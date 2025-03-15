package main

import (
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
    Id int
    Title string
    Date string
    Link string
    Description template.HTML
}

func getBlogPreview(limit int) []Post {
    var posts []Post
    rows, err := db.Query(`SELECT id, Title, Link FROM Blog LIMIT ?;`, limit)
    if err != nil {
        return nil
    }
    defer rows.Close()

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

func getBlog(link string) (Post, error) {
    var p Post
    var d string
    err := db.QueryRow(`SELECT id, Title, Date, Link, Description FROM Blog WHERE Link = ?;`, link).Scan(&p.Id, &p.Title, &p.Date, &p.Link, &d)
    if err != nil {
        return p, err
    }
    p.Description = template.HTML(d)
    return p, nil
}

func personal_blogPageHandler(w http.ResponseWriter, r *http.Request) {
    data := tplData {
        "preview": getBlogPreview(-1),
    }
    blog, err := getBlog(r.URL.Path[len("/blog/"):])
    if err == nil {
        data["blog"] = blog
    }
    
    err = tpl.ExecuteTemplate(w, "personal_blog.html", data)
    if err != nil {
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}

func getProjectPreview(limit int) []Post {
    var posts []Post
    rows, err := db.Query(`SELECT id, Title, Link FROM Project LIMIT ?;`, limit)
    if err != nil {
        return nil
    }
    defer rows.Close()
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

func personal_projectPageHandler(w http.ResponseWriter, r *http.Request) {
    data := tplData {
        "preview": getProjectPreview(-1),
    }
    project, err := getProject(r.URL.Path[len("/project/"):])
    if err == nil {
        data["project"] = project 
    }
    
    err = tpl.ExecuteTemplate(w, "personal_project.html", data)
    if err != nil {
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}

