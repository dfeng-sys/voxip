package voxip

import (
	"encoding/json"
	"io"
	"net"
	"net/http"

	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"
)

var (
	db *geoip2.Reader
)

/*
 * json request/response object structs
 */
type Request struct {
	Ip        string   `json:"ip"`
	Whitelist []string `json:"whitelist"`
}
type Response struct {
	Whitelisted bool   `json:"whitelisted"`
	Country     string `json:"country,omitempty"`
	Error       string `json:"error,omitempty"`
}

/*
 * open mmdb
 */
func Init(mmdbFile string) {
	log.Infof("opening mmdb '%s'", mmdbFile)
	reader, err := geoip2.Open(mmdbFile)
	if err != nil {
		log.Fatalf("failed to open mmdb: %v", err)
	}
	db = reader
}

/*
 * handle request to update mapping data,
 * not fully implemented
 */
func UpdateHandler(w http.ResponseWriter, r *http.Request) {

	log.Info()
	log.Info("received update request")

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(&Response{Error: "update function not yet implemented"})
}

/*
 * receive API request,
 * perform IP location check,
 * respond with result
 */
func RequestHandler(w http.ResponseWriter, r *http.Request) {

	log.Info()
	log.Info("received ip lookup request:")

	// receive request
	reqBody, _ := io.ReadAll(r.Body)
	log.Infof("request: %v", string(reqBody))

	// unmarshal request
	var req Request
	err := json.Unmarshal(reqBody, &req)

	w.Header().Set("content-type", "application/json")
	if err != nil {
		log.Error("invalid json")
		json.NewEncoder(w).Encode(&Response{Error: "invalid json in request"})
		return
	}

	response := CheckWhitelist(req)

	// send response
	respJson, _ := json.Marshal(response)
	log.Infof("sending response: %v", string(respJson))
	json.NewEncoder(w).Encode(&response)
}

/*
 * validate and process request
 */
func CheckWhitelist(req Request) Response {

	// initialize response
	resp := Response{Whitelisted: false}

	// check that ip is present
	if len(req.Ip) == 0 {
		log.Error("received null or empty IP address")
		resp.Error = "received null or empty IP address"
		return resp
	}

	// check that ip is valid
	ip := net.ParseIP(req.Ip)
	if ip == nil {
		log.Error("received invalid IP address")
		resp.Error = "received invalid IP address"
		return resp
	}

	// check that whitelist is present
	if req.Whitelist == nil || len(req.Whitelist) == 0 {
		log.Error("received null or empty whitelist")
		resp.Error = "received null or empty whitelist"
		return resp
	}

	log.Infof("received IP address '%s'", ip.String())

	// fill response object
	whitelisted, country, err := isWhitelisted(ip, req.Whitelist)
	resp.Whitelisted = whitelisted
	resp.Country = country
	if err != nil {
		resp.Error = err.Error()
	}
	return resp
}

/*
 * check if given IP address's country is in the given whitelist
 * returns true if it is
 * returns the IP address's country code
 * returns an error if the db lookup fails
 */
func isWhitelisted(ip net.IP, whitelist []string) (bool, string, error) {

	// look up country
	record, err := db.Country(ip)
	if err != nil {
		log.Errorf("error during lookup: %v", err)
		return false, "", err
	}

	log.Infof("IP's country code: %s", record.Country.IsoCode)

	// check if country is present in whitelist
	for _, country := range whitelist {
		if country == record.Country.IsoCode {
			log.Info("IP address IS on whitelist")
			return true, record.Country.IsoCode, nil
		}
	}

	log.Info("IP address IS NOT on whitelist")
	return false, record.Country.IsoCode, nil
}
