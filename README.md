# Planet Express
A set of toy applications to demonstrate the event sourcing and CQRS patterns.

The applications utilise PostgreSQL as an event store and Redis for use as a
message broker and NoSQL database for CQRS.

##  Requirements
* Go 1.13+
* Docker
* Docker-compose

## Usage
Start the databases with Docker compose:

```bash
$ docker-compose up -d
```