package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gosuri/uilive"
)

const nba_score_box_url string = "https://in.global.nba.com/stats2/scores/daily.json?countryCode=IN&locale=en&tz=%2B8"

type daily struct {
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
	Boxscore   boxscore `json:"boxscore"`
	HomeTeam   team     `json:"homeTeam"`
	AwayTeam   team     `json:"awayTeam"`
	SeriesText string   `json:"seriesText"`
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

func main() {
	readable_data := &daily{}
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()
	print_data(readable_data, writer)
	c := time.Tick(15 * time.Second)
	for range c {
		print_data(readable_data, writer)
	}
}

func print_data(readable_data *daily, writer *uilive.Writer) {
	retriveDailyData(readable_data)
	for _, v := range readable_data.Payload.All_date.Games {
		if v.SeriesText != "" {
			fmt.Fprintf(writer, "%s\n", v.SeriesText)
		}
		fmt.Fprintf(writer, "%s %d    -    %s %d  (Period:%s Time:%s)\n\n", v.AwayTeam.Profile.Abbr, v.Boxscore.AwayScore, v.HomeTeam.Profile.Abbr, v.Boxscore.HomeScore, v.Boxscore.Period, v.Boxscore.PeriodClock)
	}
}

func retriveDailyData(d *daily) {
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
