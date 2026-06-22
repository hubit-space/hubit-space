package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Right struct {
	RightCode string `bson:"right_code" json:"right_code"`
	RightName string `bson:"right_name" json:"right_name"`
	IsActive  bool   `bson:"is_active" json:"is_active"`
	CreatedBy string `bson:"created_by" json:"created_by"`
	CreatedOn string `bson:"created_on" json:"created_on"`
	UpdatedBy string `bson:"updated_by" json:"updated_by"`
	UpdatedOn string `bson:"updated_on" json:"updated_on"`
	APIUrl    string `bson:"api_url" json:"api_url"`
	PageUrl   string `bson:"page_url" json:"page_url"`
}

type Menu struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	MenuCode    string             `bson:"menu_code" json:"menu_code"`
	MenuName    string             `bson:"menu_name" json:"menu_name"`
	URL         string             `bson:"url" json:"url"`
	Icon        string             `bson:"icon" json:"icon"`
	Sort        int                `bson:"sort" json:"sort"`
	IsMobile    bool               `bson:"is_mobile" json:"is_mobile"`
	Description string             `bson:"description" json:"description"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
	CreatedOn   string             `bson:"created_on" json:"created_on"`
	UpdatedBy   string             `bson:"updated_by" json:"updated_by"`
	UpdatedOn   string             `bson:"updated_on" json:"updated_on"`
	Rights      []Right            `bson:"rights" json:"rights"`
}
