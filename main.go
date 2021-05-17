package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "encoding/json"
  "time"
  "strconv"
  "sync"

  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  "github.com/joho/godotenv"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error
var lock sync.Mutex

type Reservation struct {
  Id uint `gorm:"primary_key;auto_increment" json:"id"`
  Name string `gorm:"size:100;not null" json:"name"`
  Surname string `gorm:"size:100;not null" json:"surname"`
  Date time.Time `gorm:"not null" json:"date"`
  Email string `gorm:"size:100;not null;unique" json:"email"`
  Phone string `gorm:"size:12;not null;unique" json:"phone"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Welcome to Seat Reservation app!")
  fmt.Println("Endpoint Hit: Homepage")
}

func handleRequests() {
  log.Println("Staring development server at http://127.0.0.1:1234/")
  log.Println("Quit the server with CTRL-C")

  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", homePage)
  router.HandleFunc("/new-reservation", createReservation).Methods("POST")
  router.HandleFunc("/all-reservations", getAllReservations).Methods("GET")
  router.HandleFunc("/reservation/{id}", getReservation).Methods("GET")
  router.HandleFunc("/delete-reservation/{id}", deleteReservation).Methods("DELETE")
  router.HandleFunc("/update-reservation/{id}", updateReservation).Methods("PUT")

  log.Fatal(http.ListenAndServe(":1234", router))
}

func createReservation(w http.ResponseWriter, r *http.Request) {
  reqBody, _ := ioutil.ReadAll(r.Body)
  vars := mux.Vars(r)
  key := vars["id"] 

  var reservation Reservation

  if err := json.Unmarshal(reqBody, &reservation); err != nil {
    respondError(w, http.StatusBadRequest, err.Error())
    return
  }

  if err := db.First(&reservation, key).Error; err != nil {

    if err := db.Create(&reservation).Error; err != nil {
      respondError(w, http.StatusInternalServerError, err.Error())
      return
    }

    fmt.Println("Endpoint Hit: Creating new reservation")
    respondJSON(w, http.StatusCreated, reservation)
    return

  } else {
      respondJSON(w, http.StatusOK, reservation)
      return
  }
}

func deleteReservation(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  key := vars["id"]
  reservations := []Reservation{}

  id64, _ := strconv.ParseInt(key, 10, 64)
  idToDelete := int(id64)

  if err := db.Where("id = ?", idToDelete).Delete(&reservations).Error;
   err != nil {
    respondError(w, http.StatusInternalServerError, err.Error())
    return
  }

  fmt.Println("Endpoint Hit: Deleted reservation No:", key)
  respondJSON(w, http.StatusNoContent, nil)
}

func updateReservation(w http.ResponseWriter, r *http.Request) {
  reqBody, _ := ioutil.ReadAll(r.Body)
  vars := mux.Vars(r)
  key := vars["id"]

  var reservation Reservation

  if err := db.First(&reservation, key).Error; err != nil {
    respondError(w, http.StatusNotFound, err.Error())
    return
  }

  if err := json.Unmarshal(reqBody, &reservation); err != nil {
    respondError(w, http.StatusBadRequest, err.Error())
    return
  }

  lock.Lock()
  
  if err := db.Save(&reservation).Error; err != nil {
    respondError(w, http.StatusInternalServerError, err.Error())
    return
  }

  lock.Unlock()

  respondJSON(w, http.StatusOK, reservation)
  fmt.Println("Endpoint Hit: Updating reservation No:", key)
}

func getReservation(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  key := vars["id"]
  var reservation Reservation

  if err := db.First(&reservation, key).Error; err != nil {
    respondError(w, http.StatusNotFound, err.Error())
    return
  }

  respondJSON(w, http.StatusOK, reservation)
  fmt.Println("Endpoint Hit: Reservation No:", key)
}

func getAllReservations(w http.ResponseWriter, r *http.Request) {
  reservations := []Reservation{}
  db.Find(&reservations)
  respondJSON(w, http.StatusOK, reservations)
  fmt.Println("Endpoint Hit: Get all reservations")
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
  response, err := json.Marshal(payload)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(status)
  w.Write([]byte(response))
}

func respondError(w http.ResponseWriter, code int, message string) {
  respondJSON(w, code, map[string]string{"error": message})
}

func (ts Reservation) String() string {
  t := ts.Date
  timestamp := t.Format("02-01-2006 15:04:05")
  
  return fmt.Sprintf("%v", timestamp)
}

func main() {
  var appConfig map[string]string
  appConfig, err := godotenv.Read()
  
  if err != nil {
    log.Fatal("Error reading .env file")
  }
 
  mysqlCredentials := fmt.Sprintf(
    "%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
    appConfig["MYSQL_USER"],
    appConfig["MYSQL_PASSWORD"],
    appConfig["MYSQL_PROTOCOL"],
    appConfig["MYSQL_HOST"],
    appConfig["MYSQL_PORT"],
    appConfig["MYSQL_DBNAME"],
  )
  db, err = gorm.Open("mysql", mysqlCredentials)

  if err != nil {
     fmt.Println(err)
     panic("Failed to connect to the database!")
  } else { 
     log.Println("Connected to database")
  }

  db.AutoMigrate(&Reservation{})
  handleRequests()

}
