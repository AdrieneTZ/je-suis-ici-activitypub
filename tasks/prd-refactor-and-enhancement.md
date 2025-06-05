# Project Understanding and Refactoring Plan: Je Suis Ici ActivityPub Platform

## 1. Project Overview

### 1.1 Purpose
This project is the backend of a decentralized social platform called "Je Suis Ici", designed to comply with the ActivityPub protocol. It serves as the API provider for frontend developers and implements core social features such as user management, location-based check-ins, and media handling with federation capabilities.

### 1.2 Current Architecture Analysis

**Technology Stack:**
- **Language:** Go 1.24.2
- **Database:** CockroachDB (PostgreSQL compatible)
- **Object Storage:** MinIO
- **Authentication:** JWT with Chi Router
- **Distributed Tracing:** Jaeger
- **Protocol:** ActivityPub for federation

**Current Directory Structure:**
```
je-suis-ici-activitypub/
├── cmd/server/              # Application entry point
├── internal/                # Core business logic
│   ├── activitypub/        # ActivityPub protocol implementation
│   ├── api/                # API layer (handlers, middlewares, router)
│   ├── config/             # Configuration management
│   ├── db/                 # Database layer (models, migrations)
│   ├── services/           # Business logic services
│   ├── storage/            # File storage services
│   └── tracing/            # Distributed tracing setup
├── checkin-image/          # Media storage directory
├── go.mod & go.sum         # Dependency management
└── README.md               # Project documentation
```

### 1.3 Current Features Analysis

**Implemented Features:**
- User registration/authentication with JWT
- Location-based check-in functionality
- Single media upload per check-in
- Basic ActivityPub actor management
- Federation with other ActivityPub platforms
- WebFinger and NodeInfo support
- Database migrations
- Distributed tracing with Jaeger

**Database Schema (Current):**
- Users table (with ActivityPub actor fields)
- Checkins table (with location data)
- Media table (single file per checkin)
- Activities table (ActivityPub activities)
- Followers table (basic follow relationships)

---

## 2. Feature Gaps and Improvement Opportunities

### 2.1 Critical Missing Features

**Social Graph Functionality:**
- No follow/unfollow API endpoints
- Missing feed aggregation for followed users
- No notification system for social interactions
- Lack of privacy controls (public/private posts)

**Media System Limitations:**
- Only supports single file upload per check-in
- No support for multiple images/videos
- Missing media compression/optimization
- No media metadata extraction

**API Documentation & Maintainability:**
- Insufficient API documentation (no OpenAPI/Swagger)
- Limited code comments and inline documentation
- No architecture diagrams or visual documentation
- Inconsistent error response formats

**Testing Infrastructure:**
- No integration tests
- Missing test fixtures and test data
- No API testing automation
- No code coverage reporting

### 2.2 Infrastructure Gaps

**Message Queue System:**
- No event-driven architecture for real-time updates
- Missing notification delivery system
- No background job processing

**Monitoring & Observability:**
- No metrics collection (Prometheus)
- No performance monitoring dashboards
- Limited logging and alerting capabilities

---

## 3. Detailed Refactoring Plan

### 3.1 Phase 1: Foundation and Code Quality (Weeks 1-2)

#### Step 1.1: Code Documentation and Standards
1. **Add comprehensive function/method comments** to all exported functions in handlers, services, and models
2. **Create API documentation** using OpenAPI 3.0 specification for all existing endpoints
3. **Write architecture documentation** describing system components and their interactions
4. **Establish coding standards** document with Go best practices and project-specific guidelines
5. **Create development setup guide** with step-by-step instructions for new developers

#### Step 1.2: Code Organization and Cleanup
1. **Audit existing codebase** for unused imports, dead code, and deprecated patterns
2. **Extract common utility functions** from handlers and services into shared packages
3. **Implement consistent error handling** across all API endpoints with standardized error responses
4. **Refactor repeated validation logic** into reusable validator functions
5. **Organize imports and dependencies** following Go conventions and remove unused dependencies

