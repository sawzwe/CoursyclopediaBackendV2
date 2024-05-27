package usersvc

import (
	"BackendCoursyclopedia/model/usermodel"
	userrepo "BackendCoursyclopedia/repository/userrepository"
	"context"
	"errors"
	"fmt"

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
	GoogleLogin(ctx context.Context, email, firebaseId string) (*usermodel.User, string, error)
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

func generateJWT(user *usermodel.User) (string, error) {
	expirationTime := time.Now().Add(2 * time.Hour)

	claims := jwt.RegisteredClaims{
		Subject:   user.ID.Hex(),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWTSECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*usermodel.User, string, error) {
	user, err := s.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !CheckPasswordHash(password, user.Password) {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := generateJWT(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// func (s *UserService) GoogleLogin(ctx context.Context, email, firebaseId string) (*usermodel.User, string, error) {

// 	user, err := s.UserRepository.GetUserByEmail(ctx, email)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("user with email %s not found", email)
// 	}

// 	if user.Profile.FirebaseId != firebaseId {
// 		return nil, "", errors.New("invalid Firebase ID")
// 	}

// 	token, err := generateJWT(user)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	return user, token, nil
// }

func (s *UserService) GoogleLogin(ctx context.Context, email, firebaseId string) (*usermodel.User, string, error) {
	user, err := s.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		newUser := usermodel.User{
			Email: email,
			Profile: struct {
				FirstName  string `bson:"firstName"`
				LastName   string `bson:"lastName"`
				FirebaseId string `bson:"firebaseId"`
			}{
				FirebaseId: firebaseId,
			},
			Status: "active",
			Role: usermodel.Role{
				Name:        "user",
				Slug:        "user",
				Description: "Default user role",
				Permissions: []string{},
			},
		}
		user, err = s.CreateNewUser(ctx, newUser)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create new user: %v", err)
		}
	} else {
		// If user found, check if the provided Firebase ID matches the stored Firebase ID
		if user.Profile.FirebaseId != firebaseId {
			return nil, "", errors.New("invalid Firebase ID")
		}
	}

	// Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
