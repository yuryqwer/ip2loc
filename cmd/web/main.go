package main

import (
	"crypto/md5"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/oschwald/geoip2-golang"
	"github.com/yuryqwer/ip2loc/internal"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	db       *geoip2.Reader
	limiter  *internal.IPRateLimiter
	mmdbName string
	mmdbSum  []byte
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	mmdb := flag.String("mmdb", "./dbip-full.mmdb", "The mmdb file path")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	f, err := os.Open(*mmdb)
	if err != nil {
		errorLog.Fatal(err)
	}
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		f.Close()
		errorLog.Fatal(err)
	}
	sum := h.Sum(nil)
	f.Close()

	db, err := internal.NewDB(*mmdb)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	limiter := internal.NewIPRateLimiter(1, 5)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(filepath.Dir(*mmdb))
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		db:       db,
		limiter:  limiter,
		mmdbName: filepath.Base(*mmdb),
		mmdbSum:  sum,
	}

	go app.watchAndReload(watcher)

	srv := &http.Server{
		Addr:        *addr,
		ErrorLog:    errorLog,
		Handler:     app.routes(),
		IdleTimeout: 10 * time.Second,
		ReadTimeout: 5 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
