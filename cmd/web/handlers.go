package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/yuryqwer/ip2loc/internal"
)

func (app *application) report(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		ip = getDefaultIP(r)
	}
	address := net.ParseIP(ip)
	if address == nil {
		app.infoLog.Printf("given ip is %s, which is not valid", ip)
		app.notFound(w, "please enter the right ip")
		return
	}

	info, err := internal.GetIPInfo(address.String(), app.db)
	if err != nil {
		app.serverError(w, err)
		return
	}

	respondJsonSuccess(w, ip, info)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if _, ok := map[string]bool{
		"/":        true,
		"/json":    true,
		"/en":      true,
		"/en/json": true,
	}[r.URL.Path]; !ok {
		app.notFound(w, http.StatusText(http.StatusNotFound))
		return
	}

	ip := getDefaultIP(r)
	address := net.ParseIP(ip)
	if address == nil {
		app.infoLog.Printf("given ip is %s, which is not valid", ip)
		app.clientError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	info, err := internal.GetIPInfo(address.String(), app.db)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var ipInfo *internal.IPInfo
	if strings.Contains(r.URL.Path, "en") {
		ipInfo = internal.GetIPInfoFromLocationISP(info, "en")
	} else {
		ipInfo = internal.GetIPInfoFromLocationISP(info, "zh-CN")
	}

	if strings.Contains(r.URL.Path, "json") {
		respondJsonSuccess(w, ip, ipInfo)
	} else {
		if strings.Contains(r.URL.Path, "en") {
			fmt.Fprintf(w, "Your IP: %s\tLocation: %s\tIsp: %s\tUserType: %s\n",
				ip, ipInfo.Country+" "+ipInfo.Region+" "+ipInfo.City, ipInfo.ISP, ipInfo.UserType)
		} else {
			fmt.Fprintf(w, "当前 IP：%s\t来自于：%s\t运营商：%s\t用户类型：%s\n",
				ip, ipInfo.Country+" "+ipInfo.Region+" "+ipInfo.City, ipInfo.ISP, ipInfo.UserType)
		}
	}
}
