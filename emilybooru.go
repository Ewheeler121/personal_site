package main

import (
	"database/sql"
	"net/http"
	"path/filepath"
)

func emily_hookHandles(serve *http.ServeMux) {
    serve.HandleFunc("/favicon.ico", emily_faviconHandler)
    serve.HandleFunc("/", emily_indexHandler)
}

func emily_startDB(database *sql.DB) {  }

func emily_faviconHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "https://cdn.discordapp.com/emojis/1255615186346184796.webp", http.StatusFound)
}

func emily_indexHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path[len("/"):] != "" {
        http.ServeFile(w, r, filepath.Join("static/emilybooru/", r.URL.Path))
        return
    }
    err := tpl.ExecuteTemplate(w, "emily_index.html", nil)
    if err != nil {
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}
