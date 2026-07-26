## Technology
1. Backend: Go v1.25.12
2. Database: Postgres 14.7
3. Docker (optional)

## Project Structure
- `cmd/`: This directory contains application-specific entrypoints. It's the main program of the application.

- `domain/`: This directory holds struct object representation for each table in the database.

- `domain/enums`: This directory holds enums constants for each domain.

- `entrypoint/`: This directory is the controller layer for API endpoints. The controller layer's function is to accept and validate incoming requests before they are processed by the service layer.

- `migration/`: This directory holds SQL migrations files to create and modify tables in the database.

- `repository/`: This directory is a repository layer to handle interactions between the application and database.

- `request/`: This directory holds HTTP request structures that define the input of the application.

- `response/`: This directory contains HTTP response structures that define the output formats of the API.

- `service/`: This directory implements the core business logic layer, processing data between the controller and repository layers.

- `shared/`: This directory contains external and internal services which the application interacts with.

- `util/`: This directory holds utility functions and helper code that supports the main application functionality.

## Author
Chris Christian
