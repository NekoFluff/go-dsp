package dsp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
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
			name := ItemName(strings.ToLower(string(v.OutputItem)))
			o.recipeMap[name] = v
		}
	})
}

func (o *Optimizer) GetRecipe(itemName ItemName) (Recipe, bool) {
	name := ItemName(strings.ToLower(string(itemName)))
	result, ok := o.recipeMap[name]
	return result, ok
}

func (o *Optimizer) GetRecipes() []Recipe {
	recipes := []Recipe{}
	for _, recipe := range o.recipeMap {
		recipes = append(recipes, recipe)
	}
	return recipes
}

func (o *Optimizer) GetOptimalRecipe(itemName ItemName, craftingSpeed float64, parentItemName ItemName, seenRecipes map[ItemName]bool, depth int) []ComputedRecipe {
	computedRecipes := []ComputedRecipe{}
	// fmt.Println(itemName)

	if seenRecipes[itemName] {
		return computedRecipes
	}
	seenRecipes[itemName] = true

	recipe, ok := o.GetRecipe(itemName)
	// fmt.Println(recipe, ok)

	if ok {
		consumedMats := make(map[ItemName]float64)
		numberOfFacilitiesNeeded := recipe.Time * craftingSpeed / recipe.OutputItemCount

		for materialName, materialCount := range recipe.Materials {
			c := materialCount * numberOfFacilitiesNeeded / recipe.Time
			if math.IsNaN(float64(c)) {
				c = 0.0
			}
			consumedMats[materialName] = c
		}

		computedRecipe := ComputedRecipe{
			OutputItem:           recipe.OutputItem,
			Facility:             recipe.Facility,
			NumFacilitiesNeeded:  numberOfFacilitiesNeeded,
			ItemsConsumedPerSec:  consumedMats,
			SecondsSpentPerCraft: recipe.Time,
			CraftingPerSec:       craftingSpeed,
			UsedFor:              parentItemName,
			Depth:                depth,
		}
		computedRecipes = append(computedRecipes, computedRecipe)

		for materialName, materialCountPerSec := range computedRecipe.ItemsConsumedPerSec {
			targetCraftingSpeed := materialCountPerSec
			seenRecipesCopy := make(map[ItemName]bool)
			for k, v := range seenRecipes {
				seenRecipesCopy[k] = v
			}
			cr := o.GetOptimalRecipe(materialName, targetCraftingSpeed, recipe.OutputItem, seenRecipesCopy, depth+1)
			computedRecipes = append(computedRecipes, cr...)
		}
	}

	return computedRecipes
}
