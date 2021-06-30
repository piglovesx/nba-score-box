package main

import (
	"time"

	"github.com/gosuri/uilive"
	"github.com/piglovesx/nba-score-box/game"
	"github.com/piglovesx/nba-score-box/play"
)

func main() {
	readable_data := &game.DailyGame{}
	pbps := []*play.Play{}
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()
	print_data(readable_data, writer, pbps)
	c := time.Tick(5 * time.Second)
	for range c {
		print_data(readable_data, writer, pbps)
	}
}

func print_data(d *game.DailyGame, writer *uilive.Writer, pbps []*play.Play) {
	d.RetriveDailyData()
	for i, v := range d.Payload.All_date.Games {
		if len(pbps) <= i {
			pbps = append(pbps, &play.Play{})
			pbps[i].RetrivePlayByPlay(v.Profile.GameId, v.Boxscore.Period)
		}
		v.Print(writer)
		pbps[i].Print(writer)
	}
}
