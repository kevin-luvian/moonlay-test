package model

import (
	"database/sql"
	"mime/multipart"
	"time"
)

type List struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string `gorm:"size:100;not null"`
	Description string `gorm:"size:1000;not null"`
	FileUpload  string `gorm:"size:255"`
	Level0ID    sql.NullInt64
	Sublists    []List `gorm:"foreignkey:Level0ID"`

	FileHeader *multipart.FileHeader `gorm:"-" json:"-"`
}
