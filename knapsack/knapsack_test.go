package knapsack

import (
	"fmt"
	"testing"
)

type ts struct {
	capacity int
	weights []int
	values []int
	expMaxProfit int
	pickedItems []int
}

var tt = []ts{
	{capacity: 20, weights: []int{5, 7}, values: []int{20, 30}, expMaxProfit: 50, pickedItems: []int{0, 1}},
	{capacity: 5, weights: []int{5, 7}, values: []int{20, 30}, expMaxProfit: 20, pickedItems: []int{0}},
	{capacity: 8, weights: []int{5, 7}, values: []int{20, 30}, expMaxProfit: 30, pickedItems: []int{1}},
	{capacity: 10, weights: []int{9, 2, 1}, values: []int{20, 30, 100}, expMaxProfit: 130, pickedItems: []int{1, 2}},
}

func TestMaxProfit(t *testing.T) {
	for i, tc := range tt {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			profit, err := MaxProfit(tc.capacity, tc.weights, tc.values)
			if err != nil {
				t.Fatalf("unexpected error %+v", profit)
			}

			if profit != tc.expMaxProfit {
				t.Fatalf("expected %d, got %d", tc.expMaxProfit, profit)
			}
		})
	}
}

func TestBestItems(t *testing.T) {
	for i, tc := range tt {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			items, err := BestItems(tc.capacity, tc.weights, tc.values)
			if err != nil {
				t.Fatalf("unexpected error %+v", err)
			}

			if len(items) != len(tc.pickedItems) {
				t.Fatalf("expected %d items to be picked, got %d", tc.pickedItems, len(items))
			}

			for idx := range tc.pickedItems {
				if items[idx] != tc.pickedItems[idx] {
					t.Fatalf("expected %d at index %d, got %d", tc.pickedItems[idx], idx, items[idx])
				}
			}
		})
	}
}

func TestMaxProfitAndBestItems(t *testing.T) {
	for i, tc := range tt {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			profit, items, err := MaxProfitAndBestItems(tc.capacity, tc.weights, tc.values)
			if err != nil {
				t.Fatalf("unexpected error %+v", profit)
			}

			if profit != tc.expMaxProfit {
				t.Fatalf("expected %d, got %d", tc.expMaxProfit, profit)
			}

			if len(items) != len(tc.pickedItems) {
				t.Fatalf("expected %d items to be picked, got %d", tc.pickedItems, len(items))
			}

			for idx := range tc.pickedItems {
				if items[idx] != tc.pickedItems[idx] {
					t.Fatalf("expected %d at index %d, got %d", tc.pickedItems[idx], idx, items[idx])
				}
			}
		})
	}
}
