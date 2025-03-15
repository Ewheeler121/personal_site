package main

import (
	"database/sql"
	"net/http"
	"path/filepath"
)

func game_hookHandles(serve *http.ServeMux) {
    serve.HandleFunc("/favicon.ico", game_faviconHandler)
    serve.HandleFunc("/", game_indexHandler)
}

func game_startDB(database *sql.DB) {  }

func game_faviconHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "https://cdn.discordapp.com/emojis/1255615186346184796.webp", http.StatusFound)
}

func game_indexHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path[len("/"):] != "" {
        http.ServeFile(w, r, filepath.Join("static/game/", r.URL.Path))
        return
    }
    err := tpl.ExecuteTemplate(w, "game_index.html", nil)
    if err != nil {
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}
