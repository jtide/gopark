package api

// ConfigRates are used to update or replace existing rates.
type ConfigRates struct {
	Rates []ConfigRate `json:"rates"`
}

// ConfigRate is an individual rate for a specific time range on a list of days.
type ConfigRate struct {
	Days  string `json:"days"`
	Times string `json:"times"`
	Price uint   `json:"price"`
}
