
# Cousryclopedia Backend

This is the api endpoints for the Coursyclopedia application using go language with fiber framework. Below is the project structure. MiddleWareAdded

## directory

coursyclopediabackend/

├── handler/             # Handlers for the HTTP API endpoints 
   
├── student/             # Student-related API handlers

│   ├── user/            # User-related API handlers

│   ├── subject/                # Subject-related API handlers

│   └── auditlog/        # AuditLog-related API handlers

├── model/               # Data models for the application

│   ├── user.go          # User model definition

│   ├── subject.go       # Subject model definition

│   ├── auditlog.go      # AuditLog model definition

│   └── faculty.go       # Faculty model definition

├── service/             # Business logic layer

│   ├── userService.go

│   ├── subjectService.go

│   └── auditService.go

├── repository/          # Database interaction layer

│   ├── userRepository.go

│   ├── subjectRepository.go

│   └── auditRepository.go

├── db/                  # Database connection setup

│   └── mongodb.go

├── pkg/                 # Utility packages and common libraries

│   └── utils/

├── route/

│   └── routes/          # Routes for endpoints

└── main.go              # Entry point of the application

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`MONGODB_URI=mongodb://your_mongo_uri`

## Clone Repository

To get started with this project, make sure you have Go installed on your system. then you can clone it with following comman:

`git clone https://github.com/BXBZwe/CoursyclopediaBackend.git`

`cd CoursyclopediaBackend`

## Install the dependencies
`go mod tidy`

## Run the application

`go run main.go`
