# Task List: Je Suis Ici ActivityPub Platform Refactoring

Based on the PRD: `prd-refactor-and-enhancement.md`

## Relevant Files

- `internal/api/handlers/*.go` - All handler files need documentation and error handling updates
- `internal/services/*.go` - Service layer files need comments and refactoring
- `internal/db/models/*.go` - Model files need documentation and validation improvements
- `internal/config/config.go` - Configuration management centralization
- `internal/utils/` - New utility package for shared functions
- `internal/middleware/` - New middleware package for logging and error handling
- `docs/api/` - OpenAPI specification files
- `docs/architecture/` - Architecture documentation and diagrams
- `internal/db/migrations/` - New migration files for social graph features
- `internal/api/handlers/follow.go` - New follow/unfollow handler
- `internal/services/follow.go` - New follow service implementation
- `internal/api/handlers/feed.go` - Enhanced feed handler
- `internal/services/media.go` - Enhanced media service for multiple uploads
- `tests/integration/` - Integration test files
- `tests/fixtures/` - Test data fixtures
- `postman/` - Postman collection files
- `internal/messaging/` - RabbitMQ integration package
- `internal/metrics/` - Prometheus metrics package
- `docker/grafana/` - Grafana dashboard configurations
- `Dockerfile` - Optimized container configuration
- `docker-compose.yml` - Development environment setup

### Notes

- Integration tests should be placed in a dedicated `tests/` directory with clear naming conventions
- Use `go test ./...` to run all tests, or `go test ./internal/...` for specific packages
- Code coverage can be generated with `go test -cover ./...`
- Follow Go naming conventions for all new files and packages
- Use struct tags for validation in all model definitions

## Tasks

### Week 1: Code Documentation and Organization

- [ ] 1.0 Code Documentation and Standards Implementation
  - [ ] 1.0.1 Add godoc comments to all exported functions in `internal/api/handlers/`
  - [ ] 1.0.2 Add godoc comments to all exported functions in `internal/services/`
  - [ ] 1.0.3 Add godoc comments to all exported functions in `internal/db/models/`
  - [ ] 1.0.4 Create OpenAPI 3.0 specification file in `docs/api/openapi.yaml`
  - [ ] 1.0.5 Document all existing API endpoints with request/response schemas
  - [ ] 1.0.6 Create architecture documentation in `docs/architecture/system-overview.md`
  - [ ] 1.0.7 Create coding standards document in `docs/development/coding-standards.md`
  - [ ] 1.0.8 Update README.md with detailed setup instructions

- [ ] 1.1 Code Organization and Cleanup
  - [ ] 1.1.1 Audit all Go files for unused imports using `goimports`
  - [ ] 1.1.2 Remove dead code and commented-out functions
  - [ ] 1.1.3 Create `internal/utils/` package for shared utility functions
  - [ ] 1.1.4 Extract common validation logic into `internal/utils/validation.go`
  - [ ] 1.1.5 Extract common response formatting into `internal/utils/response.go`
  - [ ] 1.1.6 Refactor error handling patterns across all handlers
  - [ ] 1.1.7 Organize imports according to Go conventions (std, external, internal)
  - [ ] 1.1.8 Run `go mod tidy` to clean up dependencies

- [ ] 1.2 Database Layer Enhancement
  - [ ] 1.2.1 Add connection pooling configuration to `internal/db/cockroach.go`
  - [ ] 1.2.2 Implement transaction helper functions in `internal/db/transaction.go`
  - [ ] 1.2.3 Create database seeding script in `internal/db/seeds/`
  - [ ] 1.2.4 Add database health check endpoint `/health/db`
  - [ ] 1.2.5 Optimize existing queries with proper indexing
  - [ ] 1.2.6 Add query performance logging middleware
  - [ ] 1.2.7 Implement connection retry logic with exponential backoff
  - [ ] 1.2.8 Add database connection metrics collection

### Week 2: Configuration and Error Handling

- [ ] 2.0 Configuration Management Overhaul
  - [ ] 2.0.1 Centralize all configuration in `internal/config/config.go`
  - [ ] 2.0.2 Create environment-specific config files (`config/dev.yaml`, `config/prod.yaml`)
  - [ ] 2.0.3 Implement configuration validation on startup
  - [ ] 2.0.4 Add hot-reloading for development environment
  - [ ] 2.0.5 Document all configuration options in `docs/configuration.md`
  - [ ] 2.0.6 Add configuration struct tags for validation
  - [ ] 2.0.7 Implement secure configuration loading (env vars, secrets)
  - [ ] 2.0.8 Add configuration unit tests

