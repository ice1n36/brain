# Brain

‘Brain’, Penny’s dog and sidekick, is responsible for all things mobile data (e.g. storing and fetching mobile network data).

## Install

### Go
```
go get github.com/ice1n36/brain
```

See below on how to build and run

### Docker
```
docker pull ice1n36/brain:latest
```

## Run
```
docker run --rm -it -d -p8081:8081 ice1n36/brain
```

### Docker Compose
In this directory:
```
docker-compose up -d
```

and to teardown:
```
docker-compose down
```

## Development
### Build & Run

update dependencies
```
./update-deps.sh
```

build and run

```
bazel run :brain
```

### Docker
```
bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 :brain_container_image
docker run --rm -it -d -p8081:8081 bazel:brain_container_image
```

## Publish

### Docker

```
bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 :brain_container_image_push
```

## Test

### Unit
TODO

### Locally
```
curl -X GET localhost:8081/hello
```

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This is distributed under the [MIT License](http://www.opensource.org/licenses/MIT).
