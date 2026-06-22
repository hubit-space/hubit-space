package model

import "time"

type DocumentNumber struct {
	ID         int       `json:"id" gorm:"column:id"`
	Month      string    `json:"month" gorm:"column:month"`
	Year       string    `json:"year" gorm:"column:year"`
	Day        string    `json:"day" gorm:"column:day"`
	DocName    string    `json:"docname" gorm:"column:docname"`
	Format     string    `json:"format" gorm:"column:format"`
	ResetType  string    `json:"resettype" gorm:"column:resettype"`
	LastNumber int       `json:"lastnumber" gorm:"column:lastnumber"`
	IsActive   bool      `json:"is_active" gorm:"column:is_active"`
	IsDefault  bool      `json:"is_default" gorm:"column:is_default"`
	CreatedBy  string    `json:"created_by" gorm:"column:created_by"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedBy  string    `json:"updated_by" gorm:"column:updated_by"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
}
