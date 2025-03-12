# NBA Statistics Backend

## Assignment: A Comprehensive Architecture Guide

### Problem Analysis and System Requirements
The assignment asks for a backend system to log NBA player statistics and calculate aggregate metrics, with specific technical constraints. Before diving into architecture decisions, it's essential to fully understand the problem space and requirements. The system needs to track various player statistics (points, rebounds, assists, etc.) and calculate real-time aggregations at both player and team levels.

The core technical constraints include implementation in either Go or Java without frameworks, containerization for easy deployment, and usage of a relational database. Beyond these basics, the system must handle high concurrency, provide real-time data availability, and maintain solid architectural principles for long-term maintainability1. These requirements point to a system that balances performance with clean architecture while remaining extensible for future enhancements.

Understanding the data flow is crucial: external systems will submit player statistics in machine-readable formats, and the system must validate, store, and make this data immediately available for aggregation. This suggests a design that efficiently handles write operations while optimizing for the read patterns needed for aggregation queries.

## Technology Stack Selection

### Programming Language: Go (Golang)

For this assignment, I would select Go over Java for several compelling reasons. Go's lightweight concurrency model with goroutines provides an efficient approach to handling multiple simultaneous requests without the overhead of Java threads3. This aligns perfectly with the assignment's scalability requirements. Additionally, Go's simplicity and lack of reliance on frameworks makes it ideal for this "framework-less" assignment, where clear architecture is valued over framework magic.

Go also offers strong type safety while maintaining relatively simple syntax, reducing development complexity compared to Java. Its efficient memory management and fast compilation times create a smoother development experience for this time-constrained assignment. Finally, Go's standard library is robust enough to build a complete HTTP service without external dependencies, which aligns with the assignment's minimal framework approach.

### Database: PostgreSQL

PostgreSQL stands out as the optimal database choice for this system. As a battle-tested relational database, it provides the structured data storage needed for player statistics while offering advanced features that benefit this specific application. PostgreSQL's transactional integrity ensures data consistency, which is critical when updating player statistics that must immediately reflect in aggregations.

PostgreSQL's strong query optimization capabilities will be particularly valuable for calculating aggregate statistics efficiently. Additionally, its JSON support enables flexible schema evolution while maintaining relational constraints - useful for potentially storing additional metadata about games or players in the future. PostgreSQL's concurrent read/write performance also helps meet the real-time data availability requirement without sacrificing throughput.

### Deployment: Docker Compose with Minikube

For deployment, Docker Compose provides the simplest containerization approach for local development and testing, while Minikube offers a path to Kubernetes deployment for production-like environments. This dual approach allows demonstration of both simple startup and more advanced orchestration capabilities.

Docker Compose will enable quick local setup with just the database and application service defined in a compose file. Meanwhile, Minikube setup allows showcasing the solution's scalability potential, with kompose helping convert Docker Compose configurations to Kubernetes resources for more advanced deployment scenarios.

## System Architecture Design

### Core Architectural Principles

The architecture follows clean architecture principles with clear separation of concerns, domain-driven design concepts, and an emphasis on testability. The system is designed around the following layers:

1. API Layer: Handles HTTP requests and responses

1. Service Layer: Contains business logic for player statistics processing

1. Repository Layer: Manages data persistence

1. Domain Layer: Defines core entities and value objects

This model normalizes data appropriately to minimize redundancy while enabling efficient querying for aggregation purposes. The relationships between entities (Player belongs to Team, PlayerGameStats connects Player to Game) provide a clean structure for data integrity while supporting the required queries.

### API Design

The API follows RESTful principles with clearly defined endpoints:

1. POST /api/v1/player-stats: Submit statistics for a player's game performance

1. GET /api/v1/player-stats/player/{playerId}: Retrieve season averages for a specific player

1. GET /api/v1/player-stats/team/{teamId}: Retrieve season averages for all players on a team

These endpoints provide clear semantics and follow standard REST practices for resource naming and HTTP method usage. Input validation ensures data integrity, with appropriate error responses following standard HTTP status codes.

### Concurrency and Scalability Approach

The system leverages Go's goroutines and channels for efficient concurrent processing. When processing multiple incoming player statistics simultaneously, each request is handled in a separate goroutine, allowing non-blocking operation. This prevents bottlenecks during high-volume submissions, such as during multiple simultaneous games.

For database access, a connection pool manages concurrent database operations efficiently, preventing connection exhaustion while maximizing throughput. The system also implements strategic caching of frequently accessed data (like team rosters and recent aggregation results) to reduce database load.

Load balancing capabilities are built into the architecture from the start, allowing the system to scale horizontally by deploying multiple instances behind a load balancer. This approach, combined with database connection pooling, enables the system to handle significant concurrent load efficiently.

