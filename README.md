# Company Management System

## About

Company Management System is a Golang microservice application for managing company entities. It provides CRUD operations for companies. And also enables user registration and authentication. 

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Docker
- Postman application

### Running with Docker Compose

- To create build images and create containers for CompanyService, UserService, MongoDB and Kafka navigate into the /company-management-system folder and run:
```
docker-compose up
```

### API and testing
- For API testing postman collection is provided, some of the tests are executed as a part of GitHub Actions and linter.
After running the application, import postman collection and test endpoints. After registration and login use generated token in Auth header to perform Create, Patch and Delete operations.

### Future steps
- Improve integration tests. Create test database and test Kafka broker, currently storage is mocked with interface.
- Create consumer service to consume Kafka events and to store them into the MongoDB or another database as a replica for state restoration or data analytics.
- Assign roles to the users and add role based authentication. Divide user service to 2 services, authentication and user.
- Secure connections.
- Add security.
- Improve error handling with more custom errors.