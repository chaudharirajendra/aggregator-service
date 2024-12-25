# Feed Service

This is a Feed Service application built in Go using the Gin web framework. It fetches data from external services (repair and session services), combines the results, and serves the data with sorting and pagination.

## Features
- Fetches repair and session data from external services.
- Combines, sorts, and paginates the data based on timestamps.
- Provides an API to get machine feeds.
- Configurable URLs for repair and session services via environment variables.
- Simple error handling and validation for input parameters.


## Installation

### Prerequisites

- Go (1.18+ recommended)
- Git

### Clone the Repository

```bash
git clone https://github.com/yourusername/feed-service.git
cd feed-service
```

## Set Up Environment Variables

Before running the application, set the following environment variables to provide the URLs for the repair and session services.

- REPAIR_SERVICE_URL: The URL for the repair service (e.g., http://localhost:8082/repairs).
- SESSION_SERVICE_URL: The URL for the session service (e.g., http://localhost:8083/sessions).

```bash
export REPAIR_SERVICE_URL="http://localhost:8082/repairs"
export SESSION_SERVICE_URL="http://localhost:8083/sessions"
export MACHINE_SERVICE_URL="http://localhost:8081/machine"
```

## Running the Application

go run main.go


## Running the Application

Endpoint: /machine-feeds/:machineId

This endpoint fetches machine feeds, combining data from repair and session services, sorting the data by timestamp, and providing paginated results.


### Query Parameters:

- machine_id (required): The unique ID of the machine whose feed you want to retrieve.
- page (optional): The page number for pagination (default: 1).
- size (optional): The number of items per page for pagination (default: 10).

### Example Request

GET http://localhost:8080/machine-feeds?machine_id=123&page=1&size=10












