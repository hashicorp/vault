# Overview

This directory contains all artifacts associated with OpenAPI Client generation.

## How the OpenAPI generator works

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