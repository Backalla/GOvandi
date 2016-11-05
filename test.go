package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
  "time"
  
)

func checkCookie(r *http.Request) string {
  uid, err := r.Cookie("uid")
  // log.Println("Uid = ",uid)
  if err != nil{
  	return ""
  } else {
  	return uid.Value
  }

}

func checkErr(err error, s string) {
	if err != nil {
		log.Fatal(err, s)
	}
}

func get_contentType(path string) string {
	var contentType string
	_, err := ioutil.ReadFile(string(path))
	if err == nil {
		if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".html") {
			contentType = "text/html"
		} else if strings.HasSuffix(path, ".js") {
			contentType = "application/javascript"
		} else if strings.HasSuffix(path, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(path, ".gif") {
			contentType = "image/gif"
		} else {
			contentType = "text/plain"
		}
	}
	return contentType
}

func login_handler(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.Method)
	// log.Println("login_handler")
	if r.Method == "POST" {
		var uid string
		db, err := sql.Open("mysql", "root:@/godb")
		checkErr(err, "Cannot connect to db login_handler")
		username := r.FormValue("username")
		password := r.FormValue("password")
		// log.Println(username,password)
		err = db.QueryRow("select uid from users where username = ? and password=?", username, password).Scan(&uid)
		if err != nil {
			log.Println("Not found")
		} else {
      expiration := time.Now().Add(time.Hour)
      cookie := http.Cookie{Name: "uid", Value: uid, Expires: expiration}
      http.SetCookie(w, &cookie)
      http.Redirect(w, r, "/", 302)
			log.Println(uid)
		}
	} else {
		path := r.URL.Path[1:]
		// log.Println(path)
		contentType := get_contentType(path)
		data, _ := ioutil.ReadFile(string(path))
		w.Header().Add("Content Type", contentType)
		w.Write(data)
		t, _ := template.ParseFiles("pages/login.html")
		t.Execute(w, nil)
	}
}

func logout_handler(w http.ResponseWriter, r *http.Request) {
  cookie := http.Cookie{Name: "uid", Value: "", Expires: time.Now(), MaxAge: -1}
  http.SetCookie(w,&cookie)
  http.Redirect(w, r, "/", 302)
}


func home_handler(w http.ResponseWriter, r *http.Request) {
  // cookie,_ := r.Cookie("uid")
  // log.Println(r.Cookies())
	type Bus struct {
		Id          int
		Name        string
		Source      string
		Destination string
		Price       int
	}

	type User struct {
		Uid  string
		Name string
	}

	type Tdata struct {
		U     User
		Buses []Bus
	}
	// log.Println("home_handler")
	db, err := sql.Open("mysql", "root:@/godb")
	checkErr(err, "Cannot connect to db home_handler")
	// log.Println("started db")

	rows, err := db.Query("select id,name,source,destination,price from buses")
	checkErr(err, "Cannot fetch from buses")
	var (
		id          int
		name        string
		source      string
		destination string
		price       int
	)
	defer rows.Close()
	// var buses []Bus
	var tdata Tdata
	for rows.Next() {
		err := rows.Scan(&id, &name, &source, &destination, &price)
		checkErr(err, "Cannot scan rows..")
		bus := Bus{}
		bus.Id = id
		bus.Name = name
		bus.Source = source
		bus.Destination = destination
		bus.Price = price
		tdata.Buses = append(tdata.Buses, bus)
	}
	// log.Println("buses are -> ",buses)
  uid := checkCookie(r)
  log.Println("UID=",uid)
  var loggedInName string
	err = db.QueryRow("select name from users where uid = ?", uid).Scan(&loggedInName)
	user := User{Name: loggedInName, Uid: uid}
	tdata.U = user
	path := r.URL.Path[1:]
	// log.Println(path)
	contentType := get_contentType(path)
	data, err := ioutil.ReadFile(string(path))
	w.Header().Add("Content Type", contentType)
	w.Write(data)
	t, _ := template.ParseFiles("pages/home.html")
	t.Execute(w, tdata)
}

func main() {
	http.HandleFunc("/login", login_handler)
	http.HandleFunc("/logout", logout_handler)
	http.HandleFunc("/", home_handler)
	http.ListenAndServe(":8080", nil)
	// fmt.Println("Hello, Mothafaka")
}
