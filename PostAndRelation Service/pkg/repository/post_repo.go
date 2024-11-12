package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	requestmodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/requestModel"
	responsemodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/responseModel"
	interface_repo "github.com/ShahabazSulthan/friendzy_post/pkg/repository/interface"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type PostRepo struct {
	DB    *gorm.DB
	Redis *redis.Client
	Ctx   context.Context
}

func NewPostRepo(db *gorm.DB, redis *redis.Client, ctx context.Context) interface_repo.IPostRepo {
	return &PostRepo{DB: db, Redis: redis, Ctx: ctx}
}

func (p *PostRepo) AddNewPost(postData *requestmodel.AddPostData) error {
	var posdt_id string

	// Start a transaction
	tx := p.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Insert post
	query := "INSERT INTO posts (user_id, caption, created_at) VALUES ($1, $2, $3) RETURNING post_id"
	err := tx.Raw(query, postData.UserId, postData.Caption, time.Now()).Scan(&posdt_id).Error
	if err != nil {
		tx.Rollback() // Rollback if there's an error
		fmt.Println("Error in Add New Post:", err)
		return err
	}

	// Insert post media
	mediaInsQuery := "INSERT INTO post_media (post_id, media_url) VALUES ($1, $2)"
	for _, url := range postData.MediaURLs {
		errIns := tx.Exec(mediaInsQuery, posdt_id, url).Error
		if errIns != nil {
			tx.Rollback() // Rollback if there's an error
			fmt.Println("Error in inserting post media:", errIns)
			return errIns
		}
	}

	// Commit transaction if everything is successful
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (p *PostRepo) GetAllActivePostByUser(userId, limit, offset *string) (*[]responsemodel.PostData, error) {
	var response []responsemodel.PostData

	query := `SELECT posts.post_id, posts.user_id, posts.caption, post_media.media_url 
	          FROM posts 
	          LEFT JOIN post_media ON posts.post_id = post_media.post_id 
	          WHERE posts.user_id=$1 AND posts.post_status = $2 
	          ORDER BY posts.created_at DESC LIMIT $3 OFFSET $4;`

	err := p.DB.Raw(query, *userId, "normal", *limit, *offset).Scan(&response).Error
	if err != nil {
		fmt.Println("Error in GetAllActivePostByUser:", err)
		return nil, err
	}

	return &response, nil
}

func (p *PostRepo) CalculatePostAge(postID int) (string, error) {
	var age responsemodel.PostAge

	query := `
		SELECT 
			EXTRACT(EPOCH FROM NOW() - created_at) / 60 AS age_minutes,
			EXTRACT(EPOCH FROM NOW() - created_at) / 3600 AS age_hours,
			EXTRACT(EPOCH FROM NOW() - created_at) / 86400 AS age_days,
			EXTRACT(EPOCH FROM NOW() - created_at) / 604800 AS age_weeks,
			EXTRACT(EPOCH FROM NOW() - created_at) / 2592000 AS age_months,
			EXTRACT(EPOCH FROM NOW() - created_at) / 31536000 AS age_years
		FROM posts 
		WHERE post_id = $1;`

	err := p.DB.Raw(query, postID).Scan(&age).Error
	if err != nil {
		fmt.Println("Error in CalculatePostAge:", err)
		return "", err
	}

	if int(age.AgeYears) > 0 {
		return fmt.Sprintf("%d year(s) ago", int(age.AgeYears)), nil
	} else if int(age.AgeMonths) > 0 {
		return fmt.Sprintf("%d month(s) ago", int(age.AgeMonths)), nil
	} else if int(age.AgeWeeks) > 0 {
		return fmt.Sprintf("%d week(s) ago", int(age.AgeWeeks)), nil
	} else if int(age.AgeDays) > 0 {
		return fmt.Sprintf("%d day(s) ago", int(age.AgeDays)), nil
	} else if int(age.AgeHours) > 0 {
		return fmt.Sprintf("%d hour(s) ago", int(age.AgeHours)), nil
	} else if int(age.AgeMinutes) > 0 {
		return fmt.Sprintf("%d minute(s) ago", int(age.AgeMinutes)), nil
	} else {
		return "Just now", nil
	}
}

func (p *PostRepo) GetPostMediaById(postId *string) (*[]string, error) {
	var res []string

	query := "SELECT media_url FROM post_media WHERE post_id=$1 ORDER BY media_id DESC"
	err := p.DB.Raw(query, *postId).Scan(&res).Error
	if err != nil {
		fmt.Println("Error in  Get Post Media")
		return &res, err
	}

	return &res, nil
}

func (p *PostRepo) DeletePostById(postId, userId *string) error {
	query := "DELETE FROM posts WHERE post_id=$1 AND user_id=$2"
	res := p.DB.Exec(query, postId, userId)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("enter a valid post id, rows affected 0")
	}

	return nil
}

func (p *PostRepo) DeletePostMedias(postId *string) error {
	query := "DELETE FROM post_media WHERE post_id=$1"
	res := p.DB.Exec(query, postId).Error
	if res != nil {
		fmt.Println("Error in Delete Post")
		return res
	}
	return nil
}

func (p *PostRepo) EditPost(inputData *requestmodel.EditPost) error {
	query := "UPDATE posts SET caption=$1 WHERE post_id=$2 AND user_id=$3"
	res := p.DB.Exec(query, inputData.Caption, inputData.PostId, inputData.UserId)
	if res.RowsAffected == 0 {
		return errors.New("enter valid post id")
	}
	return nil
}

