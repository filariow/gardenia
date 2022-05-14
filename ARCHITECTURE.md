# Architectures

## Repository

This repository contains the different products that can be used to create an intelligent assistant for your garden.

The source code is organized as a monorepo with the following Go and TypeScript projects:
- Go
    - valve:
    - valvectl:
    - valved:
    - valvedmock:
    - schedulerd:
- TypeScript/JavaScript
    - gardenia-web

The **Go** repository follows the main ideas from [Standard Go Project Layout](https://github.com/golang-standards/project-layout), so you can in the directory `cmd` the main applications.

In the directory [fe](./fe) you can find the front-ends developed for this applications.
At the time of writing, we have [gardenia-web](./fe/gardenia-web) a SPA with Vue.JS and TypeScript.
To be able to communicate with the gRPC backend it's leveragin on an Envoy proxy and grpc-web.


