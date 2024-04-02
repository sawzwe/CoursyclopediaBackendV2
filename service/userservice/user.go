package usersvc

import (
	"BackendCoursyclopedia/model/usermodel"
	userrepo "BackendCoursyclopedia/repository/userrepository"
	"context"
	"errors"

	"time"

	"os"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	GetAllUsers(ctx context.Context) ([]usermodel.User, error)
	GetUserByID(ctx context.Context, userID string) (*usermodel.User, error)
	GetUserByEmail(ctx context.Context, email string) (*usermodel.User, error)
	CreateNewUser(ctx context.Context, user usermodel.User) (*usermodel.User, error)
	DeleteSpecificUser(ctx context.Context, userID string) error
	UpdateSpecificByID(ctx context.Context, userID string, updateUser usermodel.User) (*usermodel.User, error)
	DropAllUsers(ctx context.Context) error
	Login(ctx context.Context, email, password string) (*usermodel.User, string, error)
}

type UserService struct {
	UserRepository userrepo.IUserRepository
}

func NewUserService(userRepo userrepo.IUserRepository) IUserService {
	return &UserService{
		UserRepository: userRepo,
	}
}

// HashPassword generates a bcrypt hash of the password using a default cost of 10.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a bcrypt hashed password with its possible plaintext equivalent.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*usermodel.User, error) {
	return s.UserRepository.FindUserByID(ctx, userID)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]usermodel.User, error) {
	return s.UserRepository.FindAllUsers(ctx)

}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	return s.UserRepository.GetUserByEmail(ctx, email)

}

// func (s *UserService) CreateNewUser(ctx context.Context, user usermodel.User) (*usermodel.User, error) {
// 	return s.UserRepository.CreateUser(ctx, user)
// }

func (s *UserService) CreateNewUser(ctx context.Context, user usermodel.User) (*usermodel.User, error) {
	// Hash the password before saving it to the database
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword // Store the hashed password

	return s.UserRepository.CreateUser(ctx, user)
}

func (s *UserService) DeleteSpecificUser(ctx context.Context, userID string) error {
	return s.UserRepository.DeleteUserByID(ctx, userID)
}

func (s *UserService) UpdateSpecificByID(ctx context.Context, userID string, updateUser usermodel.User) (*usermodel.User, error) {
	// Check if the password field is set in updateUser
	if updateUser.Password != "" {
		// Hash the new password before updating
		hashedPassword, err := HashPassword(updateUser.Password)
		if err != nil {
			return nil, err
		}
		updateUser.Password = hashedPassword
	}

	return s.UserRepository.UpdateUserByID(ctx, userID, updateUser)
}

func (s *UserService) DropAllUsers(ctx context.Context) error {
	return s.UserRepository.DropAllUsers(ctx)
}

// generateJWT generates a JWT token for authenticated users
func generateJWT(user *usermodel.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token is valid for 24 hours

	claims := jwt.RegisteredClaims{
		Subject:   user.ID.Hex(), // Use the user's ID as the subject
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWTSECRET"))) // Use a secret key

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*usermodel.User, string, error) {
	user, err := s.UserRepository.GetUserByEmailLogin(ctx, email)
	if err != nil || !CheckPasswordHash(password, user.Password) {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
