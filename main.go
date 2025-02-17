package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"sync"
    "path/filepath"
    "strings"

	_ "github.com/mattn/go-sqlite3"
)

var tpl *template.Template
var db *sql.DB
var Mu sync.Mutex

func main() {
    gameMux := http.NewServeMux()
    personalMux := http.NewServeMux()
    snootMux := http.NewServeMux()
    
    tpl = template.Must(template.ParseGlob("templates/*.html"))
    
    personalMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    snootMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static_snoot"))))

    personalMux.HandleFunc("/comment-preview", commentPreviewHandler)
    personalMux.HandleFunc("/resume", resumePageHandler)
    personalMux.HandleFunc("/construction", constructionPageHandler)
    personalMux.HandleFunc("/blog/", blogPageHandler)
    personalMux.HandleFunc("/project/", projectPageHandler)
    personalMux.HandleFunc("/favicon.ico", faviconHandler)
    personalMux.HandleFunc("/submit-comment", submitComment)
    personalMux.HandleFunc("/", indexPageHandler)
    
    snootMux.HandleFunc("/", snootIndexHandler)
    snootMux.HandleFunc("/favicon.ico", snootFaviconHandler)

    gameMux.HandleFunc("/", serveStatic)
    
    startDatabase()
    defer db.Close()

    personalCert, err := tls.LoadX509KeyPair("domain.cert.pem", "private.key.pem")
    if err != nil {
        panic(err.Error())
    }
    snootCert, err := tls.LoadX509KeyPair("snoot.domain.cert.pem", "snoot.private.key.pem")
    if err != nil {
        panic(err.Error())
    }

    certMap := map[string]*tls.Certificate {
        "ewheeler121.xyz": &personalCert,
        "game.ewheeler121.xyz": &personalCert,
        "devlog.pink": &snootCert,
        "localhost": &personalCert,
    }

    tlsConfig := &tls.Config {
        GetCertificate: func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
            if cert, ok := certMap[chi.ServerName]; ok {
                return cert, nil
            }
            return nil, fmt.Errorf("No Certificate Found for %s", chi.ServerName)
        },
    }

    server := &http.Server {
        Addr: ":443",
        TLSConfig: tlsConfig,
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Host {
        case "game.ewheeler121.xyz":
            gameMux.ServeHTTP(w, r)
        case "ewheeler121.xyz":
            personalMux.ServeHTTP(w, r)
        case "devlog.pink":
            snootMux.ServeHTTP(w, r)
        default:
            gameMux.ServeHTTP(w, r)
        }
    })
    
    err = server.ListenAndServeTLS("", "")
    if err != nil {
        fmt.Println("ERROR: could not start server: ", err)
    }
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/images/favicon.ico")
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path[len("/"):] != "" && r.URL.Path[len("/"):] != "submit-comment" {
        http.Error(w, "Page Not Found", http.StatusNotFound)
        return
    }

    n, u := getHitCounter(w, r)
    data := map[string]interface{} {
        "status": getStatus(),
        "hits": n,
        "uniqueHits": u,
        "blog": getBlogPreview(5),
        "project": getProjectPreview(5),
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

    _, err = db.Query("SELECT * FROM Counter;")
    _, err = db.Exec(`
    INSERT OR IGNORE INTO Counter (Label, Count) VALUES (?, ?);
    INSERT OR IGNORE INTO Counter (Label, Count) VALUES (?, ?);
    `, "unique", 0, "normal", 0)
    if err != nil {
        panic(err)
    }
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	ext := filepath.Ext(r.URL.Path)
	if strings.HasSuffix(r.URL.Path, ".br") {
		w.Header().Set("Content-Encoding", "br")
		if ext == ".js.br" {
			w.Header().Set("Content-Type", "application/javascript")
		} else if ext == ".wasm.br" {
			w.Header().Set("Content-Type", "application/wasm")
		} else if ext == ".data.br" {
			w.Header().Set("Content-Type", "application/wasm")
		}
	}
	http.ServeFile(w, r, "game/"+r.URL.Path)
}
