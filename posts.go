package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// RSS feed structs
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	Language      string `xml:"language"`
	LastBuildDate string `xml:"lastBuildDate"`
	Items         []Item `xml:"item"`
}

type CDATA struct {
	Value string `xml:",cdata"`
}

type Item struct {
	Title       CDATA  `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description CDATA `xml:"description"`
}

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

func blogPageHandler(w http.ResponseWriter, r *http.Request) {
    data := tplData {
        "preview": getBlogPreview(-1),
    }
    blog, err := getBlog(r.URL.Path[len("/blog/"):])
    if err == nil {
        data["blog"] = blog
    }
    
    err = tpl.ExecuteTemplate(w, "blog.html", data)
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

func projectPageHandler(w http.ResponseWriter, r *http.Request) {
    data := tplData {
        "preview": getProjectPreview(-1),
    }
    project, err := getProject(r.URL.Path[len("/project/"):])
    if err == nil {
        data["project"] = project 
    }
    
    err = tpl.ExecuteTemplate(w, "project.html", data)
    if err != nil {
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}

func personal_blogRSSFeed(posts []Post) ([]byte, error) {
	const layout = "Jan 2 2006"
	now := time.Now().Format(time.RFC1123Z)

	var items []Item
	for _, p := range posts {
		t, err := time.Parse(layout, p.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date %q: %w", p.Date, err)
		}
		items = append(items, Item{
			Title:       CDATA{p.Title},
			Link:        fmt.Sprintf("https://ewheeler121.xyz/blog/%s", p.Link),
			GUID:        fmt.Sprintf("https://ewheeler121.xyz/blog/%s", p.Link),
			PubDate:     t.Format(time.RFC1123Z),
			Description: CDATA{string(p.Description)},
		})
	}

	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:         "Ewheeler121",
			Link:          "https://ewheeler121.xyz/blog",
			Description:   "Random Blog Posts from Ewheeler121",
			Language:      "en-us",
			LastBuildDate: now,
			Items:         items,
		},
	}

	buf := &bytes.Buffer{}
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(buf)
	enc.Indent("", "  ")
	if err := enc.Encode(rss); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}


func blogRSSFeedHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT Title, Date, Link, Description FROM Blog ORDER BY id DESC`)
	if err != nil {
		http.Error(w, "Failed to generate feed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.Title, &p.Date, &p.Link, &p.Description); err != nil {
			http.Error(w, "Failed to generate feed", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}

	feed, err := personal_blogRSSFeed(posts)
	if err != nil {
		http.Error(w, "Failed to generate feed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.Write(feed)
}
