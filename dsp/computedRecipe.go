package dsp

type ComputedRecipe struct {
	OutputItem           ItemName
	Facility             string
	NumFacilitiesNeeded  float64
	ItemsConsumedPerSec  map[ItemName]float64
	SecondsSpentPerCraft float64
	CraftingPerSec       float64
	UsedFor              ItemName
	Depth                int `json:"Depth,omitempty"`
}
