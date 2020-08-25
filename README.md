# geoip
Maxmind geoip IP  to region extraction.

# Getting started

> All following steps assume you are on `macOS Catalina`

## Installation of go

```
$ brew install golang
```

Assuming your `$HOME` is `/Users/<your_mac_user_name>`

```
$ mkdir -p $HOME/go
$ vi ~/.bash_profile
```
and add

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

Now go to http://localhost:8080/geolocations?ip=24.3.77.32

Change `ip` address as needed.

