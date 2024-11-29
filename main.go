package main

import (
	"GoScrapper/scraper"
	"GoScrapper/vis"
	"fmt"
)

func main() {
	fmt.Println("Running scraper...")
	scraper.RunScraper()
	fmt.Println("Running data visualization...")
	vis.RunVisualization()
}
