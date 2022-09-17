package dsp

import (
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
