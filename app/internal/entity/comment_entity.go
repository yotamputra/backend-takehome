package entity

import "time"

type Comment struct {
	ID         int       `gorm:"primaryKey;autoIncrement"`
	PostID     int       `gorm:"not null;index"`
	AuthorName string    `gorm:"type:varchar(100);not null"`
	Content    string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	Post *Blog `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
