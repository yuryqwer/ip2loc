package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/yuryqwer/ip2loc/internal"
)

type jsonResponse struct {
	Code int         `json:"code"` // 1:success 2:client error 3:server error
	Msg  string      `json:"msg"`
	IP   string      `json:"ip"`
	Data interface{} `json:"data"`
}

func respondJson(w http.ResponseWriter, r jsonResponse, status int) {
	// the content sniffing cannot distingish JSON from plain text
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonBytes, _ := json.Marshal(r)
	fmt.Fprintf(w, "%s", jsonBytes)
}

func respondJsonSuccess(w http.ResponseWriter, ip string, data interface{}) {
	httpResponse := jsonResponse{
		Code: 1,
		Msg:  "success",
		IP:   ip,
		Data: data,
	}
	respondJson(w, httpResponse, http.StatusOK)
}

func (app *application) clientError(w http.ResponseWriter, msg string, status int) {
	httpResponse := jsonResponse{
		Code: 2,
		Msg:  "success",
		Data: map[string]string{"msg": msg},
	}
	respondJson(w, httpResponse, status)
}

func (app *application) notFound(w http.ResponseWriter, msg string) {
	app.clientError(w, msg, http.StatusNotFound)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	httpResponse := jsonResponse{
		Code: 3,
		Msg:  "error",
		Data: map[string]string{"msg": http.StatusText(http.StatusInternalServerError)},
	}
	respondJson(w, httpResponse, http.StatusInternalServerError)
}

func getDefaultIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-Ip")
	}
	if ip == "" {
		ip = r.Header.Get("Remote-Addr")
	}
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	// X-Forwarded-For may get several ips separated by comma
	ip = strings.TrimSpace(strings.Split(ip, ",")[0])
	return ip
}

func (app *application) watchAndReload(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			app.infoLog.Printf("watchAndReload get notified: %#v", event)
			if !ok {
				app.errorLog.Println("watchAndReload error: watcher.Events has benn closed")
				return
			}
			if (event.Op == fsnotify.Write || event.Op == fsnotify.Create) &&
				filepath.Base(event.Name) == app.mmdbName {
				f, err := os.Open(event.Name)
				if err != nil {
					app.infoLog.Println("watchAndReload error:", err)
					continue
				}

				h := md5.New()
				if _, err := io.Copy(h, f); err != nil {
					f.Close()
					app.infoLog.Println("watchAndReload error:", err)
					continue
				}
				sum := h.Sum(nil)
				f.Close()

				if !bytes.Equal(sum, app.mmdbSum) {
					db, err := internal.NewDB(event.Name)
					if err != nil {
						app.infoLog.Println("watchAndReload error:", err)
						continue
					}
					defer db.Close()
					app.db.Close()
					app.db = db
					app.mmdbSum = sum
					app.infoLog.Println("watchAndReload change the mmdb")
				}
			}
		case _, ok := <-watcher.Errors:
			if !ok {
				app.errorLog.Println("watchAndReload error: watcher.Errors has benn closed")
				return
			}
		}
	}
}
