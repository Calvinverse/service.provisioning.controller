# service.provisioning.controller

The `service.provisioning.controller` repository is part of the group of
[provisioning services](https://github.com/Calvinverse?q=topic%3Aprovisioning&type=&language=).
This group consists of services and tools which allow users to provision a complete environment from
scratch and upgrade existing ones with new components. In this case an environment is defined as

  > An **environment** is a collection of resource instances and services that work together
  > to achieve one or more goals, e.g. to provide the ability to serve customers with
  > the ability to create, edit and store notes.

## Goals

The `service.provisioning.controller` service provides a standard REST API which will

* Process requests for environment creation based on a set of templates for the different
  services that should be present in the environment
* Process requests for the updating of one or more services in an existing environment
* Process requests for the deletion of one or more services in an existing environment
* Process requests for the deletion of an existing environment

Internally the creation, updating or deleting of an environment will be done with tools
like [Terraform](https://www.terraform.io/).

## Architecture / Design

The [architecture](./architecture.md) document contains some of the design decisions made for this
project.

### Communication

In order to communicate with the application two different paths exist.

1) There are a number of REST API's that provide information about the application itself. The REST
  API does not provide a means to send environment specific commands to the application. Those are
  expected to be send via an event bus, e.g. RabbitMQ or Kafka.
1) Commands to make changes to the environments are send via an event bus.
1) Information about the current state of the environments should be obtained directly from the
  data storage. It is expected that other services and BFFs have their own data cache.

#### REST API

When running as a service the application provides a number of REST API's which are grouped by
function. The API sub path takes the following format

    /api/<API_VERSION>/<API_CATEGORY>/<API_CALL>

Where:

* `API_VERSION` is the [version of the API](https://www.troyhunt.com/your-api-versioning-is-wrong-which-is/). Currently only a `v1` API is available.
* `API_CATEGORY` is the category of the API grouping, e.g. `self`.
* `API_CALL` is the API method.

The following REST API groupings exist

* `doc` - Methods that provide the API documentation for the application
* `self` - Methods that describe the state of the currently executing application.

##### Doc

The `doc` API group contains methods that provide the API documentation for the application.

* `/api/v1/doc` - Returns the OpenAPI document for the current application.

##### Self

The `self` API group contains methods that describe the application as it is currently executing.

* `/api/v1/self/info` - Returns information about the application consisting of version number, time
  the application was built and the GIT SHA for the commit on which the application was based.
* `/api/v1/self/liveliness?type=<REPORT_TYPE>` - Returns the liveliness report, indicating both the
  overall health of the application and which of the health checks are passing and which are failing.
  If `summary` is passed for `REPORT_TYPE` then the short version of the information is returned, containing only the pass or fail values. If `detailed` is passed for `REPORT_TYPE` then the extended
  version of information is returned containing the pass and fail values as well as other information.
  If no value is set for `type` then `summary` is assumed.
* `/api/v1/self/ping` - Returns a `pong` response with the current time.
* `/api/v1/self/readiness` - Returns a response indicating if the application is ready to start serving
  requests. Readiness is defined as being connected to all dependencies, e.g. databases and messaging systems, having all the credentials required to operate, or being able to obtain those credentials, and having registered with the service discovery system.
* `/api/v1/self/started` - Returns a response indicating if the application has started successfully.
  A successful start is defined as the application is loaded and was able to load the configuration
  information from the different configuration sources.

The `liveliness`, `readiness` and `started` API's are used for
[Kubernetes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#before-you-begin) probes.

#### Events and commands

* Commands to the service will be send asynchronously via an event bus.
  * Initially we will support a single message bus, being the AMQP / RabbitMQ
  * Might add others later on
* Service will send notifications back to the event bus when it has done work,
  e.g. created a new environment

## Running

To run the container execute the following command

    docker run -p 8080:8080 -p 8301:8301 --read-only -v d:\ops\local\docker\unbound_zones.conf:/etc/unbound.d/unbound_zones.conf --env-file d:/ops/local/docker/service.provisioning/env.txt -i -t service-provisioning-controller:0.2.0-health-checks.1
