package main

import (
	"time"

	"github.com/gosuri/uilive"
	"github.com/piglovesx/nba-score-box/game"
)

func main() {
	dailyGame := &game.DailyGame{T: time.NewTicker(5 * time.Second)}
	// plays := []*play.Play{}
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	dailyGame.Print()
	t := time.NewTicker(2 * time.Hour)
	for range t.C {
		return
	}
	// print_data(dailyGame, writer, plays)
	// c := time.Tick(5 * time.Second)
	// for range c {
	// 	print_data(dailyGame, writer, plays)
	// }
	// s := ""
	// for ; ; fmt.Scanln(&s) {
	// 	if s == "q" {
	// 		return
	// 	}
	// }
}

// func print_data(d *game.DailyGame, writer *uilive.Writer, pbps []*play.Play) {
// 	d.RetriveDailyData()
// 	for i, v := range d.Payload.All_date.Games {
// 		if len(pbps) <= i {
// 			pbps = append(pbps, &play.Play{})
// 			pbps[i].RetrivePlayByPlay(v.Profile.GameId, v.Boxscore.Period)
// 		}
// 		v.Print(writer)
// 		pbps[i].Print(writer)
// 	}
// }
