package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type tplData = map[string]interface{} 

var tpl *template.Template
var db *sql.DB
var Mu sync.Mutex

func main() {
    var err error
    personalMux := http.NewServeMux()
    emilybooruMux := http.NewServeMux()
    gameMux := http.NewServeMux()
    emilyGameMux := http.NewServeMux()

    tpl = template.Must(template.ParseGlob("templates/*.html"))
    if tpl == nil {
        panic("no tpl???")
    }
    
    //creates/starts database
    db, err = sql.Open("sqlite3", "./database.db")
    defer db.Close()
    if err != nil {
        panic(err)
    }

    //creates tables
    personal_startDB(db)
    emily_startDB(db)
    emilyGame_startDB(db)
    game_startDB(db)

    //hook handlers
    personal_hookHandles(personalMux)
    emily_hookHandles(emilybooruMux)
    emilyGame_hookHandles(emilyGameMux)
    game_hookHandles(gameMux)
    
    //load certs
    personalCert, err := tls.LoadX509KeyPair("certs/domain.cert.pem", "certs/private.key.pem")
    if err != nil {
        panic(err.Error())
    }
    snootCert, err := tls.LoadX509KeyPair("certs/snoot.domain.cert.pem", "certs/snoot.private.key.pem")
    if err != nil {
        panic(err.Error())
    }
    //certMap for config
    certMap := map[string]*tls.Certificate {
        "ewheeler121.xyz": &personalCert,
        "game.ewheeler121.xyz": &personalCert,
        "devlog.pink": &snootCert,
        "game.devlog.pink": &snootCert,
        //use this for testing
        //"localhost": &snootCert,
        "localhost": &personalCert,
    }
    tlsConfig := &tls.Config {
        GetCertificate: func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
            if cert, ok := certMap[chi.ServerName]; ok {
                return cert, nil
            }
            return nil, fmt.Errorf("No Certificate Found for %s", chi.ServerName)
        },
    }
    server := &http.Server {
        Addr: ":443",
        TLSConfig: tlsConfig,
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Host {
        case "ewheeler121.xyz":
            personalMux.ServeHTTP(w, r)
        case "game.ewheeler121.xyz":
            gameMux.ServeHTTP(w, r)
        case "devlog.pink":
            emilybooruMux.ServeHTTP(w, r)
        case "game.devlog.pink":
            emilybooruMux.ServeHTTP(w, r)
            //gameMux.ServeHTTP(w, r)
        default:
            //personalMux.ServeHTTP(w, r)
            //emilybooruMux.ServeHTTP(w, r)
            //emilyGameMux.ServeHTTP(w, r)
            gameMux.ServeHTTP(w, r)
        }
    })

    err = server.ListenAndServeTLS("", "")
    if err != nil {
        fmt.Println("ERROR: could not start server: ", err)
    }
}

