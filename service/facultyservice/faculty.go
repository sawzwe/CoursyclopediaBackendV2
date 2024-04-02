package facultyservice

import (
	"BackendCoursyclopedia/model/facultymodel"
	"BackendCoursyclopedia/model/majormodel"
	facultyrepo "BackendCoursyclopedia/repository/facultyrepository"
	"BackendCoursyclopedia/repository/majorrepository"
	"context"
)

type IFacultyService interface {
	GetAllFaculties(ctx context.Context) ([]facultymodel.Faculty, error)
	GetFacultyByID(ctx context.Context, facultyID string) (*facultymodel.Faculty, error)
	GetMajorsForFaculty(ctx context.Context, facultyId string) ([]majormodel.Major, error)
	CreateFaculty(ctx context.Context, faculty facultymodel.Faculty, image []byte) (facultymodel.Faculty, error)
	UpdateFaculty(ctx context.Context, facultyID string, faculty facultymodel.Faculty, image []byte) (facultymodel.Faculty, error)
	DeleteFaculty(ctx context.Context, facultyID string) error
}

type FacultyService struct {
	FacultyRepository facultyrepo.IFacultyRepository
	MajorRepository   majorrepository.IMajorRepository
}

func NewFacultyService(facultyRepo facultyrepo.IFacultyRepository, MajorRepo majorrepository.IMajorRepository) IFacultyService {
	return &FacultyService{
		FacultyRepository: facultyRepo,
		MajorRepository:   MajorRepo,
	}
}

func (s FacultyService) GetAllFaculties(ctx context.Context) ([]facultymodel.Faculty, error) {
	return s.FacultyRepository.FindAllFaculties(ctx)
}

func (s FacultyService) GetFacultyByID(ctx context.Context, facultyID string) (*facultymodel.Faculty, error) {
	return s.FacultyRepository.FindFacultyByID(ctx, facultyID)
}

func (s *FacultyService) GetMajorsForFaculty(ctx context.Context, facultyId string) ([]majormodel.Major, error) {
	faculty, err := s.FacultyRepository.FindFacultyByID(ctx, facultyId)
	if err != nil {
		return nil, err
	}

	majors, err := s.MajorRepository.FindMajorsByIDs(ctx, faculty.MajorIDs)
	if err != nil {
		return nil, err
	}

	return majors, nil
}

func (s *FacultyService) CreateFaculty(ctx context.Context, faculty facultymodel.Faculty, image []byte) (facultymodel.Faculty, error) {
	facultyName := faculty.FacultyName

	return s.FacultyRepository.CreateFaculty(ctx, facultyName, image)
}

func (s *FacultyService) UpdateFaculty(ctx context.Context, facultyID string, faculty facultymodel.Faculty, image []byte) (facultymodel.Faculty, error) {
	return s.FacultyRepository.UpdateFaculty(ctx, facultyID, faculty, image)
}

func (s FacultyService) DeleteFaculty(ctx context.Context, facultyID string) error {
	return s.FacultyRepository.DeleteFaculty(ctx, facultyID)
}
