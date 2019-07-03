# rpi3/
This folder contains the system software components that are intended to run in rpi3.

## API_REST/
Contains the implementation of `rpi3_api`.

## Makefile
Makefile to compile the custom light web server binary. Can be called with
```
make <options>
```
Allowed options are:
* `linux` or `amd64`: both options compiles the binary for amd64 processor architecture.
* `rpi` or `arm`: both options compiles the binary for ARM processor architecture.

## web_server.go
Implementation of a custom light web server to serve the rpi3 GUI (GatsbyJS dashboard).