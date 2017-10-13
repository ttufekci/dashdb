package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
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
	Id int64
	S  []interface{}
}

type colmeta struct {
	Ai    bool
	Name  string
	Value string
	Prim  bool
}

func main() {
	// db, err := sql.Open("mysql", "test:test123@/dashdb?charset=utf8")
	db, err := sql.Open("mysql", "test:test123@/dashdb?charset=utf8")
	checkErr(err)

	var countofrows int

	countrows, counterr := db.Query("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'dashdb'")
	checkErr(counterr)

	fmt.Println("test")

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
		c.String(http.StatusOK, "Hello my friend_t2 %s", name)
	})

	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":  "Dash Db",
			"test":   "test",
			"tables": myslice,
		})
	})

	router.POST("/deletedata", func(c *gin.Context) {
		c.Request.ParseForm()

		tablename := c.PostForm("name") // shortcut for c.Request.URL.Query().Get("lastname")
		id := c.PostForm("id")

		deleteStmt := "delete from " + tablename + " where id=?"

		// delete
		stmt, err := db.Prepare(deleteStmt)
		checkErr(err)

		res, err := stmt.Exec(id)
		checkErr(err)

		affect, err := res.RowsAffected()
		checkErr(err)

		fmt.Println(affect)

		c.Redirect(http.StatusMovedPermanently, "/tabledata?name="+tablename)
	})

	router.GET("/deletedata", func(c *gin.Context) {
		tablename := c.Query("name") // shortcut for c.Request.URL.Query().Get("lastname")
		id := c.Query("id")

		var countofcols int
		var queryCountStr string

		queryCountStr = "SELECT count(*) from information_schema.columns where table_schema='dashdb' and table_name='" + tablename + "'"

		countrows, counterr := db.Query(queryCountStr)
		checkErr(counterr)
		countrows.Scan(&countofcols)

		queryStr := "SELECT column_name, extra, column_key from information_schema.columns where table_schema='dashdb' and table_name='" + tablename + "'"

		// query show tables
		tablecols, err := db.Query(queryStr)
		checkErr(err)

		var mycols = make([]colmeta, 0)

		queryDataStr := "SELECT * from " + tablename + " where id = " + id
		dataRows, err := db.Query(queryDataStr)
		checkErr(err)

		columns, _ := dataRows.Columns()
		count := len(columns)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)

		var mydatas = make([]datarow, 0)

		var valuesStr = make([]interface{}, 0)

		for dataRows.Next() {

			for i := range columns {
				valuePtrs[i] = &values[i]
			}

			var curid int64

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

				if col == "id" {
					curids := v.(string)
					curid, err = strconv.ParseInt(curids, 10, 64)
					checkErr(err)
				}

				valuesStr = append(valuesStr, v)
			}

			var drow = datarow{curid, valuesStr}

			mydatas = append(mydatas, drow)

			break
		}

		indx := 0

		for tablecols.Next() {
			var columnName string
			var extra string
			var column_key string
			var ai bool
			var prim bool
			err = tablecols.Scan(&columnName, &extra, &column_key)
			checkErr(err)

			fmt.Println("extra ne olaki:", extra)

			if strings.HasPrefix(extra, "auto_increment") {
				fmt.Println("ai true")
				ai = true
			} else {
				fmt.Println("ai false")
				ai = false
			}

			if strings.HasPrefix(extra, "PRI") {
				fmt.Println("prim true")
				prim = true
			} else {
				fmt.Println("prim false")
				prim = false
			}

			var ivalue = valuesStr[indx].(string)

			var cmeta = colmeta{ai, columnName, ivalue, prim}
			mycols = append(mycols, cmeta)

			indx++
		}

		c.HTML(http.StatusOK, "deletedata.tmpl", gin.H{
			"title":     "Dash Db",
			"test":      "test",
			"tablename": tablename,
			"cols":      mycols,
			"tables":    myslice,
			"id":        id,
		})
	})

	// function LoadColMetadata(tablename string) []colmeta {
	// 	queryStr := "SELECT column_name, extra, column_key from information_schema.columns where table_schema='dashdb' and table_name='" + tablename + "'"

	// 	// query show tables
	// 	tablecols, err := db.Query(queryStr)
	// 	checkErr(err)

	// 	var mycols = make([]colmeta, 0)

	// 	for tablecols.Next() {
	// 		var column_name string
	// 		var extra string
	// 		var ai bool
	// 		var prim bool
	// 		var column_key string
	// 		err = tablecols.Scan(&column_name, &extra, &column_key)
	// 		checkErr(err)

	// 		fmt.Println("extra ne olaki:", extra)

	// 		if strings.HasPrefix(extra, "auto_increment") {
	// 			fmt.Println("ai true")
	// 			ai = true
	// 		} else {
	// 			fmt.Println("ai false")
	// 			ai = false
	// 		}

	// 		if strings.HasPrefix(column_key, "PRI") {
	// 			fmt.Println("prim true")
	// 			prim = true
	// 		} else {
	// 			fmt.Println("prim false")
	// 			prim = false
	// 		}

	// 		var cmeta = colmeta{ai, column_name, "", prim}
	// 		mycols = append(mycols, cmeta)
	// 	}

	// 	return mycols
	// }

	router.GET("/addnewdata", func(c *gin.Context) {
		tablename := c.Query("name") // shortcut for c.Request.URL.Query().Get("lastname")

		var countofcols int
		var queryCountStr string

		queryCountStr = "SELECT count(*) from information_schema.columns where table_schema='dashdb' and table_name='" + tablename + "'"

		countrows, counterr := db.Query(queryCountStr)
		checkErr(counterr)
		countrows.Scan(&countofcols)

		queryStr := "SELECT column_name, extra, column_key from information_schema.columns where table_schema='dashdb' and table_name='" + tablename + "'"

		// query show tables
		tablecols, err := db.Query(queryStr)
		checkErr(err)

		var mycols = make([]colmeta, 0)

		for tablecols.Next() {
			var column_name string
			var extra string
			var ai bool
			var prim bool
			var column_key string
			err = tablecols.Scan(&column_name, &extra, &column_key)
			checkErr(err)

			fmt.Println("extra ne olaki:", extra)

			if strings.HasPrefix(extra, "auto_increment") {
				fmt.Println("ai true")
				ai = true
			} else {
				fmt.Println("ai false")
				ai = false
			}

			if strings.HasPrefix(column_key, "PRI") {
				fmt.Println("prim true")
				prim = true
			} else {
				fmt.Println("prim false")
				prim = false
			}

			var cmeta = colmeta{ai, column_name, "", prim}
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

	router.POST("/addnewdata", func(c *gin.Context) {
		c.Request.ParseForm()

		tablename := c.PostForm("name") // shortcut for c.Request.URL.Query().Get("lastname")
		message := c.PostForm("message")
		fields := c.Request.Form["fields"]

		queryStr := "SELECT column_name, extra, column_key from information_schema.columns where table_schema='dashdb' and table_name='" + tablename + "'"

		// query show tables
		tablecols, err := db.Query(queryStr)
		checkErr(err)

		var mycols = make([]colmeta, 0)

		colstr := "("

		for tablecols.Next() {
			var column_name string
			var extra string
			var ai bool
			var prim bool
			var column_key string
			err = tablecols.Scan(&column_name, &extra, &column_key)
			checkErr(err)

			fmt.Println("extra ne olaki:", extra)

			if strings.HasPrefix(extra, "auto_increment") {
				fmt.Println("ai true")
				ai = true
			} else {
				fmt.Println("ai false")
				ai = false
			}

			if strings.HasPrefix(extra, "PRI") {
				fmt.Println("prim true")
				prim = true
			} else {
				fmt.Println("prim false")
				prim = false
			}

			//fmt.Print(column_name)
			//fmt.Print(" ")
			var cmeta = colmeta{ai, column_name, "", prim}
			mycols = append(mycols, cmeta)

			if column_name != "id" {
				colstr = colstr + column_name + ","
			}
		}

		colstr = strings.TrimSuffix(colstr, ",") + ")"

		insertStr := "insert into " + tablename + " " + colstr + " values ('" + strings.Join(fields[:], "','") + "')"

		// insert
		stmt, err := db.Prepare(insertStr)
		checkErr(err)

		res, err := stmt.Exec()
		checkErr(err)

		id, err := res.LastInsertId()
		checkErr(err)

		fmt.Println(id)

		c.HTML(http.StatusOK, "addnewdata.tmpl", gin.H{
			"status":    "posted",
			"message":   message,
			"title":     "Dash Db",
			"test":      "test",
			"tablename": tablename,
			"cols":      mycols,
			"fields":    fields,
			"tables":    myslice,
		})
	})

	router.POST("/editdata", func(c *gin.Context) {
		c.Request.ParseForm()

		tablename := c.PostForm("name") // shortcut for c.Request.URL.Query().Get("lastname")
		fields := c.Request.Form["fields"]
		id := c.PostForm("id")

		queryStr := "SELECT column_name, extra, column_key from information_schema.columns where table_schema='dashdb' and table_name='" + tablename + "'"

		// query show tables
		tablecols, err := db.Query(queryStr)
		checkErr(err)

		var mycols = make([]colmeta, 0)

		colstr := "set "

		indx := 0

		fmt.Println("fields:", strings.Join(fields[:], "','"))

		fmt.Println("fields[0]", fields[0])

		for tablecols.Next() {
			var column_name string
			var extra string
			var column_key string
			var ai bool
			var prim bool
			err = tablecols.Scan(&column_name, &extra, &column_key)
			checkErr(err)

			fmt.Println("extra ne olaki:", extra)

			if strings.HasPrefix(extra, "auto_increment") {
				fmt.Println("ai true")
				ai = true
			} else {
				fmt.Println("ai false")
				ai = false
			}

			if strings.HasPrefix(column_key, "PRI") {
				fmt.Println("prim true")
				prim = true
			} else {
				fmt.Println("prim false")
				prim = false
			}

			var cmeta = colmeta{ai, column_name, "", prim}
			mycols = append(mycols, cmeta)

			fmt.Println("indx", indx, "fields:", fields[indx])

			columnValue := fields[indx]

			if column_name != "id" {
				colstr = colstr + column_name + "='" + columnValue + "',"
				indx++
			}
		}

		colstr = strings.TrimSuffix(colstr, ",")

		updateStr := "update " + tablename + " " + colstr + " " + "where id=" + id

		fmt.Println("updatestr:", updateStr)

		// update
		stmt, err := db.Prepare(updateStr)
		checkErr(err)

		res, err := stmt.Exec()
		checkErr(err)

		affect, err := res.RowsAffected()
		checkErr(err)

		fmt.Println(affect)

		c.HTML(http.StatusOK, "editdata.tmpl", gin.H{
			"title":     "Dash Db",
			"test":      "test",
			"tablename": tablename,
			"cols":      mycols,
			"tables":    myslice,
		})
	})

	router.GET("/editdata", func(c *gin.Context) {
		tablename := c.Query("name") // shortcut for c.Request.URL.Query().Get("lastname")
		id := c.Query("id")

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

		queryDataStr := "SELECT * from " + tablename + " where id = " + id
		dataRows, err := db.Query(queryDataStr)
		checkErr(err)

		columns, _ := dataRows.Columns()
		count := len(columns)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)

		var mydatas = make([]datarow, 0)

		var valuesStr = make([]interface{}, 0)

		for dataRows.Next() {

			for i := range columns {
				valuePtrs[i] = &values[i]
			}

			var curid int64

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

				if col == "id" {
					curids := v.(string)
					curid, err = strconv.ParseInt(curids, 10, 64)
					checkErr(err)
				}

				valuesStr = append(valuesStr, v)
			}

			var drow = datarow{curid, valuesStr}

			mydatas = append(mydatas, drow)

			break
		}

		indx := 0

		for tablecols.Next() {
			var columnName string
			var extra string
			var column_key string
			var ai bool
			var prim bool
			err = tablecols.Scan(&columnName, &extra, &column_key)
			checkErr(err)

			fmt.Println("extra ne olaki:", extra)

			if strings.HasPrefix(extra, "auto_increment") {
				fmt.Println("ai true")
				ai = true
			} else {
				fmt.Println("ai false")
				ai = false
			}

			if strings.HasPrefix(column_key, "PRI") {
				fmt.Println("prim true")
				prim = true
			} else {
				fmt.Println("prim false")
				prim = false
			}

			var ivalue = valuesStr[indx].(string)

			var cmeta = colmeta{ai, columnName, ivalue, prim}
			mycols = append(mycols, cmeta)

			indx++
		}

		c.HTML(http.StatusOK, "editdata.tmpl", gin.H{
			"title":     "Dash Db",
			"test":      "test",
			"tablename": tablename,
			"cols":      mycols,
			"tables":    myslice,
			"id":        id,
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

		queryStr := "SELECT column_name, extra, column_key from information_schema.columns where table_schema='dashdb' and table_name='" + tablename + "'"

		// query show tables
		tablecols, err := db.Query(queryStr)
		checkErr(err)

		// var mycols = make([]string, countofcols)

		var mycols = make([]colmeta, 0)

		for tablecols.Next() {
			var column_name string
			var extra string
			var column_key string
			var ai bool
			var prim bool
			err = tablecols.Scan(&column_name, &extra, &column_key)
			checkErr(err)

			fmt.Println("extra ne olaki:", extra)

			if strings.HasPrefix(extra, "auto_increment") {
				fmt.Println("ai true")
				ai = true
			} else {
				fmt.Println("ai false")
				ai = false
			}

			if strings.HasPrefix(column_key, "PRI") {
				fmt.Println("prim true")
				prim = true
			} else {
				fmt.Println("prim false")
				prim = false
			}

			var cmeta = colmeta{ai, column_name, "", prim}
			mycols = append(mycols, cmeta)
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

			for i := range columns {
				valuePtrs[i] = &values[i]
			}

			var curid int64

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

				if col == "id" {
					curids := v.(string)
					curid, err = strconv.ParseInt(curids, 10, 64)
					checkErr(err)
				}

				valuesStr = append(valuesStr, v)
			}

			var drow = datarow{curid, valuesStr}

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

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
}
