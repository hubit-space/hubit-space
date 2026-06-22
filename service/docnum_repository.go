package model

type Parameter struct {
	ID             int    `json:"id" gorm:"primaryKey;column:id"`
	ParameterCode  string `json:"parameter_code" gorm:"column:parameter_code"`
	ParameterValue string `json:"parameter_value" gorm:"column:parameter_value"`
	Description    string `json:"description" gorm:"column:description"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
