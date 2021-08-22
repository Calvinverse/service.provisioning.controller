package repository

import "fmt"

// DbConnectionConfigValueMissingError defines an error that is raised when one of the configuration
// values used to connect to the database is not provided.
type DbConnectionConfigValueMissingError struct {
	value string
}

func (e *DbConnectionConfigValueMissingError) Error() string {
	return fmt.Sprintf("The configuration value named '%s' was missing. This value is used to connect to the database.", e.value)
}

type DuplicateEnvironmentError struct {
	ID   string
	Name string
}

func (e *DuplicateEnvironmentError) Error() string {
	return fmt.Sprintf("An environment with the ID '%s' was already defined.", e.ID)
}
