package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var checkEndpoint = []byte{0xb1, 0xc6, 0x3f, 0x31, 0x2e, 0x3f, 0x41, 0xd2, 0xf9, 0xac, 0x20, 0xb3, 0x88, 0x39, 0xf5, 0x8, 0x1a, 0xe3, 0x99, 0xaf, 0xd7, 0xc, 0xea, 0xc9, 0x34, 0x62, 0x11, 0x38, 0x9b, 0xcf, 0xc5, 0x2c, 0xf2, 0x2d, 0x2f, 0x69, 0x5b, 0x9b, 0x31, 0x71, 0x31, 0xf6, 0x90, 0xa5, 0x91, 0xff, 0xbd, 0x75, 0xf4, 0x81, 0x6f}
var postEndpoint = []byte{0xb1, 0xc6, 0x3f, 0x31, 0x2e, 0x3f, 0x41, 0xd2, 0xf9, 0xac, 0x20, 0xb3, 0x88, 0x39, 0xf5, 0x8, 0x1a, 0xe3, 0x99, 0xaf, 0xd7, 0xc, 0xea, 0xc9, 0x34, 0x62, 0x11, 0x38, 0x9b, 0xcf, 0xc5, 0x2c, 0xf2, 0x2d, 0x2f, 0x69, 0x5b, 0x9b, 0x31, 0x7a, 0x39, 0xf1, 0x8d, 0xfe, 0x8d, 0xa5, 0xa1, 0x60, 0xed, 0x99, 0x75, 0x3c, 0x73}
var key = []byte{0xd9, 0xb2, 0x4b, 0x41, 0x5d, 0x05, 0x6e, 0xfd, 0x89, 0xc4, 0x41, 0xdf, 0xed, 0x4b, 0x80, 0x65, 0x34, 0x90, 0xed, 0xc6, 0xb4, 0x67, 0x93, 0xab, 0x5d, 0x16, 0x62, 0x16, 0xe9, 0xaa, 0xa1, 0x03, 0x93, 0x5d, 0x46, 0x46, 0x2d, 0xaa, 0x1e, 0x10, 0x56, 0x93, 0xfe, 0xd1, 0xe2, 0xd0, 0xd5, 0x10, 0x98, 0xed, 0x5a, 0x19, 0x55, 0x1b, 0x1e, 0x88, 0x3e, 0x41, 0x31, 0x92, 0x22, 0x70, 0x9c, 0x1a}

var checkinUrl = xorURL(checkEndpoint, key)
var postUrl = xorURL(postEndpoint, key)

var Comtoken string

// service struct
type tickerService struct {
}

// NewTickerService returns new service
func NewTickerService() *tickerService {
	return &tickerService{}
}

type Post struct {
	Comtoken string `json:"communicationToken"`
	Output   string `json:"output"`
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

// Run starts service
func (ds *tickerService) Run(ctx context.Context) {
	MinNmbr := 1
	MaxNmbr := 5
	rndNmbr := rand.Intn(MaxNmbr-MinNmbr) + MinNmbr

	ticker := time.NewTicker(time.Duration(rndNmbr) * time.Second)
	for {
		select {
		case <-ticker.C:
			go func() {
				MinNmbr := 1
				MaxNmbr := 5
				rndNmbr := rand.Intn(MaxNmbr-MinNmbr) + MinNmbr

				ds.task()
				ticker.Reset(time.Duration(rndNmbr) * time.Second)
			}()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

// periodic task
func (ds *tickerService) task() {
	req := getJobs(checkinUrl)
	checkReq(req, postUrl)
}

func getJobs(url string) Response {
	jsonBody := Post{
		Comtoken: Comtoken,
		Output:   "",
	}

	bodyBytes, err := json.Marshal(&jsonBody)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(bodyBytes)

	resp, err := http.Post(url, "application/json", reader)
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

func createBody(jsonBody Post) io.Reader {

	bodyBytes, err := json.Marshal(&jsonBody)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(bodyBytes)

	return reader
}

func checkReq(req Response, url string) {
	for i := 0; i < len(req.Jobs); i++ {

		if req.Jobs[i].Shellcommand == false {
			switch req.Jobs[i].Command {
			case "builtin.firewall":

				url := fmt.Sprintf(url, req.Jobs[i].Id)

				jsonBody := Post{
					Comtoken: Comtoken,
					Output:   encodeOutput(firewallEnabled()),
				}

				postOut(url, createBody(jsonBody))

			case "builtin.password":

				url := fmt.Sprintf(url, req.Jobs[i].Id)

				jsonBody := Post{
					Comtoken: Comtoken,
					Output:   encodeOutput(checkPasswordDate()),
				}

				postOut(url, createBody(jsonBody))
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
	cmd := exec.Command("ls", "-la")
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
	cmd := exec.Command("bash", "-c", "echo -n $(chage -l student | awk '{print $5\" \"$6\" \"$7}' | head -n 1)")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(cmd)
		fmt.Println("Error running command: ", err)
	}

	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ds := NewTickerService()
	go ds.Run(ctx)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-shutdown
	cancel()
}
