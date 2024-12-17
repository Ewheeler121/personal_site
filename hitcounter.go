package main

import (
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func getHitCounter(w http.ResponseWriter, r *http.Request) (int, int) {
    isBrowser := false
    isUnique := false

    browsers := []string{"Mozilla", "Chrome", "Safari", "Edge", "Firefox", "Opera"}
    for _, b := range browsers {
        if strings.Contains(r.Header.Get("User-Agent"), b) {
            isBrowser = true
        }
    }
    _, err := r.Cookie("visted")
    if err != nil {
        isUnique = true
    }
    
    if isBrowser == true {
        _, err = db.Exec("UPDATE Counter SET Count = Count + 1 WHERE Label='normal'RETURNING Count;")
        if err != nil {
            return -1, -1
        }

        if isUnique == true {
            http.SetCookie(w, &http.Cookie{Name: "visted", Value: "true", Path:"/",MaxAge: 2147483647, HttpOnly: true, Secure: true})
            _, err = db.Exec("UPDATE Counter SET Count = Count + 1 WHERE Label='unique' RETURNING Count;")
            if err != nil {
                return -1, -1
            }
        }
    }

    normal := 0
    unique := 0
    err = db.QueryRow("SELECT Count FROM Counter WHERE Label='unique';").Scan(&unique)
    if err != nil {
        return -1, -1
    }
    err = db.QueryRow("SELECT Count FROM Counter WHERE Label='normal';").Scan(&normal)
    if err != nil {
        return -1, -1
    }

    return normal, unique
}
