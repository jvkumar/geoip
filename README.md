# geoip
Maxmind geoip IP  to region extraction.

# Getting started

> All following steps assume you are on `macOS Catalina`

## Installation of go

```
$ brew install golang
```

Assuming your `$HOME` is `/Users/<your_mac_user_name>`. You can validate by running command `echo $HOME`

```
$ mkdir -p $HOME/go
$ vi ~/.bash_profile
```
and add following three lines in your `~/.bash_profile` file

```
export GOPATH=$HOME/go
export GOROOT="$(brew --prefix golang)/libexec"
export PATH="$PATH:${GOPATH}/bin:${GOROOT}/bin"
```

## Application setup

```
$ cd ~/go
$ git clone https://github.com/jvkumar/geoip.git
$ cd geoip
$ go get github.com/oschwald/maxminddb-golang
```

Download `GeoIP2-City.mmdb` file from Maxmind website and save in `geoip` folder

## Start the action now

```
$ go build http.go
```
This should compile and create a go binary that you may see a new file `http` in this folder. Now time to run the http server on port `8080`. Run the following  command -

```
$ ./http
```
You may get a popup to allow the traffic. Click `Allow` 

This server exposes one endpoint `GET /geolocations` with url parameter `ip`. An example URI would be `http://localhost:8080/geolocations?ip=24.3.77.32`

Change url parameter `ip` address as needed. It is a mandatory parameter.

This REST endpoint will return result as (with status code `200` and `content-type: application/json`)

```
{
  "city": "Pittsburgh",
  "state": "Pennsylvania",
  "state_code": "PA",
  "zip_code": "15221",
  "country": "United States",
  "country_code": "US",
  "is_restricted": false,
  "is_cremia_region": false
}
```

While all errors will return `400` with a text error message. (Content type - text)

