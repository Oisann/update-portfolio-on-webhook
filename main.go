package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func serveGithubWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	logError(err)
	secret := []byte(os.Getenv("SECRET"))
	message := []byte(body)
	hash := hmac.New(sha1.New, secret)
	hash.Write(message)

	hashResult := hex.EncodeToString(hash.Sum(nil))

	remoteHash := r.Header.Get("X-Hub-Signature")

	if ("sha1=" + hashResult) == remoteHash {
		var y map[string]interface{}
		json.Unmarshal(body, &y)
		repo := y["repository"].(map[string]interface{})["full_name"]
		branch := y["repository"].(map[string]interface{})["default_branch"]
		id := y["repository"].(map[string]interface{})["id"]
		url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/PORTFOLIO", repo, branch)
		file := fmt.Sprintf("%s/%.0f.portfolio", os.Getenv("DB_FOLDER"), id)
		err := downloadFile(file, url)
		if err != nil {
			logError(err)
		}
		fmt.Println("Hashes matched, ran command!")
	} else {
		fmt.Println("Hashes didn't match, doing nothing.")
	}

	fmt.Fprintf(w, hashResult)
}

func downloadFile(filepath string, url string) (err error) {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Repo does not have a PORTFOLIO file in its root folder")
		out.Close()
		os.RemoveAll(fmt.Sprintf("%s/%s", os.Getenv("DB_FOLDER"), out.Name()))
		return nil
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func logError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	http.HandleFunc("/", serveGithubWebhook)
	log.Fatal(http.ListenAndServe(":80", nil))
}
