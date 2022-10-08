package dsp

type Recipe struct {
	OutputItem      ItemName             `json:"OutputItem"`
	OutputItemCount float64              `json:"OutputItemCount"`
	Facility        string               `json:"Facility"`
	Time            float64              `json:"Time"`
	Materials       map[ItemName]float64 `json:"Materials"`
}

type ItemName string
