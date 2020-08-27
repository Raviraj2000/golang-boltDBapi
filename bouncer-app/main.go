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

type server struct {
  db *database.Database
}

func NewServer()(*server, error){
  s := &server{}
  var err error
  s.db, err = database.OpenDB()
  if err != nil {
    log.Printf("Error opening Database Err:%s", err.Error())
    return nil, err
  }
  return s, nil
}



func HomePage(w http.ResponseWriter, r *http.Request){
  fmt.Fprintf(w, "Welcome to the HomePage!")
  fmt.Println("Endpoint Hit: homePage")

}

func (s *server)UpdateUser(w http.ResponseWriter, r *http.Request) {
  var user database.User
  bytes, _ := ioutil.ReadAll(r.Body)
  err := json.Unmarshal(bytes, &user)
  if err != nil {
    fmt.Fprintf(w, "Error:%s", err)
    return
  }

  vars := mux.Vars(r)
  key := vars["id"]

  err = s.db.UpdateUser(key, user)
  if err == nil{
    fmt.Fprintf(w, "User Data Updated Succesfully!")
  } else {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "User Does not exist. Check user key again.")
    return
  }

}

func (s *server)DeleteUser(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  key := vars["id"]
  err := s.db.DeleteUser(key)
  if err == nil {
    fmt.Fprintf(w, "User Deleted")
  } else {
    http.Error(w, "Could not delete user.", http.StatusInternalServerError)
  }
}

func (s *server)CreateUser(w http.ResponseWriter, r *http.Request) {
  var user database.User
  bytes, err := ioutil.ReadAll(r.Body)
  if err != nil{
    http.Error(w, "Please Check Post Data", http.StatusBadRequest)
    return
  }
  err = json.Unmarshal(bytes, &user)
  if err != nil {
    http.Error(w, "Json Unmarshalling Error", http.StatusInternalServerError)
    return
  }

  err = s.db.CreateUser(user)
  if err != nil {
    http.Error(w, "Could not create user", http.StatusInternalServerError)
    return
  }
  if err == nil {
    fmt.Fprintf(w, "User Created Succesfully")
  }

}

func (s *server)RetrieveAllUsers(w http.ResponseWriter, r *http.Request) {

  data, err := s.db.ListUsers()
  if data == nil {
    http.Error(w, "Database Empty.", http.StatusInternalServerError)
    return
  }
  var user database.User
  for _, users := range data {
    err = json.Unmarshal(users, &user)
    if err != nil {
      http.Error(w, "Json unmarshalling error.", http.StatusInternalServerError)
      return
    }
    err = json.NewEncoder(w).Encode(user)
    if err != nil{
      http.Error(w, "Json encoding error.", http.StatusInternalServerError)
      return
    }
  }
}


func (s *server)RetrieveSingleUser(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  key := vars["id"]
  data, err := s.db.RetrieveUser(key)
  if err != nil {
    http.Error(w, "Data Retrival Error", http.StatusInternalServerError)
    return
  }
  if data == nil {
    http.Error(w, "User does not exist. Check user key again.", http.StatusNotFound)
    return
  }

  var user database.User

  err = json.Unmarshal(data, &user)
  if err != nil {
    http.Error(w, "Json unmarshalling error.", http.StatusInternalServerError)
    return
  }
  err = json.NewEncoder(w).Encode(user)
  if err != nil{
    http.Error(w, "Json encoding error.", http.StatusInternalServerError)
    return
  }
  return
}

func handleRequests() {
  s, err := NewServer()
 if err != nil {
     panic(err)
 }
  r := mux.NewRouter().StrictSlash(true)
  r.HandleFunc("/", HomePage)
  r.HandleFunc("/users", s.RetrieveAllUsers).Methods("GET")
  r.HandleFunc("/users/{id}", s.RetrieveSingleUser).Methods("GET")
  r.HandleFunc("/users/create", s.CreateUser).Methods("POST")
  r.HandleFunc("/users/delete/{id}", s.DeleteUser).Methods("DELETE")
  r.HandleFunc("/users/update/{id}", s.UpdateUser).Methods("PUT")
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
