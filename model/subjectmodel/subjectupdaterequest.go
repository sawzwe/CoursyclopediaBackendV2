package subjectmodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubjectUpdateRequest struct {
	SubjectCode        string               `json:"subjectCode"`
	Name               string               `json:"name"`
	Professors         []primitive.ObjectID `json:"professors"`
	SubjectDescription string               `json:"subjectDescription"`
	Campus             string               `json:"campus"`
	Credit             *int                 `json:"credit"`
	PreRequisite       *[]string            `json:"preRequisite"`
	CoRequisite        *[]string            `json:"coRequisite"`
	SubjectStatus      string               `json:"subjectStatus"`
	AvailableDuration  *int                 `json:"availableDuration"`
}
