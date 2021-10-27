package knapsack

import (
	"github.com/pkg/errors"
)

var ErrInvalidArguments = errors.New("invalid arguments")

func MaxProfit(capacity int, weights []int, values []int) (int, error) {
	profit, _, err := calculate(capacity, weights, values)
	if err != nil {
		return 0, err
	}

	return profit, nil
}

func BestItems(capacity int, weights []int, values []int) ([]int, error) {
	_, items, err := calculate(capacity, weights, values)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func MaxProfitAndBestItems(capacity int, weights []int, values []int) (int, []int, error) {
	return calculate(capacity, weights, values)
}

func calculate(capacity int, weights []int, values []int) (int, []int, error) {
	if capacity < 0 || weights == nil || values == nil || len(weights) != len(values) {
		return 0, nil, ErrInvalidArguments
	}

	rows := len(weights)

	table := makeTable(rows + 1, capacity + 1)

	for row := 1; row <= rows; row++ {
		// get the value and the weight of the item
		value, weight := values[row - 1], weights[row - 1]
		for col := 1; col <= capacity; col++ {
			table[row][col] = table[row - 1][col]
			if col >= weight && (table[row - 1][col - weight] + value) > table[row][col] {
				table[row][col] = table[row - 1][col - weight] + value
			}
		}
	}

	col := capacity
	itemsSelected := make([]int, 0)

	// now we backtrack from the bottom-right element of the table to collect
	// the items that were put into the knapsack. The criteria - table[row][col] != table[row - 1][col]
	for row := rows; row > 0; row-- {
		if table[row][col] != table[row - 1][col] {
			itemIdx := row - 1
			itemsSelected = append(itemsSelected, itemIdx)
			col -= weights[itemIdx]
		}
	}

	for i, j := 0, len(itemsSelected)-1; i < j; i, j = i+1, j-1 {
		itemsSelected[i], itemsSelected[j] = itemsSelected[j], itemsSelected[i]
	}

	return table[rows][capacity], itemsSelected, nil
}

func makeTable(rows, cols int) (table [][]int) {
	table = make([][]int, rows)
	for i := range table {
		table[i] = make([]int, cols)
	}
	return
}
