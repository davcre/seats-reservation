package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
  "strconv"
)

var db *gorm.DB
var err error

type Reservation struct {
  Id int `json:"id"`
  Name string `json:"name"`
  Surname string `json:"surname"`
  No_of_seats int `json:"no_of_seats"`
  Date string `json:"date"`
  Email string `json:"email"`
  Phone string `json:"phone"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Welcome to Seat Reservation!")
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
//  router.HandleFunc("/delete-reservation/{id}", deleteReservation).Methods("DELETE")
 // router.HandleFunc("/update-reservation/{id}", updateReservation).Methods("UPDATE")
  log.Fatal(http.ListenAndServe(":1234", router))
}

func createReservation(w http.ResponseWriter, r *http.Request) {
  reqBody, _ := ioutil.ReadAll(r.Body)

  var reservation Reservation
  json.Unmarshal(reqBody, &reservation)
  db.Create(&reservation)

  fmt.Println("Endpoint Hit: Creating new reservation")
  json.NewEncoder(w).Encode(reservation)
}

//func deleteReservation(w http.ResponseWriter, r *http.Request) {

//}

//func updateReservation(w http.ResponseWriter, r *http.Request) {

//}

func getReservation(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  key := vars["id"]
  reservations := []Reservation{}
  db.Find(&reservations)

  for _, reservation := range reservations {
    s, err := strconv.Atoi(key)
    if err == nil {
      if reservation.Id == s {
        fmt.Println(reservation)
        fmt.Println("Endpoint Hit: Reservation No:", key)
        json.NewEncoder(w).Encode(reservation)
      }
    }
  }
}

func getAllReservations(w http.ResponseWriter, r *http.Request) {
  reservations := []Reservation{}
  db.Find(&reservations)
  fmt.Println("Endpoint Hit: Get all reservations")
  json.NewEncoder(w).Encode(reservations)
}

func main() {
  db, err = gorm.Open("mysql", "root:Terka01234@tcp(127.0.0.1:3306)/Cinema?charset=utf8&parseTime=True")

  if err != nil {
    log.Println("Connection Failed to Open")
  } else { 
     log.Println("Connection Established")
  }

  db.AutoMigrate(&Reservation{})
  handleRequests()
}
