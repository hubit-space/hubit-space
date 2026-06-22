package model

type Option struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	DetailID    *int   `json:"detail_id" gorm:"->;column:detail_id"`
	OptionName  string `json:"option_name" gorm:"column:option_name"`
	Description string `json:"description" gorm:"column:description"`
	CreatedAt   string `json:"created_at" gorm:"column:created_at"`
	OptionType  string `json:"option_type" gorm:"column:option_type"`
	OptionLabel string `json:"option_label" gorm:"column:option_label"`
	OptionValue string `json:"option_value" gorm:"column:option_value"`
}

type OptionDetail struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	OptionType  string `json:"option_type" gorm:"column:option_type"`
	OptionLabel string `json:"option_label" gorm:"column:option_label"`
	OptionValue string `json:"option_value" gorm:"column:option_value"`
	IsActive    *bool  `json:"is_active" gorm:"column:is_active"`
}
