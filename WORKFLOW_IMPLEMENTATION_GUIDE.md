# DIGIT 3.0 `sw-services` Go Migration: Engineering Execution Walkthrough

This document is a practical engineering walkthrough of how we migrated the DIGIT 3.0 `sw-services` module from Java Spring Boot to Go. Instead of just showing the final code, this log breaks down what we actually did step-by-step, the problems we hit along the way, and why we made specific technical calls.

---

## 1. Introduction

**The Goal:**
Rebuild the `sw-services` (Sewerage Services) backend in Go while making sure it perfectly matches the existing DIGIT API contracts.

**Why move to Go?**
* **Memory savings:** Go uses a fraction of the RAM that Spring Boot needs, which translates to immediate infrastructure cost savings.
* **Cold start speed:** Go builds into a single static binary that starts almost instantly. This is a game-changer for Kubernetes autoscaling.
* **Concurrency:** Go's goroutines make handling thousands of concurrent requests incredibly efficient without the heavy thread overhead of the JVM.

**Our "Vertical Slice" Strategy:**
We didn't just build out all the handlers, then all the services, and then the database layer. Instead, we built one complete "vertical slice." We focused entirely on getting a single HTTP request to parse, validate, hit the database, and return a proper response. We intentionally delayed adding things like Kafka or the workflow engine. Why? Because we needed to prove the core synchronous data flow and database connection were rock solid before adding moving parts.

---

## 2. Initial Setup

Before writing any code, we needed to get the repo and our local environment ready.

**Getting the Code:**
```bash
git clone https://github.com/Sushant-RM/sw-services-go.git
cd sw-services-go
git checkout swserv-go-conversion
```
*Why this matters:* We isolated this work on a dedicated migration branch to keep the main codebase clean while we experimented.

**Setting up the Go Workspace:**
```bash
go mod init github.com/BrajK111/sw-services
go get github.com/joho/godotenv
go get github.com/lib/pq
go get github.com/google/uuid
```
*Why this matters:* We initialized the Go module to track our dependencies. We kept it light: `godotenv` for local environment variables, `lib/pq` for the Postgres driver, and `uuid` for generating IDs. Notice we didn't pull in massive web frameworks—we kept it standard library where possible.

---

## 3. How the Project is Structured

With the environment ready, we laid out the code using standard Go domain-driven patterns.

* `cmd/sw-services/`: This is where the app actually boots up (`main.go`). It wires up the database, injects dependencies, and starts the server.
* `internal/transport/http/handler/`: The entry point for incoming requests. It parses JSON and formats the response.
* `internal/validator/`: We catch bad payloads here before they ever reach our business logic.
* `internal/service/`: Where the real work happens. It handles ID generation and orchestrates the database calls.
* `internal/repository/postgres/`: The data access layer. Just raw, parameterized SQL queries.
* `internal/domain/`: The blueprint for our data (DTOs, structs, and custom errors).
* `migrations/ddl/`: The raw SQL scripts to spin up the database schema.

**The Request Lifecycle:**
When a request comes in, it flows predictably:
`HTTP Request → Handler → Validator → Service → Repository → PostgreSQL → HTTP Response`

---

## 4. Setting up PostgreSQL

With the code skeleton in place, we needed a database that mirrored the legacy Java setup.

**Creating the Database:**
```bash
psql -U postgres -c "CREATE DATABASE rainmaker;"
psql -U postgres -d rainmaker -f migrations/ddl/init.sql
```
*Expected Output:* `CREATE DATABASE` followed by table creation logs.
*Why this matters:* This spins up the `rainmaker` database and creates the `eg_sw_connection` table so we have somewhere to write data.

**Why we used JSONB:**
For fields like `connectionHolders`, `plumberInfo`, and `roadCuttingInfo`, we used Postgres `JSONB` columns instead of creating separate relational tables.
*The reasoning:* These arrays belong completely to the main connection record. If we delete the connection, we delete the arrays. If we fetch the connection, we need the arrays. By using JSONB, we avoid the overhead of complex `JOIN` queries and foreign key management, keeping reads and writes fast and simple.

**Database Configuration Gotchas:**
When connecting Go and Docker to local Postgres, we hit a couple of classic hurdles:

1.  **The Listen Addresses Issue:**
    *   *What happened:* By default, local Postgres only listens on `localhost`. This means a Docker container can't talk to it.
    *   *The fix:* We opened `/etc/postgresql/<version>/main/postgresql.conf` and updated it to `listen_addresses = '*'`.
