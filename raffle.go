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
	Wins map[string]bool // Map of item names indicating if the participant has won the item
}

// Item represents an item in the raffle with its name, value, and participants
type Item struct {
	Name             string
	Value            int
	ParticipantNames []string // Array of participant names for the current item
}

// WeightedRaffle picks a winner from the provided items based on their values and participants
func WeightedRaffle(items []Item, participants map[string]*Participant, winners map[string]string) error {
	// Sort items by value in ascending order
	sortItemsByValue(items)

	// Random seed initialization
	rand.Seed(time.Now().UnixNano())

	// Iterate over items and pick a winner for each item
	for _, item := range items {
		if len(item.ParticipantNames) == 0 {
			continue // Skip items with no participants
		}

		// Display users eligible to win and their current chance of winning for the current item
		fmt.Printf("Item %s with value %d\n", item.Name, item.Value)
		displayEligibleParticipants(item.ParticipantNames, participants, item.Name)
		/*
			// Prompt user to confirm before proceeding
			fmt.Println("Press Enter to reveal the winner...")
			fmt.Scanln()
		*/
		// Pick a winner based on their chance of winning
		winner := pickWinner(item.ParticipantNames, participants, item.Name)

		// Update wins for the winner
		participants[winner].Wins[item.Name] = true

		// Print the winner for the current item
		fmt.Printf("Winner for item %s with value %d: %s\n", item.Name, item.Value, winner)
		winners[item.Name] = winner // Store the winner for the current item
		fmt.Println("-----------------------------")
		/*
			// Ask the winner if they want to remain in the following drawings
			fmt.Printf("Would you like to remain in the following drawings, %s? (yes/no): ", winner)
			var response string
			fmt.Scanln(&response)

			// If the winner chooses to opt out, remove them from future items and the participants map
			if strings.ToLower(response) != "yes" {
				for _, participantName := range item.ParticipantNames {
					delete(participants[participantName].Wins, item.Name)
				}
				fmt.Printf("%s has opted out and will be removed from future drawings.\n", winner)
			}
		*/
	}

	return nil
}

// sortItemsByValue sorts items by value in ascending order
func sortItemsByValue(items []Item) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Value < items[j].Value
	})
}

// pickWinner picks a winner from the provided participant names based on their chance of winning
func pickWinner(participantNames []string, participants map[string]*Participant, itemName string) string {
	var totalTickets float64
	tickets := make(map[string]float64)

	// Calculate total tickets and assign tickets to each participant
	for _, name := range participantNames {
		// Calculate the number of tickets for the participant
		ticket := 1.0 / float64(1+len(participants[name].Wins))

		// Add the ticket to the total
		totalTickets += ticket

		// Assign the ticket to the participant
		tickets[name] = ticket
	}

	// Generate a random number between 0 and totalTickets
	randomNumber := rand.Float64() * totalTickets

	// Determine the winner based on the random number
	var cumulativeTickets float64
	for name, ticket := range tickets {
		// Add the current participant's ticket to the cumulativeTickets
		cumulativeTickets += ticket

		// If the random number falls within the range of the current participant's tickets
		// (from cumulativeTickets - ticket to cumulativeTickets), return the participant as the winner
		if randomNumber <= cumulativeTickets {
			return name
		}
	}

	// Return an empty string if no winner is found (should never reach here)
	return ""
}

// displayEligibleParticipants displays the eligible participants for the current item and their total previous wins
func displayEligibleParticipants(participantNames []string, participants map[string]*Participant, itemName string) {
	fmt.Println("Eligible Participants:")
	for _, name := range participantNames {
		totalWins := 0
		// Count the total wins for the current participant based on the Wins map for the specified item
		for _, hasWon := range participants[name].Wins {
			if hasWon {
				totalWins++
			}
		}
		fmt.Printf("- Participant: %s, Total Previous Wins: %d\n", name, totalWins)
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
		participantsForItem := make([]string, len(participantNames))
		for i, participantName := range participantNames {
			participantName = strings.TrimSpace(participantName)
			participantsForItem[i] = participantName
			// Update total wins for the participant across all items
			if _, ok := participants[participantName]; !ok {
				participants[participantName] = &Participant{Name: participantName, Wins: make(map[string]bool)}
			}
		}

		// Create new item with participants
		item := Item{
			Name:             record[0],
			Value:            value,
			ParticipantNames: participantsForItem,
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
