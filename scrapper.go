package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"github.com/gocolly/colly"
)

type PlayerStats struct {
	// name string
	gamesPlayed string
	points string
	fieldGoalsMade string
	fieldGoalAttemps string
	fieldGoalPct string
	threePointsMade string
	threePointPct string
	freeThrowsMade string
	freeThrowPct string
	offensiveRebounds string
	defensiveRebounds string
	assists string
	turnovers string
	steals string
	blocks string
	personalFouls string
	fantasyPoints string
	doubleDoubles string
	tripleDoubles string
	plusMinus string
}

func main() {

	c := colly.NewCollector(
		colly.AllowedDomains("www.nba.com"),
	)

	var players []PlayerStats

	// c.Visit("https://www.nba.com/players")
	c.Visit("https://www.nba.com/stats/player/1630163")
	
	c.OnHTML("a.Anchor_anchor__cSc3P.RosterRow_playerLink__qw1vG", func(e *colly.HTMLElement) {
		playerLink := e.Attr("href")
		fullUrl := "www.nba.com" + playerLink
		fmt.Println("Player Link: ", fullUrl)

		e.Request.Visit(fullUrl)
	})

	c.OnHTML("tbody.Crom_body__UYOcU", func(e *colly.HTMLElement) {
		fmt.Println("hi")
		if e.Index == 0 {
			e.ForEach("tr", func(_ int, row *colly.HTMLElement) {
				player := PlayerStats {
					gamesPlayed: row.ChildText("td:nth-child(1)"),
					points: row.ChildText("td:nth-child(2)"),
					fieldGoalsMade: row.ChildText("td:nth-child(3)"),
					fieldGoalAttemps: row.ChildText("td:nth-child(4)"),
					fieldGoalPct: row.ChildText("td:nth-child(5)"),
					threePointsMade: row.ChildText("td:nth-child(6)"),
					threePointPct: row.ChildText("td:nth-child(7)"),
					freeThrowsMade: row.ChildText("td:nth-child(8)"),
					freeThrowPct: row.ChildText("td:nth-child(9)"),
					offensiveRebounds: row.ChildText("td:nth-child(10)"),
					defensiveRebounds: row.ChildText("td:nth-child(11)"),
					assists: row.ChildText("td:nth-child(12)"),
					turnovers: row.ChildText("td:nth-child(13)"),
					steals: row.ChildText("td:nth-child(14)"),
					blocks: row.ChildText("td:nth-child(15)"),
					personalFouls: row.ChildText("td:nth-child(16)"),
					fantasyPoints: row.ChildText("td:nth-child(17)"),
					doubleDoubles: row.ChildText("td:nth-child(18)"),
					tripleDoubles: row.ChildText("td:nth-child(19)"),
					plusMinus: row.ChildText("td:nth-child(20)"),
				}
				players = append(players, player)
			})
		}
	})

	// c.OnHTML("div.PlayerSummary_mainInnerBio__JQkoj", func(e *colly.HTMLElement) {
	// 	for i := range players {
	// 		players[i].name = e.ChildText("p.PlayerSummary_playerNameText___MhqC:nth-child(1)") + " " +
	// 						  e.ChildText("p.PlayerSummary_playerNameText___MhqC:nth-child(1)")
	// 	}
	// })

	// c.OnHTML("button.Pagination_button__sqGoH", func(e *colly.HTMLElement) {
	// 	nextPage := e.Request
	// 	if nextPage != nil {
	// 		e.Request.Visit("https://nba.com")
	// 	}
	// })

	c.OnScraped(func(r *colly.Response) {
		file, err := os.Create("products.csv")

		if err != nil {
			log.Fatalln("Faled to create CSV file", err)
			defer file.Close()
		}

		writer := csv.NewWriter(file)

		headers := []string {
			"Name", "Games Played", "Points", "FGM", "FGA", "FG%", "3PM", "3P%", "FTM",
			"FT%", "OR", "DR", "AST", "TO", "STL", "BLK", "PF", "FP", "DD", "TD", "+/-",
		}

		writer.Write(headers)

		for _, player := range players {
			record := []string {
				// player.name,
				player.gamesPlayed,
				player.points,
				player.fieldGoalsMade,
				player.fieldGoalAttemps,
				player.fieldGoalPct,
				player.threePointsMade,
				player.threePointPct,
				player.freeThrowsMade,
				player.freeThrowPct,
				player.offensiveRebounds,
				player.defensiveRebounds,
				player.assists,
				player.turnovers,
				player.steals,
				player.blocks,
				player.personalFouls,
				player.fantasyPoints,
				player.doubleDoubles,
				player.tripleDoubles,
				player.plusMinus,
			}
			writer.Write(record)
		}
		defer writer.Flush()
	})	

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request failed with status: ", r, "Error: ", err)
	})

}
