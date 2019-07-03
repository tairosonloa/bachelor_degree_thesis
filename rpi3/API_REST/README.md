# rpi3/API_REST/
This folder contains the implementation of `rpi3_api`.

## app/
Contains the implementation of the main features of `rpi3_api`.

## main.go
Contains the `main` method which initializes the API object and starts the server.

## Makefile
Makefile to compile the API binary. Can be called with
```
make <options>
```
Allowed options are:
* `linux` or `amd64`: both options compiles the binary for amd64 processor architecture.
* `rpi` or `arm`: both options compiles the binary for ARM processor architecture.