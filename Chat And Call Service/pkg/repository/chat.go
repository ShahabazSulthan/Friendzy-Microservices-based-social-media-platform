package repository

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ShahabazSulthan/friendzy_post/pkg/db"
	"github.com/ShahabazSulthan/friendzy_post/pkg/model/requestmodel"
	"github.com/ShahabazSulthan/friendzy_post/pkg/model/responsemodel"
	interface_chat "github.com/ShahabazSulthan/friendzy_post/pkg/repository/interface"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepo struct {
	MongoCollections db.MongoDbCollection
	LocationInd      *time.Location
}

func NewChatRepo(db db.MongoDbCollection) interface_chat.IChatRepo {
	locationInd, _ := time.LoadLocation("Asia/Kolkata")
	return &ChatRepo{
		MongoCollections: db,
		LocationInd:      locationInd,
	}
}

func (c *ChatRepo) StoreOneToOneChatToDB(chatData *requestmodel.OneToOneChatRequest) (*string, error) {
	// Set up a context with a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to insert chat data into the OneToOneChats collection
	insertResult, err := c.MongoCollections.OneToOneChats.InsertOne(ctx, chatData)
	if err != nil {
		log.Printf("Error inserting one-to-one chat data into MongoDB: %v", err)
		return nil, fmt.Errorf("failed to store chat data: %w", err)
	}

	// Extract and convert the inserted ID to a hex string
	messageID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println("Failed to convert inserted ID to ObjectID")
		return nil, fmt.Errorf("inserted ID is not an ObjectID")
	}
	hexMessageID := messageID.Hex()

	return &hexMessageID, nil
}


func (c *ChatRepo) UpdateChatStatus(senderId, recipientId *string) error {
	// Create a context with a 10-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define the filter to match chats where the recipient and sender IDs are reversed
	filter := bson.M{
		"senderid":    *recipientId,
		"recipientid": *senderId,
	}

	// Define the update operation to set the status to "sent"
	update := bson.M{
		"$set": bson.M{
			"status": "sent",
		},
	}

	// Perform the update operation on all matching documents
	_, err := c.MongoCollections.OneToOneChats.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println("Error updating chat status in MongoDB:", err)
		return fmt.Errorf("failed to update chat status: %w", err)
	}

	return nil
}

func (c *ChatRepo) GetOneToOneChats(senderId, recipientId, limit, offset *string) (*[]responsemodel.OneToOneChatResponse, error) {
	// Create a context with a 10-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var chatSlice []responsemodel.OneToOneChatResponse

	// Define the filter to fetch chats between the sender and recipient
	filter := bson.M{
		"senderid":    bson.M{"$in": bson.A{*senderId, *recipientId}},
		"recipientid": bson.M{"$in": bson.A{*senderId, *recipientId}},
	}

	// Convert limit and offset to integers and handle potential conversion errors
	limitInt, err := strconv.Atoi(*limit)
	if err != nil {
		return nil, fmt.Errorf("invalid limit value: %w", err)
	}
	offsetInt, err := strconv.Atoi(*offset)
	if err != nil {
		return nil, fmt.Errorf("invalid offset value: %w", err)
	}

	// Define options for pagination and sorting
	options := options.Find().
		SetLimit(int64(limitInt)).
		SetSkip(int64(offsetInt)).
		SetSort(bson.D{{Key: "timestamp", Value: -1}})

	// Query the database with the filter and options
	cursor, err := c.MongoCollections.OneToOneChats.Find(ctx, filter, options)
	if err != nil {
		log.Println("Error fetching one-to-one chats:", err)
		return nil, fmt.Errorf("failed to fetch chats: %w", err)
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor to decode documents into chatSlice
	for cursor.Next(ctx) {
		var message responsemodel.OneToOneChatResponse
		if err := cursor.Decode(&message); err != nil {
			log.Println("Error decoding chat message:", err)
			return nil, fmt.Errorf("failed to decode chat message: %w", err)
		}

		// Populate additional fields
		message.MessageID = message.ID.Hex()
		message.StringTime = message.TimeStamp.In(c.LocationInd).Format("2006-01-02 15:04:05")
		chatSlice = append(chatSlice, message)
	}

	// Check for any error that may have occurred during iteration
	if err := cursor.Err(); err != nil {
		log.Println("Cursor iteration error:", err)
		return nil, fmt.Errorf("cursor error during chat fetch: %w", err)
	}

	return &chatSlice, nil
}

func (c *ChatRepo) RecentChatProfileData(senderID, limit, offset *string) (*[]responsemodel.RecentChatProfileResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Match stage: filter for chats involving the sender
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "senderid", Value: senderID}},
				bson.D{{Key: "recipientid", Value: senderID}},
			}},
		}},
	}

	// Sort stage: order by timestamp in descending order
	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{{Key: "timestamp", Value: -1}}},
	}

	// Group stage: group by the other participant and get the latest chat details
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "$cond", Value: bson.D{
					{Key: "if", Value: bson.D{{Key: "$eq", Value: bson.A{"$senderid", senderID}}}},
					{Key: "then", Value: "$recipientid"},
					{Key: "else", Value: "$senderid"},
				}},
			}},
			{Key: "lastChat", Value: bson.D{
				{Key: "$first", Value: bson.D{
					{Key: "content", Value: "$content"},
					{Key: "timestamp", Value: "$timestamp"},
					{Key: "recipientid", Value: "$recipientid"},
					{Key: "senderid", Value: "$senderid"},
				}},
			}},
		}},
	}

	// Project stage: define the fields to include in the final output
	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "participantid", Value: "$_id"},
			{Key: "content", Value: "$lastChat.content"},
			{Key: "timestamp", Value: "$lastChat.timestamp"},
			{Key: "senderid", Value: "$lastChat.senderid"},
			{Key: "recipientid", Value: "$lastChat.recipientid"},
		}},
	}

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{matchStage, sortStage, groupStage, projectStage}

	// Execute the aggregation pipeline
	cursor, err := c.MongoCollections.OneToOneChats.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Error executing aggregation:", err)
		return nil, fmt.Errorf("failed to execute aggregation: %w", err)
	}
	defer cursor.Close(ctx)

	var chatSummaries []responsemodel.RecentChatProfileResponse

	// Process each document in the cursor
	for cursor.Next(ctx) {
		var message responsemodel.RecentChatProfileResponse
		if err := cursor.Decode(&message); err != nil {
			log.Println("Error decoding chat message:", err)
			return nil, fmt.Errorf("failed to decode chat message: %w", err)
		}

		// Format timestamp and assign participant ID correctly
		message.StringTime = message.TimeStamp.In(c.LocationInd).Format("2006-01-02 15:04:05")
		if message.UserId == *senderID {
			message.UserId = message.UserIdAlt
		}
		chatSummaries = append(chatSummaries, message)
	}

	// Check for any errors encountered during iteration
	if err := cursor.Err(); err != nil {
		log.Println("Cursor iteration error:", err)
		return nil, fmt.Errorf("cursor error during chat fetch: %w", err)
	}

	return &chatSummaries, nil
}
