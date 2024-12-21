package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

    "github.com/TwiN/go-away"
	_ "github.com/mattn/go-sqlite3"
)

type Comment struct {
    User string
    Site string
    Comment string
}


func renderComments() template.HTML {
    var comments []Comment
    rows, err := db.Query("SELECT Username, Site, Comment FROM (SELECT * FROM Comment ORDER BY id DESC LIMIT 25) AS row ORDER BY ID ASC")
    if err != nil {
        return template.HTML("<p>ERROR</p>")
    }

    for rows.Next() {
        var c Comment
        err = rows.Scan(&c.User, &c.Site, &c.Comment)
        if err != nil {
            continue
        }
        comments = append(comments, c)
    }

    var builder strings.Builder
    
    for _, comment := range comments {
        builder.WriteString("<tr>\n")
        c := strings.Split(comment.Comment, "\n")
        if comment.Site == "" {
            builder.WriteString(fmt.Sprintf("<td class=\"username\"><h4>%s:</h4></td>\n<td class=\"comment\">", comment.User))
        } else {
            builder.WriteString(fmt.Sprintf("<td class=\"username\"><h4>%s@%s:</h4></td>\n<td class=\"comment\">", comment.User, comment.Site))
        }
        for _, l := range c {
            builder.WriteString(fmt.Sprintf("<p>%s</p>\n", l))
        }
        builder.WriteString("</td>\n</tr>\n")
    }
    return template.HTML(builder.String())
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
    _, err = db.Exec("INSERT INTO Comment (Username, Site, Comment) VALUES (?, ?, ?)", goaway.Censor(c.User), goaway.Censor(c.Site), goaway.Censor(c.Comment))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    
    indexPageHandler(w, r)
}

func commentPreviewHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{} {
        "comments": renderComments(),
    }
    err := tpl.ExecuteTemplate(w, "comment-preview.html", data)
    if err != nil {
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}
