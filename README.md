YouTube GO API
======================

The YouTube Go API is a REST API built in Go which wraps the kkdai YouTube package - https://github.com/kkdai/youtube.

This tool is meant to be used to download CC0 licenced content, we do not support nor recommend using it for illegal activities.

## Dependencies
List of the primary 3rd party dependencies. There are other sub-dependencies not listed here.

| Name | URL |
| :---- | :---- |
| kkdai/youtube | https://github.com/kkdai/youtube |
| Gin Web Framework | https://github.com/gin-gonic/gin |
| gin-swagger | https://github.com/swaggo/gin-swagger |
| swag | https://github.com/swaggo/swag |

## Build and Run
These instructions assume you have cloned the repo and are at a command prompt at the root of the project.  

To generate Swagger docs, you will need to install swag - https://github.com/swaggo/swag

1. Tidy up  
```sehll
go mod tidy`
```
   
3. Generate Swagger docs (optional)  
```shell
swag init --pd --pdl 1`
```

5. Run  
```shell
go run main.go`
```
