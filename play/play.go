package play

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const nba_play_by_play_url string = "https://in.global.nba.com/stats2/game/playbyplay.json?"

type Play struct {
	Payload playload_for_pbp `json:"payload"`
	Ticker  *time.Ticker
	Started bool
	Msg     chan int
	mux     *sync.Mutex
}

type playload_for_pbp struct {
	PlayByPlays []play_by_play_data `json:"playByPlays"`
}

type play_by_play_data struct {
	Events []Event `json:"events"`
}

type Event struct {
	AwayScore   string `json:"awayScore"`
	HomeScore   string `json:"homeScore"`
	Description string `json:"description"`
}

func (p *Play) Print(gameId string, period string) {
	if p.Started {
		return
	}
	p.mux.Lock()
	defer p.mux.Unlock()

	go func() {
		for range p.Ticker.C {
			p.RetrivePlayByPlay(gameId, period)
			p.Msg <- 0

		}
	}()
	p.Started = true
}

func (p *Play) String() string {
	result := ""
	for i, v := range p.Payload.PlayByPlays[0].Events {
		if i <= 2 {
			result = result + fmt.Sprintf("%s\n", v.Description)
		}
	}
	return result
}

func (pbp *Play) RetrivePlayByPlay(gameId string, period string) {
	params := url.Values{
		"gameId": {gameId},
		"period": {period},
		"locale": {"en"},
	}
	play_url := nba_play_by_play_url + params.Encode()

	resp, err := http.Get(play_url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	resp_body, err := io.ReadAll(resp.Body)
	// resp_body, err := dataSource()
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(resp_body, pbp)

	if err != nil {
		log.Fatal(err)
	}
}
