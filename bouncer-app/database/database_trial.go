package database

import ("fmt"
        "github.com/boltdb/bolt"
        "encoding/json"
        "github.com/satori/go.uuid"
        )



type User struct {
  ID uuid.UUID `json:"ID"`
  FirstName string `json:"FirstName"`
  LastName string `json:"LastName"`
  Address string  `json:"Address"`
}

type Database struct {
  db *bolt.DB
}

func OpenDB()(d *Database, err error) {
  d = &Database{}
  d.db, err = bolt.Open("database/trial.db", 0600, nil)
  return
}

func (d *Database) CreateUser(user User) error {
  err := d.db.Update(func(tx *bolt.Tx) error {
    var bucket string
    bucket = "DB"
    b, err := tx.CreateBucketIfNotExists([]byte(bucket))
    if err != nil {
      return err
    }
    fmt.Println("Bucket Created!")

    uid := uuid.NewV4()
    user.ID = uid
    fmt.Println(uid)
    key := uid.String()
    encoded, err := json.Marshal(user)
    if err != nil {
      return err
    }
    fmt.Printf("%s", user)
    b.Put([]byte(key), encoded)
    fmt.Println("User Created Succesfully!")
    return nil
  })
return err
}


func(d *Database) RetrieveUser(key string)(data []byte, err error) {
  err = d.db.View(func(tx *bolt.Tx) error {
      var bucket string
      bucket = "DB"
      b := tx.Bucket([]byte(bucket))
      val := b.Get([]byte(key))
      if val != nil{
        data = make([]byte, len(val))
        copy(data, val)
      }
      val = b.Get([]byte(fmt.Sprintf("%s", key)))
      if val == nil {
        fmt.Println("User does not exist. Check user key again.")
        return nil
      }

      ct := string(val)
      fmt.Println(ct)

      return nil
    })

  return
}

func (d *Database) DeleteUser(key string) error {
  return d.db.Update(func(tx *bolt.Tx) error {
      var bucket string
      bucket = "DB"
      b := tx.Bucket([]byte(bucket))
      b.Delete([]byte(key))
      fmt.Println("User Deleted")
      return nil
  })
}

func (d *Database) UpdateUser(key string, updatedUser User) error {
  err := d.db.Update(func(tx *bolt.Tx) error {
    var user User
    var bucket string
    bucket = "DB"
    b := tx.Bucket([]byte(bucket))
    r := b.Get([]byte(key))
    
    err := json.Unmarshal(r, &user)
    if err != nil {
      return err
    }
    fmt.Println(user)

    user.FirstName = updatedUser.FirstName
    user.LastName = updatedUser.LastName
    user.Address = updatedUser.Address
    fmt.Println(user)

    encoded, err := json.Marshal(user)
    if err != nil {
      return err
    }
    b.Put([]byte(key), encoded)

    fmt.Print("User Data Updated Succesfully!")
    return nil
  })
return err
}

func (d *Database) ListUsers() (vals [][]byte, err error) {
  vals = make([][]byte, 0)
  err = d.db.View(func(tx *bolt.Tx) error {
          var bucket string
          bucket = "DB"
        	b := tx.Bucket([]byte(bucket))
        	c := b.Cursor()

        	for k, v := c.First(); k != nil; k, v = c.Next() {
            vals = append(vals, v)
        		fmt.Printf("key=%s, value=%s\n", k, v)
        	}

        	return nil
})
return vals, err
}