## Implementation Strategy

### Project Structure

The project follows a clean, maintainable structure that clearly separates concerns:

```
nba-stats/
├── cmd/
│   └── server/
│       └── main.go         # Application entry point
├── internal/
│   ├── api/                # HTTP handlers
│   ├── domain/             # Domain entities and value objects
│   ├── repository/         # Data access layer
│   └── service/            # Business logic
├── pkg/                    # Reusable packages
│   ├── validator/          # Input validation
│   ├── logger/             # Logging utilities
│   └── errors/             # Error handling
├── scripts/                # Build and deployment scripts
├── docker-compose.yml      # Local deployment configuration
└── Dockerfile              # Container definition
```

This structure follows Go best practices by separating internal application code from reusable packages, making the codebase navigable and maintainable.

### Error Handling Strategy

The system implements comprehensive error handling with custom error types that provide context-aware error messages. Errors are categorized into different types (validation errors, database errors, not found errors, etc.) allowing appropriate HTTP status codes to be returned based on error type.

All errors are logged with sufficient context for debugging, without exposing sensitive information in responses. This approach ensures robust error handling while maintaining good security practices and providing meaningful feedback to API consumers.

### Validation Implementation

Input validation occurs in multiple layers:

1. Request validation for JSON schema and type correctness

2. Domain validation for business rules (e.g., minutes played must be between 0 and 48.0)
   
3. Data consistency validation (e.g., player must belong to the specified team)

This multi-layered validation strategy ensures data integrity while providing clear error messages that pinpoint exactly what's wrong with a request.

## Testing Approach

### Unit Testing Strategy

The system employs comprehensive unit testing with a focus on testing each component in isolation. The use of interfaces throughout the codebase enables mock implementations for dependencies, allowing true unit tests that don't rely on external systems like databases.

Test cases cover positive scenarios, edge cases, and error conditions to ensure robust behavior. For Go implementation, the standard testing package is used along with testify for assertions, maintaining the framework-less requirement while providing sufficient testing capabilities.

### Integration Testing

Integration tests verify the interaction between components, particularly focusing on the API layer's interaction with the service and repository layers. These tests use in-memory database implementations or containerized test databases to provide realistic but controlled testing environments.

Additionally, end-to-end tests using the API endpoints verify complete workflows from data submission to aggregation retrieval, ensuring the system functions correctly as a whole.

### Performance Testing

Performance testing is critical for a system with scalability requirements. The implementation includes benchmark tests that measure throughput and response times under various load conditions. These tests specifically target:

1. Concurrent submission of player statistics

1. Real-time calculation of aggregate statistics

1. Database query performance

The results of these tests inform optimization decisions and provide confidence in the system's ability to handle the expected load.

## Deployment and Operations

### Docker Setup

The Docker configuration includes separate containers for the application and database, with proper health checks and restart policies. The Dockerfile follows best practices with multi-stage builds to minimize image size and security vulnerabilities:

1. First stage uses the Go build environment to compile the application

1. Second stage uses a minimal base image with only the compiled binary and necessary runtime files

1. Application runs as a non-root user for enhanced security

This approach results in a small, secure container image that starts quickly and efficiently uses resources.

### Kubernetes/Minikube Configuration

For more advanced deployment scenarios, Kubernetes configurations demonstrate scalability capabilities. The setup includes:

1. Deployment configurations with appropriate resource limits and requests

1. Horizontal Pod Autoscaler (HPA) for automatic scaling based on CPU/memory usage

1. Service definitions for internal and external access

1. ConfigMaps and Secrets for environment-specific configuration

These Kubernetes resources can be generated using kompose from the Docker Compose file, then enhanced with additional Kubernetes-specific features.

### Monitoring and Observability

Though not explicitly required, the implementation includes essential observability features through structured logging and metrics exposure. Each service component emits logs in a consistent JSON format with correlation IDs to track requests across the system. Additionally, the application exposes metrics endpoints compatible with Prometheus for monitoring performance and health.

These observability features simplify troubleshooting and provide visibility into system performance, demonstrating a production-ready approach to operations.

### AWS Deployment Considerations

As requested in the assignment, the system is designed for easy deployment on AWS with several key components:

1. Amazon RDS for PostgreSQL provides a managed, scalable database service with automatic backups and high availability

1. Amazon ECS or EKS for container orchestration, depending on complexity requirements

1. Application Load Balancer (ALB) for distributing traffic across multiple application instances

1. CloudWatch for centralized logging and monitoring

1. AWS Secrets Manager for secure credential management

The application configuration is designed to be environment-aware, loading database credentials and other configuration from environment variables or AWS parameter store, making it adaptable to different AWS environments.

