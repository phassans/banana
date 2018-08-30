package prefernce

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type preferenceEngine struct {
	sql    *sql.DB
	logger zerolog.Logger
}

// PreferenceEngine an interface for preference operations
type PrefernceEngine interface {
	PreferenceAdd(phoneID string, cuisine []string) error
	PreferenceDelete(phoneID string, cuisine []string) error
	PreferenceAll(phoneID string) ([]string, error)
}

// NewPrefernceEngine returns an instance of notificationEngine
func NewPrefernceEngine(psql *sql.DB, logger zerolog.Logger) PrefernceEngine {
	return &preferenceEngine{psql, logger}
}

func (*preferenceEngine) PreferenceAdd(phoneID string, cuisine []string) error {
	return nil
}

func (*preferenceEngine) PreferenceDelete(phoneID string, cuisine []string) error {
	return nil
}

func (*preferenceEngine) PreferenceAll(phoneID string) ([]string, error) {
	return []string{"BBQ", "Seafood"}, nil
}
