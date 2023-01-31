package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ZonesResponse struct {
	Result []Zone `json:"result"`
}

type PurgeResponse struct {
	Result  interface{} `json:"result"`
	Success bool        `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Messages []string `json:"messages"`
}

const (
	EMAIL  = "flamezaxaou1@gmail.com"
	APIKEY = "d78b50397907351e3b6e8ac41a0eae60a9136"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Input Domain [ALL, domain.com]: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimRight(text, "\n")
	text = strings.TrimRight(text, "\r")
	text = strings.Replace(text, "%0D", "", 1)
	text, _ = url.QueryUnescape(text)
	fmt.Println("Input domain: ", text)

	if text == "ALL" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("confrim purge cache ALL domain [Y/n] :")
		confrim, _ := reader.ReadString('\n')
		confrim = strings.TrimRight(confrim, "\n")
		confrim = strings.TrimRight(text, "\r")
		confrim, _ = url.QueryUnescape(confrim)
		if confrim != "Y" {
			fmt.Println("exit program...")
			time.Sleep(10 * time.Second)
			return
		}
	}

	// Replace with your Cloudflare API Key
	APIKEY := APIKEY

	url := "https://api.cloudflare.com/client/v4/zones?per_page=1000"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("X-Auth-Email", EMAIL)
	req.Header.Set("X-Auth-Key", APIKEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var zonesResponse ZonesResponse
	err = json.Unmarshal(body, &zonesResponse)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	fmt.Println("total domain in profile", len(zonesResponse.Result))

	for index, zone := range zonesResponse.Result {

		if text == zone.Name || text == "ALL" {
			fmt.Println("Purging cache for zone with ID:", zone.ID)
			fmt.Println("Purging cache for domain:", index+1, zone.Name)

			// Purge cache for the current zone
			purgeUrl := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/purge_cache", zone.ID)
			purgeReq, err := http.NewRequest("POST", purgeUrl, nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}

			purgeReq.Header.Set("X-Auth-Email", EMAIL)
			purgeReq.Header.Set("X-Auth-Key", APIKEY)
			purgeReq.Header.Set("Content-Type", "application/json")

			purgeReqBody, err := json.Marshal(map[string]interface{}{
				"purge_everything": true,
			})
			if err != nil {
				fmt.Println("Error marshalling request body:", err)
				return
			}
			purgeReq.Body = ioutil.NopCloser(bytes.NewReader(purgeReqBody))

			purgeResp, err := client.Do(purgeReq)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			defer purgeResp.Body.Close()

			purgeBody, err := ioutil.ReadAll(purgeResp.Body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				return
			}

			var purgeResponse PurgeResponse
			err = json.Unmarshal(purgeBody, &purgeResponse)
			if err != nil {
				fmt.Println("Error unmarshalling response:", err)
				return
			}

			if purgeResponse.Success {
				fmt.Println("Cache purged successfully for zone with ID:", zone.ID)
				fmt.Println("Cache purged successfully for domain:", zone.Name)
			} else {
				fmt.Println("Error purging cache for zone with ID:", zone.ID)
				fmt.Println("Error messages:", purgeResponse.Messages)
			}

			if text == zone.Name {
				time.Sleep(10 * time.Second)
				return
			}
		}
	}
	fmt.Println("Not found domain", text)
	fmt.Println("exit program...")
	time.Sleep(10 * time.Second)
	return
}

// GOOS=windows GOARCH=amd64 go build -o cloudflare-purgecache-v1.exe
