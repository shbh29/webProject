package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"text/template"
)

type Emp struct {
	Id   int
	Name string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	userName := "root"
	password := "Beni$on123"
	dbName := "godb"
	db, err := sql.Open(dbDriver, userName+":"+password+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func List(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close() // will execute at the end of the method
	selRs, err := db.Query("select * from emp")
	if err != nil {
		panic(err.Error())
	}
	res := []Emp{}
	emp := Emp{}
	for selRs.Next() {
		// Use direct Emp Fields in with address of Operator
		err = selRs.Scan(&emp.Id, &emp.Name)
		if err != nil {
			panic(err.Error())
		}
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "List", res)
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()
	if r.Method == "POST" {
		name := r.FormValue("name")
		stmt, err := db.Prepare("insert into emp(name) values(?)")
		if err != nil {
			panic(err.Error())
		}
		stmt.Exec(name)
		log.Println("Inserted Name:" + name)
	}
	http.Redirect(w, r, "/", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()
	id := r.URL.Query().Get("id")
	selRs, err := db.Query("select * from emp where id=?", id)
	if err != nil {
		panic(err.Error())
	}
	emp := Emp{}
	if selRs.Next() {
		var uid int
		var name string
		err := selRs.Scan(&uid, &name)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = uid
		emp.Name = name
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		stmt, err := db.Prepare("update emp set name=? where id=?")
		if err != nil {
			panic(err.Error())
		}
		stmt.Exec(name, id)
		log.Println("Updated name: " + name)
	}
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	id := r.URL.Query().Get("id")

	stmt, err := db.Prepare("delete from emp where id=?")
	if err != nil {
		panic(err.Error())
	}
	stmt.Exec(id)
	log.Println("Deleted name: " + id)

	http.Redirect(w, r, "/", 301)
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World\n")
}

func main() {
	http.HandleFunc("/hello", handleFunc)
	http.HandleFunc("/", List)
	http.HandleFunc("/new", New)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":5000", nil)
}
