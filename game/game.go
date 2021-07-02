package game

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/piglovesx/nba-score-box/play"
)

const nba_score_box_url string = "https://in.global.nba.com/stats2/scores/daily.json?countryCode=IN&locale=en&tz=%2B8"

type DailyGame struct {
	Payload playload `json:"payload"`
	T       *time.Ticker
}

func (d *DailyGame) Print() {
	go func() {
		for range d.T.C {
			d.RetriveDailyData()
			for _, g := range d.Payload.All_date.Games {
				g.Print()
			}
		}
	}()
}

type playload struct {
	All_date all_date `json:"date"`
}

type all_date struct {
	Games     []*game `json:"games"`
	GameCount string  `json:"gameCount"`
}

type game struct {
	Boxscore   boxscore     `json:"boxscore"`
	HomeTeam   team         `json:"homeTeam"`
	AwayTeam   team         `json:"awayTeam"`
	SeriesText string       `json:"seriesText"`
	Profile    game_profile `json:"profile"`
	Plays      *play.Play
}

type game_profile struct {
	GameId    string `json:"gameId"`
	StartTime string `json:"utcMillis"`
}

type boxscore struct {
	AwayScore   int    `json:"awayScore"`
	HomeScore   int    `json:"homeScore"`
	Period      string `json:"period"`
	PeriodClock string `json:"periodClock"`
}

type team struct {
	Profile profile `json:"profile"`
}

type profile struct {
	Abbr string `json:"abbr"`
}

func (g *game) Print() {
	if g.isStart() && !g.Plays.Started {
		g.Plays.Print(g.Profile.GameId, g.Boxscore.Period)
	}
}

func (g *game) String() string {
	result := ""
	if g.SeriesText != "" {
		temp, _ := strconv.Atoi(g.Profile.StartTime)
		result = result + fmt.Sprintf("%s (%s)\n", g.SeriesText, time.Unix(0, int64(temp)*int64(time.Millisecond)).Format(time.Kitchen))
	}
	result = result + fmt.Sprintf("%s %d    -    %s %d  (Period:%s Time:%s)\n", g.AwayTeam.Profile.Abbr, g.Boxscore.AwayScore, g.HomeTeam.Profile.Abbr, g.Boxscore.HomeScore, g.Boxscore.Period, g.Boxscore.PeriodClock)
	result = result + g.Plays.String()
	return result
}

func (g game) isStart() bool {
	t, _ := strconv.Atoi(g.Profile.StartTime)
	return time.Now().Unix()*1000 > int64(t)
}

func (d *DailyGame) RetriveDailyData() {
	resp, err := http.Get(nba_score_box_url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	resp_body, err := io.ReadAll(resp.Body)
	// resp_body, err := dataSource()
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(resp_body, d)

	if d.Payload.All_date.GameCount == "0" {
		log.Fatal("No Game Now")
	}

	if err != nil {
		log.Fatal(err)
	}
}
