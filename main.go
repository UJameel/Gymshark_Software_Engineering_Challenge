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
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	packSize, err := strconv.Atoi(r.FormValue("packSize"))
	if err != nil {
		http.Error(w, "Invalid pack size", http.StatusBadRequest)
		return
	}

	packSizes = append(packSizes, packSize)
	sort.Slice(packSizes, func(i, j int) bool {
		return packSizes[i] > packSizes[j]
	})

	if err := saveConfig(); err != nil {
		http.Error(w, "Error saving config", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func removePackSizeHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	packSize, err := strconv.Atoi(r.FormValue("packSize"))
	if err != nil {
		http.Error(w, "Invalid pack size", http.StatusBadRequest)
		return
	}

	for i, size := range packSizes {
		if size == packSize {
			packSizes = append(packSizes[:i], packSizes[i+1:]...)
			break
		}
	}

	if err := saveConfig(); err != nil {
		http.Error(w, "Error saving config", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
