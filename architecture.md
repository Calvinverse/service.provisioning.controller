# Implementation / Architecture

There are a number of areas that decisions need to be made on for this application.

* Deployment location
* Communication
* Storage layer
* Communication with other services
* Invocation of tooling

## Location

* Trying to be cloud agnostic as far as possible
  * For those cases where it's harder to be cloud agnostic we will fall back to
    using Azure services
* The service will be deployed in a container so in theory it can run anywhere

## Communication


* Any API calls will need to be resilient
  * https://github.com/cep21/circuit
  * https://dev.to/trongdan_tran/circuit-breaker-and-retry-2klg
  * https://callistaenterprise.se/blogg/teknik/2017/09/11/go-blog-series-part11/
  * https://github.com/eapache/go-resiliency
  * https://medium.com/@slok/goresilience-a-go-library-to-improve-applications-resiliency-14d229aee385
* There will be no commands or API calls to get the current state. It is expected that other
  services will keep their own data store containing the state

### Service bus

* Using RabbitMQ / AMQP
* Second option will be Apache Kafka

### Security

* Service to service security
* Authentication with the event system. Ideally done through Vault
* Authentication with Vault

## Storage

Demands

* How are we going to keep the data consistent with the actual infrastructure. Infrastructure can change, either
  because somebody click-op-sed or because something failed
  * Need to somehow ‘know’ what is in an environment, what is expected and what isn’t. Can keep track of those parts.
  * For now we assume that there is a permission control system that blocks users from doing things
    manually
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

## Service discovery

* Might be tricky or not always necessary because we're using event systems
* In general discovery is done with Consul
  * https://alex.dzyoba.com/blog/go-consul-service/
  * https://github.com/segmentio/consul-go
* Want to be using a service mesh

## Calling out to services

## Testing

* API testing
* Integration testing

### API testing

* https://www.sisense.com/blog/rest-api-testing-strategy-what-exactly-should-you-test/

## Templates

* Describes an environment
* Describes all actions for an environment
* Permissions?

## Observability

* Logs
* Metrics
* Tracing
  * Done via service mesh? -> Note that won't work with service bus messages so the service might
    still have to sort something out
* Audit logging

## Resilience

* Retries / backoff on all calls
*

## Elements

### Environment

An environment consists of

* ID
* Meta
  * Name
  * Description
  * Tags
  * Date of creation
  * Date of destruction
  * Date of planned destruction
* Templates -> Link to templates that create it
* Resources -> Link to resources

Other information that can be obtained for an environment

* Status
  * Deploying
  * Deployed - Waiting / validating
  * Deployed - OK
  * Deployed - Fail
  * Destroying
  * Destroyed

### Resource

A resource contains

* ID
* Link to the template that created it
* Count
* Status
* Links

### Template

A template contains

* ID
* Version
* Name
* Commit
* Dependency graph
* Resource groups
* Resources
* Tags
* Actions
  * App reference
    * URL
    * Name
    * Command to execute
  * Secrets to generate
  * Certificates to generate

Templates can have parameters which have a default value which can be overwritten when deploying the
template.
### Resource group

* ID
* Resource link
  * ID
  * Template that makes it
* Dependency graph
* Name
* Tags

### Deployment

A combination of a template (or multiple templates) and a configuration that describes what the
values should be for the parameters of the templates



### Possible executors

* Terraform
* Consul KV
* Vault secrets

## Commands


* `server` - Runs the application as a service
  * Modes:
    * primary - Indicates that the current instance belongs to the primary cluster, i.e. the active cluster
    * secondary - Indicates that the current instance belongs to the secondary cluster, i.e. the cluster
      that is the disaster recovery (DR) cluster for the primary cluster.
  * Parameters
    * `primary` - Flag indicating if the application belongs to the primary cluster. If so then it will
      assume it is part of the active cluster. If not then it assumes it is part of the disaster recovery
      (DR) cluster. The primary cluster will send signals to the DR cluster(s)
    * `primary-cluster-name` - The Consul name of the primary cluster. It is assumed that there exists
      a WAN connection between the two Consul clusters.
* `bootstrap` - Starts the application in bootstrap mode
  * Starts the application in bootstrap mode. This will create a new cluster in the selected k8s
    instance and fully initialize it, i.e. service discovery, secrets, certificates etc. Once the
    cluster is up and running the resource information for the cluster will be send with the
  * Parameters
    * `location` - Address of the k8s cluster into which the meta environment should be bootstrapped
    * `github-organisation` - The name of the organisation or person who's github account contains the
      repositories containing the configuration files to create the meta cluster
    * `repository-prefix` - The prefix for the repositories containing the configuration files to
      create the meta cluster

Default parameters

* `config` - The file path to the configuration file