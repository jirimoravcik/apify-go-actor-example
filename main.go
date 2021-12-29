package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Example actor written in Go.")
	resp, err := http.Get("https://example.com")
	if err != nil {
		log.Fatal(err)
		return
	}

	default_kvs, token := os.Getenv("APIFY_DEFAULT_KEY_VALUE_STORE_ID"), os.Getenv("APIFY_TOKEN")
	if default_kvs == "" || token == "" {
		log.Fatal("Missing required env vars")
		return
	}
	client := http.Client{}
	url := fmt.Sprintf("https://api.apify.com/v2/key-value-stores/%v/records/OUTPUT?token=%v", default_kvs, token)
	req, _ := http.NewRequest(http.MethodPut, url, resp.Body)
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Saved fetched html to OUTPUT in key-value store.")
}
