# Delta Swagger Generator

Delta Swagger Generator is a tool that generates swagger documentation for the Delta APIs.
## How to run the generator

```
make generate-swagger
```

The above command will do the following:
- scan the `./api` directory for all the `*.go` files starting with router.go
- scan all the annotations on the handlers
- generate the swagger documentation in `./docs/swagger/swagger.yaml`
