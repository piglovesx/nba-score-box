package game

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

const nba_score_box_url string = "https://in.global.nba.com/stats2/scores/daily.json?countryCode=IN&locale=en&tz=%2B8"

type DailyGame struct {
	Payload playload `json:"payload"`
}

type playload struct {
	All_date all_date `json:"date"`
}

type all_date struct {
	Games     []game `json:"games"`
	GameCount string `json:"gameCount"`
}

type game struct {
	Boxscore   boxscore     `json:"boxscore"`
	HomeTeam   team         `json:"homeTeam"`
	AwayTeam   team         `json:"awayTeam"`
	SeriesText string       `json:"seriesText"`
	Profile    game_profile `json:"profile"`
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

func (v game) Print(writer io.Writer) {
	if v.SeriesText != "" {
		temp, _ := strconv.Atoi(v.Profile.StartTime)
		fmt.Fprintf(writer, "%s (%s)\n", v.SeriesText, time.Unix(0, int64(temp)*int64(time.Millisecond)).Format(time.Kitchen))
	}
	fmt.Fprintf(writer, "%s %d    -    %s %d  (Period:%s Time:%s)\n", v.AwayTeam.Profile.Abbr, v.Boxscore.AwayScore, v.HomeTeam.Profile.Abbr, v.Boxscore.HomeScore, v.Boxscore.Period, v.Boxscore.PeriodClock)

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