#### Step 1.3: Database Layer Enhancement
1. **Add database connection pooling** configuration and optimization
2. **Implement transaction management** patterns for complex operations
3. **Create database seeding scripts** for development and testing environments
4. **Add database health checks** and connection monitoring
5. **Optimize existing queries** and add proper indexing strategies

#### Step 1.4: Configuration Management
1. **Centralize configuration management** using structured config files
2. **Add environment-specific configurations** (dev, staging, production)
3. **Implement configuration validation** on application startup
4. **Add configuration hot-reloading** capabilities for development
5. **Document all configuration options** with examples and default values

#### Step 1.5: Logging and Error Handling
1. **Implement structured logging** throughout the application using consistent log levels
2. **Add request/response logging middleware** with configurable verbosity
3. **Create centralized error handling** with proper error wrapping and context
4. **Add panic recovery middleware** with graceful error responses
5. **Implement log rotation and retention** policies

### 3.2 Phase 2: Feature Enhancement and Expansion (Weeks 3-4)

#### Step 2.1: Social Graph Implementation
1. **Design follow/unfollow API endpoints** with proper request/response schemas
2. **Implement follower/following relationships** in database and service layers
3. **Create user search and discovery** functionality with pagination
4. **Add privacy controls** for user profiles and posts (public/private/followers-only)
5. **Implement blocking and muting** features for user safety

#### Step 2.2: Enhanced Media System
1. **Redesign media upload API** to support multiple files per check-in
2. **Implement media type validation** and file size restrictions
3. **Add image/video processing** pipeline with compression and thumbnail generation
4. **Create media metadata extraction** system (EXIF data, dimensions, duration)
5. **Implement media storage optimization** with CDN integration preparation

#### Step 2.3: Feed and Timeline System
1. **Design aggregated timeline API** showing posts from followed users
2. **Implement timeline generation** algorithms (chronological and relevance-based)
3. **Add real-time feed updates** preparation for WebSocket integration
4. **Create feed pagination** with cursor-based pagination for performance
5. **Implement feed caching strategies** for improved performance

#### Step 2.4: Enhanced ActivityPub Integration
1. **Extend ActivityPub implementation** to support additional activity types (Like, Share, Comment)
2. **Improve federation reliability** with retry mechanisms and error handling
3. **Add ActivityPub object validation** and security measures
4. **Implement proper ActivityPub Collections** for followers, following, and outbox
5. **Add support for ActivityPub Mentions** and notifications

#### Step 2.5: API Consistency and Validation
1. **Standardize all API responses** with consistent JSON structure and status codes
2. **Implement comprehensive input validation** for all endpoints using struct tags
3. **Add API versioning strategy** for future compatibility
4. **Create consistent pagination patterns** across all list endpoints
5. **Implement rate limiting** and API usage throttling

### 3.3 Phase 3: Testing and Quality Assurance (Weeks 5-6)

#### Step 3.1: Integration Testing Infrastructure
1. **Set up test database** with automated schema creation and cleanup
2. **Create test fixtures** and mock data generation for consistent testing
3. **Implement database transaction rollback** for isolated test execution
4. **Add test configuration management** separate from production configs
5. **Create testing utilities** for common test operations (user creation, authentication)

#### Step 3.2: API Testing Implementation
1. **Write integration tests** for all authentication endpoints (login, register, token refresh)
2. **Create comprehensive user API tests** covering CRUD operations and edge cases
3. **Implement check-in API tests** including media upload and location validation
4. **Add ActivityPub endpoint tests** for federation functionality
5. **Write error scenario tests** for all endpoints (invalid input, unauthorized access, not found)

#### Step 3.3: Postman Collection Development
1. **Create comprehensive Postman collection** with all API endpoints organized by feature
2. **Add environment variables** for different testing environments (local, staging)
3. **Implement automated test scripts** within Postman for response validation
4. **Create test data setup scripts** for preparing test scenarios
5. **Document test execution procedures** and expected outcomes

#### Step 3.4: Code Coverage and Quality
1. **Set up code coverage reporting** using Go's built-in coverage tools
2. **Implement automated coverage checks** in CI pipeline with minimum thresholds
3. **Add static code analysis** using golangci-lint with comprehensive rule sets
4. **Create code quality gates** preventing merges below quality standards
5. **Establish code review checklist** ensuring consistent quality standards

