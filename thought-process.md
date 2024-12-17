---

### **`thought-process.md` Content**

Here’s a high-level overview of the implementation approach and design considerations:

---

#### **Thought Process and Design Considerations**

**Objective**:  
Develop a scalable REST service that processes at least 10,000 requests per second, handles deduplication of IDs across multiple instances, and logs or streams the count of unique IDs.

---

### **1. Core Features**

#### **1.1 REST API Implementation**
- **Endpoint**: `/api/verve/accept`
    - Accepts:
        - `id` (integer, mandatory) for deduplication.
        - `endpoint` (string, optional) for triggering external HTTP requests.
    - **Response**:
        - Returns `"ok"` for successful processing.
        - Returns `"failed"` in case of errors.

#### **1.2 Request Deduplication**
- **Challenge**: Deduplication across multiple instances behind a load balancer.
- **Solution**:
    - Used **Redis** as a centralized store to track unique IDs.
    - Redis commands:
        - `SETNX` ensures atomic deduplication.
        - TTL of 1 minute ensures IDs expire after the time window.
    - Benefits:
        - Scalable across multiple instances.
        - Handles high concurrency efficiently.

#### **1.3 Logging Unique Counts**
- **Requirement**: Log the count of unique IDs every minute.
- **Implementation**:
    - Periodic logging using `time.Ticker` to trigger the task every minute.
    - Leveraged Redis to calculate and clear unique counts centrally.

---

### **2. Extensions**

#### **2.1 HTTP POST Instead of GET**
- **Requirement**: Fire a POST request instead of a GET when `endpoint` is provided.
- **Implementation**:
    - Payload structured as JSON:
      ```json
      {
          "timestamp": "2024-12-17T14:30:00Z",
          "unique_request_count": 100
      }
      ```
    - Used Go’s `http.Client` to send the POST request.
    - Logged HTTP response status codes for debugging.

#### **2.2 Load Balancer Deduplication**
- **Requirement**: Ensure deduplication across multiple service instances.
- **Solution**:
    - Centralized deduplication using Redis.
    - Redis provides:
        - Atomic operations for unique checks (`SETNX`).
        - TTL for automatic cleanup.
    - No dependency on individual instance memory or logs.

#### **2.3 Streaming Unique Counts**
- **Requirement**: Send unique counts to a distributed streaming service instead of a log file.
- **Solution**:
    - Integrated **Apache Kafka** for distributed streaming.
    - Kafka producer sends JSON messages with unique counts to a specific topic:
      ```json
      {
          "timestamp": "2024-12-17T14:30:00Z",
          "unique_request_count": 42
      }
      ```
    - Benefits:
        - Decouples processing from downstream systems.
        - Supports real-time data analysis and scaling.

---

### **3. Design Considerations**
1. **Scalability**:
    - Stateless REST service, allowing horizontal scaling.
    - Centralized deduplication and logging using Redis and Kafka.

2. **Performance**:
    - Optimized for high throughput using lightweight Redis operations.
    - Periodic logging minimizes impact on request handling.

3. **Error Handling**:
    - Graceful fallback for Redis or Kafka failures:
        - Log errors locally.
        - Continue processing other requests.

4. **Extensibility**:
    - Modular design for future feature additions:
        - Easily replace Kafka with other streaming services.
        - Enhance deduplication logic as needed.

---

### **Packaging the Solution**

#### **1. Source Code**
Host the complete source code on a Git repository (e.g., GitHub or Bitbucket). Ensure the repository has:
1. **Folders**:
    - `handlers/`: For HTTP handlers.
    - `services/`: For business logic like deduplication, logging, Kafka integration.
    - `log/`: Placeholder for log files.
2. **Files**:
    - `thought-process.md`: For describing the implementation.
    - Source files (`main.go`, etc.).

#### **2. Dockerfile**
Here’s a sample `Dockerfile` for your Go application:
```dockerfile
# Build stage
FROM golang:1.20 AS build

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o main .

# Production stage
FROM alpine:latest
WORKDIR /root/

COPY --from=build /app/main .
CMD ["./main"]
EXPOSE 8080
```

Build and run the Docker image:
```bash
docker build -t go-rest-service .
docker run -p 8080:8080 go-rest-service

