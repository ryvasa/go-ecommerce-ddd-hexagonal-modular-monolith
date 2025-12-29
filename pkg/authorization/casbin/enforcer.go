package casbin

import (
	"fmt"

	lib "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// NewEnforcer creates a new Casbin enforcer.
// If db is provided, it uses a database adapter for persistent policy storage.
// Otherwise, it uses embedded file-based policies for development.
func NewEnforcer(db *gorm.DB) (*lib.Enforcer, error) {
	// Load model from embedded files
	modelBytes, err := Files.ReadFile("model.conf")
	if err != nil {
		return nil, fmt.Errorf("failed to read model.conf: %w", err)
	}

	m := model.NewModel()
	if err := m.LoadModelFromText(string(modelBytes)); err != nil {
		return nil, fmt.Errorf("failed to load model: %w", err)
	}

	// Create enforcer based on configuration
	var enforcer *lib.Enforcer

	if db != nil {
		// Use database adapter for production
		adapter, err := gormadapter.NewAdapterByDB(db)
		if err != nil {
			return nil, fmt.Errorf("failed to create gorm adapter: %w", err)
		}

		enforcer, err = lib.NewEnforcer(m, adapter)
		if err != nil {
			return nil, fmt.Errorf("failed to create enforcer with db adapter: %w", err)
		}

		// Load policies from database
		if err := enforcer.LoadPolicy(); err != nil {
			return nil, fmt.Errorf("failed to load policy from database: %w", err)
		}
	} else {
		// Use file-based policies for development
		enforcer, err = lib.NewEnforcer(m)
		if err != nil {
			return nil, fmt.Errorf("failed to create enforcer: %w", err)
		}

		// Load policies from embedded CSV file
		if err := loadPoliciesFromCSV(enforcer); err != nil {
			// Log warning but don't fail - policies can be added programmatically
			fmt.Printf("Warning: failed to load policies from CSV: %v\n", err)
		}
	}

	return enforcer, nil
}

// loadPoliciesFromCSV loads policies from the embedded policy.csv file
func loadPoliciesFromCSV(enforcer *lib.Enforcer) error {
	policyBytes, err := Files.ReadFile("policy.csv")
	if err != nil {
		return fmt.Errorf("failed to read policy.csv: %w", err)
	}

	// Parse CSV content line by line
	lines := parseCSVLines(string(policyBytes))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		// Determine policy type (p for policy, g for grouping/role)
		policyType := line[0]
		params := line[1:]

		switch policyType {
		case "p":
			// Authorization policy: p, subject, object, action
			if len(params) >= 3 {
				_, err := enforcer.AddPolicy(params[0], params[1], params[2])
				if err != nil {
					return fmt.Errorf("failed to add policy: %w", err)
				}
			}
		case "g":
			// Role/grouping policy: g, user, role
			if len(params) >= 2 {
				_, err := enforcer.AddGroupingPolicy(params[0], params[1])
				if err != nil {
					return fmt.Errorf("failed to add grouping policy: %w", err)
				}
			}
		}
	}

	return nil
}

// parseCSVLines parses CSV content into a 2D string slice
func parseCSVLines(content string) [][]string {
	var result [][]string

	// Split by newlines
	for _, line := range splitLines(content) {
		line = trimSpace(line)

		// Skip empty lines and comments
		if line == "" || (len(line) > 0 && line[0] == '#') {
			continue
		}

		// Split by comma and trim spaces
		var fields []string
		for _, field := range splitByComma(line) {
			fields = append(fields, trimSpace(field))
		}

		if len(fields) > 0 {
			result = append(result, fields)
		}
	}

	return result
}

// Helper functions for string parsing
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func splitByComma(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		parts = append(parts, s[start:])
	}
	return parts
}

func trimSpace(s string) string {
	start := 0
	end := len(s)

	// Trim leading spaces
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\r') {
		start++
	}

	// Trim trailing spaces
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}
