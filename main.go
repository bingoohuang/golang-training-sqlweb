package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"

	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var (
	port      = flag.String("port", "8080", "Port to serve.")
	rowsLimit = flag.Int("limit", 50, "Max number of rows to return.")
)

var templates *template.Template
var db *sql.DB
var err error

func main() {
	flag.Parse()

	db, err = sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}

	templates, err = template.ParseFiles("index.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", indexHtml)
	http.HandleFunc("/execute", executeServer)

	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatal(err)
	}
}

func indexHtml(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "index", nil); err != nil {
		http.Error(w, "failed to build page: "+err.Error(), http.StatusInternalServerError)
	}
}

func executeServer(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "error parsing form: "+err.Error(), http.StatusBadRequest)
	}
	data := executeSql(r.PostForm.Get("sql"))

	if err := templates.ExecuteTemplate(w, "index", data); err != nil {
		http.Error(w, "failed to build page: "+err.Error(), http.StatusInternalServerError)
	}
}

// Queries query and returns pageData with results.
func executeSql(query string) *PageData {
	log.Printf("executing: %s", query)
	rows, err := db.Query(query)
	if err != nil {
		return &PageData{
			Query:   query,
			Results: QueryResults{Error: err},
		}
	}
	return &PageData{
		Query:   query,
		Results: *NewQueryResults(rows, *rowsLimit),
	}
}

// Holds data for the webpage template.
type PageData struct {
	Query   string
	Results QueryResults
}

// Holds results from a query.
type QueryResults struct {
	Error   error
	Columns []string
	Data    [][]string
}

// Converts sql.Rows into QueryResults.
func NewQueryResults(rows *sql.Rows, rowLimit int) *QueryResults {
	columns, err := rows.Columns()
	if err != nil {
		return &QueryResults{Error: err}
	}
	data := make([][]string, 0)
	row := 1
	for rows.Next() && row < rowLimit {
		stringValues := make([]string, len(columns)+1)
		stringValues[0] = strconv.Itoa(row)
		pointers := make([]interface{}, len(columns))
		for i := 0; i < len(columns); i++ {
			pointers[i] = &stringValues[i+1]
		}
		if err := rows.Scan(pointers...); err != nil {
			return &QueryResults{Error: err}
		}
		data = append(data, stringValues)
		row++
	}
	return &QueryResults{
		Columns: append([]string{"#"}, columns...),
		Data:    data,
	}
}
