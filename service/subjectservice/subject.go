package subjectservice

import (
	"BackendCoursyclopedia/model/subjectmodel"
	"BackendCoursyclopedia/repository/majorrepository"
	"BackendCoursyclopedia/repository/subjectrepository"
	"context"
	"log"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"time"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISubjectService interface {
	GetAllSubjects(ctx context.Context) ([]subjectmodel.Subject, error)
	GetSubjectByID(ctx context.Context, subjectID string) (*subjectmodel.Subject, error)
	CreateSubject(ctx context.Context, subject subjectmodel.Subject, majorId string) (string, error)
	DeleteSubject(subjectId string) error
	UpdateSubject(ctx context.Context, subjectId string, updates subjectmodel.SubjectUpdateRequest, newMajorId string) error
	UpdateLikes(ctx context.Context, subjectID string, likes int) error
	AddLikeByEmail(ctx context.Context, subjectID string, userEmail string) error
}

type SubjectService struct {
	SubjectRepository subjectrepository.ISubjectRepository
	MajorRepository   majorrepository.IMajorRepository
}

func NewSubjectService(SubjectRepo subjectrepository.ISubjectRepository, MajorRepo majorrepository.IMajorRepository) ISubjectService {
	return &SubjectService{
		SubjectRepository: SubjectRepo,
		MajorRepository:   MajorRepo,
	}
}

func (s SubjectService) GetAllSubjects(ctx context.Context) ([]subjectmodel.Subject, error) {
	return s.SubjectRepository.FindAllSubjects(ctx)
}

func (s SubjectService) GetSubjectByID(ctx context.Context, subjectID string) (*subjectmodel.Subject, error) {
	return s.SubjectRepository.FindSubjectbyID(ctx, subjectID)
}

func (s *SubjectService) CreateSubject(ctx context.Context, subject subjectmodel.Subject, majorId string) (string, error) {
	subjectId, err := s.SubjectRepository.CreateSubject(ctx, subject)
	if err != nil {
		return "", err
	}

	if subject.SubjectStatus == "" {
		subject.SubjectStatus = "AVAILABLE"
	}

	subjectIdHex := subjectId.Hex()

	err = s.MajorRepository.AddSubjectToMajor(ctx, majorId, subjectIdHex)
	if err != nil {
		return "", err
	}

	return subjectIdHex, nil
}

func (s *SubjectService) DeleteSubject(subjectId string) error {
	ctx := context.Background()

	objId, err := primitive.ObjectIDFromHex(subjectId)
	if err != nil {
		return err
	}

	err = s.SubjectRepository.DeleteSubject(ctx, objId)
	if err != nil {
		return err
	}

	return s.MajorRepository.RemoveSubjectFromMajors(ctx, objId)
}

func (s *SubjectService) UpdateSubject(ctx context.Context, subjectId string, updates subjectmodel.SubjectUpdateRequest, newMajorId string) error {
	subjectObjId, err := primitive.ObjectIDFromHex(subjectId)
	if err != nil {

		return err
	}

	updateFields := bson.M{}

	if updates.SubjectCode != "" {
		updateFields["subjectCode"] = updates.SubjectCode
	}
	if updates.Name != "" {
		updateFields["name"] = updates.Name
	}
	if len(updates.Professors) > 0 {
		updateFields["professors"] = updates.Professors
	}
	if updates.SubjectDescription != "" {
		updateFields["subjectDescription"] = updates.SubjectDescription
	}
	if updates.Campus != "" {
		updateFields["campus"] = updates.Campus
	}
	if updates.Credit != nil {
		updateFields["credit"] = *updates.Credit
	}
	if updates.PreRequisite != nil {
		updateFields["pre_requisite"] = updates.PreRequisite
	}
	if updates.CoRequisite != nil {
		updateFields["co_requisite"] = updates.CoRequisite
	}
	if updates.SubjectStatus != "" {
		updateFields["subjectStatus"] = updates.SubjectStatus
	}
	if updates.AvailableDuration != nil {
		updateFields["available_duration"] = *updates.AvailableDuration
	}

	if newMajorId != "" {
		newmajObjId, err := primitive.ObjectIDFromHex(newMajorId)
		if err != nil {

			return err
		}

		currentmajor, err := s.MajorRepository.FindMajorBySubjectId(ctx, subjectObjId)
		if err != nil {

			return err
		}

		if currentmajor.ID != newmajObjId {
			err = s.MajorRepository.UpdatemajorforSubject(ctx, subjectObjId, currentmajor.ID, newmajObjId)
			if err != nil {
				return err
			}
		}
	}

	if len(updateFields) > 0 {
		err = s.SubjectRepository.UpdateSubject(ctx, subjectObjId, updateFields)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SubjectService) UpdateLikes(ctx context.Context, subjectID string, likes int) error {
	id, err := primitive.ObjectIDFromHex(subjectID)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return err
	}

	err = s.SubjectRepository.UpdateLikes(ctx, id, likes)
	if err != nil {
		log.Printf("Error updating subject likes: %v", err)
		return err
	}

	return nil
}

func (s *SubjectService) AddLikeByEmail(ctx context.Context, subjectID string, userEmail string) error {
	id, err := primitive.ObjectIDFromHex(subjectID)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return err
	}

	err = s.SubjectRepository.AddEmailToLikeList(ctx, id, userEmail)
	if err != nil {
		log.Printf("Error updating subject likelist: %v", err)
		return err
	}

	return nil
}
