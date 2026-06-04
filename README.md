# URL Shortener

A scalable backend REST API that shortens URLs with built-in caching, load balancing, and performance optimization. Built with Go, PostgreSQL, Nginx, and Docker.

Registered users can create short links, share them, and track click analytics. The system demonstrates practical SRE skills: horizontal scaling, load balancing, caching, and performance testing under load.

This project was developed as a backend engineering exercise focused on:
- Scalable API design and architecture
- Database optimization and caching strategies
- Concurrent request handling with goroutines
- Load balancing and horizontal scaling
- Performance testing and optimization

## Tech Stack

- **Go 1.25+** — main language
- **Gin** — HTTP framework
- **PostgreSQL** — primary database
- **Redis** — caching layer
- **Nginx** — load balancer
- **Docker + docker-compose** — containerization and orchestration
- **Locust** — load testing framework
- **JWT** — authentication

## Requirements

- [Go 1.25+](https://golang.org/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Locust](https://locust.io/) (optional, for load testing)

## Quick Start (Development)

**1. Clone the repo**
```bash
git clone https://github.com/TanyaKremnova/url-shortener.git
cd url-shortener
```

**2. Set up environment variables**
```bash
cp .env.example .env
```

Fill in `.env`:
```
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=url_shortener
POSTGRES_PORT=5432

DATABASE_URL=postgres://postgres:postgres@localhost:5432/url_shortener?sslmode=disable
APP_PORT=8080
JWT_SECRET=your_secret_key_here
APP_BASE_URL=http://localhost

CACHE_URL=redis:6379
```

**Security Note:** 
The values shown above are examples for local development only. 
In production, use strong random secrets and secure configuration management.

**3. Start all services (single command)**
```bash
docker-compose up
```

This starts:
- ✓ PostgreSQL database starts
- ✓ Redis cache starts
- ✓ 2 Go app instances start (containers: `url_shortener_app1`, `url_shortener_app2`)
- ✓ Nginx load balancer starts and distributes traffic

Server running at `http://localhost:8080`

**Verify load balancer is working:**
```bash
# Check health through load balancer (should hit both instances)
curl http://localhost:8080/health

# Make multiple requests - Nginx distributes across both instances
for i in {1..10}; do curl http://localhost:8080/health; done
```

**View container logs:**
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f nginx
docker-compose logs -f app1
docker-compose logs -f app2
```

**Stop all services:**
```bash
docker-compose down
```

## Performance & Scalability Features

### Caching Layer
- **Redirect endpoint** reads from Redis cache first, reducing database queries
- **Implementation:** Goroutines handle cache updates asynchronously
- **Benefit:** Redirect requests (most frequent) complete without waiting for click tracking updates

### Concurrent Request Handling
- **Goroutines** process click counter updates asynchronously
- **Benefit:** Improves response time and throughput
- **Trade-off:** Eventual consistency for click counts (acceptable for analytics)

### Load Balancing
- **Nginx** distributes incoming requests across 2 Go instances
- **Benefit:** Horizontal scaling, fault tolerance, improved throughput
- **Verification:** See load test results below

### Database Optimization
- Efficient schema design with proper indexing
- Connection pooling through database driver
- Optimized queries to minimize latency

## Load Testing (Locust)

Demonstrates performance under realistic load.

### Prerequisites
```bash
pip install locust
```

### Run Load Tests
```bash
# Basic test: 100 concurrent users, 50 spawn rate
locust -f locustfile.py --host=http://localhost:8080
```

### Test Scenarios
1. **Redirect (most frequent operation)**
   - Expected: 100+ requests/sec per instance
   - Metric: Response time < 50ms (cache hit)

2. **Create Short URL**
   - Expected: 50+ requests/sec per instance
   - Metric: Response time < 200ms

3. **Mixed Workload**
   - 80% redirects, 15% create URLs, 5% stats queries
   - Expected: system handles load gracefully
   - Metric: No 5xx errors under realistic load

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | No | Health check / load balancer probe |
| `POST` | `/auth/register` | No | Create account |
| `POST` | `/auth/login` | No | Login, get JWT token |
| `POST` | `/urls/` | Yes | Create short URL |
| `GET` | `/:code` | No | Redirect to original URL (cached) |
| `GET` | `/admin/urls/stats` | Yes | View your URLs and click analytics |

**Authentication:** Protected routes require header `Authorization: Bearer <token>`

## Error Responses

All errors follow the same shape:
```json
{ "error": "message here", "code": 400 }
```

| Code | Meaning |
|------|---------|
| `400` | Bad request — invalid input |
| `401` | Unauthorized — missing or invalid token |
| `404` | Not found — short code does not exist |
| `409` | Conflict — email already registered |
| `500` | Server error |

## Architecture Decisions

### Why Goroutines for Click Tracking?
- Allows redirect endpoint to respond immediately
- Click count updates happen asynchronously
- Improves perceived performance and throughput
- Trade-off: eventual consistency (acceptable for analytics)

### Why Caching?
- Redirect is the most frequent operation (80%+ of traffic)
- Cache hit (< 10ms) vs database query (50-100ms)
- Reduces database load significantly
- Redis provides distributed caching for multi-instance setup

### Why Nginx Load Balancer?
- Distributes traffic evenly across instances
- Single point of failure prevention
- Easy horizontal scaling (add more Go instances)
- Standard production practice

## Learning Outcomes

### This project demonstrates:  
✓ **Scalability:** Horizontal scaling with load balancing  
✓ **Reliability:** Concurrent request handling, graceful error handling  
✓ **Performance:** Caching strategies, asynchronous operations, optimization  
✓ **Operations:** Docker containerization, load testing, monitoring  
✓ **Backend Fundamentals:** REST API design, authentication, database design

## Author

Tetiana Kremnova
- GitHub: [github.com/TanyaKremnova](https://github.com/TanyaKremnova)
- LinkedIn: [linkedin.com/in/tetianakremnova](https://linkedin.com/in/tetianakremnova)