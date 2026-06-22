package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	RoleCode    string             `bson:"role_code" json:"role_code"`
	RoleName    string             `bson:"role_name" json:"role_name"`
	Description string             `bson:"description" json:"description"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
	CreatedOn   string             `bson:"created_on" json:"created_on"`
	UpdatedBy   string             `bson:"updated_by" json:"updated_by"`
	UpdatedOn   string             `bson:"updated_on" json:"updated_on"`
	Type        string             `bson:"type,omitempty" json:"type,omitempty"`
	Menus       []Menu             `bson:"menus" json:"menus"`
	MobileMenus []Menu             `bson:"mobile_menus" json:"mobile_menus"`
}
