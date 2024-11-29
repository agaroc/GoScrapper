package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Player struct {
	FirstName      string
	LastName       string
	Team           string
	PointsPerGame  string
	ReboundsPerGame string
	AssistsPerGame string
}

// function to write players data to CSV
func writeToCSV(players []Player, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write CSV header
	writer.Write([]string{"First Name", "Last Name", "Team", "PPG", "RPG", "APG"})

	// write player data
	for _, player := range players {
		writer.Write([]string{
			player.FirstName,
			player.LastName,
			player.Team,
			player.PointsPerGame,
			player.ReboundsPerGame,
			player.AssistsPerGame,
		})
	}
}

func main() {
	players := []Player{}

	c := colly.NewCollector()

	// temporary player instance
	var currentPlayer Player

	// extract player team, first name, and last name
	c.OnHTML("div.PlayerSummary_mainInnerBio__JQkoj", func(e *colly.HTMLElement) {
		teamText := strings.TrimSpace(e.ChildText("p.PlayerSummary_mainInnerInfo__jv3LO"))
		teamParts := strings.Fields(teamText)
		if len(teamParts) >= 2 {
			currentPlayer.Team = teamParts[0] + " " + teamParts[1]
		} else {
			currentPlayer.Team = teamText
		}

		currentPlayer.FirstName = strings.TrimSpace(e.ChildText("p.PlayerSummary_playerNameText___MhqC:nth-of-type(2)"))
		currentPlayer.LastName = strings.TrimSpace(e.ChildText("p.PlayerSummary_playerNameText___MhqC:nth-of-type(3)"))
	})

	// extract player stats
	c.OnHTML("div.PlayerSummary_playerStat__rmEOP", func(e *colly.HTMLElement) {
		statKey := strings.TrimSpace(e.ChildText("p.PlayerSummary_playerStatLabel__I3TO3"))
		statValue := strings.TrimSpace(e.ChildText("p.PlayerSummary_playerStatValue___EDg_"))

		switch statKey {
		case "PPG":
			currentPlayer.PointsPerGame = statValue
		case "RPG":
			currentPlayer.ReboundsPerGame = statValue
		case "APG":
			currentPlayer.AssistsPerGame = statValue
		}
	})

	// once per player page, add the player to the list
	c.OnScraped(func(r *colly.Response) {
		players = append(players, currentPlayer)
		currentPlayer = Player{} // Reset for next player
	})

	// visit player URLs
	playerURLs := []string{
		"https://www.nba.com/stats/player/1630163",
		"https://www.nba.com/player/201939/stephen-curry",
		"https://www.nba.com/player/203110/draymond-green",
		"https://www.nba.com/player/1627741/buddy-hield",
		"https://www.nba.com/player/1626172/kevon-looney",
		"https://www.nba.com/player/1627780/gary-payton-ii",
		"https://www.nba.com/player/203952/andrew-wiggins",
		"https://www.nba.com/player/203937/kyle-anderson",
	}

	for _, url := range playerURLs {
		c.Visit(url)
	}

	// save players to CSV
	writeToCSV(players, "players_stats.csv")
	fmt.Println("Player data successfully written to players_stats.csv")
}
