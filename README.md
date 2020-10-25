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
