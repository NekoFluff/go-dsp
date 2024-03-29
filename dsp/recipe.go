package dsp

type Recipe struct {
	OutputItem      ItemName             `json:"OutputItem"`
	OutputItemCount float64              `json:"OutputItemCount"`
	Facility        string               `json:"Facility"`
	Time            float64              `json:"Time"`
	Materials       map[ItemName]float64 `json:"Materials"`
	Image           string               `json:"Image"`
}

type RecipeRequirements map[ItemName]int

type ItemName string
