This code represents a functional implementation of a weighted raffle system in Go, handling participant information, item values, weighted selection of winners, and result storage in CSV format.
The strucutre of the input csv is:
`item,value,participant(s)`

## Data Structures:
- `Participant`: Represents a participant in the raffle with their name and a map of item names indicating if the participant has won the item.
- `Item`: Represents an item in the raffle with its name, value, and an array of participant names for the current item.

## Functions:
- `WeightedRaffle`: Conducts the weighted raffle by picking a winner for each item based on their values and participants. It iterates over the items, displays eligible participants, prompts for confirmation, picks a winner, updates wins for the winner, and handles participant opt-out.
- `sortItemsByValue`: Sorts items by value in ascending order.
- `pickWinner`: Picks a winner from the provided participant names based on their chance of winning.
- `displayEligibleParticipants`: Displays the eligible participants for the current item and their total previous wins.
- `WriteResultsToCSV`: Writes the results (item and winner) to a CSV file.

## Main Functionality:
- Reads data from a CSV file (`raffle_data.csv`), where each row represents an item in the raffle with its value and participants.
- Parses the CSV records, creating items and updating participant information.
- Displays the total list of participants.
- Prompts the user to start the raffle.
- Conducts the weighted raffle using the parsed items and participants.
- Writes the raffle results to a new CSV file (`raffle_results.csv`).

## Execution:
- The program is executed from the `main` function, where it reads the CSV data, conducts the raffle, and writes the results to a file.
