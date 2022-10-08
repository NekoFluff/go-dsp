package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/NekoFluff/go-dsp/dsp"
)

func main() {
	log.Println("Starting DSP Optimizer Program")
	optimizer := dsp.NewOptimizer(dsp.OptimizerConfig{
		DataSource: "data/items.json",
	})
	// fmt.Println(optimizer.GetRecipe(dsp.ItemName("Iron Ingot")))

	recipes := []dsp.ComputedRecipe{}
	recipeName := dsp.ItemName("Deuteron Fuel Rod")
	// recipe = recipe.concat(getRecipeForItem('Electromagnetic matrix', 2));
	// recipe = recipe.concat(getRecipeForItem('Energy matrix', 2));
	// // recipe = recipe.concat(getRecipeForItem('Plastic', 2));
	// recipe = recipe.concat(getRecipeForItem('Foundation', 2));
	// // // recipe = recipe.concat(getRecipeForItem('Conveyor belt MK.I', 1));
	// // recipe = recipe.concat(getRecipeForItem('Conveyor belt MK.II', 1));
	// recipe = recipe.concat(getRecipeForItem('Sorter MK.III', 0.5));
	// recipe = recipe.concat(getRecipeForItem('Graphene', 4));

	// recipe = append(recipe, optimizer.GetOptimalRecipe("Conveyor belt MK.II", 1, "", map[dsp.ItemName]bool{})...)
	recipes = append(recipes, optimizer.GetOptimalRecipe(recipeName, 1, "", map[dsp.ItemName]bool{}, 1)...)
	recipes = combineRecipes(recipes)

	// Sort
	sortRecipes(recipes)

	// Print out
	jsonStr, err := json.MarshalIndent(recipes, "", "\t")
	if err != nil {
		fmt.Println(err)
	}

	fileName := strings.ToLower(strings.TrimSpace(string(recipeName)))
	fileName = strings.ReplaceAll(fileName, " ", "_")

	err = os.WriteFile("dsp_output_"+fileName+".json", jsonStr, 0644)
	if err != nil {
		fmt.Println("Failed to write to output.json", err)
	}
	log.Println("Output to output.json")
}

func combineRecipes(recipes []dsp.ComputedRecipe) []dsp.ComputedRecipe {
	uniqueRecipes := make(map[dsp.ItemName]dsp.ComputedRecipe)

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
			uRecipe.UsedFor = uRecipe.UsedFor + " | " + recipe.UsedFor
			// uRecipe.UsedFor = uRecipe.UsedFor.filter((v, i, a) => a.indexOf(v) === i); // get unique values
			uRecipe.NumFacilitiesNeeded += recipe.NumFacilitiesNeeded
			uRecipe.Depth = max(uRecipe.Depth, recipe.Depth)
			uniqueRecipes[recipe.OutputItem] = uRecipe

		} else { // add recipe object
			uniqueRecipes[recipe.OutputItem] = recipe
		}
	}

	v := []dsp.ComputedRecipe{}
	for _, value := range uniqueRecipes {
		v = append(v, value)
	}
	return v
}

func sortRecipes(recipes []dsp.ComputedRecipe) {
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

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
