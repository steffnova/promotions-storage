package api

type Promotion struct {
	ID             string `json:"id"`              // ID promotion ID
	Price          string `json:"price"`           // Price is price is expressed as string but can be parsed to float64
	ExparationDate string `json:"exparation_date"` // Exparation data is a string and can be parsed to time.Time
}