#### Step 3.5: Performance Testing
1. **Create load testing scripts** using tools like k6 or Apache Bench
2. **Implement database performance tests** measuring query execution times
3. **Add memory and CPU profiling** capabilities for performance debugging
4. **Create benchmark tests** for critical business logic functions
5. **Establish performance baselines** and regression testing procedures

### 3.4 Phase 4: Infrastructure and Observability (Weeks 7-8)

#### Step 4.1: RabbitMQ Integration
1. **Add RabbitMQ client library** and connection management to the project
2. **Design event schema** for follow/unfollow, check-in creation, and media upload events
3. **Implement event publishers** in relevant service layers (user service, checkin service)
4. **Create event consumers** for handling notifications and feed updates
5. **Add dead letter queues** and retry mechanisms for failed message processing

#### Step 4.2: Metrics and Monitoring Setup
1. **Integrate Prometheus metrics** collection for API endpoints, database operations, and business metrics
2. **Add custom metrics** for user activity, check-in frequency, and federation statistics
3. **Implement health check endpoints** for application and dependency monitoring
4. **Create metrics middleware** for automatic HTTP request/response metrics collection
5. **Add resource utilization monitoring** (memory, CPU, disk, network)

#### Step 4.3: Grafana Dashboard Implementation
1. **Install and configure Grafana** with Prometheus as primary data source
2. **Create API performance dashboard** showing request rates, response times, and error rates
3. **Build database performance dashboard** with query performance and connection pool metrics
4. **Implement business metrics dashboard** showing user growth, check-in activity, and federation health
5. **Add alerting rules** for critical system failures and performance degradation

#### Step 4.4: Centralized Logging
1. **Enhance structured logging** with consistent fields across all components
2. **Add log aggregation** preparation for centralized logging systems (ELK stack)
3. **Implement log correlation** with request IDs and user context
4. **Create log-based alerting** for error patterns and anomalies
5. **Add log retention policies** and storage optimization

#### Step 4.5: Container and Deployment Preparation
1. **Create optimized Dockerfile** with multi-stage builds and minimal base images
2. **Add docker-compose configuration** for local development environment
3. **Implement graceful shutdown** handling for containerized deployments
4. **Create environment configuration** for container orchestration
5. **Add deployment readiness** and liveness probes

---

## 4. RabbitMQ Integration Detailed Steps

### 4.1 Installation and Basic Setup
1. **Add RabbitMQ dependencies** to go.mod (streadway/amqp or rabbitmq/amqp091-go)
2. **Create RabbitMQ configuration** in internal/config with connection parameters
3. **Implement connection manager** with automatic reconnection and error handling
4. **Add RabbitMQ health checks** to application monitoring
5. **Create message broker interface** for testability and future flexibility

### 4.2 Event System Design
1. **Define event types** for user actions (follow, unfollow, checkin creation, media upload)
2. **Create event payload structures** with versioning support for future compatibility
3. **Implement event serialization** using JSON with schema validation
4. **Design exchange and queue topology** with proper routing for different event types
5. **Add event metadata** (timestamp, source, correlation ID) for tracing

### 4.3 Publisher Implementation
1. **Create event publisher service** with transaction support for reliable publishing
2. **Implement follow/unfollow event publishing** in user service operations
3. **Add check-in creation events** for real-time feed updates
4. **Implement media upload events** for background processing triggers
5. **Add publisher error handling** with local event storage for retry scenarios

### 4.4 Consumer Implementation
1. **Create notification consumer** for delivering follow notifications to users
2. **Implement feed update consumer** for maintaining user timelines
3. **Add media processing consumer** for background image/video processing
4. **Create federation consumer** for ActivityPub outbound message processing
5. **Implement consumer scaling** with worker pool patterns

### 4.5 Monitoring and Management
1. **Add RabbitMQ metrics collection** for queue lengths, message rates, and consumer performance
2. **Implement message dead letter handling** with retry policies and failure analysis
3. **Create RabbitMQ dashboard** in Grafana for queue monitoring
4. **Add consumer lag monitoring** and alerting for processing delays
5. **Implement message replay** capabilities for system recovery scenarios

