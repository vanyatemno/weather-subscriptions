package models

import (
	"time"

	"gorm.io/gorm"
)

type TokenType = string

const (
	Sub   TokenType = "subscribe"
	Unsub TokenType = "unsubscribe"
)

const (
	unsubTokenExpiryTime = time.Hour * 24 * 28 * 13
	subTokenExpiryTime   = time.Hour * 24
)

type Token struct {
	Token            string            `gorm:"primaryKey;default:uuid_generate_v4()"`
	Type             TokenType         `gorm:"not null;text;uniqueIndex:uni_user_id_token_type"`
	SubscriptionType *SubscriptionType `gorm:"text"`
	ExpiryAt         time.Time         `gorm:"not null;check:expiry_at > now()"`
	UserID           string            `gorm:"not null;text;uniqueIndex:uni_user_id_token_type;"`
	User             User              `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt        time.Time         `gorm:"not null;check:created_at > now()"`
	UpdatedAt        time.Time         `gorm:"not null;check:updated_at > now()"`
	DeletedAt        gorm.DeletedAt
}

func (t *Token) SetTokenExpiry() {
	if t.Type == Sub {
		t.ExpiryAt = t.CreatedAt.Add(subTokenExpiryTime)
	}
	if t.Type == Unsub {
		t.ExpiryAt = t.CreatedAt.Add(unsubTokenExpiryTime)
	}
}