- [ ] 2.1 Logging and Error Handling Standardization
  - [ ] 2.1.1 Implement structured logging with consistent fields across all components
  - [ ] 2.1.2 Create request/response logging middleware in `internal/middleware/logging.go`
  - [ ] 2.1.3 Implement centralized error handling in `internal/middleware/error.go`
  - [ ] 2.1.4 Add panic recovery middleware with graceful responses
  - [ ] 2.1.5 Create error wrapping utilities in `internal/utils/errors.go`
  - [ ] 2.1.6 Implement log rotation and retention policies
  - [ ] 2.1.7 Add request ID generation and correlation
  - [ ] 2.1.8 Standardize HTTP status codes and error response format

### Week 3: Social Graph Core Implementation

- [ ] 3.0 Social Graph Database and API Design
  - [ ] 3.0.1 Create database migration for enhanced followers table
  - [ ] 3.0.2 Create database migration for following table
  - [ ] 3.0.3 Create database migration for user blocks table
  - [ ] 3.0.4 Design follow/unfollow API endpoint schemas
  - [ ] 3.0.5 Create user search and discovery API endpoint schema
  - [ ] 3.0.6 Add privacy controls to user model (public/private/followers-only)
  - [ ] 3.0.7 Create user blocking and muting database models
  - [ ] 3.0.8 Add database indexes for social graph queries

- [ ] 3.1 Follow/Unfollow Functionality Implementation
  - [ ] 3.1.1 Implement follow model in `internal/db/models/follow.go`
  - [ ] 3.1.2 Implement follow service in `internal/services/follow.go`
  - [ ] 3.1.3 Create follow/unfollow API handlers in `internal/api/handlers/follow.go`
  - [ ] 3.1.4 Implement user search functionality with pagination
  - [ ] 3.1.5 Add follower/following count to user profile
  - [ ] 3.1.6 Implement blocking and muting functionality
  - [ ] 3.1.7 Add privacy checks to all user-related endpoints
  - [ ] 3.1.8 Update user model with social graph relationships

### Week 4: Media System and Feed Enhancement

- [ ] 4.0 Enhanced Media System Implementation
  - [ ] 4.0.1 Redesign media table schema for multiple files per check-in
  - [ ] 4.0.2 Update media model in `internal/db/models/media.go`
  - [ ] 4.0.3 Implement multiple file upload in media service
  - [ ] 4.0.4 Add media type validation (images, videos, file sizes)
  - [ ] 4.0.5 Implement image compression and thumbnail generation
  - [ ] 4.0.6 Add media metadata extraction (EXIF, dimensions)
  - [ ] 4.0.7 Update media storage structure in MinIO
  - [ ] 4.0.8 Add media processing status tracking

- [ ] 4.1 Feed and Timeline System Development
  - [ ] 4.1.1 Design aggregated timeline database schema
  - [ ] 4.1.2 Implement timeline generation algorithm (chronological)
  - [ ] 4.1.3 Create feed API endpoints in `internal/api/handlers/feed.go`
  - [ ] 4.1.4 Implement cursor-based pagination for feeds
  - [ ] 4.1.5 Add feed caching mechanism (in-memory or Redis prep)
  - [ ] 4.1.6 Implement feed filtering by content type
  - [ ] 4.1.7 Add real-time feed update preparation
  - [ ] 4.1.8 Optimize feed queries with proper indexing

- [ ] 4.2 API Consistency and Validation
  - [ ] 4.2.1 Standardize JSON response structure across all endpoints
  - [ ] 4.2.2 Implement comprehensive input validation using struct tags
  - [ ] 4.2.3 Add API versioning strategy (/api/v1/)
  - [ ] 4.2.4 Create consistent pagination patterns for all list endpoints
  - [ ] 4.2.5 Implement rate limiting middleware
  - [ ] 4.2.6 Add API usage throttling and quotas
  - [ ] 4.2.7 Standardize HTTP status codes usage
  - [ ] 4.2.8 Update OpenAPI specification with new endpoints

### Week 5: Testing Infrastructure Setup

- [ ] 5.0 Integration Testing Framework Setup
  - [ ] 5.0.1 Set up test database with automated schema creation
  - [ ] 5.0.2 Create test configuration in `tests/config/test.yaml`
  - [ ] 5.0.3 Implement test fixtures in `tests/fixtures/`
  - [ ] 5.0.4 Create test utilities for user creation and authentication
  - [ ] 5.0.5 Implement database transaction rollback for test isolation
  - [ ] 5.0.6 Set up test data factories for consistent object generation
  - [ ] 5.0.7 Create test helper functions in `tests/helpers/`
  - [ ] 5.0.8 Add test environment setup and teardown scripts

