package links

import "time"

type Links struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	OriginalUrl string    `gorm:"not null"`
	ShortUrl    string    `gorm:"not null;unique"`
	ExpiresAt   time.Time `gorm:"not null"`
}
