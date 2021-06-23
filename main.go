package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func getMD5(url string, w chan string) {
	resp, err := http.Get(url)
	if err != nil {
		//returning empty string
		w <- ""
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w <- ""
		return
	}
	w <- getHash(content)
	return
}

func getHash(content []byte) string {
	newHash := md5.New()
	newHash.Write(content)
	return hex.EncodeToString(newHash.Sum(nil))
}

func getValidUrl(rawUrl string) string {
	u, _ := url.Parse(rawUrl)
	if u.Scheme == "" {
		return "http://"+u.Path
	}
	return rawUrl
}

func main(){
	if len(os.Args) < 2 {
		log.Fatal("need atleast 1 url")
	}
	maxP := flag.Int("parallel", 10, "to prevent exhausting local resources")
	flag.Parse()
	ch := make(chan string, *maxP)
	result := map[string]string{}
	waitForAll := make(chan bool)
	go func() {
		for i := 1; i < len(os.Args); i++ {
			validUrl := getValidUrl(os.Args[i])
			getMD5(validUrl, ch)
			hContent := <- ch
			if hContent != "" {
				result[validUrl] = hContent
			}
		}
		waitForAll <- true
	}()
	<-waitForAll
	for k,v := range result {
		fmt.Println(k, v)
	}
}
