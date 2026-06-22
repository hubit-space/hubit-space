package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RuleAccess struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	RoleCode    string             `bson:"role_code" json:"role_code"`
	MenuCode   string             `bson:"menu_code" json:"menu_code"`
	RuleCode   string             `bson:"rule_code" json:"rule_code"`
}
