package main

import (
	"crypto/tls"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/geosoft1/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// TODO application version
// https://semver.org/#semantic-versioning-200
const SW_VERSION = "1.1.1-release-build030519"

// NOTE session cookie name
const SESSID = "SESSID"

// NOTE cache manifest file
const MANIFEST = `CACHE MANIFEST
FALLBACK:
/ /static/offline.html
`

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	isActive bool   `json:"is_active"`
}

type Database struct {
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type SMTP struct {
	Server   string `json:"smtp"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// main configuration structure
type Config struct {
	Database `json:"database"`
	SMTP
}

var httpAddress = flag.String("http", ":8080", "http address")
var httpsAddress = flag.String("https", ":8090", "https address")
var httpsEnabled = flag.Bool("https-enabled", false, "enable https server")

var config = &Config{}
var templ = template.New("templ").Delims("[[", "]]") // integrate with angular
var router = mux.NewRouter()                         // main router

// NOTE this is the salt for secured passwords
const salt = "super-secret-key"

// NOTE random recovery password length
const token_len = 8

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte(salt)
	store = sessions.NewCookieStore(key)
)

func main() {
	// NOTE this will avoid processor overload in some circumstances
	runtime.GOMAXPROCS(1)

	flag.Usage = func() {
		fmt.Printf("usage: %s [options]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	folder, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Open(filepath.Join(folder, "config.json"))
	if err != nil {
		log.Fatalln(err)
	}
	json.Decode(file, &config)
	// NOTE [DEBUG] print configuration
	// log.Println(config)

	if _, err := templ.ParseGlob(filepath.Join(folder, "ui", "*.html")); err != nil {
		log.Fatalln(err)
	}

	sqlConnect()

	//https://github.com/gorilla/sessions/issues/58#issuecomment-154217736
	gob.Register(User{})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templ.ExecuteTemplate(w, "index", nil)
	})
	router.HandleFunc("/signin", signin)
	router.HandleFunc("/signup", signup)
	router.HandleFunc("/reset", reset) // reset password if you forgot
	router.HandleFunc("/signout", signout)
	router.HandleFunc("/cache.manifest", func(w http.ResponseWriter, r *http.Request) {
		// offline page is shown if network is down
		w.Header().Set("Content-Type", "text/cache-manifest")
		w.Write([]byte(MANIFEST))
	})
	application := router.PathPrefix("/a").Subrouter()
	application.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// resolve cross-site access issues
			w.Header().Set("Access-Control-Allow-Origin", "*")
			session, _ := store.Get(r, SESSID)
			if session.IsNew || session.Values["authenticated"] == false {
				http.Error(w, "Session expired", http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
	application.HandleFunc("/home", home)
	application.HandleFunc("/user", user)

	// TODO application handlers come here

	// NOTE [DEBUG] print registered routes
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		// if name, err := route.GetPathTemplate(); err == nil {
		// 	log.Println(name)
		// }
		return nil
	})

	// file server
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(folder, "static")))))
	// --- error handlers (NotFoundHandler, MethodNotAllowedHandler)
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		//http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	if *httpsEnabled {
		go func() {
			// allow you to use self signed certificates
			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			// NOTE geneate self signed certificates for tests
			// openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt
			log.Printf("start https server on %s", *httpsAddress)
			if err := http.ListenAndServeTLS(*httpsAddress, filepath.Join(folder, "server.crt"), filepath.Join(folder, "server.key"), router); err != nil {
				log.Fatalln(err)
			}
		}()
	}

	log.Printf("start http server on %s", *httpAddress)
	if err := http.ListenAndServe(*httpAddress, router); err != nil {
		log.Fatalln(err)
	}
}
