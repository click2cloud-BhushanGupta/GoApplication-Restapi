package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//build struct and how they represented in json
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	City      string `json:"city"`
}

var db *sql.DB
var err error

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User
	result, err := db.Query("SELECT * from users ORDER BY id desc")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var user User
		err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.City)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT * FROM users WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var user User
	for result.Next() {
		err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.City)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "User with ID = %s was deleted", params["id"])

}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("INSERT INTO users(firstname,lastname,email,city) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	firstname := keyVal["firstname"]
	lastname := keyVal["lastname"]
	email := keyVal["email"]
	city := keyVal["city"]
	_, err = stmt.Exec(firstname, lastname, email, city)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New User was created")
}

func deleteallUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := db.Query("DELETE  from users")
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "All user was deleted")
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT * FROM users WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var user User
	for result.Next() {
		err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.City)
		if err != nil {
			panic(err.Error())
		}
	}
	stmt, err := db.Prepare("UPDATE users SET firstname = ?, lastname=?,email=?,city=? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	fmt.Println(keyVal)
	firstname := keyVal["firstname"]
	fmt.Println(firstname)
	lastname := keyVal["lastname"]
	fmt.Println(lastname)
	email := keyVal["email"]
	fmt.Println(email)
	city := keyVal["city"]
	fmt.Println(city)
	m := []string{firstname, lastname, email, city}
	fmt.Println(m)
	for i, _ := range m {
		if m[i] == "" {
			if i == 0 {
				firstname = user.FirstName
			} else if i == 1 {
				lastname = user.LastName
			} else if i == 2 {
				email = user.Email
			} else {
				city = user.City
			}
		}
	}
	fmt.Println(r)
	_, err = stmt.Exec(firstname, lastname, email, city, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "User with ID = %s was updated", params["id"])
}

func main() {
	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//init router
	router := mux.NewRouter()

	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/users", deleteallUser).Methods("DELETE")
	//to run server
	http.ListenAndServe(":8000", router)
}
