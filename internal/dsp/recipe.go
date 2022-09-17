package dsp

type Recipe struct {
	OutputItem      ItemName             `json:"OutputItem"`
	OutputItemCount float32              `json:"OutputItemCount"`
	Facility        string               `json:"Facility"`
	Time            float32              `json:"Time"`
	Materials       map[ItemName]float32 `json:"Materials"`
}

type ItemName string
