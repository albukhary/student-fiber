package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	setupRoutes(app)

	app.Listen(":3000")

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
}

func setupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Assalamu alaykum 👋!")
	})
	app.Get("/students", getStudents)
	app.Get("/student/:id", getStudent)
	app.Post("/create/student", createStudent)
	app.Delete("/delete/student/:id", deleteStudent)
	app.Patch("update/student/:id", updateStudent)

}

// API Controllers

// swagger: route GET /students students listStudents
// Returns a list of students
// responses :
// 200: studentsListResponse
func getStudents(c *fiber.Ctx) error {
	var students []Student

	err = db.Select(&students, "SELECT * FROM student;")

	return c.JSON(students)
}

// swagger: route GET /student/{id} student getStudent
// Finds and returns a particular student with the requested ID
// responses:
// 200: studentResponse

// constroller of Person
func getStudent(c *fiber.Ctx) error {
	idParam := c.Params("id")

	var student Student

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Get(&student, "SELECT id, name, email, age FROM student WHERE id=$1", id)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(student)
}

// swagger: route POST /create/student createStudent
// creates a student of given parameters and writes into the database
//
// Consumes:
// - application/json
// Produces:
// - application/json

// Postman will send student data as JSON
// and we will put it into student struct and then into database
func createStudent(c *fiber.Ctx) error {
	var student Student

	c.BodyParser(&student)

	//	// Read requested student details and save it to Student struct
	//	student.ID, _ = strconv.Atoi(c.Params("id"))
	//	student.Name = c.Params("name")
	//	student.Email = c.Params("email")
	//	student.ID, _ = strconv.Atoi(c.Params("age"))

	// insert into Person query
	insertStudent := `INSERT INTO student (id, name, email, age) VALUES ($1, $2, $3, $4);`

	// Insert the student
	db.MustExec(insertStudent, student.ID, student.Name, student.Email, student.Age)

	// print the newly added student to web site
	return c.JSON(student)
}

// swagger: route DELETE /delete/student/{id} delete student deleteStudent
// Finds and deletes a student with the requested ID
// Consumes:
// - application/json
// Produces:
// - application/json

func deleteStudent(c *fiber.Ctx) error {
	idParam := c.Params("id")

	var student Student

	id, err1 := strconv.Atoi(idParam)
	if err1 != nil {
		log.Fatal(err1)
	}
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

// swagger: route PUT /student/{id} update student updateStudent
// Updates the details of the student with the requested ID according to the requested parameters
// Consumes:
// - application/json
// Produces:
// -application/json

// Update controller
func updateStudent(c *fiber.Ctx) error {
	var student Student

	c.BodyParser(&student)
	//	// Read requested student details and save it to Student struct
	//	student.ID, _ = strconv.Atoi(c.Params("id"))
	//	student.Name = c.Params("name")
	//	student.Email = c.Params("email")
	//	student.ID, _ = strconv.Atoi(c.Params("age"))

	// insert into Person query
	updateStudent := `UPDATE student SET name=$1, email=$2, age=$3 WHERE id=$4;`

	// Insert the student into the database
	db.MustExec(updateStudent, student.Name, student.Email, student.Age, student.ID)

	// print the newly updated student details
	return c.JSON(student)
}
