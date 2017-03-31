package main

import (
	"encoding/json"
	"fmt"
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

type Project struct {
	Id      string                 `json:"id"`
	Key     string                 `json:"k"`
	Sc      string                 `json:"sc"`
	Qu      string                 `json:"qu"`
	Lv      string                 `json:"lv"`
	Version map[string]VersionList `json:"v,string"`
}

var wg sync.WaitGroup

//type Projects []Project

func deleteKey(key string) {
    
    defer wg.Done()
    
	var dat = url.Values{"key": {key}}
	p := fmt.Println
    p("Deleteing: " + key)
	client := http.Client{}
	r, err := client.PostForm("https://****:*****@sonarqube-host.domain/api/projects/delete", dat)
	if err != nil {
		log.Fatal(err)
	}
	//rr, _ := ioutil.ReadAll(r.Body)
	//p(string(rr))
    defer r.Body.Close()
		
	p("Deleted: " + key)
	p(r.StatusCode)
}

func main() {
	maxProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcs)
	var data []Project
	// var timeout = time.Duration(60 * time.Second)
	client := http.Client{
	//Timeout: timeout,
	}

	resp, err := client.Get("https://****:*****@sonarqube-host.domain/api/projects?versions=true")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&data)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		switch v := err.(type) {
		case *json.SyntaxError:
			fmt.Println("JSON error %s", v)
		}
	}
	p := fmt.Println
	//fmt.Printf("Results: %v\n", data)

	for _, track := range data {
		for a, b := range track.Version {

			t1, _ := time.Parse("2006-01-02T15:04:05+0000", string(b.Date))

			if t1.Before(time.Now().AddDate(0, -1, 0)) {
				//if time.Now().Sub(t1).Hours>time.Duration.Hours{
				//if t1. < time.Now().AddDate(0,-1,0){
				if !strings.Contains(track.Key, ":master") {
					if !strings.Contains(track.Key, ":cake2") {
						fmt.Printf("Project:%s Version:%s Date: %s Key: %s\n", string(track.Id), string(a), string(b.Date), string(track.Key))
						p(t1)
						p("To delete\n")
						wg.Add(1)
						deleteKey(track.Key)
					}
				}
			}
		}
	}

	p("Waiting for all goroutines...")
	wg.Wait()
	p("Done")

}
