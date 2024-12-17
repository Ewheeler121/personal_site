package main

import (
	"database/sql"
	"net/http"
	"net/url"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Comment struct {
    User string
    Site string
    Comment string
}

func start_database() {
    var err error
    db, err = sql.Open("sqlite3", "./database.db");
    if err != nil {
        panic(err.Error())
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS Comment (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        Username TEXT NOT NULL,
        Site TEXT,
        Comment TEXT NOT NULL
    );`)
    if err != nil {
        panic(err.Error())
    }
}

func getComments() []Comment {
    var comments []Comment
    rows, err := db.Query("SELECT Username, Site, Comment FROM Comment WHERE 1=1 LIMIT 25")
    if err != nil {
        return nil
    }

    for rows.Next() {
        var c Comment
        err = rows.Scan(&c.User, &c.Site, &c.Comment)
        if err != nil {
            continue
        }
        comments = append(comments, c)
    }

    return comments 
}
func normalizeURL(u string) (string, error) {
    if u == "" {
        return "", nil
    }

	if !strings.Contains(u, "://") {
		u = "http://" + u
	}

	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	host := parsedURL.Hostname()
	if host == "" {
		return "", err
	}

	return host, nil
}

func submitComment(w http.ResponseWriter, r *http.Request) {
    var c Comment
    var err error

    r.ParseMultipartForm(10 << 20)
    c.User = r.Form.Get("username")
    c.Site, err = normalizeURL(r.Form.Get("website"))
    c.Comment = r.Form.Get("comment")
    
    if err != nil || c.User == "" || c.Comment == "" {
        http.Error(w, "Invalid Comment", http.StatusInternalServerError)
        return
    }

    Mu.Lock()
    defer Mu.Unlock()
    _, err = db.Exec("INSERT INTO Comment (Username, Site, Comment) VALUES (?, ?, ?)", c.User, c.Site, c.Comment)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    
    indexPageHandler(w, r)
}
