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
)

type Post struct {
	Output string `json:"output"`
}

type Response struct {
	Jobs []Job `json:"jobs"`
}

type Job struct {
	Id             string `json:"_id"`
	Jobname        string `json:"jobName"`
	Jobdescription string `json:"jobDescription"`
	Shellcommand   bool   `json:"shellCommand"`
	Command        string `json:"command"`
}

func getJobs() Response {
	resp, err := http.Post("https://phalerum.stickybits.red/api/v1/agents/hello", "application/json", nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body) // response body is []byte
	var result Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	return result
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

func checkReq(req Response) {
	for i := 0; i < len(req.Jobs); i++ {

		if req.Jobs[i].Shellcommand == false {
			switch req.Jobs[i].Command {
			case "builtin.firewall":

				url := fmt.Sprintf("https://phalerum.stickybits.red/api/v1/jobs/output/%s", req.Jobs[i].Id)

				jsonBody := Post{
					Output: encodeOutput(firewallEnabled()),
				}

				bodyBytes, err := json.Marshal(&jsonBody)
				if err != nil {
					log.Fatal(err)
				}

				reader := bytes.NewReader(bodyBytes)
				postOut(url, reader)

			case "builtin.password":

				url := fmt.Sprintf("https://phalerum.stickybits.red/api/v1/jobs/output/%s", req.Jobs[i].Id)

				jsonBody := Post{
					Output: encodeOutput(checkPasswordDate()),
				}

				bodyBytes, err := json.Marshal(&jsonBody)
				if err != nil {
					log.Fatal(err)
				}

				reader := bytes.NewReader(bodyBytes)
				postOut(url, reader)
			}
		}
	}
}

func postOut(url string, reader io.Reader) {
	_, err := http.Post(url, "application/json", reader)
	if err != nil {
		log.Fatal(err)
	}
}

func encodeOutput(output []byte) string {
	sEnc := base64.StdEncoding.EncodeToString([]byte(string(output)))

	return sEnc
}

func xorURL(url []byte, key []byte) string {
	for i := 0; i < len(url)-1; i++ {
		url[i] ^= key[i%len(key)]
	}
	return string(url)
}

func firewallEnabled() []byte {
	cmd := exec.Command("ufw", "status")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(cmd)
		fmt.Println("Error running command: ", err)
	}

	sEnc := base64.StdEncoding.EncodeToString([]byte(string(out)))

	if sEnc == "Status: inactive" {
		return []byte("false")
	} else {
		return []byte("false")
	}
}

func checkPasswordDate() []byte {
	cmd := exec.Command("bash", "-c", "echo -n $(chage -l chinou | awk '{print $5\" \"$6\" \"$7}' | head -n 1)")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(cmd)
		fmt.Println("Error running command: ", err)
	}

	return out
}

func main() {
	// eUrl := []byte{0xB1, 0xC6, 0x3F, 0x31, 0x2E, 0x3F, 0x41, 0xD2, 0xF9, 0xAC, 0x20, 0xB3, 0x88, 0x39, 0xF5, 0x8, 0x1A, 0xE3, 0x99, 0xAF, 0xD7, 0xC, 0xEA, 0xC9, 0x34, 0x62, 0x11, 0x38, 0x9B, 0xCF, 0xC5, 0x2C, 0xF2, 0x2D, 0x2F, 0x69, 0x5B, 0x9B, 0x31, 0x71, 0x31, 0xF6, 0x90, 0xA5, 0x91, 0xFF, 0xA1, 0x75, 0xEB, 0x74}
	// key := []byte{0xd9, 0xb2, 0x4b, 0x41, 0x5d, 0x05, 0x6e, 0xfd, 0x89, 0xc4, 0x41, 0xdf, 0xed, 0x4b, 0x80, 0x65, 0x34, 0x90, 0xed, 0xc6, 0xb4, 0x67, 0x93, 0xab, 0x5d, 0x16, 0x62, 0x16, 0xe9, 0xaa, 0xa1, 0x03, 0x93, 0x5d, 0x46, 0x46, 0x2d, 0xaa, 0x1e, 0x10, 0x56, 0x93, 0xfe, 0xd1, 0xe2, 0xd0, 0xd5, 0x10, 0x98, 0xed, 0x5a, 0x19, 0x55, 0x1b, 0x1e, 0x88, 0x3e, 0x41, 0x31, 0x92, 0x22, 0x70, 0x9c, 0x1a}

	// url := xorURL(eUrl, key)

	req := getJobs()
	checkReq(req)

}
