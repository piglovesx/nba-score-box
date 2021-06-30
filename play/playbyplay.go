package play

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

const nba_play_by_play_url string = "https://in.global.nba.com/stats2/game/playbyplay.json?"

type Play struct {
	Payload playload_for_pbp `json:"payload"`
}

type playload_for_pbp struct {
	PlayByPlays []play_by_play_data `json:"playByPlays"`
}

type play_by_play_data struct {
	Events []event `json:"events"`
}

type event struct {
	AwayScore   string `json:"awayScore"`
	HomeScore   string `json:"homeScore"`
	Description string `json:"description"`
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
