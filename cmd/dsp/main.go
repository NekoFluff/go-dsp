package main

import (
	"dsp/internal/dsp"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// Global
var once sync.Once
var recipeMap = make(map[dsp.ItemName]dsp.Recipe)

func GetRecipe(itemName dsp.ItemName) (dsp.Recipe, bool) {

	once.Do(func() {

		// Open up the file
		jsonFile, err := os.Open("data/items.json")
		if err != nil {
			log.Fatal(err)
		}
		defer jsonFile.Close()

		// Read and unmarshal the file
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var recipe []dsp.Recipe
		err = json.Unmarshal(byteValue, &recipe)
		if err != nil {
			log.Fatal(err)
		}

		// Map the recipe
		for _, v := range recipe {
			recipeMap[v.OutputItem] = v
		}
	})

	result, ok := recipeMap[itemName]
	return result, ok
}

func init() {
	fmt.Println("Initializing")
	fmt.Println(GetRecipe(dsp.ItemName("asdf")))
}

func main() {
	log.Println("Starting DSP Optimizer Program")

	recipe := []dsp.ComputedRecipe{}

	// recipe = recipe.concat(getRecipeForItem('Electromagnetic matrix', 2));
	// recipe = recipe.concat(getRecipeForItem('Energy matrix', 2));
	// // recipe = recipe.concat(getRecipeForItem('Plastic', 2));
	// recipe = recipe.concat(getRecipeForItem('Foundation', 2));
	// // // recipe = recipe.concat(getRecipeForItem('Conveyor belt MK.I', 1));
	// // recipe = recipe.concat(getRecipeForItem('Conveyor belt MK.II', 1));
	// recipe = recipe.concat(getRecipeForItem('Sorter MK.III', 0.5));
	// recipe = recipe.concat(getRecipeForItem('Graphene', 4));

	recipe = append(recipe, GetRecipeForItem("Conveyor belt MK.II", 1, "")...)

	// function combineRecipes(recipes) {
	//   const uniqueRecipes = {};

	//   recipes.forEach(recipe => {
	//     if (uniqueRecipes[recipe['Produce']]) { // combine recipe objects

	//       const old_num = uniqueRecipes[recipe['Produce']]['NumFacilitiesNeeded'];
	//       const new_num = recipe['NumFacilitiesNeeded'];
	//       const total_num = old_num + new_num;

	//       for (const [materialName, perSecConsumption] of Object.entries(uniqueRecipes[recipe['Produce']]['ItemsConsumedPerSec'])) {
	//         uniqueRecipes[recipe['Produce']]['ItemsConsumedPerSec'][materialName] = perSecConsumption + recipe['ItemsConsumedPerSec'][materialName]
	//       }

	//       uniqueRecipes[recipe['Produce']]['SecondsSpentPerCraft'] = (uniqueRecipes[recipe['Produce']]['SecondsSpentPerCraft'] * old_num + recipe['SecondsSpentPerCraft'] * new_num) / total_num;
	//       uniqueRecipes[recipe['Produce']]['CraftingPerSec'] = uniqueRecipes[recipe['Produce']]['CraftingPerSec'] + recipe['CraftingPerSec'];
	//       uniqueRecipes[recipe['Produce']]['For'] = uniqueRecipes[recipe['Produce']]['For'].concat(recipe['For']);
	//       uniqueRecipes[recipe['Produce']]['For'] = uniqueRecipes[recipe['Produce']]['For'].filter((v, i, a) => a.indexOf(v) === i); // get unique values
	//       uniqueRecipes[recipe['Produce']]['NumFacilitiesNeeded'] += recipe['NumFacilitiesNeeded'];

	//     } else { // add recipe object
	//       uniqueRecipes[recipe['Produce']] = recipe;
	//     }
	//   });

	//   return Object.values(uniqueRecipes);
	// }

	// const uniqueRecipes = combineRecipes(recipe);
	jsonStr, _ := json.MarshalIndent(recipe, "", "\t")
	fmt.Println(string(jsonStr))
}

func GetRecipeForItem(itemName dsp.ItemName, craftingSpeed float32, parentItemName dsp.ItemName) []dsp.ComputedRecipe {
	computedRecipes := []dsp.ComputedRecipe{}
	recipe, ok := GetRecipe(itemName)

	if ok {
		consumedMats := make(map[dsp.ItemName]float32)
		numberOfFacilitiesNeeded := recipe.Time * craftingSpeed / recipe.OutputItemCount

		for materialName, materialCount := range recipe.Materials {
			consumedMats[materialName] = materialCount * numberOfFacilitiesNeeded / recipe.Time
		}

		computedRecipe := dsp.ComputedRecipe{
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
			cr := GetRecipeForItem(materialName, targetCraftingSpeed, recipe.OutputItem)
			computedRecipes = append(computedRecipes, cr...)
		}
	}

	return computedRecipes
}
