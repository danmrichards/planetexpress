# Planet Express
Good news everyone! A set of toy applications to demonstrate the event sourcing
and CQRS patterns.

The applications utilise PostgreSQL as an event store and Redis for use as a
message broker and NoSQL database for CQRS.

##  Requirements
* Go 1.13+
* Docker
* Docker-compose

## Applications

* API - Provides a simple REST API for managing the shipment of packages
* Event Store - Listens for events and persists them to the event store
* View Builder - Listens for events and updates a read-only view, implementing the CQRS pattern

## Building From Source
Clone this repo and build the binaries:

```bash
$ make build
```

## Usage
Start the databases with Docker compose:

```bash
$ docker-compose up -d
```

### API
```bash
Usage of ./bin/api-linux-amd64:
```

### Event Store
```bash
Usage of ./bin/eventstore-linux-amd64:
```

### View Builder
```bash
Usage of ./bin/viewbuilder-linux-amd64:
```
