package models

import "time"

type Token struct {
	Token            string    `gorm:"primaryKey;default:uuid_generate_v4()"`
	Type             string    `gorm:"not null;text;uniqueIndex:uni_user_id_token_type"`
	SubscriptionType string    `gorm:"text"`
	ExpiryAt         time.Time `gorm:"not null;check:expiry_at > now()"`
	UserID           string    `gorm:"not null;text;uniqueIndex:uni_user_id_token_type;"`
	User             User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
