## Golang Micro Service Template

Simple golang micro service template. This service receives the "MessagePayload" as a JSON string from the main platform
and returns the "ResponseObject" as a JSON string. See "MessagePayload" and "ResponseObject" in models.
This service is placed in a docker container and hot loaded to the platform via NATS server

```
Config - config.toml is added when docker image is created. Samples files are available in the S3 config bucket.

Databases - Mongo, Neo4j, Postgres, Redis (local and shared)

Mongo uses a central db server. The database is setup as multi tenant and is the main database for client data.
The database name used by this service is the client ID.

Neo4j uses a central db server.  This database can be used for logistics, analytics, recommendations and biz intelligence.

Postgress uses a central db server. This database can be used for authentication, quotes, transactions, purchasing or
for other transactional data.

Redis is accessible via a local redis instance that should only be used for volitile data and accessible via a
central cache server that stores session data that is shared with all services.

```