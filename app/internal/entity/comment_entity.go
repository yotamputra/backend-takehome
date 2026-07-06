package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID         uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	PostID     uuid.UUID `gorm:"type:varchar(36);not null;index"`
	AuthorName string    `gorm:"type:varchar(100);not null"`
	Content    string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	Post *Blog `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}
