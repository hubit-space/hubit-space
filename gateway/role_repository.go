package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserCode  string             `json:"user_code" bson:"user_code"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Type      string             `json:"type" bson:"type"`
	IsActive  bool               `json:"is_active" bson:"is_active"`
	CreatedBy string             `bson:"created_by" json:"created_by"`
	CreatedOn string             `bson:"created_on" json:"created_on"`
	UpdatedBy string             `bson:"updated_by" json:"updated_by"`
	UpdatedOn string             `bson:"updated_on" json:"updated_on"`

	Role RoleUser `json:"role" bson:"role"`

	Menus []Menu `bson:"menus" json:"menus"`
}

type RoleUser struct {
	ID          string `bson:"_id,omitempty" json:"_id"`
	RoleCode    string `bson:"role_code" json:"role_code"`
	RoleName    string `bson:"role_name" json:"role_name"`
	Description string `bson:"description" json:"description"`
	IsActive    bool   `bson:"is_active" json:"is_active"`
	CreatedBy   string `bson:"created_by" json:"created_by"`
	CreatedOn   string `bson:"created_on" json:"created_on"`
	UpdatedBy   string `bson:"updated_by" json:"updated_by"`
	UpdatedOn   string `bson:"updated_on" json:"updated_on"`
}
