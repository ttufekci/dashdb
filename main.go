package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	cors "github.com/itsjamie/gin-cors"
	"github.com/theherk/viper"
)

type config struct {
	host     string
	schema   string
	user     string
	password string
}

type Table struct {
	Name string `json:"name"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func replacesc(ids string) string {
	var newids string

	newids = strings.Replace(ids, "éééé", "", -1)

	newids = strings.Replace(newids, "éé", " and ", -1)

	newids = strings.Replace(newids, "é", " = ", -1)

	return newids
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

var db *sql.DB

var tablelist []string

func main() {
	var countofrows int

	var countrows, rows *sql.Rows
	var counterr, err error

	var host, user, password, schema string
	host = "localhost"
	user = "test"
	password = "test123"
	schema = "dashdb"

	succ, _ := testConnection(host, user, password, schema)

	if !succ {
		goto SkipDBInit
	}

	db, err = sql.Open("mysql", "test:test123@/dashdb?charset=utf8")
	checkErr(err)

	countrows, counterr = db.Query("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'dashdb'")
	checkErr(counterr)

	fmt.Println("test")

	countrows.Scan(&countofrows)

	// query show tables
	rows, err = db.Query("show tables")
	checkErr(err)

	tablelist = make([]string, countofrows)

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		checkErr(err)
		fmt.Print(name)
		fmt.Print(" ")
		tablelist = append(tablelist, name)
	}

SkipDBInit:

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	// Apply the middleware to the router (works with groups too)
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

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
			"tables": tablelist,
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
			"tables":    tablelist,
			"id":        id,
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
			"tables":    tablelist,
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
			"tables":    tablelist,
		})
	})

	router.POST("/editdata", func(c *gin.Context) {
		c.Request.ParseForm()

		tablename := c.PostForm("name") // shortcut for c.Request.URL.Query().Get("lastname")
		fields := c.Request.Form["fields"]
		id := c.PostForm("id")
		ids := c.PostForm("ids")

		ids = replacesc(ids)

		fmt.Println("\nids:", ids, ",id:", id)

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

		var updateStr string

		if id == "0" {
			updateStr = "update " + tablename + " " + colstr + " " + "where " + ids
		} else {
			updateStr = "update " + tablename + " " + colstr + " " + "where id=" + id
		}

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
			"tables":    tablelist,
		})
	})

	router.GET("/editdata", func(c *gin.Context) {
		tablename := c.Query("name") // shortcut for c.Request.URL.Query().Get("lastname")
		id := c.Query("id")
		primcols := c.Query("primcols")
		ids := c.Query("ids")

		ids = replacesc(ids)

		fmt.Println("primcols:", primcols, ",ids:", ids, ",id:", id)

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
		var queryDataStr string

		if id == "0" {
			fmt.Println("0")
			queryDataStr = "SELECT * from " + tablename + " where " + ids
		} else {
			fmt.Println("0 degil")
			queryDataStr = "SELECT * from " + tablename + " where id = " + id
		}

		fmt.Println("queryDataStr", queryDataStr)

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
			"tables":    tablelist,
			"id":        id,
			"ids":       ids,
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

		primQueryStr := "SELECT column_name, extra, column_key from information_schema.columns where table_schema='dashdb' and column_key = 'PRI' and table_name='" + tablename + "'"

		// query show tables
		primarycols, err := db.Query(primQueryStr)
		checkErr(err)

		primcols := ""

		// sliceprimcols := make([]string, 5)

		primcolsmap := make(map[string]string)

		for primarycols.Next() {
			var column_name string
			var extra string
			var column_key string
			err = primarycols.Scan(&column_name, &extra, &column_key)
			checkErr(err)
			primcols = primcols + column_name + ","
			// append(sliceprimcols, column_name)
			primcolsmap[column_name] = ""
		}

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
		// var ids string

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

				if _, ok := primcolsmap[col]; ok {
					//do something here
					primcolsmap[col] = v.(string)
				}

				valuesStr = append(valuesStr, v)
			}

			var drow = datarow{curid, valuesStr}

			mydatas = append(mydatas, drow)

			indx = indx + 1
		}

		var ids string

		ids = "éé"

		for k, v := range primcolsmap {
			fmt.Println("k:", k, "v:", v)
			ids += "éé" + k + "é" + v
		}

		c.HTML(http.StatusOK, "tabledata.tmpl", gin.H{
			"title":     "Dash Db",
			"test":      "test",
			"tablename": tablename,
			"primcols":  primcols,
			"tables":    tablelist,
			"cols":      mycols,
			"datas":     mydatas,
			"ids":       ids,
		})
	})

	router.GET("/tables", func(c *gin.Context) {
		myarray := make([]Table, 0)
		mytable := Table{"table1"}
		myarray = append(myarray, mytable)
		mytable = Table{"table2"}
		myarray = append(myarray, mytable)
		mytable = Table{"table3"}
		myarray = append(myarray, mytable)

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")

		c.JSON(200, myarray)
	})

	router.GET("/testconnection", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")

		host := c.Query("host")
		user := c.Query("user")
		password := c.Query("password")
		schema := c.Query("schema")

		fmt.Println(host, user, password, schema)

		succ, testerr := testConnection(host, user, password, schema)

		if succ {
			c.JSON(200, gin.H{
				"success": "true",
				"message": "Ping to database successful, connection is still alive",
			})
		} else {
			c.JSON(200, gin.H{
				"success": "false",
				"message": testerr,
			})
		}
	})

	router.GET("/saveconfig", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")

		host := c.Query("host")
		user := c.Query("user")
		password := c.Query("password")
		schema := c.Query("schema")

		viper.SetConfigName("config") // name of config file (without extension)
		viper.AddConfigPath(".")      // optionally look for config in the working directory
		err := viper.ReadInConfig()   // Find and read the config file

		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}

		fmt.Println("")
		fmt.Println("database.user", viper.GetString("database.user"))
		fmt.Println("database.host", viper.GetString("database.host"))

		viper.Set("database.user", user)
		viper.Set("database.host", host)
		viper.Set("database.password", password)
		viper.Set("database.schema", schema)

		viper.WriteConfig()

		c.JSON(200, gin.H{
			"success": "true",
			"message": "saved successfully",
		})
	})

	router.GET("/readconfig", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")

		viper.SetConfigName("config") // name of config file (without extension)
		viper.AddConfigPath(".")      // optionally look for config in the working directory
		err := viper.ReadInConfig()   // Find and read the config file

		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

		c.JSON(200, gin.H{
			"success":  "true",
			"message":  "saved successfully",
			"user":     viper.GetString("database.user"),
			"password": viper.GetString("database.password"),
			"schema":   viper.GetString("database.schema"),
			"host":     viper.GetString("database.host"),
		})
	})

	router.GET("/tablelist", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")

		c.JSON(200, tablelist)
	})

	router.GET("/columnlist", func(c *gin.Context) {
		tablename := c.Query("name")

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

		primQueryStr := "SELECT column_name, extra, column_key from information_schema.columns where table_schema='dashdb' and column_key = 'PRI' and table_name='" + tablename + "'"

		// query show tables
		primarycols, err := db.Query(primQueryStr)
		checkErr(err)

		primcols := ""

		// sliceprimcols := make([]string, 5)

		primcolsmap := make(map[string]string)

		for primarycols.Next() {
			var column_name string
			var extra string
			var column_key string
			err = primarycols.Scan(&column_name, &extra, &column_key)
			checkErr(err)
			primcols = primcols + column_name + ","
			// append(sliceprimcols, column_name)
			primcolsmap[column_name] = ""
		}

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
		// var ids string

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

				if _, ok := primcolsmap[col]; ok {
					//do something here
					primcolsmap[col] = v.(string)
				}

				valuesStr = append(valuesStr, v)
			}

			var drow = datarow{curid, valuesStr}

			mydatas = append(mydatas, drow)

			indx = indx + 1
		}

		var ids string

		ids = "éé"

		for k, v := range primcolsmap {
			fmt.Println("k:", k, "v:", v)
			ids += "éé" + k + "é" + v
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")

		c.JSON(200, gin.H{
			"title":     "Dash Db",
			"test":      "test",
			"tablename": tablename,
			"primcols":  primcols,
			"tables":    tablelist,
			"cols":      mycols,
			"datas":     mydatas,
			"ids":       ids,
		})
	})

	router.GET("/reditdata", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		tablename := c.Query("name") // shortcut for c.Request.URL.Query().Get("lastname")
		id := c.Query("id")

		ids := c.Query("ids")

		ids = replacesc(ids)

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
		var queryDataStr string

		if id == "0" {
			queryDataStr = "SELECT * from " + tablename + " where " + ids
		} else {
			queryDataStr = "SELECT * from " + tablename + " where id = " + id
		}

		fmt.Println("queryDataStr", queryDataStr)

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

			if strings.HasPrefix(extra, "auto_increment") {
				ai = true
			} else {
				ai = false
			}

			if strings.HasPrefix(column_key, "PRI") {
				prim = true
			} else {
				prim = false
			}

			var ivalue = valuesStr[indx].(string)

			var cmeta = colmeta{ai, columnName, ivalue, prim}
			mycols = append(mycols, cmeta)

			indx++
		}

		c.JSON(200, gin.H{
			"title":     "Dash Db",
			"test":      "test",
			"tablename": tablename,
			"cols":      mycols,
			"tables":    tablelist,
			"id":        id,
			"ids":       ids,
		})
	})

	// LoginJSON stuff
	type editdata struct {
		Name   string   `json:"name" binding:"required"`
		Id     string   `json:"id" binding:"required"`
		Ids    string   `json:"ids" binding:"required"`
		Fields []string `json:"fields" binding:"required"`
	}

	router.POST("/reditdatam", func(c *gin.Context) {
		var json editdata
		c.BindJSON(&json)

		tablename := json.Name
		fields := json.Fields
		id := json.Id
		ids := json.Ids

		fmt.Println("beforeids", ids)

		ids = replacesc(ids)

		fmt.Println("\nids:", ids, ",id:", id, "tablename:", tablename)

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

			indx++

			if column_name != "id" {
				colstr = colstr + column_name + "='" + columnValue + "',"
			}
		}

		colstr = strings.TrimSuffix(colstr, ",")

		var updateStr string

		if id == "0" {
			updateStr = "update " + tablename + " " + colstr + " " + "where " + ids
		} else {
			updateStr = "update " + tablename + " " + colstr + " " + "where id=" + id
		}

		fmt.Println("updatestr:", updateStr)

		// update
		stmt, err := db.Prepare(updateStr)
		checkErr(err)

		res, err := stmt.Exec()
		checkErr(err)

		affect, err := res.RowsAffected()
		checkErr(err)

		fmt.Println(affect)

		c.JSON(200, gin.H{
			"title":     "Dash Db",
			"test":      "test",
			"tablename": tablename,
			"cols":      mycols,
			"tables":    tablelist,
		})
	})

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run(":8081")
}

func testConnection(host string, user string, password string, schema string) (succ bool, errStr string) {
	var sqlConnStr string
	var calcHost string
	if host == "" || host == "localhost" {
		calcHost = ""
	} else {
		calcHost = host
	}

	sqlConnStr += user + ":" + password + "@" + calcHost + "/" + schema + "?charset=utf8"
	fmt.Println("sqlconstr: ", sqlConnStr)

	conn, err := sql.Open("mysql", sqlConnStr)
	defer conn.Close()

	if err != nil {
		fmt.Println(err)
		return false, err.Error()
	}

	err = conn.Ping()
	if err != nil {
		fmt.Println(err)
		return false, err.Error()
	}

	return true, ""
}
