# Overview

This directory contains all artifacts associated with OpenAPI Client generation.

## How the OpenAPI generator works

The [OpenAPI Generator Project](https://github.com/OpenAPITools/openapi-generator)
works by parsing an OpenAPI Specification file into an in-memory representation
of the API, and feeding those structures to a set of mustache templates that
and generates client or server code in the target language.

An incomplete list of supported languages can be found [here](https://openapi-generator.tech/docs/generators).

The complete list of generators can be found in the GitHub repo at [this path](https://github.com/OpenAPITools/openapi-generator/tree/master/modules/openapi-generator/src/main/resources).

## Generating clients from the spec

Run `make v2`

## What output you get

## Configuring the OpenAPI generator


## Shortcomings of the OpenAPI generator

## Presentation Notes

- Languages we targeted
- Things I learned along the way
  - A single language might have more than one generator - usually splits based HTTP Framework
  - All generators are not created equally - Coarsely grained API output
  - You can override the built-in mustache templates
  - You can add additional templates to generate custom code like a service layer