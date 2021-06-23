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
	"sync"
)

func getMD5(url string) string{
	resp, err := http.Get(url)
	if err != nil {
		//returning empty string
		return ""
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return getHash(content)
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
	sem := make(chan struct{}, *maxP)
	result := map[string]string{}
	var wg sync.WaitGroup
	for i := 1; i < len(os.Args); i++ {
		wg.Add(1)
		sem <- struct {}{}
		go func(i int) {
			defer wg.Done()
			validUrl := getValidUrl(os.Args[i])
			hContent := getMD5(validUrl)
			if hContent != "" {
				result[validUrl] = hContent
			}
			//release sem
			<- sem
		}(i)
	}
	wg.Wait()
	for k,v := range result {
		fmt.Println(k, v)
	}
}
