package main

// @title Student CRUD using fiber
// @version 1.1
// @description This is a sample CRUD application implementing gofiber/fiber and /arsmn/fiber-swagger
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@novalab.uz
// @license.name Novalab 2.0
// @license.url novalab.uz
// @host localhost:8084
// @BasePath /

import (
	"fmt"
	"log"
	"os"
	"strconv"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	// docs are generated by Swag CLI, you have to import them.
	_ "github.com/albukhary/student_fiber/docs"
)

//struct Student represents body of
type Student struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Age   int    `db:"age"`
}

// variable of type pointer to a database
var db *sqlx.DB
var err error

func main() {
	app := fiber.New()

	app.Get("/swagger/*", swagger.Handler)

	setupRoutes(app)

	//Loading environment variables for DATABASE connection
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	// Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort)

	//open and connect to the database at the same time
	db, err = sqlx.Connect(dialect, dbURI)
	if err != nil {
		log.Fatal(err)
	}

	app.Listen(":8084")
}

func setupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Assalamu alaykum 👋!")
	})
	app.Get("/students", getStudents)
	app.Get("/student/:id", getStudent)
	app.Post("/create/student", createStudent)
	app.Delete("/delete/student/:id", deleteStudent)
	app.Put("update/student/:id", updateStudent)
}

// API Controllers

// getStudents godoc
// @Summary Retrieves the list of all students
// @Produce json
// @Success 200 {object} []Student
// @Router /students [get]
func getStudents(c *fiber.Ctx) error {
	var students []Student

	// Use db.Select() to write all the rows in a slice
	err := db.Select(&students, "SELECT * FROM student")
	if err != nil {
		log.Fatal(err)
	}

	//return the slice of students to http
	return c.JSON(students)
}

// getStudent godoc
// @Summary Retrieves user based on given ID
// @Produce json
// @Param id path integer true "Student ID"
// @Success 200 {object} Student
// @Router /student/{id} [get]
func getStudent(c *fiber.Ctx) error {
	idParam := c.Params("id")

	var student Student

	// ID is initially a string when we get it from JSON
	// convert into int to use in a query
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Fatal(err)
	}

	// query the database using db.Get()
	err = db.Get(&student, "SELECT id, name, email, age FROM student WHERE id=$1", id)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(student)
}

// createStudent godoc
// @Summary Creates a student record with user input details and writes into database
// @Accept json
// @Produce json
// @Param details body Student true "Student details"
// @Success 200 {object} Student
// @Router /create/student [post]
func createStudent(c *fiber.Ctx) error {
	var student Student

	// parse JSON to a student struct
	c.BodyParser(&student)

	fmt.Println(student)

	insertStudent := `INSERT INTO student (id, name, email, age) VALUES ($1, $2, $3, $4);`

	// Insert the student
	db.MustExec(insertStudent, student.ID, student.Name, student.Email, student.Age)

	// print the newly added student to web site
	return c.JSON(student)
}

// deleteStudent godoc
// @Summary Deletes a student with the specified ID
// @Produce json
// @Param id path integer true "Student ID"
// @Success 200 {object} Student
// @Router /delete/student/{id} [delete]
func deleteStudent(c *fiber.Ctx) error {
	idParam := c.Params("id")

	var student Student

	// ID is initially a string when we get it from JSON
	// convert into int to use in a query
	id, err1 := strconv.Atoi(idParam)
	if err1 != nil {
		log.Fatal(err1)
	}

	// find the requested student from database
	err = db.Get(&student, "SELECT id, name, email, age FROM student WHERE id=$1", id)
	if err != nil {
		fmt.Println("There is no student with such ID")
		log.Fatal(err)
	}

	deleteQuery := `DELETE FROM student WHERE id=$1`

	//execute deletion
	db.MustExec(deleteQuery, student.ID)

	return c.JSON(student)
}

// updateStudent godoc
// @Summary Updates a student record with user input details and writes into database
// @Accept json
// @Produce json
// @Param details body Student true "Updated Student Details"
// @Success 200 {object} Student
// @Router /update/student/{id} [put]
func updateStudent(c *fiber.Ctx) error {
	var student Student

	// parses JSON to struct
	c.BodyParser(&student)

	updateStudent := `UPDATE student SET name=$1, email=$2, age=$3 WHERE id=$4;`

	// Insert the student into the database
	db.MustExec(updateStudent, student.Name, student.Email, student.Age, student.ID)

	// print the newly updated student details
	return c.JSON(student)
}
