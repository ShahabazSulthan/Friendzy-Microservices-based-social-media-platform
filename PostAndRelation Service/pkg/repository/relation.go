package repository

import (
	"context"
	"database/sql"
	"fmt"

	interface_repo "github.com/ShahabazSulthan/friendzy_post/pkg/repository/interface"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type RelationRepo struct {
	DB    *gorm.DB
	Redis *redis.Client
	Ctx   context.Context
}

func NewRelationRepo(db *gorm.DB, redis *redis.Client, ctx context.Context) interface_repo.IRelationRepo {
	return &RelationRepo{DB: db, Redis: redis, Ctx: ctx}
}

func (r *RelationRepo) GetFollowerAndFollowingCountofUser(userId *string) (*uint, *uint, *error) {
	var count struct {
		FollowersCount uint `gorm:"column:followers_count"`
		FollowingCount uint `gorm:"column:following_count"`
	}

	query := "SELECT (SELECT COUNT(*) FROM relationships WHERE following_id=$1 AND relation_type=$2) AS followers_count, (SELECT COUNT(*) FROM relationships WHERE follower_id = $1 AND relation_type=$2) AS following_count"
	err := r.DB.Raw(query, userId, "follows").Scan(&count).Error
	if err != nil {
		return nil, nil, &err
	}
	return &count.FollowersCount, &count.FollowingCount, nil
}

func (r *RelationRepo) InitiateFollowRelationship(userId, userBId *string) (bool, error) {
	var inserted bool

	query := "INSERT INTO relationships (follower_id, following_id) VALUES ($1, $2) ON CONFLICT (follower_id, following_id) DO NOTHING RETURNING true;"
	err := r.DB.Raw(query, userId, userBId).Scan(&inserted).Error
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return inserted, nil
}

func (r *RelationRepo) InitiateUnFollowRelationship(userId, userBId *string) error {

	// Dereference the pointers before using them in the query
	query := "DELETE FROM relationships WHERE follower_id=$1 AND following_id=$2 AND relation_type=$3"
	err := r.DB.Exec(query, *userId, *userBId, "follows").Error
	if err != nil {
		return fmt.Errorf("failed to execute unfollow relationship query: %v", err)
	}

	return nil
}

func (r *RelationRepo) GetFollowersIdsOfUser(userId *string) (*[]uint64, error) {
	var userIds []uint64 = []uint64{}

	query := "SELECT follower_id FROM relationships WHERE following_id=$1 AND relation_type=$2"
	err := r.DB.Raw(query, userId, "follows").Scan(&userIds).Error
	if err != nil {
		return nil, err
	}
	return &userIds, nil
}

func (r *RelationRepo) GetFollowingsIdsOfUser(userId *string) (*[]uint64, error) {
	var userIds []uint64 = []uint64{}

	query := "SELECT following_id FROM relationships WHERE follower_id=$1 AND relation_type=$2"
	err := r.DB.Raw(query, userId, "follows").Scan(&userIds).Error
	if err != nil {
		return nil, err
	}
	return &userIds, nil
}

func (r *RelationRepo) UserAFollowingUserBorNot(userId, userBId *string) (bool, error) {
	var count uint

	query := "SELECT COUNT(*) FROM relationships WHERE follower_id = ? AND following_id = ? AND relation_type=?"
	err := r.DB.Raw(query, userId, userBId, "follows").Scan(&count).Error
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil

}
