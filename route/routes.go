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

	app.Post("/api/users/login", userHandler.Login)

	protectedUserGroup := app.Group("/api/users", middleware.JWTMiddleware)
	protectedUserGroup.Get("/getallusers", userHandler.GetUsers)
	protectedUserGroup.Get("/getoneuser/:id", userHandler.GetOneUser)
	protectedUserGroup.Get("/getuserbyemail/:email", userHandler.GetUserByEmail)
	protectedUserGroup.Post("/createoneuser", userHandler.CreateOneUser)
	protectedUserGroup.Delete("/deleteoneuser/:id", userHandler.DeleteOneUser)
	protectedUserGroup.Put("/updateoneuser/:id", userHandler.UpdateOneUser)
	protectedUserGroup.Delete("/dropallusers", userHandler.DropAllUsers)

	// app.Post("/api/users/login", userHandler.Login)

	// userGroup := app.Group("/api/users")
	// // userGroup := app.Group("/api/users", jwtmiddleware.JWTAuthMiddleware)
	// // userGroup.Get("/getallusers", middleware.JWTMiddleware, userHandler.GetUsers)
	// userGroup.Get("/getallusers", userHandler.GetUsers)
	// userGroup.Get("/getoneuser/:id", userHandler.GetOneUser)
	// userGroup.Get("/getuserbyemail/:email", userHandler.GetUserByEmail)
	// userGroup.Post("/createoneuser", userHandler.CreateOneUser)
	// userGroup.Delete("/deleteoneuser/:id", userHandler.DeleteOneUser)
	// userGroup.Put("/updateoneuser/:id", userHandler.UpdateOneUser)
	// userGroup.Delete("/dropallusers", userHandler.DropAllUsers)
	// userGroup.Post("/login", userHandler.Login)

	faculyGroup := app.Group("/api/faculties")
	faculyGroup.Get("/getallfaculties", facultyHandler.GetFaculties)
	faculyGroup.Get("/geteachfaculty/:id", facultyHandler.GetEachFaculty)
	faculyGroup.Get("/getamjorforfaculty/:id", facultyHandler.GetMajorsForeachFaculty)
	faculyGroup.Post("/createfaculty", facultyHandler.CreateFaculty)
	faculyGroup.Put("/updatefaculty/:id", facultyHandler.UpdateFaculty)
	faculyGroup.Delete("/deletefaculty/:id", facultyHandler.DeleteFaculty)

	majorGroup := app.Group("api/majors")
	majorGroup.Get("/getallmajors", majorHandler.GetMajors)
	majorGroup.Get("/geteachmajor/:id", majorHandler.Geteachmajor)
	majorGroup.Get("getsubjectsforeachmajor/:id", majorHandler.GetSubjectsForeachMajor)
	majorGroup.Post("/createmajor", majorHandler.CreateMajor)
	majorGroup.Delete("/deletemajor/:id", majorHandler.DeleteMajor)
	majorGroup.Put("/updatemajor/:id", majorHandler.UpdateMajor)

	auditlogGroup := app.Group("/api/auditlogs")
	auditlogGroup.Get("/getallauditlogs", auditlogHandler.GetAuditLogs)

	subjectGroup := app.Group("api/subjects")
	subjectGroup.Get("/getallsubjects", subjectHandler.GetSubjects)
	subjectGroup.Get("/geteachsubject/:id", subjectHandler.GetEachSubject)
	subjectGroup.Post("/createsubject", subjectHandler.CreateSubject)
	subjectGroup.Delete("/deletesubject/:id", subjectHandler.DeleteSubject)
	subjectGroup.Put("/updatesubject/:id", subjectHandler.UpdateSubject)
	subjectGroup.Put("/updatelikes/:id", subjectHandler.AddLikeByEmailHandler)

}
