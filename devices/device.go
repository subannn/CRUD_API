package devices

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	
	_ "github.com/lib/pq"

	"golang.org/x/crypto/bcrypt"

)

var DB *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1"
	dbname   = "go_crud"
)

func OpenDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	CheckError(err)

	DB = db

	defer db.Close()
}

func RequestHTTP() {
	http.HandleFunc("/jsonPost", JsonPost) 

	http.HandleFunc("/tableGet", TableGet)
	http.HandleFunc("/jsonGetByID", JsonGet)

	http.HandleFunc("/jsonPut", JsonPut)

	http.HandleFunc("/jsonDelate", JsonDelate) 

	http.ListenAndServe(":8080", nil)
}

type UserUpd struct { 
	Id string 	   `json:"id"` 
	Name  string   `json:"name"` 
	Surname string `json:"surname"` 
	Mail string	   `json:"mail"` 
	Phone string   `json:"phone"` 
	Password string`json:"password"` 
}

type User struct { 
	Name  string   `json:"name"` 
	Surname string `json:"surname"` 
	Mail string	   `json:"mail"` 
	Phone string   `json:"phone"` 
	Password string`json:"password"` 
}

func TableGet(w http.ResponseWriter, r *http.Request) {  
	if r.Method != http.MethodGet {
		panic("Wrong HTTP method.")
	}   
	rows, err := DB.Query("SELECT * FROM users")
	CheckError(err)

	defer rows.Close()

	var buf bytes.Buffer
	for rows.Next() {
		var id int
    	var name string
    	var surname string
		var mail string
		var pnumber string
		var pass string
    	if err := rows.Scan(&id, &name, &surname, &mail, &pnumber, &pass); err != nil {
       		CheckError(err)
    	}
    	buf.WriteString(fmt.Sprintf("%d: %s: %s: %s: %s: %s\n", id, name, surname, mail, pnumber, pass))
	}
	if err := rows.Err(); err != nil {
		CheckError(err)
	}

	fmt.Fprint(w, buf.String())			
}

func JsonGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		panic("Wrong HTTP method.")
	}
	
	var u User

	body, err := ioutil.ReadAll(r.Body)
	CheckError(err)

	var id string

	err = json.Unmarshal(body, &id)
	CheckError(err)

	psqlInfo := fmt.Sprintf("SELECT * FROM users WHERE id=%s", id)
	CheckError(err)

	rows, err := DB.Query(psqlInfo)
	CheckError(err)

	for rows.Next() {
		var name string
    	var surname string
		var mail string
		var phone string
		var password string
		if err := rows.Scan(&id, &name, &surname, &mail, &phone, &password); err != nil {
			CheckError(err)
		}
		u.Name = name
		u.Surname = surname
		u.Mail = mail
		u.Phone = phone
		u.Password = password
	
	}
	
	jsonU,err := json.Marshal(u)
	CheckError(err)
	
	fmt.Fprint(w, string(jsonU))	
}

func JsonPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		panic("Wrong HTTP method.")
	} 

    var u User

    body, err := ioutil.ReadAll(r.Body)
	CheckError(err)

    err = json.Unmarshal(body, &u)
    CheckError(err)

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 20)
	CheckError(err)
	
	_, err = DB.Exec("INSERT INTO users(name, surname, mail, phone, password) VALUES ($1, $2, $3, $4, $5)", u.Name, u.Surname, u.Mail, u.Phone, hash)
	CheckError(err)
}
func JsonPut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		panic("Wrong HTTP method.")
	}

	var u UserUpd

	body, err := ioutil.ReadAll(r.Body)
	CheckError(err)

	err = json.Unmarshal(body, &u)
    CheckError(err)
	
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 20)
	CheckError(err)

	_, err = DB.Exec("UPDATE users SET name=$2 , surname=$3, mail=$4 , phone=$5 , password=$6 WHERE id=$1", u.Id, u.Surname, u.Mail, u.Phone, hash)
	if err != nil{
		panic(err)
	}
	
}

func JsonDelate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		panic("Wrong HTTP method.")
	}   

	body, err := ioutil.ReadAll(r.Body)
	CheckError(err)

	var id string

	err = json.Unmarshal(body, &id)
	CheckError(err)

	_, err = DB.Exec("DELETE FROM users WHERE id = $1", id)
	CheckError(err)

}

func CheckError(err error) {
	if(err != nil) {
		panic(err)
	}
}