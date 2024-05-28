package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	WellKnownOIDCConf = ".well-known/openid-configuration"
)

func main() {
	envs := os.Environ()
	fmt.Println("Printing all Environmental Variables")
	for _, e := range envs {
		fmt.Println("   " + e)
	}
	//Retrieve base url
	baseUrl, ok := os.LookupEnv("OIDC_BASE_URL")
	if !ok {
		fmt.Println("OIDC_BASE_URL not defined.  Terminating")
		os.Exit(1)
	}

	fmt.Println("Staring endless for loop")
	for {
		fmt.Println("OIDC_BASE_URL = " + baseUrl)
		fmt.Println("Retrieving well-known/openid-configuration")

		oidcConf, err := url.Parse(WellKnownOIDCConf)
		if err != nil {
			fmt.Println("can't parse well known oidc")
			fmt.Println("err:", err)
		}
		base, err := url.Parse(baseUrl)
		if err != nil {
			fmt.Println("can't parse issuer url")
			fmt.Println("err:", err)
		}
		discoveryURL := base.ResolveReference(oidcConf)
		fmt.Println("discoveryUrl = " + discoveryURL.String())

		req, err := http.NewRequest("GET", discoveryURL.String(), nil)
		if err != nil {
			fmt.Println("error creating http.NewRequest")
			fmt.Println("err:", err)
			fmt.Println("req:", req)
			fmt.Println("Sleeping for 60 seconds")
			time.Sleep(60 * time.Second)
			fmt.Println("Sleep Over.....")
			continue
		}
		fmt.Println("req:", req)
		client := http.DefaultClient
		client.Timeout = 5 * time.Minute
		fmt.Println("client:", client)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("error client.Do")
			fmt.Println("err:", err)
			fmt.Println("req:", req)
			fmt.Println("client:", client)
			fmt.Println("Sleeping for 60 seconds")
			time.Sleep(60 * time.Second)
			fmt.Println("Sleep Over.....")
			continue
		}
		fmt.Println("resp:", resp)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("request failed with code " + resp.Status)
		} else {
			// read the data
			reader := &io.LimitedReader{R: resp.Body, N: 1024 * 1024 * 10}
			data, err := io.ReadAll(reader)
			if err != nil {
				fmt.Println("error reaing data")
				fmt.Println("err:", err)
				fmt.Println("Sleeping for 60 seconds")
				time.Sleep(60 * time.Second)
				fmt.Println("Sleep Over.....")
				continue
			}

			if reader.N <= 0 {
				fmt.Println("max body length exceeded")
			}

			fmt.Println(string(data))
		}
		fmt.Println("Sleeping for 60 seconds")
		time.Sleep(60 * time.Second)

		// Printed after sleep is over
		fmt.Println("Sleep Over.....")
	}
}
