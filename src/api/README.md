# API

## Development Environment

### Install dependencies

#### Install Realize

```
go get github.com/oxequa/realize
```

### Install Redis
See: https://redis.io/topics/quickstart

### Setup ENV file

```
cp env.example .env
```

Edit ENV file as necessary to set API port and NOMSU binary location

### Start Redis daemon
```
$ redis-server
```

### Run local server

```
cd src/api
realize start
```

| Description | URI |
| ------------| --- |
| API Server | http://localhost:1323/ |
| Swagger Docs | http://localhost:1323/swagger/index.html |
| Realize Web UI | http://localhost:5002/ |

## Production Environment

### Coming Soon
