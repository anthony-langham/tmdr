package acronym

// Acronym represents a medical acronym with its full form and definition
type Acronym struct {
	Acronym    string
	FullForm   string
	Definition string
}

// Repository defines the interface for acronym storage
type Repository interface {
	Find(acronym string) (*Acronym, error)
	Random() (*Acronym, error)
	All() ([]Acronym, error)
}