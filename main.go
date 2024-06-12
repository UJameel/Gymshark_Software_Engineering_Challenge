package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type Config struct {
	PackSizes []int `json:"packSizes"`
}

type Order struct {
	ItemsOrdered int `json:"items_ordered"`
}

type Pack struct {
	Size  int `json:"size"`
	Count int `json:"count"`
}

type Response struct {
	Message string `json:"message"`
	Packs   []Pack `json:"packs,omitempty"`
}

var packSizes []int

func main() {
	loadConfig()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add-pack-size", addPackSizeHandler)
	http.HandleFunc("/remove-pack-size", removePackSizeHandler)
	http.HandleFunc("/calculate-packs", calculatePacksHandler)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func loadConfig() {
	file, err := os.Open("packSizeConfig.json")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var config Config
	if err := json.Unmarshal(byteValue, &config); err != nil {
		fmt.Println("Error parsing config file:", err)
		return
	}

	packSizes = config.PackSizes
	if len(packSizes) == 0 {
		fmt.Println("No pack sizes defined in config file")
		os.Exit(1)
	}
}

func saveConfig() error {
	config := Config{PackSizes: packSizes}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile("packSizeConfig.json", data, 0644)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		PackSizes []int
	}{
		PackSizes: packSizes,
	}

	tmpl.Execute(w, data)
}

func addPackSizeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request to add pack size")
	packSizeStr := r.FormValue("packSize")
	fmt.Println("Form value packSize:", packSizeStr)

	packSize, err := strconv.Atoi(packSizeStr)
	if err != nil {
		fmt.Println("Invalid pack size:", packSizeStr)
		http.Error(w, "Invalid pack size", http.StatusBadRequest)
		return
	}

	// Check if the pack size already exists
	for _, size := range packSizes {
		if size == packSize {
			response := Response{Message: "Pack size already exists"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	fmt.Println("Adding pack size:", packSize)
	packSizes = append(packSizes, packSize)
	sort.Slice(packSizes, func(i, j int) bool {
		return packSizes[i] > packSizes[j]
	})

	if err := saveConfig(); err != nil {
		fmt.Println("Error saving config:", err)
		http.Error(w, "Error saving config", http.StatusInternalServerError)
		return
	}

	response := Response{Message: "Pack size added successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func removePackSizeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request to remove pack size")
	packSizeStr := r.FormValue("packSize")
	fmt.Println("Form value packSize:", packSizeStr)

	packSize, err := strconv.Atoi(packSizeStr)
	if err != nil {
		fmt.Println("Invalid pack size:", packSizeStr)
		http.Error(w, "Invalid pack size", http.StatusBadRequest)
		return
	}

	fmt.Println("Removing pack size:", packSize)
	found := false
	for i, size := range packSizes {
		if size == packSize {
			packSizes = append(packSizes[:i], packSizes[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		response := Response{Message: "Pack size not found"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := saveConfig(); err != nil {
		fmt.Println("Error saving config:", err)
		http.Error(w, "Error saving config", http.StatusInternalServerError)
		return
	}

	response := Response{Message: "Pack size removed successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func calculatePacksHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request to calculate packs")
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	itemsOrderedStr := r.PostFormValue("itemsOrdered")
	fmt.Println("PostForm value itemsOrdered:", itemsOrderedStr)

	if itemsOrderedStr == "" {
		fmt.Println("itemsOrderedStr is empty")
		http.Error(w, "Invalid number of items", http.StatusBadRequest)
		return
	}

	var order Order
	order.ItemsOrdered, err = strconv.Atoi(itemsOrderedStr)
	if err != nil {
		fmt.Println("Error converting itemsOrderedStr to integer:", err)
		http.Error(w, "Invalid number of items", http.StatusBadRequest)
		return
	}

	fmt.Println("Calculating packs for:", order.ItemsOrdered)
	packs := calculatePacks(order.ItemsOrdered)
	response := Response{
		Message: "Calculation successful",
		Packs:   packs,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("Error encoding response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
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

	// Optimize the packs
	return optimizePacks(packs)
}

func optimizePacks(packs []Pack) []Pack {
	totalItems := 0
	for _, pack := range packs {
		totalItems += pack.Size * pack.Count
	}

	fmt.Println("Total items for optimization:", totalItems)

	// Recalculate the packs to see if we can optimize without recursion
	optPacks := []Pack{}
	remainingItems := totalItems

	for _, size := range packSizes {
		if remainingItems == 0 {
			break
		}
		if remainingItems >= size {
			count := remainingItems / size
			optPacks = append(optPacks, Pack{Size: size, Count: count})
			remainingItems %= size
		}
	}

	// If there are any remaining items, use the smallest pack to cover the remainder
	if remainingItems > 0 {
		smallestPack := packSizes[len(packSizes)-1]
		optPacks = append(optPacks, Pack{Size: smallestPack, Count: 1})
	}

	// Return the more optimal solution
	if len(optPacks) < len(packs) {
		return optPacks
	}
	return packs
}