2.  **The Authentication Issue:**
    *   *What happened:* Postgres usually expects local "peer" authentication (matching your Linux username), which breaks when connecting via standard username/password from Go.
    *   *The fix:* We updated `pg_hba.conf` to use `md5` or `scram-sha-256` for local connections, and we made sure our Go connection string explicitly included `sslmode=disable`.

---

## 5. Building the Go Application

After sorting out the database, we needed a clean build pipeline to compile our code.

**1. Cleaning Dependencies**
```bash
go mod tidy
```
*Expected Output:* Downloads missing packages and cleans up unused ones from `go.mod` and `go.sum`.
*Why this matters:* Keeps the dependency tree lean and secure.

**2. Formatting Code**
```bash
go fmt ./...
```
*Expected Output:* A list of files that got auto-formatted.
*Why this matters:* Ends arguments about tabs vs. spaces. Everyone's code looks exactly the same.

**3. Catching Bugs Early**
```bash
go vet ./...
```
*Expected Output:* Nothing (if the code is clean).
*Why this matters:* `vet` catches subtle logical errors and shadowed variables that the compiler might let slip by.

**4. Compiling the Binary**
```bash
go build -o sw-services ./cmd/sw-services
```
*Expected Output:* A fresh `sw-services` binary file drops into the directory.
*Why this matters:* This proves our code actually compiles into the static binary we'll deploy.

---

## 6. Developing the Core APIs

With the binary building successfully, we implemented the three APIs needed for our vertical slice.

**1. The Health API** (`GET /sw-services/actuator/health`)
*   *What it does:* Simply returns `{"status": "UP"}`.
*   *Why we built it first:* Kubernetes and Docker need a way to check if the pod actually booted up correctly before routing traffic to it.

**2. The Create API** (`POST /sw-services/swc/_create`)
*   *Validation:* We built `validator.ValidateCreateConnection` to reject bad requests (like missing a `TenantId`) before hitting the database.
*   *Service Logic:* We generate the `ApplicationNumber` here.
*   *Repository SQL:* We write the actual `INSERT` statement. We specifically used `$1, $2` parameterized placeholders to completely block SQL injection. We also marshal our Go structs into JSON before saving them into the `JSONB` columns.
*   *Response Formatting:* We wrap the new record in a standard DIGIT `ResponseInfo` block so the frontend doesn't even know it's talking to a Go backend instead of Java.

**3. The Search API** (`POST /sw-services/swc/_search`)
*   *What it does:* Pulls `tenantId` and `applicationNumber` out of the URL query params.
*   *Repository SQL:* Runs a `SELECT` query, grabs the rows, and unmarshals the `JSONB` data back into native Go structs so we can serve it as clean JSON.

---

## 7. Containerizing with Docker

Once the code ran locally, we needed to package it up for deployment.

**The Dockerfile Approach:**
We used a Multi-stage Docker build.
*   *Stage 1 (Builder):* Uses `golang:1.24-alpine`. It downloads dependencies and runs `go build`.
*   *Stage 2 (Runner):* Uses a barebones `alpine:latest` image. We just copy the binary from Stage 1 into this image.
*   *Why this matters:* The final image only contains the compiled binary, not the Go compiler or source code. It keeps the image tiny and secure.

**Building the Image:**
```bash
docker build -t sw-services-go .
```
*Expected Output:* Docker steps through the multi-stage build and tags the final image.

**Running the Container:**
```bash
docker run -d --name sw-app -p 3468:3468 -e DB_HOST=host.docker.internal sw-services-go
```
*Expected Output:* A long container ID string.
*Why this matters:* This boots the app in the background, maps port 3468, and injects `host.docker.internal` so the container knows how to route database traffic back to the host machine's Postgres.

**Checking the Status:**
```bash
docker ps
docker logs sw-app
```
*Expected Output:* You should see the container marked "Up" and the logs should say "Server starting on :3468".

---

## 8. Real Problems We Faced (and Fixed)

Migrations are never perfect on the first try. Here is a breakdown of the specific bugs we hit and how we solved them.

**1. The Docker Permission Denied Error**
*   *What went wrong:* When the Alpine container tried to run the binary, it threw a permission error.
*   *How we fixed it:* We explicitly set `CGO_ENABLED=0` during the build phase and made sure the binary had execute rights when copied over.

