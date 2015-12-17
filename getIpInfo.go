package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage:\n Please Input Ip Or Hostname\n")
		os.Exit(0)
	}

	host := os.Args[1]

	if host == "127.0.0.1" || host == "localhost" {
		fmt.Fprintf(os.Stderr, "Locahost")
		os.Exit(0)
	}

	if ip := IsIP(host); ip {
		getIpInfo(host)
		return
	}

	u, err := url.Parse(host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Domain Parse Is Error ")
		os.Exit(0)
	}

	if strings.Contains(host, "http://") {
		getDomainIp(u.Host)
	} else {
		getDomainIp(host)
	}
}

func getDomainIp(domain string) {
	addrs, err := net.LookupHost(domain)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(2)
	}

	for _, s := range addrs {
		getIpInfo(s)
	}

	os.Exit(0)
}

func getIpInfo(ip string) {

	u, _ := url.Parse("http://ip.taobao.com/service/getIpInfo.php")
	q := u.Query()
	q.Set("ip", ip)
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	js, err := simplejson.NewJson(result)
	if err != nil {
		panic(err.Error())
	}

	if v, err := js.Get("code").Int(); v != 0 {
		fmt.Println("Get Ip Info Error: " + err.Error())
		return
	}

	jsonIpInfo := js.Get("data")
	country, _ := jsonIpInfo.Get("country").String()
	area, _ := jsonIpInfo.Get("area").String()
	region, _ := jsonIpInfo.Get("region").String()
	city, _ := jsonIpInfo.Get("city").String()
	isp, _ := jsonIpInfo.Get("isp").String()

	fmt.Println("[" + ip + "][" + isp + "] " + country + "-" + area + "-" + region + "-" + city)
}

func IsIP(ip string) (b bool) {
	if m, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", ip); !m {
		return false
	}
	return true
}
