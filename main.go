package main

import (
    "database/sql"
    "fmt"
    "html/template"
    "net/http"
    "sync"

    _ "github.com/mattn/go-sqlite3"
)

var tpl *template.Template
var db *sql.DB
var Mu sync.Mutex

func main() {
    startDatabase()
    defer db.Close()

    tpl = template.Must(template.ParseGlob("templates/*.html"))

    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.HandleFunc("/comment-preview", commentPreviewHandler)
    http.HandleFunc("/resume", resumePageHandler)
    http.HandleFunc("/construction", constructionPageHandler)
    http.HandleFunc("/blog/", blogPageHandler)
    http.HandleFunc("/favicon.ico", faviconHandler)
    
    http.HandleFunc("/submit-comment", submitComment)

    http.HandleFunc("/", indexPageHandler)
    
    
    err := http.ListenAndServeTLS(":443", "domain.cert.pem", "private.key.pem", nil)
    if err != nil {
        fmt.Println("ERROR: could not start server: ", err)
    }
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/images/favicon.ico")
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
    n, u := getHitCounter(w, r)
    data := map[string]interface{} {
        "status": getStatus(),
        "hits": n,
        "uniqueHits": u,
        "blog": getPostPreview(5),
    }
    err := tpl.ExecuteTemplate(w, "index.html", data)
    if err != nil {
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}

func constructionPageHandler(w http.ResponseWriter, r *http.Request) {
    err := tpl.ExecuteTemplate(w, "construction.html", nil)
    if err != nil {
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}

func resumePageHandler(w http.ResponseWriter, r *http.Request) {
    err := tpl.ExecuteTemplate(w, "resume.html", nil)
    if err != nil {
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}

func startDatabase() {
    var err error
    db, err = sql.Open("sqlite3", "./database.db");
    if err != nil {
        panic(err)
    }

    _, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS Comment (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        Username TEXT NOT NULL,
        Site TEXT,
        Comment TEXT NOT NULL
    );
    CREATE TABLE IF NOT EXISTS Counter (
        Label TEXT UNIQUE NOT NULL,
        Count INTEGER
    );
    CREATE TABLE IF NOT EXISTS Blog (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        Title TEXT NOT NULL,
        Date TEXT NOT NULL,
        Link TEXT NOT NULL,
        Description TEXT NOT NULL
    );
    CREATE TABLE IF NOT EXISTS Projects (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        Title TEXT NOT NULL,
        Date TEXT NOT NULL,
        Description TEXT NOT NULL
    );
    `)
    if err != nil {
        panic(err)
    }

    _, err = db.Query("SELECT * FROM Counter;")
    _, err = db.Exec(`
    INSERT OR IGNORE INTO Counter (Label, Count) VALUES (?, ?);
    INSERT OR IGNORE INTO Counter (Label, Count) VALUES (?, ?);
    `, "unique", 0, "normal", 0)
    if err != nil {
        panic(err)
    }
}
