package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"

	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	date    string
	name    string
	phone   string
	message string
	score   string
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// begin web page
	var htmlHeader = "<!DOCTYPE html><html><head><style>table {font-family: arial, sans-serif;border-collapse: collapse;width: 100%;}td, th {border: 1px solid #dddddd;text-align: left;padding: 8px;}tr:nth-child(even) {background-color: #dddddd;}</style></head><body>"
	fmt.Fprintf(w, htmlHeader)
	var hostname = getHostname()
	var gitSHA = os.Getenv("GIT_SHA")
	var appversion = "2.5.0"
	fmt.Fprintf(w, "<h1>Golang Guestbook (v%s)</h1><p>Hostname: %s</p><p>Git: %s</p><table><tr><th>Date</th><th>Name</th><th>Phone</th><th>Sentiment</th><th>Message</th></tr>", appversion, hostname, gitSHA)

	// query DB and loop through rows
	var connString = getConnectString()
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()

	rows, err := conn.Query("select * from guestlog order by entrydate DESC")
	if err != nil {
		log.Fatal("Cannot query: ", err.Error())
		return
	}
	defer rows.Close()

	// loop through result and build table
	for rows.Next() {
		err := rows.Scan(&date, &name, &phone, &message, &score)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "<tr><td>"+date+"</td><td>"+name+"</td><td>"+phone+"</td><td>"+score+"</td><td>"+message+"</td></tr>")
	}
	fmt.Fprintf(w, "</table>")
}

// HealthCheckHandler e.g. http.HandleFunc("/health-check", HealthCheckHandler)
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func testHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/html")
	resp.WriteHeader(http.StatusOK)
	fmt.Fprint(resp, "OK")
}

func testdbHandler(resp http.ResponseWriter, req *http.Request) {
	var connString = getConnectString()
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()
	resp.Header().Add("Content-Type", "text/html")
	resp.WriteHeader(http.StatusOK)
	fmt.Fprint(resp, "OK")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/testdb", testdbHandler)
	http.ListenAndServe(":8080", nil)
}

func getHostname() string {
	var result string
	localhostname, err := os.Hostname()

	if err != nil {
		result = "ERROR: Cannot find server hostname"
	} else {
		result = localhostname
	}
	return result
}

func getConnectString() string {
	var result string

	var sqlserver = os.Getenv("SQLSERVER")
	if sqlserver == "" {
		sqlserver = "13.82.145.183"
	}
	var sqlport = os.Getenv("SQLPORT")
	if sqlport == "" {
		sqlport = "10433"
	}
	var sqlid = os.Getenv("SQLID")
	if sqlid == "" {
		sqlid = "sa"
	}
	var sqlpwd = os.Getenv("SQLPWD")
	if sqlpwd == "" {
		sqlpwd = "Your@Password"
	}
	var sqldb = os.Getenv("SQLDB")
	if sqldb == "" {
		sqldb = "sql_guestbook"
	}
	result = "server=" + sqlserver + ";port=" + sqlport + ";user id=" + sqlid + ";password=" + sqlpwd + ";database=" + sqldb + ";connection timeout=45"
	return result
}
