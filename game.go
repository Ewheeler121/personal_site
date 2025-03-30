package main

import (
    "strings"
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
    http.Redirect(w, r, "https://ewheeler121.xyz/favicon.ico", http.StatusFound)
}

func game_indexHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path[len("/"):] != "" {
        filePath := filepath.Join("static/game/", r.URL.Path)
        
        if strings.HasSuffix(filePath, ".br") {
            w.Header().Set("Content-Encoding", "br")
        }

        if strings.HasSuffix(filePath, ".gz") {
            w.Header().Set("Content-Encoding", "gzip")
        }
    
        if strings.HasSuffix(filePath, ".wasm.br") {
            w.Header().Set("Content-Type", "application/wasm")
        }
        if strings.HasSuffix(filePath, ".wasm.gz") {
            w.Header().Set("Content-Type", "application/wasm")
        }

        http.ServeFile(w, r, filePath)
        return
    }

    err := tpl.ExecuteTemplate(w, "game_index.html", nil)
    if err != nil {
        http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
    }
}
