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
	once      sync.Once // Global
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
	o.once.Do(o.LoadRecipes)
	return o
}

func (o *Optimizer) LoadRecipes() {
	// Open up the file
	jsonFile, err := os.Open(o.config.DataSource)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	// Read and unmarshal the file
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var recipes []Recipe
	err = json.Unmarshal(byteValue, &recipes)
	if err != nil {
		log.Fatal(err)
	}

	// Map the recipe
	for _, recipe := range recipes {
		name := ItemName(strings.ToLower(string(recipe.OutputItem)))
		o.recipeMap[name] = append(o.recipeMap[name], recipe)
	}
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

func (o *Optimizer) GetRecipes() [][]Recipe {
	recipes := [][]Recipe{}
	for _, recipe := range o.recipeMap {
		recipes = append(recipes, recipe)
	}
	return recipes
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
	numberOfFacilitiesNeeded := guardInf(float64(recipe.Time * craftingSpeed / recipe.OutputItemCount))

	for materialName, materialCount := range recipe.Materials {
		consumedMats[materialName] = guardInf(float64(materialCount * numberOfFacilitiesNeeded / recipe.Time))
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

			uRecipe.SecondsSpentPerCraft = guardInf(float64(uRecipe.SecondsSpentPerCraft*old_num+recipe.SecondsSpentPerCraft*new_num) / total_num)
			uRecipe.CraftingPerSec = uRecipe.CraftingPerSec + recipe.CraftingPerSec
			uRecipe.UsedFor = fmt.Sprintf("%s | %s (Uses %0.2f/s)", uRecipe.UsedFor, recipe.UsedFor, recipe.CraftingPerSec)
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

func guardInf(x float64) float64 {
	if math.IsNaN(x) {
		return 0.0
	}
	return x
}
