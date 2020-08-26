package main

import (
  "github.com/oschwald/maxminddb-golang"
  "net/http"
  "log"
  "net"
  "encoding/json"
  "fmt"
  "strings"
  "time"
)

func track() (time.Time) {
  return time.Now()
}

func duration(start time.Time) {
  elapsed := time.Since(start)
  log.Printf("Time elapsed in microseconds: %v\n==\n", elapsed)
}

func parse(w http.ResponseWriter, req *http.Request) {
  //Track execution time
  defer duration(track())
  keys        := req.URL.Query()
  ip_to_parse := keys.Get("ip")
  caller_ip   := req.RemoteAddr

  caller_ip, err := getIP(req)
  if err != nil {
    w.WriteHeader(400)
    w.Write([]byte("Bad caller ip"))
    return
  }

  log.Printf("Caller IP = %s and IP to parse = %v\n", caller_ip, ip_to_parse)

  if len(ip_to_parse) < 7 {
    log.Printf("Bad IP to parse: %v\n", ip_to_parse)
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("Bad ip or that is missing: " + ip_to_parse + "\n"))
    return
  }

  ip := net.ParseIP(ip_to_parse)

  db, err := maxminddb.Open("GeoIP2-City.mmdb")
  if err != nil {
    log.Printf("Error: %v\n", err)
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("Maxmind DB reading failed: " + ip_to_parse + "\n"))
    return
  }
  defer db.Close()

  var record struct {
    Country struct {
      IsoCode string            `maxminddb:"iso_code"`
      Names   map[string]string `maxminddb:"names"`
    } `maxminddb:"country"`

    City struct {
      Names   map[string]string `maxminddb:"names"`
    } `maxminddb:"city"`

    State []struct {
      IsoCode string            `maxminddb:"iso_code"`
      Names   map[string]string `maxminddb:"names"`
    } `maxminddb:"subdivisions"`

    Zipcode struct {
      Code    string `maxminddb:"code"`
    } `maxminddb:"postal"`
  }

  type Result struct { 
     City         string  `json:"city"` 
     State        string  `json:"state"` 
     StateCode    string  `json:"state_code"` 
     ZipCode      string  `json:"zip_code"` 
     Country      string  `json:"country"` 
     CountryCode  string  `json:"country_code"` 
     Restricted   bool    `json:"is_restricted"` 
     Cremia       bool    `json:"is_cremia_region"` 
  }

  err = db.Lookup(ip, &record)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("Maxmind lookup failed: " + ip_to_parse + "\n"))
    return
  }

  if len(record.State) == 0 {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("Maxmind lookup failed to retrieve any information of : " + ip_to_parse + "\n"))
    return
  }

  location := Result {
    City:       record.City.Names["en"],    
    State:      record.State[0].Names["en"],             
    StateCode:  record.State[0].IsoCode,
    ZipCode:    record.Zipcode.Code, 
    Country:    record.Country.Names["en"], 
    CountryCode:record.Country.IsoCode, 
    Restricted: false,
    Cremia:     false,
  } 

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(location)
}

func getIP(r *http.Request) (string, error) {
  //Get IP from the X-REAL-IP header
  ip := r.Header.Get("X-REAL-IP")
  netIP := net.ParseIP(ip)
  if netIP != nil {
    return ip, nil
  }

  //Get IP from X-FORWARDED-FOR header
  ips := r.Header.Get("X-FORWARDED-FOR")
  splitIps := strings.Split(ips, ",")
  for _, ip := range splitIps {
    netIP := net.ParseIP(ip)
    if netIP != nil {
      return ip, nil
    }
  }

  //Get IP from RemoteAddr
  ip, _, err := net.SplitHostPort(r.RemoteAddr)
  if err != nil {
    return "", err
  }
  netIP = net.ParseIP(ip)
  if netIP != nil {
    return ip, nil
  }
  return "", fmt.Errorf("No valid ip found")
}

func main() {
  http.HandleFunc("/geolocations", parse)

  log.Fatal(http.ListenAndServe(":8080", nil))
}