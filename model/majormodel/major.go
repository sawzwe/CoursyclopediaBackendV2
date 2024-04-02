package majormodel

import "go.mongodb.org/mongo-driver/bson/primitive"

type Major struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	MajorName  string               `bson:"majorName"`
	SubjectIDs []primitive.ObjectID `bson:"subjectIDs"`
}
