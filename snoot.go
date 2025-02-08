package main

import (
	"net/http"
)

func snootIndexHandler(w http.ResponseWriter, r *http.Request) {
    err := tpl.ExecuteTemplate(w, "snoot.html", nil)
    if err != nil {
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}

func snootFaviconHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "https://cdn.discordapp.com/emojis/1255615186346184796.webp", http.StatusFound)
}