- [ ] 5.1 API Testing Implementation
  - [ ] 5.1.1 Write integration tests for authentication endpoints (`tests/integration/auth_test.go`)
  - [ ] 5.1.2 Write integration tests for user management endpoints (`tests/integration/user_test.go`)
  - [ ] 5.1.3 Write integration tests for check-in endpoints (`tests/integration/checkin_test.go`)
  - [ ] 5.1.4 Write integration tests for media upload endpoints (`tests/integration/media_test.go`)
  - [ ] 5.1.5 Write integration tests for follow/unfollow endpoints (`tests/integration/follow_test.go`)
  - [ ] 5.1.6 Write integration tests for feed endpoints (`tests/integration/feed_test.go`)
  - [ ] 5.1.7 Write integration tests for ActivityPub endpoints (`tests/integration/activitypub_test.go`)
  - [ ] 5.1.8 Add error scenario testing for all endpoints

### Week 6: Testing Coverage and Quality

- [ ] 6.0 Postman Collection and Performance Testing
  - [ ] 6.0.1 Create comprehensive Postman collection in `postman/je-suis-ici.postman_collection.json`
  - [ ] 6.0.2 Add environment variables for local, staging, production
  - [ ] 6.0.3 Implement automated test scripts within Postman
  - [ ] 6.0.4 Create test data setup scripts for Postman
  - [ ] 6.0.5 Document test execution procedures
  - [ ] 6.0.6 Create load testing scripts using k6 or Apache Bench
  - [ ] 6.0.7 Implement database performance tests
  - [ ] 6.0.8 Add memory and CPU profiling tests

- [ ] 6.1 Code Coverage and Quality Gates
  - [ ] 6.1.1 Set up code coverage reporting with `go test -cover`
  - [ ] 6.1.2 Configure automated coverage checks in CI pipeline
  - [ ] 6.1.3 Add static code analysis with golangci-lint
  - [ ] 6.1.4 Create code quality gates for minimum coverage thresholds
  - [ ] 6.1.5 Establish code review checklist
  - [ ] 6.1.6 Add benchmark tests for critical functions
  - [ ] 6.1.7 Implement performance baseline testing
  - [ ] 6.1.8 Add security testing (input validation, auth bypass)

### Week 7: Message Queue and Metrics Integration

- [ ] 7.0 RabbitMQ Integration Implementation
  - [ ] 7.0.1 Add RabbitMQ client library to go.mod
  - [ ] 7.0.2 Create RabbitMQ configuration in `internal/config/`
  - [ ] 7.0.3 Implement connection manager in `internal/messaging/rabbitmq.go`
  - [ ] 7.0.4 Add RabbitMQ health checks
  - [ ] 7.0.5 Create message broker interface for testability
  - [ ] 7.0.6 Define event types and payload structures
  - [ ] 7.0.7 Implement event publisher service
  - [ ] 7.0.8 Create event consumers for notifications and feed updates

- [ ] 7.1 Prometheus Metrics and Monitoring Setup
  - [ ] 7.1.1 Add Prometheus client library to project
  - [ ] 7.1.2 Create metrics middleware in `internal/middleware/metrics.go`
  - [ ] 7.1.3 Implement business metrics collectors
  - [ ] 7.1.4 Add database metrics for connection pool and query performance
  - [ ] 7.1.5 Create custom metrics for ActivityPub federation
  - [ ] 7.1.6 Implement health check endpoints for all dependencies
  - [ ] 7.1.7 Add resource utilization monitoring
  - [ ] 7.1.8 Create metrics exposition endpoint `/metrics`

### Week 8: Grafana Dashboards and Deployment

- [ ] 8.0 Grafana Dashboard Implementation
  - [ ] 8.0.1 Install and configure Grafana with Prometheus data source
  - [ ] 8.0.2 Create API performance dashboard (`docker/grafana/dashboards/api-performance.json`)
  - [ ] 8.0.3 Build database performance dashboard (`docker/grafana/dashboards/database.json`)
  - [ ] 8.0.4 Implement business metrics dashboard (`docker/grafana/dashboards/business-metrics.json`)
  - [ ] 8.0.5 Create system resources dashboard (`docker/grafana/dashboards/system-resources.json`)
  - [ ] 8.0.6 Build RabbitMQ monitoring dashboard (`docker/grafana/dashboards/rabbitmq.json`)
  - [ ] 8.0.7 Configure alerting rules for critical failures
  - [ ] 8.0.8 Set up alert channels (email, Slack integration)

- [ ] 8.1 Container and Deployment Preparation
  - [ ] 8.1.1 Create optimized Dockerfile with multi-stage builds
  - [ ] 8.1.2 Add docker-compose.yml for local development environment
  - [ ] 8.1.3 Implement graceful shutdown handling
  - [ ] 8.1.4 Create environment configuration for container orchestration
  - [ ] 8.1.5 Add deployment readiness and liveness probes
  - [ ] 8.1.6 Create deployment documentation
  - [ ] 8.1.7 Add container security scanning configuration
  - [ ] 8.1.8 Document monitoring and alerting setup procedures 