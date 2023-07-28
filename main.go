package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Resp struct {
	Data []Data
}
type Data struct {
	Id       int
	FileName string
}

func main() {
	var version *string = flag.String("version", "latest", "")
	var modpackId *int = flag.Int("pack", -1, "")
	flag.Parse()

	if *modpackId == -1 {
		fmt.Println("Please provide modpack id with '--pack <id>' option")
		os.Exit(1)
	}

	url := fmt.Sprintf("https://www.curseforge.com/api/v1/mods/%d/files", *modpackId)
	modpackFile := getData(url, version)

	additionalFilesUrl := fmt.Sprintf("https://www.curseforge.com/api/v1/mods/%d/files/%d/additional-files", *modpackId, modpackFile.Id)
	serverFileData := getData(additionalFilesUrl, version)

	downloadUrl := fmt.Sprintf("https://www.curseforge.com/api/v1/mods/%d/files/%d/download", *modpackId, serverFileData.Id)

	file, err := os.Create("server.zip")
	if err != nil {
		fmt.Printf("error creating file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	resp, err := http.Get(downloadUrl)
	if err != nil {
		fmt.Printf("error downloading server file: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("error saving server file: %s\n", err)
		os.Exit(1)
	}

}

func getData(url string, modpackVersion *string) Data {
	res, err := http.Get(url)

	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var jsonresp Resp

	json.Unmarshal(body, &jsonresp)

	if *modpackVersion == "latest" {
		return jsonresp.Data[len(jsonresp.Data)-1]
	}

	var file Data
	for i, e := range jsonresp.Data {
		if strings.Contains(e.FileName, *modpackVersion) {
			file = jsonresp.Data[i]
			break
		}
	}

	return file
}
