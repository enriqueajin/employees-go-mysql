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
)

// Responses to client API is Ok
func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

// Set up the logger settings
func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

// Employees struct with public access
type EmployeesModel struct {
	Id int `gorm:"primary_key"`
	FirstName string
	LastName string
	Salary float64
}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("firstname")
	log.WithFields(log.Fields{"firstname": firstName}).Info("First name added")
	lastName := r.FormValue("lastname")
	log.WithFields(log.Fields{"lastname": lastName}).Info("Last name added")
	var salary string = r.FormValue("salary")
	log.WithFields(log.Fields{"salary": salary}).Info("Salary added")
	sal, _ := strconv.ParseFloat(salary, 64)
	employee := &EmployeesModel{FirstName: firstName, LastName: lastName, Salary: sal }
	db.Create(&employee)
	result := db.Last(&employee)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Value)
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
	router.HandleFunc("/create", CreateEmployee).Methods("POST")
	http.ListenAndServe(":8000", router)
}