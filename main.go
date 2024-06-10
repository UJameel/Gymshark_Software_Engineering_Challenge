package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

type Order struct {
	ItemsOrdered int `json:"items_ordered"`
}

type Pack struct {
	Size  int `json:"size"`
	Count int `json:"count"`
}

var packSizes = []int{5000, 2000, 1000, 500, 250}

func main() {
	http.HandleFunc("/calculate-packs", calculatePacksHandler)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func calculatePacksHandler(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	packs := calculatePacks(order.ItemsOrdered)
	optimizedPacks := optimizePacks(packs)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(optimizedPacks)
}

func calculatePacks(itemsOrdered int) []Pack {
	sort.Slice(packSizes, func(i, j int) bool {
		return packSizes[i] > packSizes[j]
	})

	packs := []Pack{}
	remainingItems := itemsOrdered

	for _, size := range packSizes {
		if remainingItems == 0 {
			break
		}
		if remainingItems >= size {
			count := remainingItems / size
			packs = append(packs, Pack{Size: size, Count: count})
			remainingItems %= size
		}
	}

	// If there are any remaining items, use the smallest pack to cover the remainder
	if remainingItems > 0 {
		smallestPack := packSizes[len(packSizes)-1]
		packs = append(packs, Pack{Size: smallestPack, Count: 1})
	}

	return packs
}

func optimizePacks(packs []Pack) []Pack {
	totalItems := 0
	for _, pack := range packs {
		totalItems += pack.Size * pack.Count
	}

	// Recalculate the packs to see if we can optimize
	return calculatePacks(totalItems)
}