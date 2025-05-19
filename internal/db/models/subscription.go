package models

type Subscription struct {
	ID        string `gorm:"primaryKey;default:uuid_generate_v4()"`
	Frequency string `gorm:"text;not null;index"`
	UserID    string `gorm:"text;not null"`
	User      User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type SubscriptionType string

const (
	DAILY  SubscriptionType = "DAILY"
	HOURLY SubscriptionType = "HOURLY"
)

type TokenType string

const (
	Sub   TokenType = "subscribe"
	Unsub TokenType = "unsubscribe"
)
