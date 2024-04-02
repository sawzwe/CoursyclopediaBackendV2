package majorrepository

import (
	"BackendCoursyclopedia/db"
	"BackendCoursyclopedia/model/majormodel"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IMajorRepository interface {
	FindAllMajors(ctx context.Context) ([]majormodel.Major, error)
	FindmajorbyID(ctx context.Context, majorId string) (*majormodel.Major, error)
	FindMajorsByIDs(ctx context.Context, majorIDs []primitive.ObjectID) ([]majormodel.Major, error)
	CreateMajor(ctx context.Context, majorName string) (string, error)
	DeleteMajor(ctx context.Context, majorId primitive.ObjectID) error
	UpdateMajor(ctx context.Context, majorId primitive.ObjectID, newName string) error
	AddSubjectToMajor(ctx context.Context, majorId string, subjectId string) error
	RemoveSubjectFromMajors(ctx context.Context, subjectId primitive.ObjectID) error
	FindMajorBySubjectId(ctx context.Context, subjectId primitive.ObjectID) (majormodel.Major, error)
	UpdatemajorforSubject(ctx context.Context, subjectId primitive.ObjectID, currentmajorId primitive.ObjectID, newmajorId primitive.ObjectID) error
}

type MajorRepository struct {
	DB *mongo.Client
}

func NewMajorRepository(db *mongo.Client) IMajorRepository {
	return &MajorRepository{
		DB: db,
	}
}

func (r MajorRepository) FindAllMajors(ctx context.Context) ([]majormodel.Major, error) {
	collection := db.GetCollection("majors")
	var majors []majormodel.Major

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var major majormodel.Major
		if err := cursor.Decode(&major); err != nil {
			return nil, err
		}
		majors = append(majors, major)
	}

	return majors, nil
}

func (r *MajorRepository) FindmajorbyID(ctx context.Context, majorId string) (*majormodel.Major, error) {
	collection := db.GetCollection("majors")
	var major majormodel.Major

	objID, err := primitive.ObjectIDFromHex(majorId)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&major)
	if err != nil {
		return nil, err
	}
	return &major, nil
}

func (r *MajorRepository) FindMajorsByIDs(ctx context.Context, majorIDs []primitive.ObjectID) ([]majormodel.Major, error) {
	collection := db.GetCollection("majors")
	var majors []majormodel.Major

	filter := bson.M{"_id": bson.M{"$in": majorIDs}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var major majormodel.Major
		if err := cursor.Decode(&major); err != nil {
			return nil, err
		}
		majors = append(majors, major)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return majors, nil
}

func (r *MajorRepository) CreateMajor(ctx context.Context, majorName string) (string, error) {
	collection := db.GetCollection("majors")
	major := majormodel.Major{
		ID:         primitive.NewObjectID(),
		MajorName:  majorName,
		SubjectIDs: []primitive.ObjectID{},
	}
	_, err := collection.InsertOne(ctx, major)
	if err != nil {
		return "", err
	}
	return major.ID.Hex(), nil
}

func (r *MajorRepository) DeleteMajor(ctx context.Context, majorId primitive.ObjectID) error {
	collection := db.GetCollection("majors")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": majorId})
	return err
}

func (r *MajorRepository) UpdateMajor(ctx context.Context, majorId primitive.ObjectID, newName string) error {
	collection := db.GetCollection("majors")
	update := bson.M{"$set": bson.M{"majorName": newName}}

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": majorId},
		update,
	)

	return err
}

func (r *MajorRepository) AddSubjectToMajor(ctx context.Context, majorId string, subjectId string) error {
	collection := db.GetCollection("majors")

	mid, err := primitive.ObjectIDFromHex(majorId)
	if err != nil {

		return err
	}

	sid, err := primitive.ObjectIDFromHex(subjectId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": mid}
	update := bson.M{"$addToSet": bson.M{"subjectIDs": sid}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *MajorRepository) RemoveSubjectFromMajors(ctx context.Context, subjectId primitive.ObjectID) error {
	collection := db.GetCollection("majors")

	_, err := collection.UpdateMany(
		ctx,
		bson.M{"subjectIDs": subjectId},
		bson.M{"$pull": bson.M{"subjectIDs": subjectId}},
	)
	return err
}

func (r *MajorRepository) FindMajorBySubjectId(ctx context.Context, subjectId primitive.ObjectID) (majormodel.Major, error) {

	collection := db.GetCollection("majors")
	var major majormodel.Major

	filter := bson.M{"subjectIDs": subjectId}
	err := collection.FindOne(ctx, filter).Decode(&major)
	if err != nil {
		return majormodel.Major{}, err
	}
	return major, nil
}

func (r *MajorRepository) UpdatemajorforSubject(ctx context.Context, subjectId primitive.ObjectID, currentmajorId primitive.ObjectID, newmajorId primitive.ObjectID) error {
	collection := db.GetCollection("majors")

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": currentmajorId},
		bson.M{"$pull": bson.M{"subjectIDs": subjectId}},
	)
	if err != nil {
		return err
	}
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": newmajorId},
		bson.M{"$addToSet": bson.M{"subjectIDs": subjectId}},
	)
	return err
}
