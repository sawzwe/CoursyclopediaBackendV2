package majorservice

import (
	"BackendCoursyclopedia/model/majormodel"
	"BackendCoursyclopedia/model/subjectmodel"
	"BackendCoursyclopedia/repository/facultyrepository"
	majorrepo "BackendCoursyclopedia/repository/majorrepository"
	"BackendCoursyclopedia/repository/subjectrepository"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IMajorService interface {
	GetAllMajors(ctx context.Context) ([]majormodel.Major, error)
	GetMajorByID(ctx context.Context, majorID string) (*majormodel.Major, error)
	GetSubjectsForMajor(ctx context.Context, majorId string) ([]subjectmodel.Subject, error)
	CreateMajor(majorName string, facultyId string) error
	DeleteMajor(majorId string) error
	UpdateMajor(ctx context.Context, majorId string, newMajorName string, newFacultyId string) error
}

type MajorService struct {
	MajorRepository   majorrepo.IMajorRepository
	FacultyRepository facultyrepository.IFacultyRepository
	SubjectRepository subjectrepository.ISubjectRepository
}

func NewMajorService(MajorRepo majorrepo.IMajorRepository, FacultyRepo facultyrepository.IFacultyRepository, SubjectRepo subjectrepository.ISubjectRepository) IMajorService {
	return &MajorService{
		MajorRepository:   MajorRepo,
		FacultyRepository: FacultyRepo,
		SubjectRepository: SubjectRepo,
	}
}

func (s MajorService) GetAllMajors(ctx context.Context) ([]majormodel.Major, error) {
	return s.MajorRepository.FindAllMajors(ctx)
}

func (s *MajorService) GetMajorByID(ctx context.Context, majorID string) (*majormodel.Major, error) {
	return s.MajorRepository.FindmajorbyID(ctx, majorID)
}

func (s *MajorService) GetSubjectsForMajor(ctx context.Context, majorId string) ([]subjectmodel.Subject, error) {
	major, err := s.MajorRepository.FindmajorbyID(ctx, majorId)
	if err != nil {
		return nil, err
	}

	subjects, err := s.SubjectRepository.FindSubjectsByIDs(ctx, major.SubjectIDs)
	if err != nil {
		return nil, err
	}

	return subjects, nil
}
func (s *MajorService) CreateMajor(majorName string, facultyId string) error {
	ctx := context.Background()

	majorId, err := s.MajorRepository.CreateMajor(ctx, majorName)
	if err != nil {
		return err
	}

	return s.FacultyRepository.AddMajorToFaculty(ctx, facultyId, majorId)
}

func (s *MajorService) DeleteMajor(majorId string) error {
	ctx := context.Background()

	objId, err := primitive.ObjectIDFromHex(majorId)
	if err != nil {
		return err
	}

	err = s.MajorRepository.DeleteMajor(ctx, objId)
	if err != nil {
		return err
	}

	return s.FacultyRepository.RemoveMajorFromFaculty(ctx, objId)
}

func (s *MajorService) UpdateMajor(ctx context.Context, majorId string, newMajorName string, newFacultyId string) error {
	majorObjId, err := primitive.ObjectIDFromHex(majorId)
	if err != nil {
		return err
	}

	if newMajorName != "" {
		err = s.MajorRepository.UpdateMajor(ctx, majorObjId, newMajorName)
		if err != nil {
			return err
		}
	}

	if newFacultyId != "" {
		newFacObjId, err := primitive.ObjectIDFromHex(newFacultyId)
		if err != nil {
			return err
		}

		currentFaculty, err := s.FacultyRepository.FindFacultyByMajorId(ctx, majorObjId)
		if err != nil {
			return err
		}

		if currentFaculty.ID != newFacObjId {
			err = s.FacultyRepository.UpdateFacultyForMajor(ctx, majorObjId, currentFaculty.ID, newFacObjId)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
