module github.com/yuryqwer/ip2loc

go 1.20

replace github.com/oschwald/geoip2-golang v1.9.0 => ./internal/geoip2-golang

require (
	github.com/fsnotify/fsnotify v1.6.0
	github.com/oschwald/geoip2-golang v1.9.0
	golang.org/x/time v0.3.0
)

require (
	github.com/oschwald/maxminddb-golang v1.11.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
)
