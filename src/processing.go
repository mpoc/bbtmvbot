package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type post struct {
	url         string
	phone       string
	description string
	address     string
	heating     string
	floor       int
	floorTotal  int
	area        int
	price       int
	rooms       int
	year        int
}

// Must be lowercase!!!
var exclusionKeywords = []string{
	" bus taikomas vienkartinis agentūros mokestis",
	" bus taikomas vienkartinis agentūrinis mokestis",
	" bus taikomas vienkartinis agenturos mokestis",
	" bus taikomas vienkartinis agenturinis mokestis",
	" bus taikomas vienkartinis tarpininkavimo mokestis",
	",bus taikomas vienkartinis agentūros mokestis",
	",bus taikomas vienkartinis agentūrinis mokestis",
	",bus taikomas vienkartinis agenturos mokestis",
	",bus taikomas vienkartinis agenturinis mokestis",
	",bus taikomas vienkartinis tarpininkavimo mokestis",
	"yra taikomas vienkartinis agentūros mokestis",
	"yra taikomas vienkartinis agentūrinis mokestis",
	"yra taikomas vienkartinis agenturos mokestis",
	"yra taikomas vienkartinis agenturinis mokestis",
	"yra taikomas vienkartinis tarpininkavimo mokestis",
	"vienkartinis agentūros mokestis jei",
	"vienkartinis agentūrinis mokestis jei",
	"vienkartinis agenturos mokestis jei",
	"vienkartinis agenturinis mokestis jei",
	"vienkartinis tarpininkavimo mokestis jei",
	" tiks, bus taiko",
	" tiks bus taiko",
	"taikomas vienkartinis tarpininkavimo mokestis",
	"tiks vienkartinis tarpininkavimo mokestis",
	"tarpininkavimo mokestis-",
	"tarpininkavimo mokestis -",
	"(yra mokestis)",
	" bus imamas vienkartinis",
	" bus imamas tarpininkavimo",
	" bus taikomas vienkartinis",
	".bus taikomas vienkartinis",
	",bus taikomas vienkartinis",
	" bus taikomas tarpininkavimo",
	"mokestis (jei butas tiks)",
	"ir imamas vienkartinis mokestis",
	",yra vienkartinis agent",
	" yra vienkartinis agent",
	".yra vienkartinis agent",
}

var regexExclusion1 = regexp.MustCompile(`(agenturos|agenturinis|agentūrinis|tarpininkavimo) mokestis[\s:]{0,3}\d+`)

// Note that post is already checked against DB in parsing functions!
func processPost(p post) {

	// Check if description contains exclusion keyword
	desc := strings.ToLower(p.description)
	for _, v := range exclusionKeywords {
		if strings.Contains(desc, v) {
			fmt.Println(">> Excluding", p.url, "reason:", v)
			databaseAddPost(p)
			return
		}
	}

	// Passed blacklisted keywords test, so let's do some regex tests
	arr := regexExclusion1.FindStringSubmatch(desc)
	if len(arr) >= 1 {
		fmt.Println(">> Excluding", p.url, "reason: /regex1/")
		databaseAddPost(p)
		return
	}

	// Skip posts without price and let user know:
	if p.price == 0 {
		fmt.Println(">> 0eur price", p.url)
		return
	}

	// Add to database, so it won't be sent again:
	insertedRowID := databaseAddPost(p)

	// Send to users
	databaseGetUsersAndSendThem(p, insertedRowID)

	// Show debug info:
	p.description = strconv.Itoa(len(p.description))
	fmt.Println(p)
}

func getCompiledMessage(p post, ID int64) string {
	message := strconv.FormatInt(ID, 10) + ". " + p.url + "\n"

	if p.phone != "" {
		message += "» *Tel. numeris:* [" + p.phone + "](tel:" + p.phone + ")\n"
	}

	if p.address != "" {
		message += "» *Adresas:* [" + p.address + "](https://maps.google.com/?q=" + url.QueryEscape(p.address) + ")\n"
	}

	if p.price != 0 && p.area != 0 {
		message += "» *Kaina:* `" + strconv.Itoa(p.price) + "€ (" + strconv.Itoa(int(p.price/p.area)) + "€/m²)`\n"
	} else if p.price != 0 {
		message += "» *Kaina:* `" + strconv.Itoa(p.price) + "€`\n"
	}

	if p.rooms != 0 && p.area != 0 {
		message += "» *Kambariai:* `" + strconv.Itoa(p.rooms) + " (" + strconv.Itoa(p.area) + "m²)`\n"
	} else if p.rooms != 0 {
		message += "» *Kambariai:* `" + strconv.Itoa(p.rooms) + "`\n"
	}

	if p.year != 0 {
		message += "» *Statybos metai:* `" + strconv.Itoa(p.year) + "`\n"
	}

	if p.heating != "" {
		message += "» *Šildymo tipas:* `" + p.heating + "`\n"
	}

	if p.floor != 0 && p.floorTotal != 0 {
		message += "» *Aukštas:* `" + strconv.Itoa(p.floor) + "/" + strconv.Itoa(p.floorTotal) + "`\n"
	} else if p.floor != 0 {
		message += "» *Aukštas:* `" + strconv.Itoa(p.floor) + "`\n"
	}

	return strings.TrimSpace(message)
}
