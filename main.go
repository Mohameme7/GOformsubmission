package main

import (
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const ADMINPASSWORD = "admin1234"

func checkpassword(a string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(ADMINPASSWORD)) == 1
}

var db *sql.DB

// Inits a connection to the database
func init() {
	db, _ = sql.Open("sqlite3", "./database.db")
}
func main() {
	r := chi.NewRouter()
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	r.Get("/", DefaultPage)
	r.Get("/panel", AdminPage)
	r.Get("/admin", AdminLoginPage)
	r.Get("/admin/getdata", SendForms)
	r.Post("/sendfourm", ReceiveData)
	r.Post("/admin/login", LoginEndPoint)
	http.ListenAndServe(":5000", r)
}

// returns the default page
func DefaultPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

// Returns a page for login to admin panel(very advanced one fr)
func AdminLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/AdminLogin.html")
}
func AdminPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/adminformpage.html")
}

func LoginEndPoint(w http.ResponseWriter, r *http.Request) {
	var ReceivedInfo password

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	json.Unmarshal(body, &ReceivedInfo)

	if !checkpassword(ReceivedInfo.Password) {
		http.Error(w, "Wrong Password", http.StatusUnauthorized)

	}
}

// Receive Data from users after submission and insert into the Database
func ReceiveData(w http.ResponseWriter, r *http.Request) {
	var data formdata
	body, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(body, &data)

	cookie := gencookie()
	db.Exec("INSERT INTO FORMS (Subject, Email, Body, Cookie, Answered) VALUES (?,?,?,?,?)", data.Subject, data.Email, data.Body, cookie.Value, "No")
	http.SetCookie(w, &cookie)
}

// Sends all the received forms from the database to the admin page, will add verification later to avoid abuse
func SendForms(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT Subject, Email, Body FROM FORMS WHERE Answered = 'No'")
	defer rows.Close()

	var data []formdata
	for rows.Next() {
		var form formdata
		rows.Scan(&form.Subject, &form.Email, &form.Body)

		data = append(data, form)
	}
	json.NewEncoder(w).Encode(data)

}

type formdata struct {
	Subject string `json:"Subject"`
	Email   string `json:"Email"`
	Body    string `json:"Body"`
}

type password struct {
	Password string `json:"Password"`
}

func gencookie() http.Cookie {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	cookie := make([]byte, 15)
	for i := range cookie {
		cookie[i] = charset[rand.Intn(len(charset))]
	}
	return http.Cookie{Name: "ID", Value: string(cookie)}
}
