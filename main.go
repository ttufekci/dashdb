package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := sql.Open("mysql", "test:test123@/dashdb?charset=utf8")
	checkErr(err)

	// query
	rows, err := db.Query("SELECT * FROM testtbl")
	checkErr(err)

	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		checkErr(err)
		fmt.Print(id, " ")
		fmt.Println(name)
	}

	// query show tables
	rows, err = db.Query("show tables")
	checkErr(err)

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		checkErr(err)
		fmt.Print(name)
		fmt.Print(" ")
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
