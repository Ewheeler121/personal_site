package main

import (
	"database/sql"
	"net/http"
	"path/filepath"
)

func emilyGame_hookHandles(serve *http.ServeMux) {
    serve.HandleFunc("/favicon.ico", emilyGame_faviconHandler)
    serve.HandleFunc("/", emilyGame_indexHandler)
}

func emilyGame_startDB(database *sql.DB) {  }

func emilyGame_faviconHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "https://cdn.discordapp.com/emojis/1255615186346184796.webp", http.StatusFound)
}

func emilyGame_indexHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path[len("/"):] != "" {
        http.ServeFile(w, r, filepath.Join("static/emily_game/", r.URL.Path))
        return
    }
    err := tpl.ExecuteTemplate(w, "emily_game_index.html", nil)
    if err != nil {
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}
