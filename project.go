package main

import (
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)


func getProjectPreview(limit int) []Post {
    var posts []Post
    rows, err := db.Query(`SELECT id, Title, Link FROM Project LIMIT ?;`, limit)
    for rows.Next() {
        var p Post
        err = rows.Scan(&p.Id, &p.Title, &p.Link)
        if err != nil {
            continue
        }
        posts = append(posts, p)
    }
    return posts
}

func getProject(link string) (Post, error) {
    var p Post
    var d string
    err := db.QueryRow(`SELECT id, Title, Date, Link, Description FROM Project WHERE Link = ?;`, link).Scan(&p.Id, &p.Title, &p.Date, &p.Link, &d)
    if err != nil {
        return p, err
    }
    p.Description = template.HTML(d)
    return p, nil
}

func projectPageHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{} {
        "preview": getProjectPreview(-1),
    }
    project, err := getProject(r.URL.Path[len("/project/"):])
    if err == nil {
        data["project"] = project 
    }
    
    err = tpl.ExecuteTemplate(w, "project.html", data)
    if err != nil {
        fmt.Println(err)
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}
