package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	requestmodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/requestModel"
	responsemodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/responseModel"
	interface_repo "github.com/ShahabazSulthan/friendzy_post/pkg/repository/interface"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type CommentRepo struct {
	DB    *gorm.DB
	Redis *redis.Client
	Ctx   context.Context
}

func NewCommentRepo(db *gorm.DB, redis *redis.Client, ctx context.Context) interface_repo.ICommentRepo {
	return &CommentRepo{
		DB:    db,
		Redis: redis,
		Ctx:   ctx,
	}
}

func (d *CommentRepo) CheckingCommentHierarchy(input *uint64) (bool, error) {
	var parentCommentId uint64
	query := "SELECT parent_comment_id FROM comments WHERE comment_id=?"
	err := d.DB.Raw(query, input).Scan(&parentCommentId)

	if err.Error != nil {
		return true, err.Error
	}
	if err.RowsAffected == 0 {
		return true, errors.New("no parent-comment found in this id,enter a valid parent-comment-id")
	}

	if parentCommentId != 0 {
		return true, nil
	}

	return false, nil

}

func (c *CommentRepo) AddComment(input *requestmodel.CommentRequest) error {
	
	query := "INSERT INTO comments (post_id,user_id,comment_text,created_at) VALUES ($1,$2,$3,$4)"
	if input.ParentCommentId != 0 {
		query = "INSERT INTO comments (post_id,user_id,comment_text,parent_comment_id,created_at) VALUES ($1,$2,$3,$4,$5)"
	}

	var err error
	if input.ParentCommentId == 0 {
		err = c.DB.Exec(query, input.PostId, input.UserId, input.CommentText, time.Now()).Error
	} else {
		err = c.DB.Exec(query, input.PostId, input.UserId, input.CommentText, input.ParentCommentId, time.Now()).Error
	}

	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return errors.New("foreign key constraint violation: Post not found")
		}

		return err
	}
	return nil
}

func (c *CommentRepo) DeleteCommentAndReturnIsParentStat(userId, commentId *string) (bool, error) {
	var parentCommentId uint64
	query := "DELETE FROM comments WHERE user_id=$1 AND comment_id=$2 RETURNING parent_comment_id"
	result := c.DB.Raw(query, userId, commentId).Scan(&parentCommentId)
	if result.RowsAffected == 0 {
		return false, errors.New("no comment found with this id ,enter a valid commentid")
	}
	if result.Error != nil {
		return false, result.Error
	}

	if parentCommentId != 0 {
		return false, nil
	}

	return true, nil
}

func (c *CommentRepo) DeleteChildComments(parentCommentId *string) error {
	query := "DELETE FROM comments WHERE parent_comment_id=?"
	err := c.DB.Exec(query, parentCommentId).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *CommentRepo) EditComment(userId, commentText *string, commentId *uint64) error {
	query := "UPDATE comments SET comment_text=$1 WHERE user_id=$2 AND comment_id=$3"
	result := c.DB.Exec(query, commentText, userId, commentId)
	if result.RowsAffected == 0 {
		return errors.New("no comment found with this id,enter a valid commentId ")
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (c *CommentRepo) FetchParentCommentsOfPost(userId, postId, limit, offset *string) (*[]responsemodel.ParentComments, error) {
	var ScannerStruct []responsemodel.ParentComments

	query := "SELECT * FROM comments WHERE post_id=$1 AND parent_comment_id=0 LIMIT $2 OFFSET $3"
	result := c.DB.Raw(query, postId, limit, offset).Scan(&ScannerStruct).Error
	if result != nil {
		return nil, result
	}
	return &ScannerStruct, nil
}

func (c *CommentRepo) FetchChildCommentsOfComment(parentCommentId *uint) (*[]responsemodel.ChildComments, error) {
	var ScannerStruct []responsemodel.ChildComments
	query := "SELECT * FROM comments WHERE parent_comment_id=$1"
	err := c.DB.Raw(query, parentCommentId).Scan(&ScannerStruct).Error
	if err != nil {
		return nil, err
	}
	return &ScannerStruct, nil
}

func (c *CommentRepo) FindCommentCreatorId(CommentId *uint64) (*string, error) {
	var parentId string

	query := "SELECT user_id FROM comments WHERE comment_id=?"
	err := c.DB.Raw(query, CommentId).Scan(&parentId).Error
	if err != nil {
		log.Println("--------", err)
		return nil, err
	}
	return &parentId, nil
}

func (p *CommentRepo) CalculateCommntAge(CommentsID int) (string, error) {
	var age responsemodel.PostAge

	query := `
		SELECT 
			EXTRACT(EPOCH FROM NOW() - created_at) / 60 AS age_minutes,
			EXTRACT(EPOCH FROM NOW() - created_at) / 3600 AS age_hours,
			EXTRACT(EPOCH FROM NOW() - created_at) / 86400 AS age_days,
			EXTRACT(EPOCH FROM NOW() - created_at) / 604800 AS age_weeks,
			EXTRACT(EPOCH FROM NOW() - created_at) / 2592000 AS age_months,
			EXTRACT(EPOCH FROM NOW() - created_at) / 31536000 AS age_years
		FROM comments 
		WHERE comment_id = $1;`

	err := p.DB.Raw(query, CommentsID).Scan(&age).Error
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
