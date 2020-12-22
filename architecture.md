# Implementation / Architecture

There are a number of areas that decisions need to be made on for this application.

* Deployment location
* Communication
* Storage layer
* Communication with other services
* Invocation of tooling

Bootstrapping will be important. How are we going to do that?

* No services will be available
* Bootstrapping can be a command with some options that determine the location
  where the local container is deployed to using a terraform script of some kind
  to generate the initial resources
* Deploy to a k8s cluster somewhere. From there we can have all the actual services.
  * Need some way to record what is being created and store that in the database
  * Also need some way to migrate data from one database to another
* Bootstrap should work just from a repo, i.e. all you have done is `git clone`
  so the bootstrap script should be able to
  * Build the application, possibly into a docker container
  * Run the code(?)
  * Deploy the base infrastructure
    * k8s cluster(?)
    * Consul masters
    * storage
    * event system
    * api front-end
    * Service
  * Send data describing the environment that was just created
  * send commands to create the other environments that should be created(?)
* What do we do if there is no k8s cluster
  * microk8s
  * Azure k8s

## Location

* Trying to be cloud agnostic as far as possible
  * For those cases where it's harder to be cloud agnostic we will fall back to
    using Azure services
* The service will be deployed in a container so in theory it can run anywhere

## Communication

* Will provide a (somewhat pragmatic) REST API for obtaining information about the service,
  e.g. health checks etc.
* Commands to the service will be send asynchronously via an event bus.
  * Initially we will support a single message bus, being the AMQP / RabbitMQ
  * Might add others later on
* Service will send notifications back to the event bus when it has done work,
  e.g. created a new environment

### REST API Versioning

* Many different ways to do versioning. All of them have advantages / disadvantages
  and there doesn't seem to be a consensus about the best method.
* Based on this [post and the comments](https://www.troyhunt.com/your-api-versioning-is-wrong-which-is/)
	the the `api/v1` approach is selected

### Service bus

* Using RabbitMQ / AMQP

### Security

## Storage

Demands

* How are we going to keep the data consistent with the actual infrastructure. Infrastructure can change, either
  because somebody click-op-sed or because something failed
  * Need to somehow ‘know’ what is in an environment, what is expected and what isn’t. Can keep track of those parts.
* Storing
  * Environment
    * Resources
    * Tags
    * Name
    * Description
    * Entrypoint(s)
    * Status endpoints
  * Resource
    * Tags
    * Name
    * Description
    * Status endpoints
    * Dependencies -> Resource IDs / External resources
* Do we need to store a dependency graph for quick retrieval?
* How are we going to deal with search?
* Database
  * Prefer faster reads over faster writes – We won’t be doing many writes compared to the number of reads we will do
  * How much data do we store – Not very much, definitely less than a Gb
  * How many users do we have – Not very many. Kinda 1 to start with
  * Data
    * Some is nested – Tags / Environment -> Resource links
    * Some of the data forms a graph – Resource dependencies
    * Some of the data needs to be searchable - All text fields like names and descriptions. Additionally also might
      want to search the dependency tree
  * Do we care about multi-master / geographic distribution? Probably not at the start
* Keep the data for a service together. No other services should touch this data directly
* Use events to communicate with other services
  * How are we going to report progress etc.

Potentially looking at using [arangodb]() because
* OSS available
* Multi-node available in OSS
* Multiple database types
* Has some form of security (LDAP isn't available in OSS, but there's a simple basic auth)


Also need

* Some way of keeping things up to date. Ideally we would get notified when things change, but we may have to poll. Can
  we link to Consul and keep track of the health status?


Options

* Azure: Cosmos
* OSS
  * ArangoDB <-- probably the most suitable because GraphQL and distributed by defaul
  * OrientDB

## Calling out to services