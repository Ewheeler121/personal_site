package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

type tplData = map[string]interface{} 

var tpl *template.Template
var db *sql.DB
var Mu sync.Mutex

func main() {
    var err error

    tpl = template.Must(template.ParseGlob("templates/*.html"))
    if tpl == nil {
        panic("no tpl???")
    }
    
    //creates/starts database
    db, err = sql.Open("sqlite3", "./database.db")
    defer db.Close()
    if err != nil {
        panic(err)
    }

    //creates tables
    personal_startDB(db)

    //hook handlers
    http.HandleFunc("/comment-preview", commentPreviewHandler)
    http.HandleFunc("/resume", resumePageHandler)
    http.HandleFunc("/construction", constructionPageHandler)
    http.HandleFunc("/blog/", blogPageHandler)
    http.HandleFunc("/project/", projectPageHandler)
    http.HandleFunc("/submit-comment", submitComment)
    http.HandleFunc("/rss", blogRSSFeedHandler)
    http.HandleFunc("/favicon.ico", faviconHandler)
    http.HandleFunc("/", indexHandler)

	err = http.ListenAndServeTLS("127.0.0.1:3001", "certs/domain.cert.pem", "certs/private.key.pem", nil)
    if err != nil {
        fmt.Println("ERROR: could not start server: ", err)
    }
}

