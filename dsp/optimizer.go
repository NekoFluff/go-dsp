package dsp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"sync"
)

type Optimizer struct {
	// Global
	once      sync.Once
	recipeMap map[ItemName][]Recipe
	config    OptimizerConfig
}

type OptimizerConfig struct {
	DataSource string
}

func NewOptimizer(config OptimizerConfig) *Optimizer {
	o := &Optimizer{
		recipeMap: make(map[ItemName][]Recipe),
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
			o.recipeMap[name] = append(o.recipeMap[name], v)
		}
	})
}

func (o *Optimizer) GetRecipe(itemName ItemName, recipeIdx int) (Recipe, bool) {
	name := ItemName(strings.ToLower(string(itemName)))
	recipes, ok := o.recipeMap[name]

	if !ok {
		return Recipe{}, ok
	}

	if len(recipes) > recipeIdx {
		recipe := recipes[recipeIdx]
		return recipe, true
	}

	return recipes[0], true
}

func (o *Optimizer) GetRecipes() map[ItemName][]Recipe {
	return o.recipeMap
}

func (o *Optimizer) GetOptimalRecipe(itemName ItemName, craftingSpeed float64, parentItemName ItemName, seenRecipes map[ItemName]bool, depth int, recipeRequirements RecipeRequirements) []ComputedRecipe {
	computedRecipes := []ComputedRecipe{}

	if seenRecipes[itemName] {
		return computedRecipes
	}
	seenRecipes[itemName] = true

	rRequirement, ok := recipeRequirements[itemName]
	recipeIdx := 0
	if ok {
		recipeIdx = rRequirement
	}

	recipe, ok := o.GetRecipe(itemName, recipeIdx)
	if !ok {
		return computedRecipes
	}

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
		UsedFor:              string(parentItemName),
		Depth:                depth,
	}
	computedRecipes = append(computedRecipes, computedRecipe)

	for materialName, materialCountPerSec := range computedRecipe.ItemsConsumedPerSec {
		targetCraftingSpeed := materialCountPerSec
		seenRecipesCopy := make(map[ItemName]bool)
		for k, v := range seenRecipes {
			seenRecipesCopy[k] = v
		}
		cr := o.GetOptimalRecipe(materialName, targetCraftingSpeed, recipe.OutputItem, seenRecipesCopy, depth+1, recipeRequirements)
		computedRecipes = append(computedRecipes, cr...)
	}

	return computedRecipes
}

func (o *Optimizer) SortRecipes(recipes []ComputedRecipe) {
	sort.SliceStable(recipes, func(i, j int) bool {
		if recipes[i].Depth != recipes[j].Depth {
			return recipes[i].Depth < recipes[j].Depth
		} else if recipes[i].OutputItem != recipes[j].OutputItem {
			return recipes[i].OutputItem < recipes[j].OutputItem
		} else if recipes[i].UsedFor != recipes[j].UsedFor {
			return recipes[i].UsedFor < recipes[j].UsedFor
		} else {
			return recipes[i].CraftingPerSec < recipes[j].CraftingPerSec
		}
	})
}

func (o *Optimizer) CombineRecipes(recipes []ComputedRecipe) []ComputedRecipe {
	uniqueRecipes := make(map[ItemName]ComputedRecipe)

	for _, recipe := range recipes {
		if uRecipe, ok := uniqueRecipes[recipe.OutputItem]; ok { // combine recipe objects

			old_num := uRecipe.NumFacilitiesNeeded
			new_num := recipe.NumFacilitiesNeeded
			total_num := old_num + new_num
			for materialName, perSecConsumption := range uRecipe.ItemsConsumedPerSec {
				uRecipe.ItemsConsumedPerSec[materialName] = perSecConsumption + recipe.ItemsConsumedPerSec[materialName]
			}

			sspc := (uRecipe.SecondsSpentPerCraft*old_num + recipe.SecondsSpentPerCraft*new_num) / total_num
			if math.IsNaN(float64(sspc)) {
				sspc = 0.0
			}
			uRecipe.SecondsSpentPerCraft = sspc

			uRecipe.CraftingPerSec = uRecipe.CraftingPerSec + recipe.CraftingPerSec
			uRecipe.UsedFor = fmt.Sprintf("%s | %s (Uses %0.2f/s)", uRecipe.UsedFor, recipe.UsedFor, recipe.CraftingPerSec)
			// uRecipe.UsedFor = uRecipe.UsedFor.filter((v, i, a) => a.indexOf(v) === i); // get unique values
			uRecipe.NumFacilitiesNeeded += recipe.NumFacilitiesNeeded
			uRecipe.Depth = max(uRecipe.Depth, recipe.Depth)
			uniqueRecipes[recipe.OutputItem] = uRecipe

		} else { // add recipe object
			if recipe.UsedFor != "" {
				recipe.UsedFor = fmt.Sprintf("%s (Uses %0.2f/s)", recipe.UsedFor, recipe.CraftingPerSec)
			}
			uniqueRecipes[recipe.OutputItem] = recipe
		}
	}

	v := []ComputedRecipe{}
	for _, value := range uniqueRecipes {
		v = append(v, value)
	}
	return v
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
