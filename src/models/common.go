package models

import (
	"database/sql/driver"

	"github.com/ochom/gutils/helpers"
)

// MetaData the metadata for an mpesa express request
type MetaData map[string]any

// Scan implements the sql.Scanner interface.
func (m *MetaData) Scan(value any) error {
	v, ok := value.([]byte)
	if !ok {
		return nil
	}

	*m = helpers.FromBytes[map[string]any](v)
	return nil
}

// Value implements the driver.Valuer interface.
func (m MetaData) Value() (driver.Value, error) {
	return helpers.ToBytes(m), nil
}

// Get returns a key in the metadata
func (m MetaData) Get(key string) any {
	return m[key]
}
