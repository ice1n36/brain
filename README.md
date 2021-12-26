# Brain

‘Brain’, Penny’s dog and sidekick, is responsible for doing the leg work (e.g. storing data, fetching data, processing data, etc).

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
docker run --rm -it -d -p8081:80 ice1n36/brain
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
CONFIG_DIR=<path to config directory> bazel run :brain
```

### Docker
```
bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 :brain_container_image
docker run --rm -it -d -p8081:80 bazel:brain_container_image
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

note: if running with just bazel, use port 80

hello
```
curl -X GET localhost:8081/hello
```

write network traffic

```
curl -X POST localhost:80/api/v1/writeNetworkTraffic -d '{"app_id":"com.test.app"}'
```

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This is distributed under the [MIT License](http://www.opensource.org/licenses/MIT).
