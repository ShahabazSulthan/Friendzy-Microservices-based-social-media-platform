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

func (c *ChatRepo) CreateNewGroup(groupInfo *requestmodel.NewGroupInfo) error {
	// Set a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the new group information into the Chatgroup collection
	if _, err := c.MongoCollections.Chatgroup.InsertOne(ctx, groupInfo); err != nil {
		log.Printf("Error inserting data into MongoDB Chatgroup collection: %v", err)
		return fmt.Errorf("could not create new group: %w", err)
	}
	return nil
}


func (c *ChatRepo) GetGroupMembersList(groupId *string) (*[]uint64, error) {
	objGroupID, err := primitive.ObjectIDFromHex(*groupId)
	if err != nil {
		return nil, fmt.Errorf("invalid group ID: %v", err)
	}

	fmt.Println("objectGroupId", objGroupID)

	var group struct {
		Members []uint64 `bson:"groupmembers"`
	}

	filter := bson.M{"_id": objGroupID}
	err = c.MongoCollections.Chatgroup.FindOne(context.Background(), filter).Decode(&group)
	if err != nil {
		return nil, fmt.Errorf("could not find group: %v", err)
	}

	return &group.Members, nil
}

func (c *ChatRepo) StoreOneToManyChatToDB(msg *requestmodel.OneToManyMessageRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := c.MongoCollections.Chatgroup.InsertOne(ctx, msg)
	if err != nil {
		log.Println("-----error: from chatrepo:StoreOneToManyChatToDB() failed to store chat data")
		fmt.Println("--------", err)
		return err
	}
	return nil
}

func (c *ChatRepo) GetRecentGroupProfilesOfUser(userId, limit, offset *string) (*[]responsemodel.GroupInfoLite, error) {
	userIdInt, _ := strconv.Atoi(fmt.Sprint(*userId))
	limitInt, _ := strconv.Atoi(*limit)
	offsetInt, _ := strconv.Atoi(*offset)

	filter := bson.M{"groupmembers": userIdInt}
	findOptions := options.Find()
	findOptions.SetLimit(int64(limitInt))
	findOptions.SetSkip(int64(offsetInt))

	cursor, err := c.MongoCollections.Chatgroup.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var groups []responsemodel.GroupInfoLite
	for cursor.Next(context.TODO()) {
		var group responsemodel.GroupInfoLite
		if err = cursor.Decode(&group); err != nil {
			return nil, err
		}

		group.GroupID = group.ID.Hex()
		groups = append(groups, group)
	}

	return &groups, nil
}

func (c *ChatRepo) GetGroupLastMessageDetailsByGroupId(groupid *string) (*responsemodel.OneToManyMessageLite, error) {
	filter := bson.M{"groupid": *groupid}
	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{Key: "timestamp", Value: -1}}) // Sort by timestamp in descending order

	var chat responsemodel.OneToManyMessageLite
	err := c.MongoCollections.OneToManyChats.FindOne(context.TODO(), filter, findOptions).Decode(&chat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No message found for the given groupID
		}
		return nil, err
	}

	chat.StringTime = fmt.Sprint(chat.TimeStamp.In(c.LocationInd))
	return &chat, nil
}

