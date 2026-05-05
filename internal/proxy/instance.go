package proxy

import (
	"errors"
	"fmt"
	"strings"
)

// CloudSQLInstanceParts represents a parsed Cloud SQL connection name.
// Canonical format: project:region:instance.
type CloudSQLInstanceParts struct {
	Project  string
	Region   string
	Instance string
}

var ErrInvalidInstanceFormat = errors.New("invalid instance format")

// ParseCloudSQLInstance parses `project:region:instance` and validates each segment.
func ParseCloudSQLInstance(instance string) (CloudSQLInstanceParts, error) {
	instance = strings.TrimSpace(instance)
	parts := strings.Split(instance, ":")
	if len(parts) != 3 {
		return CloudSQLInstanceParts{}, fmt.Errorf("%w: expected project:region:instance", ErrInvalidInstanceFormat)
	}

	project := strings.TrimSpace(parts[0])
	region := strings.TrimSpace(parts[1])
	name := strings.TrimSpace(parts[2])
	if project == "" || region == "" || name == "" {
		return CloudSQLInstanceParts{}, fmt.Errorf("%w: project, region, and instance must be non-empty", ErrInvalidInstanceFormat)
	}

	return CloudSQLInstanceParts{
		Project:  project,
		Region:   region,
		Instance: name,
	}, nil
}
