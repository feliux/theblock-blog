package models

type Relationship struct {
	UserId         string `bson:"userid" json:"userid"`
	UserRelationId string `bson:"userrelationid" json:"userrelationid"`
}

type ResponseRelationship struct {
	Status bool `json:"status"`
}