func (c *ChatRepo) CheckUserIsGroupMember(userid, groupid *string) (bool, error) {
	objGroupID, err := primitive.ObjectIDFromHex(*groupid)
	if err != nil {
		return false, fmt.Errorf("invalid group ID: %v", err)
	}
	userIdInt, _ := strconv.Atoi(*userid)

	filter := bson.M{
		"_id":          objGroupID,
		"groupmembers": userIdInt,
	}
	var result bson.M
	err = c.MongoCollections.Chatgroup.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *ChatRepo) GetOneToManyChats(groupId, limit, offset *string) (*[]responsemodel.OneToManyChatResponse, error) {
	var chatSlice []responsemodel.OneToManyChatResponse

	// Filter setup
	filter := bson.M{"groupid": *groupId}

	// Convert limit and offset safely
	limitInt, err := strconv.Atoi(*limit)
	if err != nil {
		return nil, fmt.Errorf("invalid limit: %v", err)
	}
	offsetInt, err := strconv.Atoi(*offset)
	if err != nil {
		return nil, fmt.Errorf("invalid offset: %v", err)
	}

	// Set options for the query
	options := options.Find().
		SetLimit(int64(limitInt)).
		SetSkip(int64(offsetInt)).
		SetSort(bson.D{{Key: "timestamp", Value: -1}})

	// Execute the find query
	cursor, err := c.MongoCollections.Chatgroup.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Iterate over the cursor
	for cursor.Next(context.TODO()) {
		var message responsemodel.OneToManyChatResponse
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		// Additional processing
		message.MessageID = message.ID.Hex()
		message.StringTime = fmt.Sprint(message.TimeStamp.In(c.LocationInd))
		chatSlice = append(chatSlice, message)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &chatSlice, nil
}


func (c *ChatRepo) AddNewMembersToGroupByGroupId(inputData *requestmodel.AddNewMemberToGroup) error {
	objGroupID, err := primitive.ObjectIDFromHex(inputData.GroupID)
	if err != nil {
		return fmt.Errorf("invalid group ID: %v", err)
	}

	filter := bson.M{"_id": objGroupID}
	var group struct {
		GroupMembers []uint64 `bson:"groupmembers"`
	}

	err = c.MongoCollections.Chatgroup.FindOne(context.TODO(), filter).Decode(&group)
	if err != nil {
		return fmt.Errorf("failed to find group: %v", err)
	}

	// Create a map to ensure uniqueness
	memberSet := make(map[uint64]struct{})
	for _, member := range group.GroupMembers {
		memberSet[member] = struct{}{}
	}

	// Add new members if they are not already present
	for _, newMember := range inputData.GroupMembers {
		if _, exists := memberSet[newMember]; !exists {
			memberSet[newMember] = struct{}{}
		}
	}

	// Convert the map back to a slice
	updatedMembers := make([]uint64, 0, len(memberSet))
	for member := range memberSet {
		updatedMembers = append(updatedMembers, member)
	}

	// Update the group document with the new members
	update := bson.M{"$set": bson.M{"groupmembers": updatedMembers}}
	_, err = c.MongoCollections.Chatgroup.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update group members: %v", err)
	}

	return nil
}

func (c *ChatRepo) RemoveGroupMember(inputData *requestmodel.RemoveMemberFromGroup) error {
	objGroupID, err := primitive.ObjectIDFromHex(inputData.GroupID)
	if err != nil {
		return fmt.Errorf("invalid group ID: %v", err)
	}
	memberIdInt, _ := strconv.Atoi(inputData.MemberID)

	filter := bson.M{"_id": objGroupID}
	update := bson.M{"$pull": bson.M{"groupmembers": memberIdInt}}

	result, err := c.MongoCollections.Chatgroup.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to remove user from group: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("group not found")
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("user with this member id is not in this group")
	}

	return nil
}

func (c *ChatRepo) CountMembersInGroup(groupId string) (int, error) {
	objGroupID, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		return 0, fmt.Errorf("invalid group ID: %v", err)
	}

	filter := bson.M{"_id": objGroupID}
	var group struct {
		GroupMembers []uint64 `bson:"groupmembers"`
	}

	err = c.MongoCollections.Chatgroup.FindOne(context.TODO(), filter).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, fmt.Errorf("group not found")
		}
		return 0, err
	}
	return len(group.GroupMembers), nil
}

func (c *ChatRepo) DeleteOneToManyChatsByGroupId(groupId string) error {
	deleteFilter := bson.M{"groupid": groupId}
	_, err := c.MongoCollections.OneToManyChats.DeleteMany(context.TODO(), deleteFilter)
	if err != nil {
		return fmt.Errorf("failed to delete group chats: %v", err)
	}

	return nil
}

func (c *ChatRepo) DeleteGroupDataByGroupId(groupId string) error {
	objGroupID, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		return fmt.Errorf("invalid group ID: %v", err)
	}
	filter := bson.M{"_id": objGroupID}
	_, err = c.MongoCollections.Chatgroup.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
