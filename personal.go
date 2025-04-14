package main

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func personal_hookHandles(serve *http.ServeMux) {
    serve.HandleFunc("/comment-preview", personal_commentPreviewHandler)
    serve.HandleFunc("/resume", personal_resumePageHandler)
    serve.HandleFunc("/construction", personal_constructionPageHandler)
    serve.HandleFunc("/blog/", personal_blogPageHandler)
    serve.HandleFunc("/project/", personal_projectPageHandler)
    serve.HandleFunc("/submit-comment", personal_submitComment)
    serve.HandleFunc("/favicon.ico", personal_faviconHandler)
    serve.HandleFunc("/", personal_indexHandler)

}

func personal_startDB(database *sql.DB) {
    _, err := database.Exec(`
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
    CREATE TABLE IF NOT EXISTS Project (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        Title TEXT NOT NULL,
        Date TEXT NOT NULL,
        Link TEXT NOT NULL,
        Description TEXT NOT NULL
    );
    `)
    if err != nil {
        panic(err)
    }
}

func getStatus() string {
    //gonna do some wonky Steam + Discord API things here
    status := os.Getenv("STATUS")
    if status == "" {
        status = "No Status Found"
    }
    return status
}

func getHitCounter(w http.ResponseWriter, r *http.Request) (int, int) {
    isBrowser := false
    isUnique := false

    browsers := []string{"Mozilla", "Chrome", "Safari", "Edge", "Firefox", "Opera"}
    for _, b := range browsers {
        if strings.Contains(r.Header.Get("User-Agent"), b) {
            isBrowser = true
        }
    }
    _, err := r.Cookie("visted")
    if err != nil {
        isUnique = true
    }
    
    if isBrowser == true {
        _, err = db.Exec("UPDATE Counter SET Count = Count + 1 WHERE Label='normal'RETURNING Count;")
        if err != nil {
            return -1, -1
        }

        if isUnique == true {
            http.SetCookie(w, &http.Cookie{Name: "visted", Value: "true", Path:"/",MaxAge: 2147483647, HttpOnly: true, Secure: true})
            _, err = db.Exec("UPDATE Counter SET Count = Count + 1 WHERE Label='unique' RETURNING Count;")
            if err != nil {
                return -1, -1
            }
        }
    }

    normal := 0
    unique := 0
    err = db.QueryRow("SELECT Count FROM Counter WHERE Label='unique';").Scan(&unique)
    if err != nil {
        return -1, -1
    }
    err = db.QueryRow("SELECT Count FROM Counter WHERE Label='normal';").Scan(&normal)
    if err != nil {
        return -1, -1
    }

    return normal, unique
}


func personal_indexHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path[len("/"):] != "" {
        http.ServeFile(w, r, filepath.Join("static/personal/", r.URL.Path))
        return
    }

    n, u := getHitCounter(w, r)
    data := tplData {
        "status": getStatus(),
        "hits": n,
        "uniqueHits": u,
    }
    
    blogs := getBlogPreview(5)
    if blogs != nil {
        data["blog"] = blogs
    }
    projects := getProjectPreview(5)
    if projects != nil {
        data["project"] = projects
    }

    err := tpl.ExecuteTemplate(w, "personal_index.html", data)
    if err != nil {
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}

func personal_faviconHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/personal/images/favicon.ico")
}

func personal_constructionPageHandler(w http.ResponseWriter, r *http.Request) {
    err := tpl.ExecuteTemplate(w, "personal_construction.html", nil)
    if err != nil {
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}

func personal_resumePageHandler(w http.ResponseWriter, r *http.Request) {
    err := tpl.ExecuteTemplate(w, "personal_resume.html", nil)
    if err != nil {
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}
