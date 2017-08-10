package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

type datarow struct {
	S []interface{}
}

type colmeta struct {
	Ai   bool
	Name string
}

func main() {
	// db, err := sql.Open("mysql", "test:test123@/dashdb?charset=utf8")
	db, err := sql.Open("mysql", "test:test123@/dashdb?charset=utf8")
	checkErr(err)

	var countofrows int

	countrows, counterr := db.Query("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'dashdb'")
	checkErr(counterr)

	countrows.Scan(&countofrows)

	// query show tables
	rows, err := db.Query("show tables")
	checkErr(err)

	var myslice = make([]string, countofrows)

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		checkErr(err)
		fmt.Print(name)
		fmt.Print(" ")
		myslice = append(myslice, name)
	}

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.Static("/fonts", "./fonts")
	router.Static("/scripts", "./scripts")
	router.Static("/styles", "./styles")

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":  "Dash Db",
			"test":   "test",
			"tables": myslice,
		})
	})

	router.GET("/addnewdata", func(c *gin.Context) {
		tablename := c.Query("name") // shortcut for c.Request.URL.Query().Get("lastname")

		var countofcols int
		var queryCountStr string

		queryCountStr = "SELECT count(*) from information_schema.columns where table_schema='dashdb' and table_name='" + tablename + "'"

		countrows, counterr := db.Query(queryCountStr)
		checkErr(counterr)
		countrows.Scan(&countofcols)

		queryStr := "SELECT column_name, extra from information_schema.columns where table_schema='dashdb' and table_name='" + tablename + "'"

		// query show tables
		tablecols, err := db.Query(queryStr)
		checkErr(err)

		var mycols = make([]colmeta, 0)

		for tablecols.Next() {
			var column_name string
			var extra string
			var ai bool
			err = tablecols.Scan(&column_name, &extra)
			checkErr(err)

			fmt.Println("extra ne olaki:", extra)

			if strings.HasPrefix(extra, "auto_increment") {
				fmt.Println("ai true")
				ai = true
			} else {
				fmt.Println("ai false")
				ai = false
			}

			//fmt.Print(column_name)
			//fmt.Print(" ")
			var cmeta = colmeta{ai, column_name}
			mycols = append(mycols, cmeta)
		}

		c.HTML(http.StatusOK, "addnewdata.tmpl", gin.H{
			"title":     "Dash Db",
			"test":      "test",
			"tablename": tablename,
			"cols":      mycols,
			"tables":    myslice,
		})
	})

	router.GET("/tabledata", func(c *gin.Context) {
		tablename := c.Query("name") // shortcut for c.Request.URL.Query().Get("lastname")

		var countofcols int
		var queryCountStr string

		queryCountStr = "SELECT count(*) from information_schema.columns where table_schema='dashdb' and table_name='" + tablename + "'"

		countrows, counterr := db.Query(queryCountStr)
		checkErr(counterr)
		countrows.Scan(&countofcols)

		queryStr := "SELECT column_name from information_schema.columns where table_schema='dashdb' and table_name='" + tablename + "'"

		// query show tables
		tablecols, err := db.Query(queryStr)
		checkErr(err)

		var mycols = make([]string, countofcols)

		for tablecols.Next() {
			var column_name string
			err = tablecols.Scan(&column_name)
			checkErr(err)
			//fmt.Print(column_name)
			//fmt.Print(" ")
			mycols = append(mycols, column_name)
		}

		queryDataStr := "SELECT * from " + tablename
		dataRows, err := db.Query(queryDataStr)
		checkErr(err)

		columns, _ := dataRows.Columns()
		count := len(columns)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)

		var mydatas = make([]datarow, 0)

		var indx int

		for dataRows.Next() {

			var valuesStr = make([]interface{}, 0)

			for i, _ := range columns {
				valuePtrs[i] = &values[i]
			}

			dataRows.Scan(valuePtrs...)

			for i, col := range columns {

				var v interface{}

				val := values[i]

				b, ok := val.([]byte)

				if ok {
					v = string(b)
				} else {
					v = val
				}

				fmt.Println("valvecol degerleri", i, col, v, b, val, valuePtrs, values)

				valuesStr = append(valuesStr, v)
			}

			var drow = datarow{valuesStr}

			mydatas = append(mydatas, drow)

			indx = indx + 1
		}

		c.HTML(http.StatusOK, "tabledata.tmpl", gin.H{
			"title":     "Dash Db",
			"test":      "test",
			"tablename": tablename,
			"tables":    myslice,
			"cols":      mycols,
			"datas":     mydatas,
		})
	})

	//http.HandleFunc("/", handler)
	//http.ListenAndServe(":8080", nil)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
}