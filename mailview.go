package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	// Some useful packages for go newbies looking to for ideas:
	/*
		"bufio"
		"html"
		"regexp"
		"time"
	*/
)

func main() {
	config, ini_err := ini.Load("settings.ini")
	if ini_err != nil {
		fmt.Printf("Failed to read configuration file: %v", ini_err)
		os.Exit(1)
	}
	//tls_cert := config.Section("server").Key("tls_cert").String()
	//tls_key := config.Section("server").Key("tls_key").String()
	// TODO: Should there be an error check for these? ^^

	port, port_err := config.Section("server").Key("http_port").Int()

	if port_err != nil {
		fmt.Printf("Invalid value for 'http_port' in 'server'.\nError: ", port_err)
		fmt.Printf("\nNo port set. We're using 8086 as a default.\n")
		port = 8086
	} else {
		fmt.Printf("Server listening on port: " + strconv.Itoa(port) + "\n")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		
		//// Enable this to output the request URL to the console:
		//fmt.Printf(r.URL.String())
		
		if r.URL.String() == "/" {
	/* EXAMPLE SECTION 1: This is what's returned at the root of your web server */
			response := "<!DOCTYPE html>\n<html>\n<head>\n</head>\n<body>\nHey there! You're running in the following mode:\n<br>\n"
			response += config.Section("").Key("app_mode").String()
			response += "\n<hr>\n<bold>Feel free to change that in your settings.ini</bold>\n</body>\n</html>"
			w.Write([]byte(response))
	/* End of SECTION 1 (Make sure to set a response variable and then w.Write([]byte(response)) ) */
	/* EXAMPLE SECTION 2+: This can be duplicated as many times as needed (or even removed) */
		} else if r.URL.String() == "/logs" {
			content, file_read_err := ioutil.ReadFile("logs.html")
			if file_read_err != nil {
				if config.Section("").Key("app_mode").String() == "production" {
					content = strings.Replace(content, "[production.min|development]", "production.min", -1)
				}
				if config.Section("").Key("app_mode").String() == "development" {
					content = strings.Replace(content, "[production.min|development]", "development", -1)
				}
				response := string(content)
			} else {
				response := "Error loading logs.html"
			}
			w.Write([]byte(response))
	/* End of SECTION 2+ (Make sure to change "/your-url-example/here" to a valid web path string) */
		} else {
			http.NotFound(w, r)
		}

	})

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))

}
