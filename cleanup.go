package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"
)

type VersionList struct {
	Sid  string `json:"sid"`
	Date string `json:"d"`
}

type Page struct {
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
	Total     int `json:"total"`
}

type Response struct {
	Paging     Page      `json:"paging"`
	Components []Project `json:"components"`
}

type Project struct {
	Organisation string `json:"organisation"`
	Id           string `json:"id"`
	Key          string `json:"key"`
	Name         string `json:"name"`
	Qualifier    string `json:"qualifier"`
}

type Component struct {
	Organization string   `json:"organization"`
	Id           string   `json:"id"`
	Key          string   `json:"key"`
	Name         string   `json:"name"`
	Qualifier    string   `json:"qualifier"`
	AnalysisDate string   `json:"analysisDate"`
	Tags         []string `json:"tags"`
	Visibility   string   `json:"visibility"`
}

type ComponentResp struct {
	Component Component `json:"component"`
	Ancestors []string  `json:"ancestors"`
}

var wg sync.WaitGroup

//type Projects []Project

func deleteKey(key string) {

	defer wg.Done()

	var dat = url.Values{"project": {key}}
	p := fmt.Println
	p("Deleteing: " + key)
	client := http.Client{}
	r, err := client.PostForm("https://user:pass@sonar.example.com/api/projects/delete", dat)
	if err != nil {
		log.Fatal(err)
	}
	rr, _ := ioutil.ReadAll(r.Body)
	p(string(rr))
	defer r.Body.Close()

	p("Deleted: " + key)
	p(r.StatusCode)
}

func main() {
	maxProcs := runtime.NumCPU()
	p := fmt.Println
	runtime.GOMAXPROCS(maxProcs)
	var data Response
	// var timeout = time.Duration(60 * time.Second)
	client := http.Client{
	//Timeout: timeout,
	}

	resp, err := client.Get("https://user:pass@sonar.example.com/api/components/search?qualifiers=TRK&pageSize=1000")

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//rr, _ := ioutil.ReadAll(resp.Body)
	//p(string(rr))

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		switch v := err.(type) {
		case *json.SyntaxError:
			fmt.Println("JSON error %s", v)
		}
	}

	for _, component := range data.Components {
		resp, err = client.Get("https://user:pass@sonar.example.com/api/components/show?component=" + component.Key)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var comp ComponentResp
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&comp)
		if err != nil {
			fmt.Printf("%T\n%s\n%#v\n", err, err, err)
			switch v := err.(type) {
			case *json.SyntaxError:
				fmt.Println("JSON error %s", v)
			}
		}

		t1, _ := time.Parse("2006-01-02T15:04:05+0000", string(comp.Component.AnalysisDate))
		if t1.Before(time.Now().AddDate(0, -3, 0)) {
			if !strings.Contains(comp.Component.Key, ":master") && !strings.Contains(comp.Component.Key, ":cake2") && !strings.Contains(comp.Component.Key, ":production") && !strings.Contains(comp.Component.Key, ":live") {
				fmt.Printf("Project:%s Anylysis Date: %s Key: %s\n", string(comp.Component.Id), string(comp.Component.AnalysisDate), string(comp.Component.Key))
				p(t1)
				p("To delete\n")
				wg.Add(1)
				go deleteKey(comp.Component.Key)
			}
		}

	}

	p("Waiting for all goroutines...")
	wg.Wait()
	p("Done")

}
