package repository

import (
	"context"
	"time"

	driver "github.com/arangodb/go-driver"
	log "github.com/sirupsen/logrus"
)

var ()

// Environment describes a logical collection of resource instances and services that work together
// to achieve one or more goals, e.g. to provide the ability to serve customers with
// the ability to create, edit and store notes.
type Environment struct {
	// ID provides the unique identity of the environment.
	ID string `json:"_key"`

	// DataVersion is the version of the schema that was active when the data was stored.
	DataVersion int `json:"version"`

	// Name is the human readable name of the environment, e.g. 'production'.
	Name string `json:"name"`

	// Description is the human readable description of the environment.
	Description string `json:"description"`

	// CreatedOn is the date and time the environment was created. May be set to nil
	// if the environment has not been created yet.
	CreatedOn *time.Time `json:"created_on"`

	// DestroyedON is the date and time the environment was destroyed. May be set to
	// nil if the environment has not been destroyed yet.
	DestroyedOn *time.Time `json:"destroyed_on"`

	// DestructionPlannedOn is the date and time the environment is planned to be
	// destroyed. If the environment is not planned to be destroyed then it will
	// be set to nil.
	DestructionPlannedOn *time.Time `json:"destruction_planned_on"`
}

// SchemaVersion defines the schema version of the Environment data
func (e *Environment) SchemaVersion() int {
	return 1
}

// FetchEnvironmentByID returns the information about the environment with the given ID.
func FetchEnvironmentByID(ctx context.Context, storage *Storage, id string) (*Environment, error) {
	graph := storage.Graph()

	vertices, err := graph.VertexCollection(ctx, EnvironmentVertex)
	if err != nil {
		log.WithError(err).Error("Failed to get the environment vertex collection")
		return nil, err
	}

	if ok, _ := vertices.DocumentExists(ctx, id); !ok {
		return nil, &UnknownEnvironmentError{
			ID: id,
		}
	}

	result := &Environment{}
	_, err = vertices.ReadDocument(ctx, id, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FetchEnvironments returns the information about all the environments.
func FetchEnvironments(ctx context.Context, storage *Storage) ([]*Environment, error) {
	graph := storage.Graph()

	vertices, err := graph.VertexCollection(ctx, EnvironmentVertex)
	if err != nil {
		log.WithError(err).Error("Failed to get the environment vertex collection")
		return nil, err
	}

	result := []*Environment{}
	_, _, err = vertices.ReadDocuments(ctx, []string{"a"}, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// StoreEnvironment stores the environment data for the given environment in the database.
func StoreEnvironment(ctx context.Context, storage *Storage, e *Environment) error {
	graph := storage.Graph()

	vertices, err := graph.VertexCollection(ctx, EnvironmentVertex)
	if err != nil {
		log.WithError(err).Error("Failed to get the environment vertex collection")
		return err
	}

	e.DataVersion = e.SchemaVersion()
	_, err = vertices.CreateDocument(ctx, *e)
	if err != nil {
		if driver.IsConflict(err) {
			log.
				WithError(err).
				WithFields(log.Fields{
					"ID":   e.ID,
					"name": e.Name,
				}).
				Error("An environment with the given ID already exists")
			return &DuplicateEnvironmentError{
				ID:   e.ID,
				Name: e.Name,
			}
		} else {
			log.Errorf("Failed to store the new environment, %v", err)
		}

		return err
	}

	return nil
}
