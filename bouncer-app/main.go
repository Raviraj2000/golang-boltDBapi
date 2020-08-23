package main


import ("fmt"
        "net/http"
        "log"
        "encoding/json"
        "github.com/bouncer-app/database"
        "github.com/gorilla/mux"
        //"github.com/satori/go.uuid"
        "io/ioutil"
      )



func HomePage(w http.ResponseWriter, r *http.Request){
  fmt.Fprintf(w, "Welcome to the HomePage!")
  fmt.Println("Endpoint Hit: homePage")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
  db, err := database.OpenDB()
  if err != nil {
    log.Fatal("Error: %s", err)
  }
  var user database.User
  bytes, _ := ioutil.ReadAll(r.Body)
  err = json.Unmarshal(bytes, &user)
  if err != nil {
    fmt.Printf("%s", err)
  }

  vars := mux.Vars(r)
  bucket := vars["bucket"]
  key := vars["id"]

  err = db.UpdateUser(bucket, key, user)
  if err == nil{
    fmt.Fprintf(w, "User Data Updated Succesfully!")
  }

}

func DeleteUser (w http.ResponseWriter, r *http.Request) {
  db, err := database.OpenDB()
  if err != nil {
    log.Fatal("Error: %s", err)
  }

  vars := mux.Vars(r)
  bucket := vars["bucket"]
  key := vars["id"]

  err = db.DeleteUser(bucket, key)
  if err == nil {
    fmt.Fprintf(w, "User Deleted")
  }
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
  db, err := database.OpenDB()
  if err != nil {
    log.Fatal("Error: %s", err)
  }
  var user database.User
  bytes, _ := ioutil.ReadAll(r.Body)

  err = json.Unmarshal(bytes, &user)
  if err != nil {
    fmt.Printf("%s", err)
  }
  err = db.CreateUser(user)
  if err != nil {
    fmt.Printf("%s", err)
  }
  if err == nil {
    fmt.Fprintf(w, "User Created Succesfully")
  }


}

func RetrieveAllUsers(w http.ResponseWriter, r *http.Request) {
  db, err := database.OpenDB()
  if err != nil {
    log.Fatal("Error: %s", err)
  }

  vars := mux.Vars(r)
  bucket := vars["bucket"]

  data, err := db.ListUsers(bucket)
  var user database.User

  for _, users := range data {
    err = json.Unmarshal(users, &user)
    if err != nil {
      fmt.Printf("%s", users)
    }
    json.NewEncoder(w).Encode(user)
  }
}


func RetrieveSingleUser(w http.ResponseWriter, r *http.Request){
  db, err := database.OpenDB()
  if err != nil {
    log.Fatal("Error: %s", err)
  }

  vars := mux.Vars(r)
  bucket := vars["bucket"]
  key := vars["id"]
  data, err := db.RetrieveUser(bucket, key)
  if err != nil {
    fmt.Printf("%s", err)
  }
  if data == nil {
    fmt.Fprintf(w, "User does not exist. Check name of bucket and key again.")
    return
  }

  var user database.User

  err = json.Unmarshal(data, &user)

  if err != nil {
    fmt.Printf("%s", err)
  }
  err = json.NewEncoder(w).Encode(user)
  if err != nil{
    fmt.Printf("%s", err)
  }
}


func handleRequests() {
  r := mux.NewRouter().StrictSlash(true)
  r.HandleFunc("/", HomePage)
  r.HandleFunc("/users/{bucket}", RetrieveAllUsers).Methods("GET")
  r.HandleFunc("/users/{bucket}/{id}", RetrieveSingleUser).Methods("GET")
  r.HandleFunc("/users/create", CreateUser).Methods("POST")
  r.HandleFunc("/users/delete/{bucket}/{id}", DeleteUser).Methods("DELETE")
  r.HandleFunc("/users/update/{bucket}/{id}", UpdateUser).Methods("PUT")
  log.Fatal(http.ListenAndServe(":8000", r))
}



func main(){
  handleRequests()

  /*db, err := database.OpenDB()
  if err != nil {
    log.Fatal("Error: %s", err)
  }*/
  //db.CreateUser()
  //db.RetrieveUser()
  //db.DeleteUser()
  //db.UpdateUser()
  //db.ListUsers("DB")

}
