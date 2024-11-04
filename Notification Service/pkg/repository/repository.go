package repository

import (
	"fmt"
	"time"

	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/models/requestmodel"
	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/models/responsemodel"
	interface_notification "github.com/ShahabazSulthan/Friendzy_Notification/pkg/repository/interface"
	"gorm.io/gorm"
)

type NotifRepo struct {
	DB *gorm.DB
}

func NewNotifRepo(db *gorm.DB) interface_notification.INotifRepo {
	return &NotifRepo{DB: db}
}

func (n *NotifRepo) CreateNewNotification(msg *requestmodel.KafkaNotification) error {
	query := "INSERT INTO notifications (user_id,actor_id,action_type,target_id,target_type,comment_text,created_at) VALUES($1,$2,$3,$4,$5,$6,$7)"
	err := n.DB.Exec(query, msg.UserID, msg.ActorID, msg.ActionType, msg.TargetID, msg.TargetType, msg.CommentText, msg.CreatedAt).Error
	if err != nil {
		return err
	}
	return nil
}

func (n *NotifRepo) GetNotificationsForUser(userId, limit, offset *string) (*[]responsemodel.NotificationModel, error) {
	var respModel []responsemodel.NotificationModel

	query := "SELECT * FROM notifications WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"
	err := n.DB.Raw(query, userId, limit, offset).Scan(&respModel).Error
	if err != nil {
		return nil, err
	}
	return &respModel, nil
}

func (n *NotifRepo) CalculatePostAge(createdAt time.Time) (string, error) {
	var age responsemodel.PostAge

	query := `
		SELECT 
			EXTRACT(EPOCH FROM NOW() - $1) / 60 AS age_minutes,
			EXTRACT(EPOCH FROM NOW() - $1) / 3600 AS age_hours,
			EXTRACT(EPOCH FROM NOW() - $1) / 86400 AS age_days,
			EXTRACT(EPOCH FROM NOW() - $1) / 604800 AS age_weeks,
			EXTRACT(EPOCH FROM NOW() - $1) / 2592000 AS age_months,
			EXTRACT(EPOCH FROM NOW() - $1) / 31536000 AS age_years
		FROM notifications;`

	err := n.DB.Raw(query, createdAt).Scan(&age).Error
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
