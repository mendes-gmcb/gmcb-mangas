# GMCB Manga API

GMCB Manga API is a Go-based RESTful API that provides manga-related functionalities, including manga information, chapters, images. It leverages AWS S3 for image storage.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
- [Usage](#usage)
- [Docker](#docker)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Manga Management:** Create, retrieve, retrieve deleted, update, update image, and delete manga information.
- **Chapter Management:** Manage manga chapters, including the number of pages and images.
- **Image Handling:** Upload images to AWS S3 and perform image size validation.

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- AWS S3 account and configuration

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/mendes-gmcb/gmcb-manga-api.git
   cd gmcb-manga-api
   ```

2. Build the application:
  ```bash
  go build
  ```

3. Run the migrations:
  ```bash
  go run migrate/migrate.go
  ```

4. Run the application:
  ```bash
  ./gmcb-manga-api
  ```

### Configuration

Configure the application using environment variables or a configuration file. The following environment variables are required:

- DB_HOST: Database host to will connect
- DB_USER= Database user to will connect
- DB_PASS: Database password to will connect
- DB_PORT: Database port to will connect
- DB_NAME: Database name to will connect

<br>

- AWS_ACCESS_KEY: AWS access key
- AWS_SECRET_KEY: AWS secrete key 
- AWS_REGION: AWS region of bucket
- AWS_BUCKET_NAME: AWS bucket name

## Usage

1. Define your manga entries, chapters, and images using the provided API routes.
2. Explore the API routes for manga, chapters, images, and cover images.

## Docker

To run the application using Docker:

```bash
docker-compose up --build
```

This will start the application, Redis server, and Nginx load balancer.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