func (p *PostRepo) GetPostCountOfUser(userId *string) (*uint, *error) {
	var count uint
	query := "SELECT COUNT(*) FROM posts WHERE user_id=$1 AND post_status=$2"
	err := p.DB.Raw(query, userId, "normal").Scan(&count).Error
	if err != nil {
		fmt.Println("Error in get post count")
		return nil, &err
	}
	return &count, nil
}

func (p *PostRepo) LikePost(postId, userId *string) (bool, error) {
	var inserted bool
	query := "INSERT INTO post_likes(user_id, post_id, created_at) VALUES (?, ?, ?) ON CONFLICT (user_id, post_id) DO NOTHING RETURNING (xmax = 0);"
	err := p.DB.Raw(query, userId, postId, time.Now()).Scan(&inserted).Error
	if err != nil && err != sql.ErrNoRows {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return false, errors.New("enter a valid postid: Post not found")
		}
		return false, err
	}

	return inserted, nil
}

func (p PostRepo) UnLikePost(postId, userId *string) error {
	query := "DELETE FROM post_likes WHERE user_id =? AND post_id=?"
	err := p.DB.Exec(query, userId, postId).Error
	if err != nil {
		return nil
	}

	return nil
}

func (p *PostRepo) RemovePostLikesByPostId(postId *string) error {
	query := "DELETE FROM post_likes WHERE post_id=$1"
	err := p.DB.Exec(query, postId).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *PostRepo) GetPostCreatorId(postId *string) (*string, error) {
	var id string
	query := "SELECT user_id FROM posts WHERE post_id=?"
	result := p.DB.Raw(query, postId).Scan(&id)
	if result.RowsAffected == 0 {
		return nil, errors.New("no post found with this id,enter a valid postid")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &id, nil
}

func (p PostRepo) GetAllActiveRelatedPostsForHomeScreen(userId, limit, offset *string) (*[]responsemodel.PostData, error) {
	var response []responsemodel.PostData
	query := `
		SELECT posts.*, 
		CASE WHEN post_likes.user_id IS NULL THEN FALSE ELSE TRUE END AS is_liked 
		FROM posts 
		LEFT JOIN relationships AS r1 ON posts.user_id = r1.following_id AND r1.follower_id = $1
		LEFT JOIN relationships AS r2 ON posts.user_id = r2.follower_id AND r2.following_id = $1
		LEFT JOIN (SELECT post_id, user_id FROM post_likes WHERE user_id = $1) AS post_likes ON posts.post_id = post_likes.post_id  
		WHERE (r1.follower_id IS NOT NULL OR r2.following_id IS NOT NULL) 
		AND posts.post_status = 'normal' 
		ORDER BY posts.created_at DESC 
		LIMIT $2 OFFSET $3
	`

	err := p.DB.Raw(query, userId, limit, offset).Scan(&response).Error
	if err != nil {
		return &response, err
	}

	return &response, nil
}

func (p *PostRepo) GetPostLikeAndCommentsCount(postId *string) (*responsemodel.LikeCommentCounts, error) {
	//cacheKey := fmt.Sprintf("postLikeCommentCounts:%s", *postId)

	// Try to fetch cached data
	var result responsemodel.LikeCommentCounts
	// err := cache.CacheGet(p.Ctx, p.Redis, cacheKey, &result)
	// if err == nil {
	// 	return &result, nil // Cache hit
	// }

	// Cache miss - query database
	query := `SELECT 
                 (SELECT COUNT(*) FROM post_likes WHERE post_id = $1) AS likes_count, 
                 (SELECT COUNT(*) FROM comments WHERE post_id = $1 AND parent_comment_id = 0) AS comments_count`
	err := p.DB.Raw(query, postId).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	// Cache the result with a TTL
	//cache.CacheSet(p.Ctx, p.Redis, cacheKey, result, 1*time.Minute)

	return &result, nil
}

func (p *PostRepo) GetMostLovedPostsFromGlobalUser(userId, limit, offset string) (*[]responsemodel.PostData, error) {
	var response []responsemodel.PostData

	// SQL query to retrieve the most-loved posts
	query := `
		SELECT 
			posts.*, 
			COUNT(post_likes.post_id) AS like_count, 
			EXISTS (
				SELECT 1 
				FROM post_likes 
				WHERE post_likes.post_id = posts.post_id AND post_likes.user_id = ?
			) AS is_liked
		FROM posts 
		LEFT JOIN post_likes ON posts.post_id = post_likes.post_id 
		WHERE posts.post_status = 'normal' 
		GROUP BY posts.post_id 
		ORDER BY like_count DESC, posts.created_at DESC 
		LIMIT ? OFFSET ?
	`

	// Parse limit and offset as integers for the query
	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	// Execute query and scan results
	err := p.DB.Raw(query, userId, limitInt, offsetInt).Scan(&response).Error
	if err != nil {
		return nil, err
	}

	return &response, nil
}


func (p *PostRepo) GetRandomPosts(limit, offset string) (*[]responsemodel.PostData, error) {
	var response []responsemodel.PostData
	query := `
        SELECT posts.*
        FROM posts 
        WHERE posts.post_status = 'normal' 
        ORDER BY RANDOM() 
        LIMIT ? OFFSET ?`

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)
	err := p.DB.Raw(query, limitInt, offsetInt).Scan(&response).Error
	if err != nil {
		return nil, err
	}

	return &response, nil
}
