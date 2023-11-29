# go-ticket

Development exercise to explore language and library capabilities by building a simple ticketing application.

## Running

Set required environment variables:

```
source .env
```

Run the server:

```
go run ./cmd/server
```

Get the API documentation:

```
curl http://localhost:8080/openapi.yml
```

### Explore API with Swagger UI

Get Swagger UI project and build

```
git clone https://github.com/swagger-api/swagger-ui.git
cd swagger-ui
npm install
npm run dev
```

Start a browser with CORS checking disabled. For example on MacOS:

```
open -na /Applications/Google\ Chrome.app --args --user-data-dir="/var/tmp/insecure" --disable-web-security
```

Open the Swagger UI found at `http://localhost:3200`.

Change the OpenAPI document used by Swagger UI to `http://localhost:8080/openapi.yml`.
