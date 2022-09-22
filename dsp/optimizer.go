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
	config    OptimizerConfig
}

type OptimizerConfig struct {
	DataSource string
}

func NewOptimizer(config OptimizerConfig) *Optimizer {
	o := &Optimizer{
		recipeMap: make(map[ItemName]Recipe),
		config:    config,
	}
	o.loadRecipes()
	return o
}

func (o *Optimizer) loadRecipes() {
	o.once.Do(func() {
		// Open up the file
		jsonFile, err := os.Open(o.config.DataSource)
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
}

func (o *Optimizer) GetRecipe(itemName ItemName) (Recipe, bool) {
	result, ok := o.recipeMap[itemName]
	return result, ok
}

func (o *Optimizer) GetRecipes() []Recipe {
	recipes := []Recipe{}
	for _, recipe := range o.recipeMap {
		recipes = append(recipes, recipe)
	}
	return recipes
}

func (o *Optimizer) GetOptimalRecipe(itemName ItemName, craftingSpeed float32, parentItemName ItemName, seenRecipes map[ItemName]bool) []ComputedRecipe {
	computedRecipes := []ComputedRecipe{}
	if seenRecipes[itemName] {
		return computedRecipes
	}
	seenRecipes[itemName] = true

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
			seenRecipesCopy := make(map[ItemName]bool)
			for k, v := range seenRecipes {
				seenRecipesCopy[k] = v
			}
			cr := o.GetOptimalRecipe(materialName, targetCraftingSpeed, recipe.OutputItem, seenRecipesCopy)
			computedRecipes = append(computedRecipes, cr...)
		}
	}

	return computedRecipes
}
