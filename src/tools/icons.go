package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const apiUrl = "https://rektdeckard-fontpack.web.val.run"
const srcFile = "./icons.json"
const dstFile = "../web/assets/phosphor.css"

func main() {
	fmt.Println("CucinAssistant Icons Packer")
	fmt.Println("===========================")

	// Loads the source file
	var icons map[string][]string
	fmt.Printf("Loading icons ids from %s...\n", srcFile)
	src, err := ioutil.ReadFile(srcFile)
	check(err, "reading file")
	check(json.Unmarshal(src, &icons), "parsing file")

	// Prepares the request
	data := make(map[string]any)
	data["formats"] = []string{"woff"}
	data["inline"] = true
	data["icons"] = icons
	buffer, err := json.Marshal(data)
	check(err, "preparing the request data")
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(buffer))
	check(err, "preparing the request")

	// Sends the request
	client := &http.Client{}
	res, err := client.Do(req)
	check(err, "during the request")
	defer res.Body.Close()

	// Checks the status code
	if res.StatusCode != 200 {
		fmt.Printf("Expected status code 200, got %d.\n", err.Error())
		os.Exit(1)
	}

	// Reads and parses the response
	body, err := io.ReadAll(res.Body)
	check(err, "reading the response")
	check(json.Unmarshal(body, &data), "parsing the response")

	// Writes the response to the file
	fmt.Printf("Saving the icons pack to %s...\n", dstFile)
	err = os.WriteFile(dstFile, []byte(data["css"].(string)), 0644)
	if err != nil {
		fmt.Printf("Error writing the file: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Done!")
}

func check(err error, msg string) {
	if err != nil {
		fmt.Printf("Error %s: %s\n", msg, err.Error())
		os.Exit(1)
	}
}
