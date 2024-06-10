# IntelligenzGo - Hacker News Scraper

IntelligenzGo is a Golang-based web scraper designed to fetch and sort entries from the Hacker News API, as well as other news sources like Lobsters. The project is designed to be modular, adhering to best practices and principles such as SOLID and clean code. This README explains the thought process behind the development, key technical points, blocking points encountered, and how to test the code.

## Thought Process

### Project Goals

1. **Fetch and Sort Entries**: Scrape the first 30 entries from Hacker News and sort them based on title length and other criteria.
2. **Modular Design**: Ensure the scraper can be easily extended to support other news sources like Lobsters.
3. **Clean Code and SOLID Principles**: Write clean, maintainable code that adheres to SOLID principles.
4. **Automation**: Implement scripts for testing, building, and running the project efficiently.

### Key Technical Points

1. **Modularity**:
    - The project is structured into packages to separate concerns (e.g., `retriever`, `aggregator`, `items`).
    - Interfaces are used to abstract different scrapers (`ApiConnector`, `WebScrapperConnector`), making the code extensible and testable.

2. **Error Handling**:
    - Comprehensive error handling ensures that network failures and data inconsistencies do not crash the program.
    - Errors are propagated and logged appropriately to provide meaningful diagnostics.

3. **Sorting Logic**:
    - Entries with more than 5 words in their title are sorted by the number of comments.
    - Entries with shorter titles are sorted by points.

4. **Concurrency**:
    - The scrapers fetch data concurrently to improve performance (individual items data endpoints, different sources).

5. **Testing**:
    - Extensive unit and integration tests are written to ensure the functionality and reliability of the code.
    - Mocks are used for external API calls, web contents and retrievers aggregations to isolate tests.

### Blocking Points and Solutions

1. **Data Inconsistencies**:
    - Issue: Inconsistent data formats and missing fields in API responses.
    - Solution: Implement robust error handling and data validation to manage inconsistencies.

2. **Concurrency Issues**:
    - Issue: Handling concurrency safely while fetching and processing data.
    - Solution: Use goroutines and channels to manage concurrent tasks effectively.

## Project Structure

 * Main function: Creates a http server that handlers the following endpoints:
   * `hacker-news-items`: retrieves, sorts and return Hacker News items through it's API connections
   * `lobsters-items`: retrieves, sorts and return Lobsters items through Lobsters front web scrapping
   * `combine-sources-items`: Combines items fetched by all sources and returns sorted items
 * Services: Retrieving items interface and specific implementation for different sources
 * Data: Sources entities and responses


## Usage

### Prerequisites

- Go 1.22.2 or higher
- `make` installed

### Installation

Clone the repository:

```sh
git clone https://github.com/fburilloUCM/intelligenzGo.git
cd intelligenzGo
```

### Running the Scraper
To build and run the scraper:

```sh
make all
make run
```

### Calling endpoints 

* `hacker-news-items`:
  ```sh
  curl -s http://localhost:8080/hacker-news-items
  ```
* `lobsters-items`:
  ```sh
  curl -s http://localhost:8080/lobsters-items
  ```
* `combine-sources-items`
  ```sh
  curl -s http://localhost:8080/combine-sources-items
  ```

Responses for these calls will contain items sorted by required parameters and also will be logged in standard output the sorted list of items. 

## Testing

### Unit Testing
The project includes unit tests and integration tests to ensure comprehensive test coverage. Mocks are used to simulate external API calls.

Unit tests and integration tests can be executed using the go test command. The tests are located in the _test.go files within their respective packages.

```sh
go test ./...
````

Or to run using make:

```sh
make test
```

## Publishing

### Create container image

```sh
make build-container
```

### Test Coverage

To run the tests and generate a coverage report:

```sh
go test -v ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Coverage report generated in `coverage.html` document

## Automation workflows

## GitHub actions

### Triggering Events
The workflow is triggered on push to the main branch and feature branches, as well as on pull requests targeting the main branch.
### Build Job
1. Checkout code: Uses the actions/checkout@v2 action to clone the repository.
2. Set up: Install dependencies and build the project and run unit tests, also generate a coverage report.
3. Upload coverage report for use in later steps.
4. Coverage Check Job: Runs a script to check the coverage percentage. If the coverage is below 80%, the job fails.
5. Functional Test Job: Runs the Docker container with the built image, waits for the server to start, then uses curl to make a request to the endpoint and checks the response.

### Summary
This GitHub Actions workflow ensures that the project is built, tested, and verified on every push and pull request. It includes:

Build: Compiles the code to ensure there are no build errors.
Unit Tests and Coverage: Runs unit tests and checks for sufficient test coverage.
Functional Tests: Verifies the functionality of the deployed application using Docker.
By following this workflow, the project maintains high code quality and reliability through automated continuous integration and continuous deployment (CI/CD) practices.