package auditlogmodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditLog struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	OperationType string             `bson:"operationType"`
	Subject       primitive.ObjectID `bson:"subject"`
	OperatedBy    primitive.ObjectID `bson:"operatedBy"`
	Timestamp     primitive.DateTime `bson:"timestamp"`
	PreviousState interface{}        `bson:"previousState"`
	NewState      interface{}        `bson:"newState"`
}
