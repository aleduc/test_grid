package grid

//go:generate go run github.com/mailru/easyjson/easyjson -all -no_std_marshalers structs.go

type Intensity string

const (
	VeryLow  Intensity = "very low"
	Low      Intensity = "low"
	Moderate Intensity = "moderate"
	High     Intensity = "high"
	VeryHigh Intensity = "very high"
)

//easyjson:json
type Response struct {
	Data []struct {
		From      string `json:"from"`
		To        string `json:"to"`
		Intensity struct {
			Forecast int       `json:"forecast"`
			Actual   int       `json:"actual"`
			Index    Intensity `json:"index"`
		} `json:"intensity"`
	} `json:"data"`
}
