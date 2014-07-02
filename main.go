package main

import (
	"fmt"
	"github.com/hoisie/redis"
	"net/http"
	"html/template"
	"log"
)

var (
	client redis.Client
)

type Profile struct {
	User string
}

func setupDB() {
    client.Addr = "192.168.3.241:16696"
    client.Db = 0
    client.Password = "ajul827gtcj4z9hbd27wl1jgpu69vogd"
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In new handler!")
	t := template.New("new.html")
	t.ParseFiles("templates/new.html")
	t.Execute(w, t)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In save handler!")
	user := r.FormValue("user")
	err := client.Rpush("user", []byte(user))

	fmt.Println("Redis Error:", err)
	http.Redirect(w, r, "/", http.StatusFound)
}

func indexHanlder(w http.ResponseWriter, r *http.Request){
	fmt.Println("In get data.")
	username, _ := client.Lrange("user", 0, 10000)
	fmt.Println("User--->", username)

	profiles := []Profile{}
	for i, v := range username {
		println("Name", i,":",string(v))

		profile := Profile{}
		profile.User = string(v)
		profiles = append(profiles, profile)
	}

	t := template.New("index.html")
	t.ParseFiles("templates/index.html")
	fmt.Println("Final user================>", profiles)
	t.Execute(w, profiles)
}

func main() {
	setupDB()
	fmt.Println("Client: ", client)
	http.HandleFunc("/", indexHanlder)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/new/", newHandler)

	http.Handle("/public/css/", http.StripPrefix("/public/css/", http.FileServer(http.Dir("public/css"))))
	http.Handle("/public/images/", http.StripPrefix("/public/images/", http.FileServer(http.Dir("public/images"))))
	if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		log.Fatalf("Error in listening:", err)
	}
}
