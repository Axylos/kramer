package kramer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type KramerPoller struct {
	uri     url.URL
	done    bool
	channel chan PollResp
	token   string
}

type pollResp struct {
	Done        bool
	AccessToken string
}

func (poller *KramerPoller) pollCheck() {
	resp, err := http.Get(poller.uri.String())
	if err != nil {
		fmt.Println(err)
		println("oopsie")
	}

	var pResp pollResp
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &pResp)

	poller.done = pResp.Done
	poller.token = pResp.AccessToken
}

func (poller *KramerPoller) Poll() {
	for !poller.done {
		time.Sleep(2000 * time.Millisecond)
		poller.pollCheck()
	}

	poller.channel <- PollResp{poller.done, poller.token, nil}
}

func GenPoller(id string, c chan PollResp) *KramerPoller {
	uri, _ := url.Parse("https://heytherekramer.com/poll")
	q := uri.Query()
	q.Set("id", id)
	uri.RawQuery = q.Encode()

	fmt.Println(uri)
	return &KramerPoller{*uri, false, c, ""}
}
