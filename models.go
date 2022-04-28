package pagocore

// Ok is boolean response
type Ok struct {
	Ok bool `json:"ok"`
}

// Input is a filterable and validatable input
type Input interface {
	// Filter filters input values
	Filter()

	// Validate checks if input values are valid
	Validate() error
}
