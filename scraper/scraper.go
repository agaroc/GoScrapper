package scraper

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Player struct {
	FirstName       string
	LastName        string
	Team            string
	PointsPerGame   string
	ReboundsPerGame string
	AssistsPerGame  string
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

func RunScraper() {

	var playerURLs []string
	fmt.Println("Enter NBA player profile URLS from NBA websiter (1 per line). Enter 'done' when ready: ")

	// read user input for the URLS
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter URL: ")
		scanner.Scan()
		url := strings.TrimSpace(scanner.Text())
		if url == "done" {
			break
		}
		if url != "" {
			playerURLs = append(playerURLs, url)
		}
	}
	//check if nothing was entered
	if len(playerURLs) == 0 {
		fmt.Println("No URLs provided. Exiting.")
		return
	}

	players := []Player{}

	c := colly.NewCollector()

	// temporary player instance
	var currentPlayer Player

	// extract player team, first name, and last name
	c.OnHTML("div.PlayerSummary_mainInnerBio__JQkoj", func(e *colly.HTMLElement) {
		teamText := strings.TrimSpace(e.ChildText("p.PlayerSummary_mainInnerInfo__jv3LO"))

		teamParts := strings.SplitN(teamText, "|", 2)
		if len(teamParts) > 0 {
			teamName := strings.SplitN(teamParts[0], "#", 2)[0]
			currentPlayer.Team = strings.TrimSpace(teamName)
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

	

	// Visit player URLs provided by the user
	for _, url := range playerURLs {
		fmt.Printf("Scraping %s...\n", url)
		c.Visit(url)
	}

	// save players to CSV
	writeToCSV(players, "players_stats.csv")
	fmt.Println("Player data successfully written to players_stats.csv")
}
