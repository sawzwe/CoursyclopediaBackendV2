package facultyrepository

import (
	"BackendCoursyclopedia/db"
	"BackendCoursyclopedia/model/facultymodel"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IFacultyRepository interface {
	FindAllFaculties(ctx context.Context) ([]facultymodel.Faculty, error)
	FindFacultyByID(ctx context.Context, facultyID string) (*facultymodel.Faculty, error)
	CreateFaculty(ctx context.Context, facultyName string, image []byte) (facultymodel.Faculty, error)
	UpdateFaculty(ctx context.Context, facultyID string, faculty facultymodel.Faculty, image []byte) (facultymodel.Faculty, error)
	DeleteFaculty(ctx context.Context, facultyID string) error
	AddMajorToFaculty(ctx context.Context, facultyId string, majorId string) error
	RemoveMajorFromFaculty(ctx context.Context, majorId primitive.ObjectID) error
	FindFacultyByMajorId(ctx context.Context, majorId primitive.ObjectID) (facultymodel.Faculty, error)
	UpdateFacultyForMajor(ctx context.Context, majorId primitive.ObjectID, currentFacultyId primitive.ObjectID, newFacultyId primitive.ObjectID) error
}

type FacultyRepository struct {
	DB *mongo.Client
}

func NewFacultyRepository(db *mongo.Client) IFacultyRepository {
	return &FacultyRepository{
		DB: db,
	}
}

func (r FacultyRepository) FindAllFaculties(ctx context.Context) ([]facultymodel.Faculty, error) {
	collection := db.GetCollection("faculties")
	var faculties []facultymodel.Faculty

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var faculty facultymodel.Faculty
		if err := cursor.Decode(&faculty); err != nil {
			return nil, err
		}
		faculties = append(faculties, faculty)
	}

	return faculties, nil
}

func (r *FacultyRepository) FindFacultyByID(ctx context.Context, facultyID string) (*facultymodel.Faculty, error) {
	collection := db.GetCollection("faculties")
	var faculty facultymodel.Faculty

	objID, err := primitive.ObjectIDFromHex(facultyID)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&faculty)
	if err != nil {
		return nil, err
	}
	return &faculty, nil
}

func (r FacultyRepository) CreateFaculty(ctx context.Context, facultyName string, image []byte) (facultymodel.Faculty, error) {
	collection := db.GetCollection("faculties")
	Faculty := facultymodel.Faculty{
		ID:          primitive.NewObjectID(),
		FacultyName: facultyName,
		Image:       image,
		MajorIDs:    []primitive.ObjectID{},
	}
	_, err := collection.InsertOne(ctx, Faculty)
	if err != nil {
		return facultymodel.Faculty{}, err
	}

	return Faculty, nil
}

func (r FacultyRepository) UpdateFaculty(ctx context.Context, facultyID string, faculty facultymodel.Faculty, image []byte) (facultymodel.Faculty, error) {
	collection := db.GetCollection("faculties")
	objID, err := primitive.ObjectIDFromHex(facultyID)
	if err != nil {
		return facultymodel.Faculty{}, err
	}

	updateData := bson.M{"$set": faculty}
	if image != nil {
		updateData["$set"].(bson.M)["image"] = image
	}

	filter := bson.M{"_id": objID}
	result, err := collection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		return facultymodel.Faculty{}, err
	}
	if result.MatchedCount == 0 {
		return facultymodel.Faculty{}, errors.New("no faculty found with given ID")
	}

	return faculty, nil
}

func (r FacultyRepository) DeleteFaculty(ctx context.Context, facultyID string) error {
	collection := db.GetCollection("faculties")
	objID, err := primitive.ObjectIDFromHex(facultyID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no faculty found with given ID")
	}

	return nil
}

func (r *FacultyRepository) AddMajorToFaculty(ctx context.Context, facultyId string, majorId string) error {
	collection := db.GetCollection("faculties")

	fid, err := primitive.ObjectIDFromHex(facultyId)
	if err != nil {
		return err
	}
	mid, err := primitive.ObjectIDFromHex(majorId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": fid}
	update := bson.M{"$addToSet": bson.M{"majorIDs": mid}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *FacultyRepository) RemoveMajorFromFaculty(ctx context.Context, majorId primitive.ObjectID) error {
	collection := db.GetCollection("faculties")

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"majorIDs": majorId},
		bson.M{"$pull": bson.M{"majorIDs": majorId}},
	)
	return err
}

func (r *FacultyRepository) FindFacultyByMajorId(ctx context.Context, majorId primitive.ObjectID) (facultymodel.Faculty, error) {
	collection := db.GetCollection("faculties")
	var faculty facultymodel.Faculty

	filter := bson.M{"majorIDs": majorId}
	err := collection.FindOne(ctx, filter).Decode(&faculty)
	if err != nil {
		return facultymodel.Faculty{}, err
	}

	return faculty, nil
}

func (r *FacultyRepository) UpdateFacultyForMajor(ctx context.Context, majorId primitive.ObjectID, currentFacultyId primitive.ObjectID, newFacultyId primitive.ObjectID) error {
	collection := db.GetCollection("faculties")

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": currentFacultyId},
		bson.M{"$pull": bson.M{"majorIDs": majorId}},
	)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": newFacultyId},
		bson.M{"$addToSet": bson.M{"majorIDs": majorId}},
	)
	return err
}
