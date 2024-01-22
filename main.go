package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Users struct {
	Users []User `json:"users"`
}

var usersFile = "data/users.json"
var lstUsers Users
var tmpl *template.Template

func init() {
	// Corrected template path
	tmpl = template.Must(template.ParseFiles("templates/home.gohtml", "templates/p1.gohtml"))
}

func ReadJSON() (Users, error) {
	jsonFile, err := os.ReadFile("data/users.json")
	if err != nil {
		fmt.Println("Error reading", err.Error())
		return Users{}, err
	}
	var jsonData Users
	err = json.Unmarshal(jsonFile, &jsonData)
	return jsonData, err
}

func EditJSON(ModifiedArticle Users) {
	modifiedJSON, errMarshal := json.Marshal(ModifiedArticle)
	if errMarshal != nil {
		fmt.Println("Error encoding ", errMarshal.Error())
		return
	}

	// Write the modified JSON to the file
	err := os.WriteFile("data/users.json", modifiedJSON, 0644)
	if err != nil {
		fmt.Println("Error writing modified JSON file:", err)
		return
	}
}

func loadUsers() (Users, error) {
	file, err := ioutil.ReadFile(usersFile)
	if err != nil {
		log.Printf("Error reading file %s: %v", usersFile, err)
		return Users{}, err
	}
	var users Users
	err = json.Unmarshal(file, &users)
	if err != nil {
		log.Printf("Error unmarshaling users: %v", err)
		return Users{}, err
	}
	return users, nil
}

func saveUsers(users Users) error {
	file, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(usersFile, file, 0644)
}

func addUser(username, password string) error {
	users, err := loadUsers()
	if err != nil {
		return err
	}
	users.Users = append(users.Users, User{Username: username, Password: password})
	return saveUsers(users)
}

func checkCredentials(username, password string) bool {
	users, err := loadUsers()
	if err != nil {
		log.Println(err)
		return false
	}
	for _, user := range users.Users {
		if user.Username == username && user.Password == password {
			return true
		}
	}
	return false
}

func homePage(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "login", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if checkCredentials(username, password) {
		// Redirect to p1.gohtml
		http.Redirect(w, r, "/p1", http.StatusSeeOther)
	} else {
		fmt.Fprintf(w, "Login failed.")
	}
}

func renderP1Page(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var person User

	person.Username = r.FormValue("new_username")
	person.Password = r.FormValue("new_password")
	hashedPassword, _ := hashPassword(person.Password)

	// Vérifier si le nom d'utilisateur existe déjà
	if usernameExists(person.Username) {
		fmt.Fprintf(w, "Nom d'utilisateur existe déjà.")
		return
	}

	err := addUser(person.Username, hashedPassword)
	if err != nil {
		fmt.Fprintf(w, "Error during signup: %v", err)
		return
	}
	fmt.Fprintf(w, "Signup successful!")
}

func usernameExists(username string) bool {
	users, err := loadUsers()
	if err != nil {
		log.Println(err)
		return false
	}
	for _, user := range users.Users {
		if user.Username == username {
			return true
		}
	}
	return false
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func login(username, password string) bool {

	storedHash := ""

	return checkPasswordHash(password, storedHash)
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/p1", renderP1Page)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/signup", handleSignup)
	rootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(rootDoc + "/styles"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))
	log.Fatal(http.ListenAndServe(":8084", nil))
}
