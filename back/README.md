# gin-template
Template GO api with gin as framework. This project is a template for a GO api with gin as framework. It includes a basic structure for a REST api with a user model and a basic CRUD. It also includes a basic authentication system with JWT.

## Dependencies

- [Gin](https://github.com/gin-gonic/gin) Web framework
- [GORM](https://gorm.io/) ORM
- [JWT](https://github.com/golang-jwt/jwt) JWT library

## Getting Started

### Prerequisites

- [Golang](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Installing

- Create repository from template with Github Cli

```bash
gh repo create <repo-name> --template gin-template
```

- Create repository from template with [Github UI](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-repository-from-a-template)
- Clone repository as template
```bash
git clone git@github:JulesRobineau/gin-template.git <repo-name>
rm -rf .git
git init
git add .
git commit -m "Initial commit"
git remote add origin <new-repo-name>
git push -u origin develop
```
Once you have created the repository, you can install the dependencies with the following command:
- First Golang
- Then the dependencies

Golang dependencies and install Swag CLI
```bash
make install
```

## Running the tests

Launch tests with tool go test
```bash
make test
```

## Running the server

Launch server with docker compose
```bash
make run
```

---

## TODO
- [ ] Functional tests
- [ ] Add more tests
- [ ] Add more documentation
- [ ] Add more examples
- [ ] Add more features
  - [ ] Add email verification
  - [ ] Add password reset

