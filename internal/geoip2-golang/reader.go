// Package geoip2 provides an easy-to-use API for the MaxMind GeoIP2 and
// GeoLite2 databases; this package does not support GeoIP Legacy databases.
//
// The structs provided by this package match the internal structure of
// the data in the MaxMind databases.
//
// See github.com/oschwald/maxminddb-golang for more advanced used cases.
package geoip2

import (
	"fmt"
	"net"

	"github.com/oschwald/maxminddb-golang"
)

// The LocationISP struct corresponds to the data in the dbip's
// `IP to Location + ISP` database.
type LocationISP struct {
	City struct {
		GeoNameID uint              `maxminddb:"geoname_id" json:"geoname_id"`
		Names     map[string]string `maxminddb:"names" json:"names"`
	} `maxminddb:"city" json:"city"`
	Continent struct {
		Code      string            `maxminddb:"code" json:"code"`
		GeoNameID uint              `maxminddb:"geoname_id" json:"geoname_id"`
		Names     map[string]string `maxminddb:"names" json:"names"`
	} `maxminddb:"continent" json:"continent"`
	Country struct {
		GeoNameID         uint              `maxminddb:"geoname_id" json:"geoname_id"`
		IsoCode           string            `maxminddb:"iso_code" json:"iso_code"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union" json:"is_in_european_union"`
		Names             map[string]string `maxminddb:"names" json:"names"`
	} `maxminddb:"country" json:"country"`
	Location struct {
		Latitude    float64 `maxminddb:"latitude" json:"latitude"`
		Longitude   float64 `maxminddb:"longitude" json:"longitude"`
		TimeZone    string  `maxminddb:"time_zone" json:"time_zone"`
		WeatherCode string  `maxminddb:"weather_code" json:"weather_code"`
	} `maxminddb:"location" json:"location"`
	Postal struct {
		Code string `maxminddb:"code" json:"code"`
	} `maxminddb:"postal" json:"postal"`
	Subdivisions []struct {
		GeoNameID uint              `maxminddb:"geoname_id" json:"geoname_id"`
		IsoCode   string            `maxminddb:"iso_code" json:"iso_code"`
		Names     map[string]string `maxminddb:"names" json:"names"`
	} `maxminddb:"subdivisions" json:"subdivisions"`
	Traits struct {
		AutonomousSystemNumber       uint   `maxminddb:"autonomous_system_number" json:"autonomous_system_number"`
		AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization" json:"autonomous_system_organization"`
		ConnectionType               string `maxminddb:"connection_type" json:"connection_type"`
		UserType                     string `maxminddb:"user_type" json:"user_type"`
		ISP                          string `maxminddb:"isp" json:"isp"`
		Organization                 string `maxminddb:"organization" json:"organization"`
		Network                      string `json:"network"`
	} `maxminddb:"traits" json:"traits"`
}

// The ConnectionType struct corresponds to the data in the GeoIP2
// Connection-Type database.
type ConnectionType struct {
	ConnectionType string `maxminddb:"connection_type"`
}

// The Domain struct corresponds to the data in the GeoIP2 Domain database.
type Domain struct {
	Domain string `maxminddb:"domain"`
}

type databaseType int

const (
	isAnonymousIP = 1 << iota
	isASN
	isCity
	isConnectionType
	isCountry
	isDomain
	isEnterprise
	isISP
)

// Reader holds the maxminddb.Reader struct. It can be created using the
// Open and FromBytes functions.
type Reader struct {
	mmdbReader   *maxminddb.Reader
	databaseType databaseType
}

// InvalidMethodError is returned when a lookup method is called on a
// database that it does not support. For instance, calling the ISP method
// on a City database.
type InvalidMethodError struct {
	Method       string
	DatabaseType string
}

func (e InvalidMethodError) Error() string {
	return fmt.Sprintf(`geoip2: the %s method does not support the %s database`,
		e.Method, e.DatabaseType)
}

// UnknownDatabaseTypeError is returned when an unknown database type is
// opened.
type UnknownDatabaseTypeError struct {
	DatabaseType string
}

func (e UnknownDatabaseTypeError) Error() string {
	return fmt.Sprintf(`geoip2: reader does not support the %q database type`,
		e.DatabaseType)
}

