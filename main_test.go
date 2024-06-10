package main

import (
	"reflect"
	"testing"
)

func TestCalculatePacks(t *testing.T) {
	testCases := []struct {
		order    int
		expected []Pack
	}{
		{7500, []Pack{{Size: 5000, Count: 1}, {Size: 2000, Count: 1}, {Size: 500, Count: 1}}},
		{1, []Pack{{Size: 250, Count: 1}}},
		{250, []Pack{{Size: 250, Count: 1}}},
		{251, []Pack{{Size: 250, Count: 1}, {Size: 250, Count: 1}}},
		{501, []Pack{{Size: 500, Count: 1}, {Size: 250, Count: 1}}},
		{12001, []Pack{{Size: 5000, Count: 2}, {Size: 2000, Count: 1}, {Size: 250, Count: 1}}},
		{800, []Pack{{Size: 500, Count: 1}, {Size: 250, Count: 1}, {Size: 250, Count: 1}}},
	}

	for _, tc := range testCases {
		result := calculatePacks(tc.order)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("calculatePacks(%v) = %v; want %v", tc.order, result, tc.expected)
		}
	}
}

func TestOptimizePacks(t *testing.T) {
	testCases := []struct {
		packs    []Pack
		expected []Pack
	}{
		{
			[]Pack{{Size: 5000, Count: 1}, {Size: 2000, Count: 1}, {Size: 500, Count: 1}},
			[]Pack{{Size: 5000, Count: 1}, {Size: 2000, Count: 1}, {Size: 500, Count: 1}},
		},
		{
			[]Pack{{Size: 5000, Count: 2}, {Size: 2000, Count: 1}, {Size: 250, Count: 1}},
			[]Pack{{Size: 5000, Count: 2}, {Size: 2000, Count: 1}, {Size: 250, Count: 1}},
		},
		{
			[]Pack{{Size: 500, Count: 1}, {Size: 250, Count: 2}},
			[]Pack{{Size: 1000, Count: 1}},
		},
		{
			[]Pack{{Size: 500, Count: 1}, {Size: 250, Count: 1}, {Size: 250, Count: 1}},
			[]Pack{{Size: 1000, Count: 1}},
		},
	}

	for _, tc := range testCases {
		result := optimizePacks(tc.packs)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("optimizePacks() = %v; want %v", result, tc.expected)
		}
	}
}

func TestCalculateAndOptimizePacks(t *testing.T) {
	testCases := []struct {
		order    int
		expected []Pack
	}{
		{7500, []Pack{{Size: 5000, Count: 1}, {Size: 2000, Count: 1}, {Size: 500, Count: 1}}},
		{1, []Pack{{Size: 250, Count: 1}}},
		{250, []Pack{{Size: 250, Count: 1}}},
		{251, []Pack{{Size: 500, Count: 1}}},
		{501, []Pack{{Size: 500, Count: 1}, {Size: 250, Count: 1}}},
		{12001, []Pack{{Size: 5000, Count: 2}, {Size: 2000, Count: 1}, {Size: 250, Count: 1}}},
		{800, []Pack{{Size: 1000, Count: 1}}},
	}

	for _, tc := range testCases {
		result := optimizePacks(calculatePacks(tc.order))
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("TestCalculateAndOptimizePacks(%v) = %v; want %v", tc.order, result, tc.expected)
		}
	}
}
