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

var JSONDefaultRateConfig = []byte(
	`{
    "rates": [
        {
            "days": "mon,tues,thurs",
            "times": "0900-2100",
            "price": 1500
        },
        {
            "days": "fri,sat,sun",
            "times": "0900-2100",
            "price": 2000
        },
        {
            "days": "wed",
            "times": "0600-1800",
            "price": 1750
        },
        {
            "days": "mon,wed,sat",
            "times": "0100-0500",
            "price": 1000
        },
        {
            "days": "sun,tues",
            "times": "0100-0700",
            "price": 925
        }
    ]
}`)
