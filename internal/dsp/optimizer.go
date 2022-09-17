package dsp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Optimizer struct {
	// Global
	once      sync.Once
	recipeMap map[ItemName]Recipe
}

func NewOptimizer() *Optimizer {
	return &Optimizer{
		recipeMap: make(map[ItemName]Recipe),
	}
}

func (o *Optimizer) GetRecipe(itemName ItemName) (Recipe, bool) {

	o.once.Do(func() {
		// Open up the file
		jsonFile, err := os.Open("data/items.json")
		if err != nil {
			log.Fatal(err)
		}
		defer jsonFile.Close()

		// Read and unmarshal the file
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var recipe []Recipe
		err = json.Unmarshal(byteValue, &recipe)
		if err != nil {
			log.Fatal(err)
		}

		// Map the recipe
		for _, v := range recipe {
			o.recipeMap[v.OutputItem] = v
		}
	})

	result, ok := o.recipeMap[itemName]
	return result, ok
}

func (o *Optimizer) GetOptimalRecipe(itemName ItemName, craftingSpeed float32, parentItemName ItemName) []ComputedRecipe {
	computedRecipes := []ComputedRecipe{}
	recipe, ok := o.GetRecipe(itemName)

	if ok {
		consumedMats := make(map[ItemName]float32)
		numberOfFacilitiesNeeded := recipe.Time * craftingSpeed / recipe.OutputItemCount

		for materialName, materialCount := range recipe.Materials {
			consumedMats[materialName] = materialCount * numberOfFacilitiesNeeded / recipe.Time
		}

		computedRecipe := ComputedRecipe{
			OutputItem:           recipe.OutputItem,
			Facility:             recipe.Facility,
			NumFacilitiesNeeded:  numberOfFacilitiesNeeded,
			ItemsConsumedPerSec:  consumedMats,
			SecondsSpentPerCraft: recipe.Time,
			CraftingPerSec:       craftingSpeed,
			UsedFor:              parentItemName,
		}
		computedRecipes = append(computedRecipes, computedRecipe)

		for materialName, materialCountPerSec := range computedRecipe.ItemsConsumedPerSec {
			targetCraftingSpeed := materialCountPerSec
			cr := o.GetOptimalRecipe(materialName, targetCraftingSpeed, recipe.OutputItem)
			computedRecipes = append(computedRecipes, cr...)
		}
	}

	return computedRecipes
}
