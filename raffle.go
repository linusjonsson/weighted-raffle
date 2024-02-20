package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Participant represents a participant in the raffle with their name and wins
type Participant struct {
	Name string
	Wins int
}

// Item represents an item in the raffle with its name, value, and participants
type Item struct {
	Name         string
	Value        int
	Participants map[string]int // Map of participants and their wins for the current item
}

// WeightedRaffle picks a winner from the provided items based on their values and participants
func WeightedRaffle(items []Item, participants map[string]*Participant, winners map[string]string) error {
	// Sort items by value in ascending order
	sortItemsByValue(items)

	// Random seed initialization
	rand.Seed(time.Now().UnixNano())

	// Iterate over items and pick a winner for each item
	for _, item := range items {
		if len(item.Participants) == 0 {
			continue // Skip items with no participants
		}

		// Display users eligible to win and their current chance of winning for the current item
		fmt.Printf("Item %s with value %d\n", item.Name, item.Value)
		displayEligibleParticipants(item.Participants, participants)

		// Prompt user to confirm before proceeding
		fmt.Println("Press Enter to reveal the winner...")
		fmt.Scanln()

		// Pick a winner based on their chance of winning
		winner := pickWinner(item.Participants)

		// Update wins for the winner
		participants[winner].Wins++

		// Print the winner for the current item
		fmt.Printf("Winner for item %s with value %d: %s\n", item.Name, item.Value, winner)
		winners[item.Name] = winner // Store the winner for the current item
		fmt.Println("-----------------------------")

		// Ask the winner if they want to remain in the following drawings
		fmt.Printf("Would you like to remain in the following drawings, %s? (yes/no): ", winner)
		var response string
		fmt.Scanln(&response)

		// If the winner chooses to opt out, remove them from future items and the participants map
		if strings.ToLower(response) != "yes" {
			for _, futureItem := range items {
				delete(futureItem.Participants, winner)
			}
			delete(participants, winner)
			fmt.Printf("%s has opted out and will be removed from future drawings.\n", winner)
		}
	}

	return nil
}

// sortItemsByValue sorts items by value in ascending order
func sortItemsByValue(items []Item) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Value < items[j].Value
	})
}

// pickWinner picks a winner based on their chance of winning
func pickWinner(participants map[string]int) string {
	totalWins := calculateTotalWins(participants)
	var totalChance float64

	// Calculate total chance considering previous wins
	for _, wins := range participants {
		totalChance += 1.0/1 + totalWins - float64(wins)
	}

	// Generate a random number between 0 and totalChance
	randomNumber := rand.Float64() * totalChance

	// Iterate over participants to determine the winner
	for participant, wins := range participants {
		chance := 1.0/1 + totalWins - float64(wins)
		if randomNumber <= chance {
			return participant
		}
		randomNumber -= chance
	}

	return "" // Should never reach here
}

// calculateTotalWins calculates the total wins across all participants
func calculateTotalWins(participants map[string]int) float64 {
	totalWins := 0.0
	for _, wins := range participants {
		totalWins += float64(wins)
	}
	return totalWins
}

// displayEligibleParticipants displays the eligible participants for the current item and their previous wins
func displayEligibleParticipants(participants map[string]int, allParticipants map[string]*Participant) {
	fmt.Println("Eligible Participants:")
	for participant := range participants {
		totalWins := allParticipants[participant].Wins
		fmt.Printf("- Participant: %s, Previous Wins: %d\n", participant, totalWins)
	}
}

// WriteResultsToCSV writes the results to a CSV file
func WriteResultsToCSV(items []Item, winners map[string]string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"Item", "Winner"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write item and winner for each entry
	for _, item := range items {
		winner, ok := winners[item.Name]
		if !ok {
			winner = "No winner"
		}
		entry := []string{item.Name, winner}
		if err := writer.Write(entry); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// Open CSV file
	file, err := os.Open("raffle_data.csv")
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow varying number of fields per row

	// Read all CSV records
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV records:", err)
		return
	}

	// Map to store items
	items := make([]Item, 0)

	// Map to store participants and their total wins across all items
	participants := make(map[string]*Participant)

	// Iterate over CSV records
	for _, record := range records {
		if len(record) < 3 {
			continue // Skip records with less than three fields
		}

		// Parse item value
		value, err := strconv.Atoi(record[1])
		if err != nil {
			fmt.Println("Error parsing item value:", err)
			continue
		}

		// Parse participants for the current item
		participantNames := record[2:]
		participantsForItem := make(map[string]int)
		for _, participantName := range participantNames {
			participantName = strings.TrimSpace(participantName)
			participantsForItem[participantName] = 0 // Initialize wins for the participant for the current item
			// Update total wins for the participant across all items
			if _, ok := participants[participantName]; !ok {
				participants[participantName] = &Participant{Name: participantName, Wins: 0}
			}
		}

		// Create new item with participants
		item := Item{
			Name:         record[0],
			Value:        value,
			Participants: participantsForItem,
		}

		// Append item to items slice
		items = append(items, item)
	}

	// Display total list of participants
	fmt.Println("Total List of Participants:")
	for participant := range participants {
		fmt.Printf("- Participant: %s\n", participant)
	}

	fmt.Println("Press Enter to start the raffle:")
	fmt.Scanln()

	// Run weighted raffle with the parsed items and participants
	winners := make(map[string]string)
	err = WeightedRaffle(items, participants, winners)
	if err != nil {
		fmt.Println("Error running weighted raffle:", err)
		return
	}

	// Write results to CSV file
	err = WriteResultsToCSV(items, winners, "raffle_results.csv")
	if err != nil {
		fmt.Println("Error writing results to CSV:", err)
		return
	}
}
