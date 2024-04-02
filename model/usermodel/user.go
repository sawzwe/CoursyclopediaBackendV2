// // package usermodel

// // import (
// // 	"go.mongodb.org/mongo-driver/bson/primitive"
// // )

// // type User struct {
// // 	ID          primitive.ObjectID   `bson:"_id,omitempty"`
// // 	Email       string               `bson:"email"`
// // 	Roles       []string             `bson:"roles"`
// // 	Wishlists   []primitive.ObjectID `bson:"wishlists"`
// // 	PhoneNumber string               `bson:"phoneNumber"`
// // 	Profile     struct {
// // 		FirstName string `bson:"firstName"`
// // 		LastName  string `bson:"lastName"`
// // 	} `bson:"profile"`
// // 	FacultyID primitive.ObjectID `bson:"facultyId,omitempty"`
// // }

// package usermodel

// import (
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type User struct {
// 	ID          primitive.ObjectID   `bson:"_id,omitempty"`
// 	Email       string               `bson:"email"`
// 	Password    string               `bson:"password"` // Add this line to store the hashed password
// 	Roles       []string             `bson:"roles"`
// 	Wishlists   []primitive.ObjectID `bson:"wishlists"`
// 	PhoneNumber string               `bson:"phoneNumber"`
// 	Profile     struct {
// 		FirstName string `bson:"firstName"`
// 		LastName  string `bson:"lastName"`
// 	} `bson:"profile"`
// 	FacultyID primitive.ObjectID `bson:"facultyId,omitempty"`
// }

package usermodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	Name        string   `bson:"name"`
	Slug        string   `bson:"slug"`
	Description string   `bson:"description"`
	Permissions []string `bson:"permissions"`
}

type User struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Email       string               `bson:"email"`
	Password    string               `bson:"password,omitempty"`
	Role        Role                 `bson:"role"`
	Wishlists   []primitive.ObjectID `bson:"wishlists"`
	PhoneNumber string               `bson:"phoneNumber"`
	Profile     struct {
		FirstName string `bson:"firstName"`
		LastName  string `bson:"lastName"`
	} `bson:"profile"`
	FacultyID primitive.ObjectID `bson:"facultyId,omitempty"`
	Status    string             `bson:"status"`
}
