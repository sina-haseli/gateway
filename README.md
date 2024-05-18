# Project Documentation

## Overview

This project consists of a gateway service and two mock services: `service1_mock` (HTTP) and `service2_mock` (gRPC). The gateway service provide load balancing, and circuit breaker functionalities.

## Running the Project

### Prerequisites

- Docker
- Docker Compose

### Steps to Run the Project

1. **Clone the repository**:

    ```sh
    git clone https://github.com/your-repo/myproject.git
    cd myproject
    ```

2. **Build the Docker images**:

    ```sh
    docker-compose build
    ```

3. **Run the containers**:

    ```sh
    docker-compose up
    ```

### Accessing the Services

- **Gateway Service**: `http://localhost:9090`
- **Service1 Mock Instances**:
    - `http://localhost:8081`
    - `http://localhost:8082`
    - `http://localhost:8083`
- **Service2 Mock (gRPC)**: `localhost:50051`

## Gateway Service

The gateway service serves as the entry point for all incoming requests. It performs the following key functions:

### Load Balancer

The load balancer distributes incoming HTTP requests across multiple instances of `service1_mock` using a round-robin algorithm. This ensures an even distribution of traffic, improving fault tolerance and scalability. Each instance of `service1_mock` receives an equal share of requests, maximizing resource utilization and minimizing response time.

### Circuit Breaker

The circuit breaker pattern is implemented to handle failures gracefully. If a backend service fails repeatedly, the circuit breaker trips, and the gateway stops sending requests to the failing service for a certain period, preventing further strain on the service and allowing it time to recover.

## Handling gRPC Services

The gateway also handles gRPC requests. It communicates with `service2_mock` using gRPC, leveraging the load balancer and circuit breaker functionalities to manage gRPC service interactions.

## Summary

This project demonstrates how to set up a gateway service with advanced functionalities like load balancing, reverse proxying, and circuit breaking. It also showcases handling both HTTP and gRPC services. By using Docker and Docker Compose, the services are containerized, making it easy to build, run, and scale the application.
