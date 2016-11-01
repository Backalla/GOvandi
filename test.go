package main
import (
      "html/template"
      "io/ioutil"
      "net/http"
      _ "github.com/go-sql-driver/mysql"
      "database/sql"
      "log"
      "strings"
)
type Bus struct{
  Id int
  Name string
  Source string
  Destination string
  Price int
}

func checkErr(err error,s string) {
  if err != nil {
    log.Fatal(err,s)
  }
}

func home_handler(w http.ResponseWriter, r *http.Request) {
  db, err := sql.Open("mysql", "root:@/godb")
  checkErr(err,"Cannot connect to db")
  log.Println("started db")

  rows, err := db.Query("select id,name,source,destination,price from buses")
  checkErr(err,"Cannot fetch from buses")
  var (
    id int
    name string
    source string
    destination string
    price int
  )
  defer rows.Close()
  var buses []Bus

  for rows.Next() {
    err := rows.Scan(&id, &name,&source,&destination,&price)
    checkErr(err,"Cannot scan rows..")
    bus := Bus{}
    bus.Id=id
    bus.Name=name
    bus.Source=source
    bus.Destination=destination
    bus.Price=price
    buses=append(buses,bus)
  }
  log.Println("buses are -> ",buses)


  path := r.URL.Path[1:]
  log.Println(path)
  var contentType string
  data,err := ioutil.ReadFile(string(path))
  if err==nil {
    if strings.HasSuffix(path,".css"){
      contentType="text/css"
    } else if strings.HasSuffix(path,".html"){
      contentType="text/html"
    } else if strings.HasSuffix(path,".js"){
      contentType="application/javascript"
    } else if strings.HasSuffix(path,".png"){
      contentType="image/png"
    } else  if strings.HasSuffix(path,".gif"){
      contentType="image/gif"
    } else {
      contentType="text/plain"
    }
  }


  w.Header().Add("Content Type",contentType)
  w.Write(data)
  t, _ := template.ParseFiles("pages/home.html")
  t.Execute(w, buses)
}




func main(){
  
  http.HandleFunc("/",home_handler)
  http.ListenAndServe(":8080",nil )
  // fmt.Println("Hello, Mothafaka")
}