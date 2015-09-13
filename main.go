package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) != 4 {
		log.Fatal("Usage: relup <user/repo> <tagname> <assetfile>\n  i.e. relup calmh/relup v1.0.0 relup-binary.tar.gz")
	}
	repo := os.Args[1]
	tag := os.Args[2]
	file := os.Args[3]
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Please export GITHUB_TOKEN=\"<your token here>\"")
	}

	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases", repo))
	if err != nil {
		log.Fatal(err)
	}

	var res []map[string]interface{}
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	err = dec.Decode(&res)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	var uploadUrl string
	for _, rel := range res {
		if rel["tag_name"] == tag {
			uploadUrl = rel["upload_url"].(string)
			break
		}
	}

	if uploadUrl == "" {
		log.Fatalln("Found no release with that tag")
	}

	fd, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	fi, _ := fd.Stat()

	log.Println("Uploading", path.Base(file))
	url := strings.Replace(uploadUrl, "{?name,label}", "?name="+path.Base(file), 1)
	req, err := http.NewRequest("POST", url, fd)
	if err != nil {
		log.Fatal(err)
	}

	req.ContentLength = fi.Size()
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", "token "+token)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.Status)
}
