package main

import (
	"io"
	"net/http"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
       "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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

var db, _ = gorm.Open("mysql", "root:root@/employees?charset=utf8&parseTime=True&loc=Local")

// Main function
func main() {
	// Close database connection
	defer db.Close()

	log.Info("Starting employees API server")
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	http.ListenAndServe(":8000", router)
}