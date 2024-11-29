package vis

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func RunVisualization() {
	// Load the CSV file
	csvFile := "players_stats.csv"
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
	}

	// Prepare data
	var labels []string
	var ppgValues, rpgValues, apgValues plotter.Values

	for i, record := range records {
		if i == 0 {
			// Skip the header row
			continue
		}

		// Combine first name and last name for the label
		playerName := record[0] + " " + record[1]
		labels = append(labels, playerName)

		// Parse PPG
		ppg, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Printf("Invalid PPG value for %s: %v", playerName, err)
			ppg = 0
		}
		ppgValues = append(ppgValues, ppg)

		// Parse RPG
		rpg, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Printf("Invalid RPG value for %s: %v", playerName, err)
			rpg = 0
		}
		rpgValues = append(rpgValues, rpg)

		// Parse APG
		apg, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			log.Printf("Invalid APG value for %s: %v", playerName, err)
			apg = 0
		}
		apgValues = append(apgValues, apg)
	}

	// Create and save the bar charts
	createBarChart(labels, ppgValues, "Points Per Game (PPG)", "ppg_chart.png")
	createBarChart(labels, rpgValues, "Rebounds Per Game (RPG)", "rpg_chart.png")
	createBarChart(labels, apgValues, "Assists Per Game (APG)", "apg_chart.png")
}

// createBarChart generates and saves a bar chart
func createBarChart(labels []string, values plotter.Values, title, filename string) {
	// Create a new plot
	p := plot.New()
	p.Title.Text = title
	p.Y.Label.Text = "Value"
	p.X.Label.Text = "Players"

	// Create the bar chart
	bars, err := plotter.NewBarChart(values, vg.Points(20))
	if err != nil {
		log.Fatalf("Failed to create bar chart: %v", err)
	}
	bars.Color = plotutil.Color(1) // Assign a color for the bars

	// Set X-axis labels
	p.Add(bars)
	p.NominalX(labels...)

	// Save the plot to a file
	if err := p.Save(12*vg.Inch, 8*vg.Inch, filename); err != nil {
		log.Fatalf("Failed to save plot: %v", err)
	}
	log.Printf("%s chart saved as %s\n", title, filename)
}
