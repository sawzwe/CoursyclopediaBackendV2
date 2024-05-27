package route

import (
	"BackendCoursyclopedia/db"
	"BackendCoursyclopedia/handler/auditloghandler"
	"BackendCoursyclopedia/handler/facultyhandler"
	"BackendCoursyclopedia/handler/majorhandler"
	"BackendCoursyclopedia/handler/subjecthandler"
	"BackendCoursyclopedia/handler/userhandler"

	"BackendCoursyclopedia/middleware"
	"BackendCoursyclopedia/repository/facultyrepository"
	"BackendCoursyclopedia/repository/majorrepository"
	"BackendCoursyclopedia/repository/subjectrepository"
	userrepo "BackendCoursyclopedia/repository/userrepository"
	"BackendCoursyclopedia/service/facultyservice"
	"BackendCoursyclopedia/service/majorservice"
	"BackendCoursyclopedia/service/subjectservice"

	auditlogrepo "BackendCoursyclopedia/repository/auditlogrepository"

	auditlogsvc "BackendCoursyclopedia/service/auditlogservice"
	usersvc "BackendCoursyclopedia/service/userservice"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	db.ConnectDB()

	userRepository := userrepo.NewUserRepository(db.DB)
	majorRepository := majorrepository.NewMajorRepository(db.DB)
	facultyRepository := facultyrepository.NewFacultyRepository(db.DB)
	auditlogRepository := auditlogrepo.NewAuditLogRepository(db.DB)
	subjectRepository := subjectrepository.NewSubjectRepository(db.DB)

	userService := usersvc.NewUserService(userRepository)
	facultyService := facultyservice.NewFacultyService(facultyRepository, majorRepository)
	majorService := majorservice.NewMajorService(majorRepository, facultyRepository, subjectRepository)
	auditlogService := auditlogsvc.NewAuditLogService(auditlogRepository)
	subjectService := subjectservice.NewSubjectService(subjectRepository, majorRepository)

	userHandler := userhandler.NewUserHandler(userService)
	facultyHandler := facultyhandler.NewFacultyHandler(facultyService)
	majorHandler := majorhandler.NewMajorHandler(majorService)
	auditlogHandler := auditloghandler.NewAuditLogHandler(auditlogService)
	subjectHandler := subjecthandler.NewSubjectHandler(subjectService)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the API")
	})

	//Auth
	app.Post("/api/auth/login", userHandler.Login)
	app.Post("/api/auth/googlelogin", userHandler.GoogleLogin)
	app.Post("/api/auth/createoneuser", userHandler.CreateOneUser)

	protectedUserGroup := app.Group("/api/users", middleware.JWTMiddleware)
	protectedUserGroup.Get("/getallusers", userHandler.GetUsers)
	protectedUserGroup.Get("/getoneuser/:id", userHandler.GetOneUser)
	protectedUserGroup.Get("/getuserbyemail/:email", userHandler.GetUserByEmail)
	protectedUserGroup.Post("/createoneuser", userHandler.CreateOneUser)
	protectedUserGroup.Delete("/deleteoneuser/:id", userHandler.DeleteOneUser)
	protectedUserGroup.Put("/updateoneuser/:id", userHandler.UpdateOneUser)
	protectedUserGroup.Delete("/dropallusers", userHandler.DropAllUsers)

	protectedFacultyGroup := app.Group("/api/faculties", middleware.JWTMiddleware)
	protectedFacultyGroup.Get("/getallfaculties", facultyHandler.GetFaculties)
	protectedFacultyGroup.Get("/geteachfaculty/:id", facultyHandler.GetEachFaculty)
	protectedFacultyGroup.Get("/getamjorforfaculty/:id", facultyHandler.GetMajorsForeachFaculty)
	protectedFacultyGroup.Post("/createfaculty", facultyHandler.CreateFaculty)
	protectedFacultyGroup.Put("/updatefaculty/:id", facultyHandler.UpdateFaculty)
	protectedFacultyGroup.Delete("/deletefaculty/:id", facultyHandler.DeleteFaculty)

	protectedMajorGroup := app.Group("/api/majors", middleware.JWTMiddleware)
	protectedMajorGroup.Get("/getallmajors", majorHandler.GetMajors)
	protectedMajorGroup.Get("/geteachmajor/:id", majorHandler.Geteachmajor)
	protectedMajorGroup.Get("getsubjectsforeachmajor/:id", majorHandler.GetSubjectsForeachMajor)
	protectedMajorGroup.Post("/createmajor", majorHandler.CreateMajor)
	protectedMajorGroup.Delete("/deletemajor/:id", majorHandler.DeleteMajor)
	protectedMajorGroup.Put("/updatemajor/:id", majorHandler.UpdateMajor)

	protectedAuditlogGroup := app.Group("/api/auditlogs", middleware.JWTMiddleware)
	protectedAuditlogGroup.Get("/getallauditlogs", auditlogHandler.GetAuditLogs)

	protectedSubjectGroup := app.Group("/api/subjects", middleware.JWTMiddleware)
	protectedSubjectGroup.Get("/getallsubjects", subjectHandler.GetSubjects)
	protectedSubjectGroup.Get("/geteachsubject/:id", subjectHandler.GetEachSubject)
	protectedSubjectGroup.Post("/createsubject", subjectHandler.CreateSubject)
	protectedSubjectGroup.Delete("/deletesubject/:id", subjectHandler.DeleteSubject)
	protectedSubjectGroup.Put("/updatesubject/:id", subjectHandler.UpdateSubject)
	protectedSubjectGroup.Put("/updatelikes/:id", subjectHandler.AddLikeByEmailHandler)

}
