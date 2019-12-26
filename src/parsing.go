package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

var userAgent = "Mozilla/5.0 (Linux; Android 9; SAMSUNG GT-I9505 Build/LRX22C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.93 Mobile Safari/537.36"

func compileAddressWithStreet(state, street, houseNumber string) (address string) {
	if state == "" {
		address = "Vilnius"
	} else if street == "" {
		address = "Vilnius, " + state
	} else if houseNumber == "" {
		address = "Vilnius, " + state + ", " + street
	} else {
		address = "Vilnius, " + state + ", " + street + " " + houseNumber
	}
	return
}

func compileAddress(state, street string) (address string) {
	if state == "" {
		address = "Vilnius"
	} else if street == "" {
		address = "Vilnius, " + state
	} else {
		address = "Vilnius, " + state + ", " + street
	}
	return
}

func getBytes(link string) ([]byte, error) {
	res, err := sendGetRequest(link)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func getGoqueryDocument(link string) (*goquery.Document, error) {
	res, err := sendGetRequest(link)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(res.Body)
}

func sendGetRequest(link string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		u, _ := url.Parse(link)
		return nil, fmt.Errorf("status code error: %s (from %s)", res.Status, u.Host)
	}
	return res, nil
}