// Open takes a string path to a file and returns a Reader struct or an error.
// The database file is opened using a memory map. Use the Close method on the
// Reader object to return the resources to the system.
func Open(file string) (*Reader, error) {
	reader, err := maxminddb.Open(file)
	if err != nil {
		return nil, err
	}
	dbType, err := getDBType(reader)
	return &Reader{reader, dbType}, err
}

// FromBytes takes a byte slice corresponding to a GeoIP2/GeoLite2 database
// file and returns a Reader struct or an error. Note that the byte slice is
// used directly; any modification of it after opening the database will result
// in errors while reading from the database.
func FromBytes(bytes []byte) (*Reader, error) {
	reader, err := maxminddb.FromBytes(bytes)
	if err != nil {
		return nil, err
	}
	dbType, err := getDBType(reader)
	return &Reader{reader, dbType}, err
}

func getDBType(reader *maxminddb.Reader) (databaseType, error) {
	switch reader.Metadata.DatabaseType {
	case "GeoIP2-Anonymous-IP":
		return isAnonymousIP, nil
	case "DBIP-ASN-Lite (compat=GeoLite2-ASN)",
		"GeoLite2-ASN":
		return isASN, nil
	// We allow City lookups on Country for back compat
	case "DBIP-City-Lite",
		"DBIP-Country-Lite",
		"DBIP-Country",
		"DBIP-Location (compat=City)",
		"GeoLite2-City",
		"GeoIP2-City",
		"GeoIP2-City-Africa",
		"GeoIP2-City-Asia-Pacific",
		"GeoIP2-City-Europe",
		"GeoIP2-City-North-America",
		"GeoIP2-City-South-America",
		"GeoIP2-Precision-City",
		"GeoLite2-Country",
		"GeoIP2-Country":
		return isCity | isCountry, nil
	case "GeoIP2-Connection-Type":
		return isConnectionType, nil
	case "GeoIP2-Domain":
		return isDomain, nil
	case "DBIP-ISP (compat=Enterprise)",
		"DBIP-Location-ISP (compat=Enterprise)",
		"GeoIP2-Enterprise",
		"IPCC-Location-ISP-Enterprise":
		return isEnterprise | isCity | isCountry, nil
	case "GeoIP2-ISP", "GeoIP2-Precision-ISP":
		return isISP | isASN, nil
	default:
		return 0, UnknownDatabaseTypeError{reader.Metadata.DatabaseType}
	}
}

// LocationISP takes an IP address as a net.IP struct and returns an LocationISP
// struct and/or an error. This is intended to be used with the dbip's
// `IP to Location + ISP` database.
func (r *Reader) LocationISP(ipAddress net.IP) (*LocationISP, error) {
	if isEnterprise&r.databaseType == 0 {
		return nil, InvalidMethodError{"LocationISP", r.Metadata().DatabaseType}
	}
	var locationISP LocationISP
	network, _, err := r.mmdbReader.LookupNetwork(ipAddress, &locationISP)
	if err != nil {
		return nil, err
	}
	locationISP.Traits.Network = network.String()
	return &locationISP, err
}

// ConnectionType takes an IP address as a net.IP struct and returns a
// ConnectionType struct and/or an error.
func (r *Reader) ConnectionType(ipAddress net.IP) (*ConnectionType, error) {
	if isConnectionType&r.databaseType == 0 {
		return nil, InvalidMethodError{"ConnectionType", r.Metadata().DatabaseType}
	}
	var val ConnectionType
	err := r.mmdbReader.Lookup(ipAddress, &val)
	return &val, err
}

// Domain takes an IP address as a net.IP struct and returns a
// Domain struct and/or an error.
func (r *Reader) Domain(ipAddress net.IP) (*Domain, error) {
	if isDomain&r.databaseType == 0 {
		return nil, InvalidMethodError{"Domain", r.Metadata().DatabaseType}
	}
	var val Domain
	err := r.mmdbReader.Lookup(ipAddress, &val)
	return &val, err
}

// Metadata takes no arguments and returns a struct containing metadata about
// the MaxMind database in use by the Reader.
func (r *Reader) Metadata() maxminddb.Metadata {
	return r.mmdbReader.Metadata
}

// Close unmaps the database file from virtual memory and returns the
// resources to the system.
func (r *Reader) Close() error {
	return r.mmdbReader.Close()
}
