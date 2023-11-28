# go-ticket

Development exercise to explore language and library capabilities by building a simple ticketing application.

## Running

Set required environment variables:

```
source .env
```

Run the application:

```
go run ./cmd/main
```

## Testing

### Swagger UI

Get Swagger UI project and build

```
git clone https://github.com/swagger-api/swagger-ui.git
cd swagger-ui
npm install
npm run dev
```

Start a browser with CORS checking disabled and open Swagger UI at location
`http://localhost:3200`.

Change the OpenAPI document used by Swagger UI to `http://localhost:8080/openapi.yml`.
