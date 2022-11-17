# Library Portal

A portal for students to browse books in a library with their availability. And for a librarian to record borrows and returns.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

You need to have `docker` installed. To install it on Linux, follow the instructions [here](https://docs.docker.com/engine/install/).

Note: Enable root-less access to docker (*Recommended*)

### Installing

First, download the source code

```bash
git clone https://github.com/sujay1844/library-portal.git
```

And then setup the demo environment

```bash
docker compose up -d
```

To test if the server is running, run

```bash
curl http://localhost:8080/all
```

The above command should return a list of books in `json` format.


## Built With

- [MongoDB](https://mongodb.com) - NoSQL database
- [Go](https://go.dev) - Backend language

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details
