package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/labstack/gommon/log"
)

const (
	// TODO fill this in directly or through environment variable
	// Build a DSN e.g. postgres://postgres:gnsecret@choukasherp.cud7jbsftjfi.ap-southeast-1.rds.amazonaws.com:5432/dbName
	DB_DSN = "postgres://postgres:password@gononeterp.cud7jbsftjfi.ap-southeast-1.rds.amazonaws.com:5432/testDB"
)

type User struct {
	tenant_id       string
	tenant_email    string
	tenant_username string
}
type User2 struct {
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

func main() {

	// Middleware

	// Create DB pool
	db, err := sql.Open("postgres", DB_DSN)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	// Create an empty user and make the sql query (using $1 for the parameter)
	var myUser User
	userSql := "SELECT tenant_id, tenant_username, tenant_email FROM imraan_auth.tenants"

	err = db.QueryRow(userSql).Scan(&myUser.tenant_id, &myUser.tenant_email, &myUser.tenant_username)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
	}

	fmt.Printf("Hi %s, welcome back!\n", myUser.tenant_email)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.POST("/users", test)

	// Start server
	e.Logger.Fatal(e.Start(":1325"))

}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

//Handler
func test(c echo.Context) (err error) {

	u := new(User2)
	if err = c.Bind(u); err != nil {
		return
	}
	// To avoid security flaws try to avoid passing binded structs directly to other methods
	// if these structs contain fields that should not be bindable.

	return c.JSON(http.StatusOK, u)
}
