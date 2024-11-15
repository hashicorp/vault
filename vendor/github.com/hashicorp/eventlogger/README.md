# go-eventlogger   [![Go Reference](https://pkg.go.dev/badge/github.com/hashicorp/eventlogger.svg)](https://pkg.go.dev/github.com/hashicorp/eventlogger)


`go-eventlogger` is a flexible event system libray implemented as a pub/sub model supporting middleware. 

The library's clients submit events to a Broker that routes them through a pipeline. Each node in pipline can modifying the event, filtering it, persisting it, etc.  

# Stability Notice 

While this library is fully open source and HashiCorp will be maintaining it (since we are and will be making extensive use of it), the API and output format is subject to minor changes as we fully bake and vet it in our projects. This notice will be removed once it's fully integrated into our major projects and no further changes are anticipated.


# Usage

An Event is a collection of data, analogous to a log entry, that we want to
process via pipelines. A pipeline is a graph composed of nodes.  The client
provides an event type and payload, and any other fields are generated as part
of processing. The library will not attempt to discover whether configured
formatter/marshaller nodes can actually handle the arbitrary payloads; it is up
to the encapsulating  program to put any such constraints on the user via its API.

The library's clients submit events to a Broker that routes them through 
pipelines (graphs) based on their type.  A pipeline is a graph composed of
Nodes.  A Node processes an Event in some way -- modifying it, filtering it,
persisting it, etc.  A Sink is a Node that persists an Event.

## Broker

Clients interact with the library via the Broker by sending events. A Broker
processes incoming Events, by sending them to the pipelines (graphs) associated
with the Event's type.  A given Broker, along with its associated set of
pipelines (graphs), will be configured programmatically. 


## Nodes 
A Node is a node in a Pipeline, that can perform operations on an Event.  A node
has a Type, one of: Filter, Formatter, Sink.

![Node example](img/pipeline.jpg)


Examples of things that a Node might do to an Event include:

Modify the Event, by storing a change description in Mutations.  Changes could
be described as a  (jsonpointer, interface{}) key-value pair. Filter the Event
out of the pipeline, by  returning nil. Get the Event ready for a sink by
rendering (formatting) it in someway, e.g. as JSON, so that downstream Sinks in 
the pipeline can then write it without any extra work.  Rendered events will be 
stored in the Formatted map.


## Pipeline 

A Pipeline is a pointer to the root of an interconnected sequence of Nodes. 

All pipelines with a Sink must contain a Formatter that precedes the
Sink and formats events in the Sink's required format.  

When using a FileSink without a specified format it will default to JSON and a
JSONFormatter must be in the pipeline before the FileSink.

All pipelines must end with a sink node.


# Contributing 

First: if you're unsure or afraid of anything, just ask or submit the issue or pull request anyways. You won't be yelled at for giving your best effort. The worst that can happen is that you'll be politely asked to change something. We appreciate any sort of contributions, and don't want a wall of rules to get in the way of that.

That said, if you want to ensure that a pull request is likely to be merged, talk to us! A great way to do this is in issues themselves. When you want to work on an issue, comment on it first and tell us the approach you want to take.

## Build

If you have the following requirements met locally:

* Golang v1.16 or greater

Please note that development may require other tools; to install the set of tools at the versions used by the Boundary team, run:

`make tools`

Before opening a PR, please run:

`make fmt` 