package subjectmodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subject struct {
	ID                 primitive.ObjectID   `bson:"_id,omitempty"`
	SubjectCode        string               `bson:"subjectCode"`
	Name               string               `bson:"name"`
	Professors         []primitive.ObjectID `bson:"professors"`
	SubjectDescription string               `bson:"subjectDescription"`
	Campus             string               `bson:"campus"`
	Credit             int                  `bson:"credit"`
	PreRequisite       []string             `bson:"pre_requisite"`
	CoRequisite        []string             `bson:"co_requisite"`
	Likes              int                  `bson:"likes"`
	Likelist           []string             `bson:"likelist"`
	SubjectStatus      string               `bson:"subjectStatus"`
	LastUpdated        primitive.DateTime   `bson:"last_updated"`
	AvailableDuration  int                  `bson:"available_duration"`
}
