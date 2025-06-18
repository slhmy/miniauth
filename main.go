package main

import (
	"fmt"
	"miniauth/database"
	"miniauth/routers"
	"miniauth/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// Main
//
//	@title						MiniAuth API
//	@version					1.0
//	@description				A simple authentication service with user and organization management
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.url				http://www.swagger.io/support
//	@contact.email				support@swagger.io
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//	@host						localhost:8080
//	@BasePath					/api
//	@securityDefinitions.basic	BasicAuth
func main() {
	// Initialize database
	databaseConfig := database.GetDatabaseConfig()
	db, err := database.NewDatabase(databaseConfig)
	if err != nil {
		panic(err)
	}
	err = database.SetupDatabase(db)
	if err != nil {
		panic(err)
	}
	println("Database setup complete")

	// Initialize default admin user
	err = database.InitializeDefaultAdmin(db)
	if err != nil {
		panic(err)
	}

	// Initialize services
	serviceManager := service.NewServiceManager(db)

	// Initialize Echo server
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Setup routes
	routers.SetupRoutes(e, serviceManager)

	// Start server
	fmt.Println("Starting server on :8080...")
	e.Logger.Fatal(e.Start(":8080"))
}