**2. The PostgreSQL Peer Auth Crash**
*   *What went wrong:* Go's `lib/pq` driver defaults to wanting an SSL connection, and local Postgres defaults to `peer` authentication. The app refused to connect.
*   *How we fixed it:* We forced `sslmode=disable` in the `database.go` connection string and updated the host's `pg_hba.conf`.

**3. The Docker ↔ Host Database Network Trap**
*   *What went wrong:* Inside the Docker container, the app tried to connect to `localhost:5432`. But to a container, "localhost" means the container itself, not the host machine running Postgres.
*   *How we fixed it:* We used `-e DB_HOST=host.docker.internal` to explicitly route the connection out of the container and onto the host's network.

**4. The MySQL vs Postgres Syntax Bug**
*   *What went wrong:* We initially wrote raw SQL queries using `?` placeholders, which is standard for MySQL. Postgres threw syntax errors.
*   *How we fixed it:* We refactored all queries to use Postgres positional arguments (`$1, $2, $3`).

**5. The JSON Capitalization Issue**
*   *What went wrong:* By default, Go capitalizes exported struct fields. This meant our JSON responses had keys like `TenantId` instead of `tenantId`, breaking the DIGIT frontend.
*   *How we fixed it:* We manually mapped every single field using struct tags: ``json:"tenantId"``.

**6. The "Null vs Empty Array" UI Bug**
*   *What went wrong:* If a user didn't provide any connection holders, our Go slice was `nil`. When serialized, this turned into `null` in the JSON response. The DIGIT frontend expects an empty array `[]` and crashed.
*   *How we fixed it:* We added a quick check in the handler. If the slice is `nil`, we initialize it to an empty slice `[]dto.ConnectionHolder{}` before sending the response.

---

## 9. The Live Demo Script

If you're presenting this work, run this exact sequence of commands to prove the architecture is fully functional.

**Step 1: Prove the container is running**
```bash
docker ps
```
*   *What it proves:* The Docker container didn't crash on startup and is actively bound to port 3468.

**Step 2: Prove the database connected**
```bash
docker logs sw-app
```
*   *What it proves:* Shows the application boot logs. If the DB connection failed, it would have panicked here.

**Step 3: Prove HTTP routing works**
```bash
curl -X GET http://localhost:3468/sw-services/actuator/health
```
*   *Expected output:* `{"status":"UP"}`
*   *What it proves:* The web server is successfully handling traffic.

**Step 4: Prove the Create API & Database Insert**
```bash
curl -X POST http://localhost:3468/sw-services/swc/_create \
-H "Content-Type: application/json" \
-d '{
  "RequestInfo": { "apiId": "Rainmaker", "ver": ".01", "msgId": "20170310130900|en_IN" },
  "SewerageConnection": {
    "tenantId": "pb.amritsar",
    "propertyId": "PROP-123",
    "connectionType": "Non-Metered",
    "connectionExecutionDate": 1616568453000
  }
}'
```
*   *Expected output:* A massive JSON payload with `status: successful` and a new `applicationNumber`.
*   *What it proves:* The request successfully passed validation, the service generated an ID, and the repository successfully executed the `$1, $2` SQL insert (along with the JSONB marshaling).

**Step 5: Prove the Search API & JSONB Unmarshaling**
```bash
curl -X POST "http://localhost:3468/sw-services/swc/_search?tenantId=pb.amritsar&applicationNumber=<INSERT_APP_NUM_HERE>" \
-H "Content-Type: application/json" \
-d '{ "RequestInfo": { "apiId": "Rainmaker" } }'
```
*   *Expected output:* The exact connection data you just created.
*   *What it proves:* The repository successfully queried Postgres and cleanly unmarshaled the JSONB data back into our Go structs.

---

## 10. What's Next?

We intentionally delayed several features to make sure the synchronous request-response cycle was completely stable first. Now that the vertical slice is proven, the next phases will introduce:

*   **Kafka Integration:** We need to publish audit logs and billing events asynchronously so we don't slow down the main API response.
*   **Workflow Engine:** We need to hook into the DIGIT `eg-workflow` service to handle state transitions (e.g., moving a connection from "Pending" to "Approved").
*   **API Middleware:** We need to add standard interceptors for checking auth tokens and injecting correlation IDs for tracing.
*   **Testing & Observability:** Adding table-driven unit tests, Swagger documentation, and structured logging so the service is actually maintainable in production.

By securing this foundation first, we ensure that adding these complex integrations later won't destabilize the core application.
