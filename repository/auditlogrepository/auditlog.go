package auditlogrepo

import (
	"BackendCoursyclopedia/db"
	"BackendCoursyclopedia/model/auditlogmodel"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IAuditLogRepository interface {
	FindAllAuditLogs(ctx context.Context) ([]auditlogmodel.AuditLog, error) // Removed the id parameter as it's not used
	FindAuditLogByID(ctx context.Context, auditlogId string) (*auditlogmodel.AuditLog, error)
}

type AuditLogRepository struct {
	DB *mongo.Client
}

func NewAuditLogRepository(db *mongo.Client) IAuditLogRepository {
	return &AuditLogRepository{
		DB: db,
	}
}

func (r *AuditLogRepository) FindAllAuditLogs(ctx context.Context) ([]auditlogmodel.AuditLog, error) {
	collection := db.GetCollection("auditlogs")
	var auditlogs []auditlogmodel.AuditLog
	cursor, err := collection.Find(ctx, bson.M{}) // Using an empty filter to fetch all documents
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var auditlog auditlogmodel.AuditLog
		if err := cursor.Decode(&auditlog); err != nil {
			return nil, err
		}
		auditlogs = append(auditlogs, auditlog)
	}

	return auditlogs, nil
}

func (r *AuditLogRepository) FindAuditLogByID(ctx context.Context, auditlogId string) (*auditlogmodel.AuditLog, error) {
	collection := db.GetCollection("auditlogs")

	var auditlog auditlogmodel.AuditLog

	objID, err := primitive.ObjectIDFromHex(auditlogId)
	if err != nil {
		return nil, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&auditlog)
	if err != nil {
		return nil, err
	}
	return &auditlog, nil

}
