# Project Title

Choose and leave only one of the following badge

![REPO-TYPE](https://img.shields.io/badge/repo--type-backend-critical?style=for-the-badge&logo=github)

This project has the objective of giving the possibility of updating the image in the edge module by pulling the new docker image and update the running container.

The project expose the following endpoints:

- /current: retrieve the version of the installed module or tell if it is not present
- /update: which update the module to the latest version
- /update/:version: which update the module to the specific version
- /listversions: which list the available versions
- /updateplugin/:plugin name: which update the plugin selected

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them

```
golang version 1.21
```

### Installing

To run locally populate the .env file with AWS credentials and then

```
go run .
```

## Deployment

The docker image should be run using the following options in order to mak the image use the host machine docker daemon:

 docker run -p 8080:8080 -v /var/run/docker.sock:/var/run/docker.sock $IMAGE

## Built With

* [Go](https://go.dev/) - The programming language used

## Authors

* **Alessandro Lucani** - *Cloud Engineer* - [AleLuc98](https://github.com/AleLuc98)
