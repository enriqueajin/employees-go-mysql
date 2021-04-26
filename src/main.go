package main

import (
	"io"
	"net/http"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
       "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
	"encoding/json"
	"github.com/rs/cors"
)

// Employees struct with public access
type EmployeesModel struct {
	Id int `gorm:"primary_key"`
	FirstName string
	LastName string
	Salary float64
}

// Set up the logger settings
func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

// Responses to client API is Ok
func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

// Function to create a new employee
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("firstname")
	log.WithFields(log.Fields{"firstname": firstName}).Info("First name added")

	lastName := r.FormValue("lastname")
	log.WithFields(log.Fields{"lastname": lastName}).Info("Last name added")

	salary, _ := strconv.ParseFloat(r.FormValue("salary"), 64)
	log.WithFields(log.Fields{"salary": salary}).Info("Salary added")

	employee := &EmployeesModel{FirstName: firstName, LastName: lastName, Salary: salary }
	db.Select("FirstName", "LastName", "Salary").Create(&employee)
	result := db.Last(&employee)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Value)
}

// Function to get an employee by id
func GetEmployeeById(Id int) bool {
	employee := &EmployeesModel{}
	result := db.First(&employee, Id)
	if result.Error != nil{
			log.Warn("Employee not found in database")
			return false
	}
	return true
}

// Function to get all employees (Read)
func GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	var employees []EmployeesModel
	FindEmployees := db.Find(&employees).Value

	log.Info("Get all employees")
	allEmployees := FindEmployees
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allEmployees)
}
 
// func GetTodoItems(completed bool) interface{} {
// 	var todos []TodoItemModel
// 	TodoItems := db.Where("completed = ?", completed).Find(&todos).Value
// 	return TodoItems
// }

// Function to update an employee
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	// Get URL parameter from mux
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Test if an employee exist in DB
	err := GetEmployeeById(id)
	if !err {
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"updated": false, "error": "Record Not Found"}`)
	} else {
			firstname := r.FormValue("firstname")
			log.WithFields(log.Fields{"firstname": firstname}).Info("Updating first name")

			lastname := r.FormValue("lastname")
			log.WithFields(log.Fields{"lastname": lastname}).Info("Updating last name")

			salary, _ := strconv.ParseFloat(r.FormValue("salary"), 64)
			log.WithFields(log.Fields{"salary": salary}).Info("Updating salary")

			employee := &EmployeesModel{}
			db.First(&employee, id)
			employee.FirstName = firstname
			employee.LastName = lastname
			employee.Salary = salary
			db.Save(&employee)
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"updated": true}`)
	}
}

// Function to delete an employee
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	// Get URL parameter from mux
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Test if an employee exist in DB
	err := GetEmployeeById(id)
	if !err {
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"deleted": false, "error": "Record Not Found"}`)
	} else {
			log.WithFields(log.Fields{"Id": id}).Info("Deleting employee")
			employee := &EmployeesModel{}
			db.First(&employee, id)
			db.Delete(&employee)
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"deleted": true}`)
	}
}

// Open database connection
var db, _ = gorm.Open("mysql", "root:root@/employees?charset=utf8&parseTime=True&loc=Local")

// Main function
func main() {
	// Close database connection
	defer db.Close()

	// Automigrates MySQL database after starting our API Server
	db.Debug().DropTableIfExists(&EmployeesModel{})
	db.Debug().AutoMigrate(&EmployeesModel{})

	log.Info("Starting employees API server")
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	router.HandleFunc("/all-employees", GetAllEmployees).Methods("GET")
	router.HandleFunc("/create", CreateEmployee).Methods("POST")
	router.HandleFunc("/employee/{id}", UpdateEmployee).Methods("POST")
	router.HandleFunc("/employee/{id}", DeleteEmployee).Methods("DELETE")

	handler := cors.New(cors.Options{
			AllowedMethods: []string{"GET", "POST", "DELETE"},
	}).Handler(router)

	http.ListenAndServe(":8000", handler)
}	