package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/NekoFluff/go-dsp/dsp"
)

func main() {
	log.Println("Starting DSP Optimizer Program")
	optimizer := dsp.NewOptimizer(dsp.OptimizerConfig{
		DataSource: "data/items.json",
	})

	rRequirements := dsp.RecipeRequirements{
		"Carbon Nanotube":    0,
		"Crystal Silicon":    0,
		"Diamond":            0,
		"Deuterium":          0,
		"Particle Container": 0,
		"Graphene":           1,
		"Hydrogen":           0,
		"Fire Ice":           1,
		"Casimir Crystal":    0,
	}
	// fmt.Println(optimizer.GetRecipe(dsp.ItemName("Iron Ingot")))

	recipes := []dsp.ComputedRecipe{}
	recipeName := dsp.ItemName("Universe Matrix")
	recipes = append(recipes, optimizer.GetOptimalRecipe(recipeName, 1, "", map[dsp.ItemName]bool{}, 1, rRequirements)...)
	optimizer.SortRecipes(recipes)

	recipes = optimizer.CombineRecipes(recipes)

	// Sort
	optimizer.SortRecipes(recipes)

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
	_ = os.WriteFile("dsp_output.json", jsonStr, 0644)

	log.Println("Output to output.json")
}
