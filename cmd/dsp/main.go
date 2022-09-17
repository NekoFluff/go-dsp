package main

import (
	"dsp/internal/dsp"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	log.Println("Starting DSP Optimizer Program")
	fmt.Println("Initializing")
	optimizer := dsp.NewOptimizer()
	fmt.Println(optimizer.GetRecipe(dsp.ItemName("asdf")))

	recipe := []dsp.ComputedRecipe{}

	// recipe = recipe.concat(getRecipeForItem('Electromagnetic matrix', 2));
	// recipe = recipe.concat(getRecipeForItem('Energy matrix', 2));
	// // recipe = recipe.concat(getRecipeForItem('Plastic', 2));
	// recipe = recipe.concat(getRecipeForItem('Foundation', 2));
	// // // recipe = recipe.concat(getRecipeForItem('Conveyor belt MK.I', 1));
	// // recipe = recipe.concat(getRecipeForItem('Conveyor belt MK.II', 1));
	// recipe = recipe.concat(getRecipeForItem('Sorter MK.III', 0.5));
	// recipe = recipe.concat(getRecipeForItem('Graphene', 4));

	recipe = append(recipe, optimizer.GetOptimalRecipe("Conveyor belt MK.II", 1, "")...)

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
