package models

import "time"

type Subscription struct {
	ID     string `gorm:"primaryKey;default:uuid_generate_v4()"`
	UserID string `gorm:"text;not null"`

	Frequency string `gorm:"text;not null;index"`

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SubscriptionType string

const (
	DAILY  SubscriptionType = "daily"
	HOURLY SubscriptionType = "hourly"
)

type TokenType string

const (
	Sub   TokenType = "subscribe"
	Unsub TokenType = "unsubscribe"
)
