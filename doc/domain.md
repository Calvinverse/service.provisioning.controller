# Domain models

Service handles deployments

## Environment

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

## Service

## Resource

A resource contains

* ID
* Link to the template that created it
* Count
* Status
* Links

## Template

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

## Resource group

* ID
* Resource link
  * ID
  * Template that makes it
* Dependency graph
* Name
* Tags

## Deployment

A combination of a template (or multiple templates) and a configuration that describes what the
values should be for the parameters of the templates
