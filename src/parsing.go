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
	address = compileAddress(state, street+" "+houseNumber)
	return
}

func compileAddress(state, street string) (address string) {
	address = "Vilnius"
	if state != "" {
		address += ", " + state
	}
	if street != "" {
		address += ", " + street
	}
	return
}

func getBytes(link string) ([]byte, error) {

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
	}
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
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func getGoqueryDocument(link string) (*goquery.Document, error) {

	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
	}

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
	return goquery.NewDocumentFromReader(res.Body)
}
