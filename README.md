# dating
Tool to help people find more friends who have a similar hoppy, passion..., then we can connect to make a date

## Prerequisites

Make sure you have the development environment matches with these notes below so we can mitigate any problems of version mismatch.

- Backend:
  - Go SDK: 1.16.
    Make sure to set `$GOROOT` and `$GOPATH` correctly.
    You can check those environment variable by typing: `go env`.
  - MongoDB: 4.4.

- Commons:
  - Install [Docker CE](https://docs.docker.com/install/) and [docker-compose](https://docs.docker.com/compose/install/).
  - Install [git](https://git-scm.com/) for manage source code.
  - IDE of your choice, recommended `Goland` or `VS Code`.

## Development

#### 1. Clone code to local

```shell
$ go get -u -v github.com/paphuc/dating
or
$ cd $GOPATH/src
$ git clone https://github.com/paphuc/dating.git
```
After this step, source code must be available at `$GOPATH/src/github.com/paphuc/dating`.

#### 2. Start development environment manually

- Start MongoDB service at localhost:27017. The easiest way is to run the Docker as below:

  ```shell
  $ docker run -p 27017:27017 -v /opt/data/mongo_home:/data/db --name mongo -d mongo:4.1.8
  ```

- Start backend API service (Go):

  ```shell
  $ go run main.go
  # Backend service will start on port :8080.
  ```
  
#### 3. Start development environment with Docker

Instead of manually start services like step 2. You can use Docker to start all services at once.

```shell
$ docker-compose up
```

After started, services will be available at `localhost` with ports as below:
```
MongoDB: 27017
Backend: 8080
```

## Notes

- Make sure to run `go fmt`, `go vet`, `go test`, and `go build / go install` before pushing your code to Github.
  Or you can just run `make` before pushing.
- Never commit directly to `master` or `develop` branches (you don't have permission to do so, anyway). Instead, checkout from `develop` branch to a separated branch then work on that.  
  Whenever you finish your work, you can create a Pull Request (PR) / Merge Request (MR) to ask for code review and merging your branch back to `develop`.   
  `master` branch will be reserved when administrator decide to release a stable version of application.

## Tech requirements:

- Authentication:
  - JWT (frontend authentication)
- Database
  - MongoDB

