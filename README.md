# fc-retrieval-gateway

Filecoin Retrieval Gateway

NOTE: There is a private key in the settings.json file. This is ONLY being used for testing. It does not post a security issues.

## Start the service

### Create a config file

Create a `.env` file, using [.env.example](./.env.example) as a reference:

```
cp .env.example .env
```

## Development

### Start the service with Docker

Start the service:

```
make dev
```

The APIs should be available at

- Client Api: `http://localhost:9010`
- Provider Api: `http://localhost:9011`
- Gateway Api: `http://localhost:9012`
- Admin Api: `http://localhost:9013`

To rebuild containers, following option can be used:

```
make dev arg=--build
```
