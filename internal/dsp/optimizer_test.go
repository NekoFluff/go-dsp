package dsp

import (
	"encoding/json"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDSP_Optimizer_DoesNotInfinitelyLoop(t *testing.T) {
	o := NewOptimizer(OptimizerConfig{
		DataSource: "test_data/loopy_items.json",
	})
	recipe := o.GetOptimalRecipe("Loopy item", 1, "", map[ItemName]bool{})
	assert.Equal(t, []ComputedRecipe{
		{
			OutputItem:          "Loopy item",
			Facility:            "Smelting facility",
			NumFacilitiesNeeded: 1,
			ItemsConsumedPerSec: map[ItemName]float32{
				"Loopy item": 1,
			},
			SecondsSpentPerCraft: 1,
			CraftingPerSec:       1,
			UsedFor:              "",
		},
	}, recipe)
}

func sortRecipes(recipes []ComputedRecipe) {
	sort.SliceStable(recipes, func(i, j int) bool {
		if recipes[i].OutputItem != recipes[j].OutputItem {
			return recipes[i].OutputItem < recipes[j].OutputItem
		} else {
			return recipes[i].UsedFor < recipes[j].UsedFor
		}
	})
}

func TestDSP_Optimizer_E2E_ConveyorBeltMKII(t *testing.T) {
	o := NewOptimizer(OptimizerConfig{
		DataSource: "../../data/items.json",
	})
	expectedRecipes := []ComputedRecipe{}
	f, err := os.ReadFile("test_data/computed_recipe_conveyor_belt_mk_2.json")
	assert.Equal(t, nil, err)
	err = json.Unmarshal(f, &expectedRecipes)
	assert.Equal(t, nil, err)
	sortRecipes(expectedRecipes)

	recipes := o.GetOptimalRecipe("Conveyor belt MK.II", 1, "", map[ItemName]bool{})
	sortRecipes(recipes)

	for k := range expectedRecipes {
		assert.Equal(t, expectedRecipes[k], recipes[k])
	}
}
