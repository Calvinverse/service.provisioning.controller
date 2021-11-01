package repository

import (
	"context"
	"time"

	"github.com/AppsFlyer/go-sundheit/checks"
	log "github.com/sirupsen/logrus"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/calvinverse/service.provisioning.controller/internal/config"
	"github.com/calvinverse/service.provisioning.controller/internal/info"
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

//
// STORAGE
//

type Storage struct {
	cfg config.Configuration

	url      string
	user     string
	password string // this should really be an encrypted string or something

	client driver.Client
	db     driver.Database
	graph  driver.Graph
}

func NewStorage(cfg config.Configuration) (*Storage, error) {
	url := "http://localhost:8529"
	if cfg.IsSet("db.url") {
		url = cfg.GetString("db.url")
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

	storage := &Storage{
		cfg:      cfg,
		url:      url,
		user:     user,
		password: password,
	}

	err := registerHealthCheck(storage)
	if err != nil {
		log.WithError(err).Error("Failed to register the health check.")
		return nil, err
	}

	// Should do this in a way that if we get disconnected we can re-connect
	err = storage.init()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

// Database returns the database.
func (s *Storage) Database() driver.Database {
	return s.db
}

// Graph returns the graph.
func (s *Storage) Graph() driver.Graph {
	return s.graph
}

func (s *Storage) init() error {
	if s.client == nil {
		conn, err := http.NewConnection(http.ConnectionConfig{
			Endpoints: []string{s.url},
		})

		if err != nil {
			log.WithError(err).Error("Failed to create HTTP connection")
			return err // Should try again
		}

		log.WithFields(log.Fields{
			"url":  s.url,
			"user": s.user,
		}).Info("Connecting to database engine")
		c, err := driver.NewClient(driver.ClientConfig{
			Connection:     conn,
			Authentication: driver.BasicAuthentication(s.user, s.password),
		})
		if err != nil {
			log.WithError(err).Error("Failed to connect to the database")
			return err // Try again, except if it is an authentication error, then just give up
		}

		s.client = c
	}

	// Create database
	ctx := context.Background()
	if s.db == nil {
		db, err := createOrLoadDatabase(ctx, s.cfg, s.client)
		if err != nil {
			return err
		}

		s.db = db
	}

	// Create or load the graph
	if s.graph == nil {
		graph, err := createOrLoadGraph(ctx, s.cfg, s.db)
		if err != nil {
			return err
		}

		s.graph = graph
	}

	return nil
}

//
// DB PRIVATE FUNCTIONS
//

func createOrLoadDatabase(ctx context.Context, cfg config.Configuration, client driver.Client) (driver.Database, error) {
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

func createOrLoadGraph(ctx context.Context, cfg config.Configuration, db driver.Database) (driver.Graph, error) {
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
		graph, err := db.Graph(ctx, graphName)
		if err != nil {
			log.
				WithError(err).
				WithFields(log.Fields{
					"database": db.Name(),
					"graph":    graphName,
				}).
				Error("Failed to load the existing graph from the database.")
			return nil, err
		}

		return graph, nil
	}
}

func registerHealthCheck(s *Storage) error {
	check := DbLivelinessCheck(s)

	center := info.GetHealthCenter()
	err := center.RegisterLivelinessCheck(
		check,
		30*time.Second,
		5*time.Second,
		false,
	)

	return err
}

//
// DB HEALTH CHECK
//

// DbLivelinessCheck returns a liveliness check for the database connection.
func DbLivelinessCheck(s *Storage) checks.Check {
	t := time.NewTicker(5 * time.Second)
	check := &dbLivelinessCheck{
		storage: s,
	}
	check.configureTicker(t)

	return check
}

type dbLivelinessCheck struct {
	details   string
	lastError error

	storage *Storage

	ticker *time.Ticker
}

func (s *dbLivelinessCheck) Execute() (details interface{}, err error) {
	return s.details, s.lastError
}

func (s *dbLivelinessCheck) Name() string {
	return "db-liveliness"
}

func (s *dbLivelinessCheck) configureTicker(ticker *time.Ticker) {
	s.ticker = ticker
	go func() {
		for {
			select {
			case <-ticker.C:
				// Check that the database is connected
				s.details = "database connected"
				s.lastError = nil
			}
		}
	}()
}
