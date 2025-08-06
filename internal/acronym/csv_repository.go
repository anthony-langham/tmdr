package acronym

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

// CSVRepository implements Repository using a CSV file
type CSVRepository struct {
	data map[string]Acronym
	list []Acronym
}

// NewCSVRepository creates a new CSV-based repository
func NewCSVRepository(path string) (*CSVRepository, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	repo := &CSVRepository{
		data: make(map[string]Acronym),
		list: []Acronym{},
	}

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}

		if len(record) < 2 {
			continue
		}

		// Parse the definition to extract full form and description
		parts := strings.SplitN(record[1], "–", 2)
		fullForm := strings.TrimSpace(parts[0])
		definition := ""
		if len(parts) > 1 {
			definition = strings.TrimSpace(parts[1])
		}

		acronym := Acronym{
			Acronym:    strings.ToUpper(record[0]),
			FullForm:   fullForm,
			Definition: definition,
		}

		repo.data[strings.ToUpper(record[0])] = acronym
		repo.list = append(repo.list, acronym)
	}

	return repo, nil
}

// Find looks up an acronym by its abbreviation
func (r *CSVRepository) Find(acronym string) (*Acronym, error) {
	a, exists := r.data[strings.ToUpper(acronym)]
	if !exists {
		return nil, fmt.Errorf("acronym '%s' not found", acronym)
	}
	return &a, nil
}

// Random returns a random acronym
func (r *CSVRepository) Random() (*Acronym, error) {
	if len(r.list) == 0 {
		return nil, fmt.Errorf("no acronyms available")
	}
	idx := rand.Intn(len(r.list))
	return &r.list[idx], nil
}

// All returns all acronyms
func (r *CSVRepository) All() ([]Acronym, error) {
	return r.list, nil
}