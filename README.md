# Secret Server

The simple service to manage secrets

[![Go Report Card](https://goreportcard.com/badge/github.com/ilyakaznacheev/secret)](https://goreportcard.com/report/github.com/ilyakaznacheev/secret) 
[![GoDoc](https://godoc.org/github.com/ilyakaznacheev/secret?status.svg)](https://godoc.org/github.com/ilyakaznacheev/secret)
[![Build Status](https://travis-ci.org/ilyakaznacheev/secret.svg?branch=master)](https://travis-ci.org/ilyakaznacheev/secret)
[![Heroku](https://pyheroku-badge.herokuapp.com/?app=secret&root=api&style=flat)](https://secret-web.herokuapp.com/)
[![Coverage Status](https://codecov.io/github/ilyakaznacheev/secret/coverage.svg?branch=master)](https://codecov.io/gh/ilyakaznacheev/secret)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

## Overview

Secret Server helps to securely store secrets and get them back when needed.

## Contents

- [About Secret Server](#about-secret-server)
- [Requirements](#requirements)
- [Usage](#usage)
    - [Download](#download)
    - [Run Local](#run-local)
    - [Docker Compose](#docker-compose)
    - [Deployment](#deployment)
        - [Heroku](#heroku)
- [API documentation](#api-documentation)
- [Contributing](#contributing)

## About Secret Server

Secret server is designed to mace secret storage easier. Simple API allowes to save and get back your secrets, and also set expiration time and view limit.

## Requirements

If you run the app using docker-compose you only need it to start and a git to download a package.

Otherwise you need Golang 1.11+ with modules enabled to compile the app and a Redis to store the secrets.

## Usage

Here is a short manual how to download and use the app.

### Download

Create a new folder and run:

```bash
git clone https://github.com/ilyakaznacheev/secret.git
```

or download it from this page manually.

### Run Local

You need to start Redis first. The default connection path is `localhost:5050` but you can overwrite it using `REDIS_URL` environment variable. Run ```go run main.go -h` for more info.

### Docker Compose

To start the whole environment run

```bash
docker-compose up
```

It will start the service with Redis DB, and also Prometheus metrics with Grafana UI.

Service will be on `localhost:8080` and Grafana will be on `localhost:3000`.

### Deployment

#### Heroku

## API documentation

API is documented using Swagger. Check the [swagger.yml](/swagger.yml) file.

Try it in [Swagger Editor](https://editor.swagger.io/)!

## Contributing

The application is open-sourced under the [MIT](/LICENSE) license.

If you will find some error, want to add something or ask a question - feel free to create an issue and/or make a pull request.

Any contribution is welcome.