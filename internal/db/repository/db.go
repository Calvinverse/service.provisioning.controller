package repository

import (
	"context"

	log "github.com/sirupsen/logrus"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/calvinverse/service.provisioning.controller/internal/config"
)

const (
	ProvisioningDb    = "provisioning"
	ProvisioningGraph = "graph-provisioning"

	EnvironmentVertex = "vertex-environment"
	ResourceVertex    = "vertex-resource"
	TemplateVertex    = "vertex-template"
	ServiceVertex     = "vertex-service"

	EnvironmentEdges                     = "edge-environment"
	EnvironmentToResourceOutgoingVertex  = "vertex-environment-to-resource-outgoing"
	EnvironmentToResourceIncomingVertext = "vertex-environment-to-resource-incoming"

	EnvironmentToTemplateOutgoingVertex = "vertex-environment-to-template-outgoing"
	EnvironmentToTemplateIncomingVertex = "vertex-environment-to-template-incoming"

	EnvironmentToServiceOutgoingVertex = "vertex-environment-to-service-outgoing"
	EnvironmentToServiceIncomingVertex = "vertex-environment-to-service-incoming"

	ResourceEdges                    = "edge-resource"
	ResourceToTemplateOutgoingVertex = "vertex-resource-to-template-outgoing"
	ResourceToTemplateIncomingVertex = "vertex-resource-to-template-incoming"

	ServiceEdges                    = "edge-service"
	ServiceToResourceOutgoingVertex = "vertex-service-to-resource-outgoing"
	ServiceToResourceIncomingVertex = "vertex-service-to-resource-incoming"
)

type Storage struct {
	client *driver.Client
	db     *driver.Database
	graph  *driver.Graph
}

func Init(cfg config.Configuration) (*Storage, error) {
	url := "http://localhost:8529"
	if cfg.IsSet("db.url") {
		url = cfg.GetString("db.url")
	}

	// Create an HTTP connection to the database
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{url},
	})

	if err != nil {
		log.WithError(err).Error("Failed to create HTTP connection")
		return nil, err
	}

	// Create a client
	if !cfg.IsSet("db.user") {
		log.Error("User name for connection to database not specified.")
		return nil, &DbConnectionConfigValueMissingError{
			value: "db.user",
		}
	}
	user := cfg.GetString("db.user")

	// SHOULD REALLY GET THIS FROM VAULT OR SOMETHING
	if !cfg.IsSet("db.password") {
		log.Error("User password for connection to database not specified.")
		return nil, &DbConnectionConfigValueMissingError{
			value: "db.password",
		}
	}
	password := cfg.GetString("db.password")

	log.WithFields(log.Fields{
		"url":  url,
		"user": user,
	}).Info("Connecting to database engine")
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(user, password),
	})
	if err != nil {
		log.WithError(err).Error("Failed to connect to the database")
		return nil, err
	}

	// Create database
	ctx := context.Background()
	db, err := createDatabase(ctx, c, cfg)
	if err != nil {
		return nil, err
	}

	graph, err := createGraph(ctx, db, cfg)
	if err != nil {
		return nil, err
	}

	return &Storage{
		client: &c,
		db:     &db,
		graph:  &graph,
	}, nil

}

func createDatabase(ctx context.Context, client driver.Client, cfg config.Configuration) (driver.Database, error) {
	databaseName := ProvisioningDb
	if cfg.IsSet("db.name") {
		databaseName = cfg.GetString("db.name")
	}

	ok, err := client.DatabaseExists(ctx, databaseName)
	if err != nil {
		log.
			WithError(err).
			WithFields(log.Fields{
				"database": databaseName,
			}).
			Error("Failed to determine if database exist")
		return nil, err
	}

	if ok {
		log.WithFields(log.Fields{
			"database": databaseName,
		}).Info("Connecting to to existing database instance.")
		retrievedDb, err := client.Database(ctx, databaseName)
		if err != nil {
			log.
				WithError(err).
				WithFields(log.Fields{
					"database": databaseName,
				}).
				Error("Failed to retrieve database")
			return nil, err
		}

		return retrievedDb, nil
	} else {
		log.WithFields(log.Fields{
			"database": databaseName,
		}).Info("Database does not exist. Creating new database.")
		createdDb, err := client.CreateDatabase(ctx, databaseName, nil)
		if err != nil {
			log.
				WithError(err).
				WithFields(log.Fields{
					"database": databaseName,
				}).
				Error("Failed to create database")
			return nil, err
		}

		return createdDb, nil
	}
}

func createEdgeDefinition(name string, outgoing []string, incoming []string) driver.EdgeDefinition {
	// define the edgeCollection to store the edges
	var edgeDefinition driver.EdgeDefinition
	edgeDefinition.Collection = name

	// define a set of collections where an edge is going out...
	edgeDefinition.From = outgoing

	// repeat this for the collections where an edge is going into
	edgeDefinition.To = incoming

	return edgeDefinition
}

func createGraph(ctx context.Context, db driver.Database, cfg config.Configuration) (driver.Graph, error) {
	graphName := ProvisioningGraph
	if cfg.IsSet("db.graph") {
		graphName = cfg.GetString("db.graph")
	}

	ok, err := db.GraphExists(ctx, graphName)
	if err != nil {
		log.
			WithError(err).
			WithFields(log.Fields{
				"database": db.Name(),
				"graph":    graphName,
			}).
			Error("Failed to determine if graph exists")
		return nil, err
	}

	if !ok {
		log.
			WithFields(log.Fields{
				"database": db.Name(),
				"graph":    graphName,
			}).
			Info("The database does not have a graph with the appropriate name. Creating a new one.")

		// define the edgeCollection to store the edges
		environmentEdgeDefinition := createEdgeDefinition(
			EnvironmentEdges,
			[]string{EnvironmentToResourceOutgoingVertex, EnvironmentToTemplateOutgoingVertex, EnvironmentToServiceOutgoingVertex},
			[]string{EnvironmentToResourceIncomingVertext, EnvironmentToTemplateIncomingVertex, EnvironmentToServiceIncomingVertex})

		resourceEdgeDefinition := createEdgeDefinition(
			ResourceEdges,
			[]string{ResourceToTemplateOutgoingVertex},
			[]string{ResourceToTemplateIncomingVertex})

		serviceEdgeDefinition := createEdgeDefinition(
			ServiceEdges,
			[]string{ServiceToResourceOutgoingVertex},
			[]string{ServiceToResourceIncomingVertex})

		// A graph can contain additional vertex collections, defined in the set of orphan collections
		var options driver.CreateGraphOptions
		options.OrphanVertexCollections = []string{EnvironmentVertex, ResourceVertex, ServiceVertex, TemplateVertex}
		options.EdgeDefinitions = []driver.EdgeDefinition{environmentEdgeDefinition, resourceEdgeDefinition, serviceEdgeDefinition}

		graph, err := db.CreateGraph(ctx, graphName, &options)
		if err != nil {
			log.
				WithError(err).
				WithFields(log.Fields{
					"database": db.Name(),
					"graph":    graphName,
				}).
				Error("Failed to create a new graph in the database.")
			return nil, err
		}

		return graph, nil
	} else {
		return db.Graph(ctx, graphName)
	}
}

// Graph returns the graph.
func (s *Storage) Graph() *driver.Graph {
	return s.graph
}
