package kramer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type GhResp struct {
	Login string `json:login`
	Id    int    `json:id`
}

type PollResp struct {
	done        bool
	accessToken string
	err         error
}

var kramerUri = "https://heytherekramer.com"

type kramerInfo struct {
	Uri string
	Id  string
}

type Kramer struct {
	Win        fyne.Window
	Buttons    map[string]*widget.Button
	Tapped     bool
	App        fyne.App
	data       kramerInfo
	poller     *KramerPoller
	pollerChan chan PollResp
	token      string
}

func buildButton(k *Kramer) *widget.Button {
	b := widget.NewButton("Sign up via Github", func() {
		uri, _ := url.Parse(k.data.Uri)

		go k.poller.Poll()
		k.Win.SetContent(widget.NewProgressBarInfinite())
		k.App.OpenURL(uri)
		resp := <-k.pollerChan
		if resp.err != nil {
			println("oh noez")
		}

		k.token = resp.accessToken
		ghUri, _ := url.Parse("https://api.github.com/user")
		client := &http.Client{}

		req, _ := http.NewRequest("GET", ghUri.String(), nil)
		req.Header.Add("Authorization", fmt.Sprintf("token %s", k.token))
		clientResp, err := client.Do(req)
		if err != nil {
			println(err)
			panic(err)
		}
		defer clientResp.Body.Close()
		body, error := ioutil.ReadAll(clientResp.Body)
		if error != nil {
			println(error)
			panic("oops again")
		}

		var ghData GhResp
		json.Unmarshal([]byte(body), &ghData)

		if resp.done {
			k.Win.SetContent(widget.NewLabel(fmt.Sprintf("it's: %s", ghData.Login)))
		}
	})

	k.Tapped = true

	return b
}

func NewKramer(a fyne.App) *Kramer {
	w := a.NewWindow("Kramer Pager")
	// for some reason constructor does not properly save title
	w.SetTitle("Kramer Pager")
	fmt.Println(w.Title())

	w.Resize(fyne.NewSize(300, 300))

	buttons := make(map[string]*widget.Button)

	k := &Kramer{
		Win:     w,
		Buttons: buttons,
		App:     a,
	}
	buttons["signup"] = buildButton(k)
	w.SetContent(widget.NewVBox(
		buttons["signup"],
	))

	return k
}

func (kram *Kramer) fetchToken() {
	resp, err := http.Get(kramerUri)
	if err != nil {
		println(err)
		panic("oops")
	}
	defer resp.Body.Close()
	body, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		println(err)
		panic("oops again")
	}

	var data kramerInfo
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	kram.data = data

	c := make(chan PollResp)
	poller := GenPoller(kram.data.Id, c)
	kram.poller = poller
	kram.pollerChan = c
}

func (kram *Kramer) Run() {
	go kram.fetchToken()

	kram.Win.ShowAndRun()
}