---

## 5. Grafana Integration Detailed Steps

### 5.1 Prometheus Metrics Implementation
1. **Add Prometheus client library** to the project dependencies
2. **Create metrics middleware** for HTTP requests with automatic labeling
3. **Implement business metrics collectors** for user registrations, check-ins, and follows
4. **Add database metrics** for connection pool usage and query performance
5. **Create custom metrics** for ActivityPub federation health and message queue performance

### 5.2 Dashboard Design and Creation
1. **Create API Performance Dashboard** with panels for:
   - Request rate by endpoint
   - Response time percentiles (P50, P95, P99)
   - Error rate by status code
   - Concurrent users and sessions
   - API usage by user agent/client

2. **Build Database Performance Dashboard** with panels for:
   - Query execution time by operation type
   - Connection pool utilization
   - Transaction success/failure rates
   - Database lock contention
   - Storage usage and growth trends

3. **Implement Business Metrics Dashboard** with panels for:
   - Daily/weekly active users
   - Check-in creation rates
   - Media upload statistics
   - Follow/unfollow activity
   - Federation message success rates

4. **Create System Resources Dashboard** with panels for:
   - CPU and memory utilization
   - Disk I/O and storage usage
   - Network traffic and bandwidth
   - Container resource usage
   - Application uptime and restarts

5. **Build RabbitMQ Monitoring Dashboard** with panels for:
   - Queue message rates and depths
   - Consumer performance and lag
   - Message acknowledgment rates
   - Dead letter queue activity
   - Connection and channel statistics

### 5.3 Alerting Configuration
1. **Set up critical alerts** for system failures (database down, high error rates)
2. **Create performance alerts** for response time degradation and resource exhaustion
3. **Implement business alerts** for unusual activity patterns or security concerns
4. **Add integration alerts** for RabbitMQ and external service failures
5. **Configure alert channels** (email, Slack, PagerDuty) with appropriate escalation

### 5.4 Data Source Configuration
1. **Configure Prometheus data source** with proper authentication and retention settings
2. **Add Jaeger data source** for distributed tracing correlation
3. **Set up log data source** if using centralized logging (Loki, Elasticsearch)
4. **Configure alert manager** for notification routing and silencing rules
5. **Implement data source health monitoring** and failover capabilities

### 5.5 Advanced Visualization Features
1. **Create template variables** for dynamic dashboard filtering (environment, user, endpoint)
2. **Implement drill-down capabilities** from high-level metrics to detailed traces
3. **Add annotation support** for deployment and incident correlation
4. **Create shared dashboard templates** for consistent visualization across teams
5. **Implement dashboard version control** and automated backup procedures

---

## 6. Comprehensive Testing Strategy

### 6.1 Integration Testing Framework
1. **Set up test environment** with isolated database and message queue instances
2. **Create test data factories** for generating consistent test objects (users, check-ins, media)
3. **Implement test authentication** helpers for simulating different user roles and permissions
4. **Add database transaction management** for test isolation and cleanup
5. **Create test configuration** separate from production with appropriate test-specific settings

### 6.2 API Endpoint Testing
1. **Authentication API Tests:**
   - User registration with valid/invalid data
   - Login with correct/incorrect credentials
   - JWT token validation and expiration
   - Password reset functionality
   - Account activation and deactivation

2. **User Management API Tests:**
   - Profile creation and updates
   - User search and discovery
   - Privacy settings management
   - Account deletion and data cleanup
   - User avatar upload and management

3. **Check-in API Tests:**
   - Check-in creation with location validation
   - Multiple media upload testing
   - Check-in listing with pagination
   - Location-based check-in filtering
   - Check-in privacy and visibility controls

4. **Social Features API Tests:**
   - Follow/unfollow operations
   - Timeline generation and filtering
   - Notification creation and delivery
   - User blocking and muting
   - Feed aggregation and caching

5. **ActivityPub API Tests:**
   - Actor creation and management
   - Activity publishing and consumption
   - Federation message handling
   - WebFinger service responses
   - NodeInfo service functionality

