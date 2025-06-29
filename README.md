# 🛠️ Data Ingestion Service

This is a Go-based data ingestion service that fetches posts from a remote API, transforms the data, and uploads the output to AWS S3. The application also provides HTTP endpoints to trigger ingestion and retrieve the most recent ingested data.

---

## 📦 Features

- Fetches data from a JSON placeholder API
- Transforms and normalizes the data
- Uploads transformed data to AWS S3 bucket as a JSON file
- REST API to manually trigger ingestion and view latest data
- Clean and modular structure

---

## ⚙️ Setup Instructions

### ✅ Prerequisites

- Go 1.20 or higher
- AWS credentials with access to S3

---

### 📁 Environment Configuration

Create a `.env` file based on the example provided:

```bash
cp .env.example .env
```

Update it with your AWS details::

```
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_REGION=us-east-1
AWS_BUCKET=your-bucket-name
```

## 🏃 Running the Application
### ▶️ Start the HTTP server
```bash
go run cmd/main.go
```
Server will be available at: http://localhost:8080

### 🔄 Run ingestion via CLI
```bash
go run cmd/main.go ingest
```
This will fetch, transform, and upload the data to your configured AWS S3 bucket.

## 📡 API Documentation
### 🔁 POST /ingest
Triggers the ingestion pipeline: fetch → transform → upload.

```bash
curl -X POST http://localhost:8080/ingest
```
Response:

```
Ingestion completed successfully
```

### 📥 GET /ingested-posts
Retrieves the most recent ingested JSON file from S3.

```bash
curl http://localhost:8080/ingested-posts
```

Response:

```json
[
  {
    "userId": 1,
    "id": 1,
    "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
    "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto",
    "ingested_at": "2025-06-24T18:31:47Z",
    "source": "placeholder_api"
  }
]
```

## 🔄 Transformation Logic
The transformation logic extracts the relevant fields (userId, id, title, body) from each post, append fields ingested_at, source and clean JSON structure suitable for downstream processing or analytics.

## 🧪 Running Tests
### ✅ Run unit and integration tests
```bash
go test ./... -v -coverprofile=coverage.out
```

### ✅ Check coverage summary
```bash
go tool cover -func=coverage.out
```

### ✅ Open detailed coverage in browser
```bash
go tool cover -html=coverage.out
```

Test coverage includes:
- API route behavior
- Upload logic (via mocks)
- Failure scenarios (e.g., fetch errors, upload errors)
- Configuration loading

## ⚖️ Trade-offs Considered
| Decision                         | Reason                                                  | Trade-off                             |
|----------------------------------|---------------------------------------------------------|----------------------------------------|
| AWS S3 only (no MinIO)          | Simplifies deployment and aligns with real use case     | No local offline testing               |
| No interfaces for fetcher/uploader | Keeps code simple and readable                         | Harder to mock deeply in unit tests    |
| Light transformation logic      | Fast implementation                                     | Not extensible for complex schemas yet |


## 🚧 Hardest Parts
- Mocking internal functions (fetcher.FetchPosts, uploader.NewUploader) in Go is non-trivial without changing them to overrideable vars or using interfaces. To avoid altering production code for tests, we tested behavior using indirect methods and log capture.
- Testing fatal errors like log.Fatal() required careful use of defer/recover to prevent test crashes.

## 💡 Improvements If Given More Time
- Use proper interfaces for better mocking and testing
- Add retry/backoff for failed uploads or fetches
- Add structured logging and error tracking
- Add Swagger/OpenAPI documentation for the HTTP API
- Add metrics/health endpoints for observability

## 🙌 Author
- Anshul Agrawal
- GitHub: [@aanshul543](https://github.com/aanshul543)
