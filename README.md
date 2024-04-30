# Testcontainers demo

This repo was created to show how to create integration tests avoiding the use of mocks for external dependencies(such as databases, caching and message streaming), using `tescontainers`.   

It's a simple REST application, with 2 endpoints, that have the responsibility to create and list players from a `postgres` database.

The main reason to write integration tests is to validate that all components in a system are working properly together. Something that can be easily overlooked if we make use of mocks. 

## Running the project
```sh 
$ make app.run
```

## Executing tests
```sh 
$ make test
```
