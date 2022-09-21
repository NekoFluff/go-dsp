package dsp

type ComputedRecipe struct {
	OutputItem           ItemName
	Facility             string
	NumFacilitiesNeeded  float32
	ItemsConsumedPerSec  map[ItemName]float32
	SecondsSpentPerCraft float32
	CraftingPerSec       float32
	UsedFor              ItemName
}
