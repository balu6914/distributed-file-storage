package models

import (
	_ "gorm.io/gorm"
)

type FilePart struct {
	ID     uint   `gorm:"primaryKey"`
	FileID string `gorm:"index"`
	Part   int
	Data   []byte
}
