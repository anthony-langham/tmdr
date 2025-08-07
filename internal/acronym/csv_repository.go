package acronym

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

//go:embed data/acronyms.csv
var embeddedCSV string

// CSVRepository implements Repository using a CSV file
type CSVRepository struct {
	data map[string]Acronym
	list []Acronym
}

// NewEmbeddedCSVRepository creates a new CSV-based repository from embedded data
func NewEmbeddedCSVRepository() (*CSVRepository, error) {
	reader := csv.NewReader(strings.NewReader(embeddedCSV))
	reader.FieldsPerRecord = -1 // Allow variable number of fields
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

// NewCSVRepository creates a new CSV-based repository from a file path (for backwards compatibility)
func NewCSVRepository(path string) (*CSVRepository, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields
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

// FindFuzzy performs fuzzy search on acronyms and returns top matches
func (r *CSVRepository) FindFuzzy(acronym string, maxResults int) ([]Acronym, error) {
	if maxResults <= 0 {
		maxResults = 3
	}
	
	acronymUpper := strings.ToUpper(acronym)
	type scoredMatch struct {
		acronym Acronym
		score   int
	}
	
	var matches []scoredMatch
	
	// Calculate similarity scores for all acronyms
	for key, acr := range r.data {
		score := calculateSimilarity(acronymUpper, key)
		if score > 0 {
			matches = append(matches, scoredMatch{acr, score})
		}
	}
	
	// Sort by score (higher is better)
	for i := 0; i < len(matches)-1; i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[j].score > matches[i].score {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}
	
	// Return top matches
	results := []Acronym{}
	for i := 0; i < len(matches) && i < maxResults; i++ {
		results = append(results, matches[i].acronym)
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("no fuzzy matches found for '%s'", acronym)
	}
	
	return results, nil
}

// calculateSimilarity calculates a similarity score between two strings
// Higher score means more similar
func calculateSimilarity(s1, s2 string) int {
	// If strings are equal, perfect score
	if s1 == s2 {
		return 100
	}
	
	// Calculate Levenshtein distance
	dist := levenshteinDistance(s1, s2)
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}
	
	// If distance is too large, no match
	// Allow up to 2 edits for short acronyms, 3 for longer ones
	maxDist := 2
	if maxLen > 4 {
		maxDist = 3
	}
	if dist > maxDist {
		return 0
	}
	
	// Convert distance to similarity score
	score := 100 - (dist * 100 / maxLen)
	
	// Boost score for common patterns
	if strings.Contains(s2, s1) || strings.Contains(s1, s2) {
		score += 20
	}
	
	// Boost for same prefix
	minLen := len(s1)
	if len(s2) < minLen {
		minLen = len(s2)
	}
	for i := 0; i < minLen && i < 3; i++ {
		if s1[i] == s2[i] {
			score += 5
		} else {
			break
		}
	}
	
	return score
}

// levenshteinDistance calculates the edit distance between two strings
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}
	
	// Create matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}
	
	// Initialize first column and row
	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}
	
	// Fill matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			
			matrix[i][j] = min(
				matrix[i-1][j]+1,     // deletion
				matrix[i][j-1]+1,     // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}
	
	return matrix[len(s1)][len(s2)]
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}