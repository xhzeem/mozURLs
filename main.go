package main

import (
    "bytes"
    "fmt"
    "flag"
    "os"
    "bufio"
    "encoding/json"
    "io/ioutil"
    "net/http"
)

func main() {

    target := flag.String("t", "", "Domain of website to get URLs")
    moz_api_key := flag.String("k", os.Getenv("MOZ_API_KEY"), "MOZ API key for basic Authorization b64(u:p)")
    
    flag.Parse()

    if *target == "" {
        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        *target = scanner.Text()
    }

    if *moz_api_key == "" {
        fmt.Println("Please provide API key using -k flag or via MOZ_API_KEY env variable")
        os.Exit(1)
    }
    
    nextToken := ""

    for {
        reqBody, _ := json.Marshal(map[string]interface{}{
            "target": target,
            "limit": 50,
            "next_token": nextToken,
        })

        req, _ := http.NewRequest("POST", "https://lsapi.seomoz.com/v2/top_pages", bytes.NewBuffer(reqBody))
        req.Header.Add("Authorization", "Basic " + *moz_api_key)
        req.Header.Add("Content-Type", "application/json")

        client := &http.Client{}
        resp, _ := client.Do(req)

        defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)

        var data map[string]interface{}
        json.Unmarshal([]byte(body), &data)

        results, _ := json.Marshal(data["results"])
        var pages []map[string]interface{}
        json.Unmarshal(results, &pages)

        for _, page := range pages {
            fmt.Println("http://" + page["page"])
        }
                
        nextToken = data["next_token"].(string)
        if nextToken == "" {
            break
        }
    }
}