### 6.3 Postman Collection Development
1. **Create environment configurations** for local, staging, and production testing
2. **Implement authentication flows** with automatic token management and refresh
3. **Add pre-request scripts** for test data setup and variable management
4. **Create comprehensive test scenarios** covering happy paths and error conditions
5. **Implement automated test execution** with Newman for CI/CD integration

### 6.4 Performance and Load Testing
1. **Create baseline performance tests** for all API endpoints measuring response times
2. **Implement load testing scenarios** simulating realistic user traffic patterns
3. **Add database performance tests** measuring query execution under various loads
4. **Create memory and CPU profiling** tests for identifying performance bottlenecks
5. **Implement concurrent user testing** for validating system behavior under load

### 6.5 Security Testing
1. **Add authentication bypass tests** ensuring proper access control enforcement
2. **Implement input validation tests** for SQL injection and XSS prevention
3. **Create authorization tests** verifying proper permission enforcement
4. **Add rate limiting tests** ensuring API protection against abuse
5. **Implement data privacy tests** verifying proper data access restrictions

---

## 7. Implementation Timeline and Milestones

### Week 1-2: Foundation Phase
- **Milestone 1:** Complete code documentation and API specification
- **Milestone 2:** Implement consistent error handling and logging
- **Deliverables:** Updated README, API docs, coding standards document

### Week 3-4: Feature Enhancement Phase
- **Milestone 3:** Complete social graph implementation (follow/unfollow)
- **Milestone 4:** Enhanced media system with multiple file support
- **Deliverables:** New API endpoints, database migrations, updated documentation

### Week 5-6: Testing and Quality Phase
- **Milestone 5:** Complete integration testing suite
- **Milestone 6:** Postman collection and automated testing
- **Deliverables:** Test coverage reports, automated test execution

### Week 7-8: Infrastructure Phase
- **Milestone 7:** RabbitMQ integration and event system
- **Milestone 8:** Grafana dashboards and monitoring
- **Deliverables:** Monitoring dashboards, alerting configuration, deployment documentation

---

## 8. Success Metrics and KPIs

### Technical Metrics
- Code coverage: >80% for new code, >60% overall
- API response time: <200ms for 95th percentile
- System uptime: >99.9% availability
- Test execution time: <5 minutes for full test suite

### Quality Metrics
- Documentation coverage: 100% for public APIs
- Code review coverage: 100% for all changes
- Security vulnerability count: 0 critical, <5 medium
- Performance regression: 0 tolerance for degradation

### Business Metrics
- Developer onboarding time: <2 hours for new developers
- Feature delivery velocity: 20% improvement in development speed
- Bug report reduction: 50% fewer production issues
- API adoption: Improved developer experience metrics

---

## 9. Risk Assessment and Mitigation

### Technical Risks
1. **Database migration complexity** - Mitigate with thorough testing and rollback procedures
2. **RabbitMQ integration complexity** - Implement gradual rollout with feature flags
3. **Performance impact of new features** - Continuous monitoring and performance testing
4. **ActivityPub compatibility issues** - Extensive testing with other platforms
5. **Testing environment stability** - Implement robust test infrastructure

### Timeline Risks
1. **Scope creep** - Strict adherence to defined milestones and regular reviews
2. **Dependency delays** - Parallel development where possible and contingency planning
3. **Resource availability** - Cross-training and documentation for knowledge sharing
4. **Integration challenges** - Early proof-of-concept testing for critical integrations
5. **Quality gates** - Automated quality checks preventing low-quality merges

---

## 10. Conclusion

This comprehensive refactoring plan transforms the Je Suis Ici ActivityPub platform from a basic social backend into a robust, scalable, and maintainable system. The plan emphasizes:

1. **Maintainability** through comprehensive documentation, code organization, and testing
2. **Scalability** through proper architecture, monitoring, and event-driven design
3. **Quality** through automated testing, code coverage, and consistent development practices

The phased approach ensures minimal disruption to existing functionality while systematically improving all aspects of the system. Each phase builds upon the previous one, creating a solid foundation for future development and scaling.

By following this plan, the development team will have a production-ready social platform backend that can scale with user growth, integrate seamlessly with frontend applications, and maintain high code quality standards for long-term success. 