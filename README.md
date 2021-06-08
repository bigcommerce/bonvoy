# Bonvoy

[![CircleCI](https://circleci.com/gh/bigcommerce/bonvoy/tree/main.svg?style=svg&circle-token=df991e7eb7eb4c38c2ece44f81cc259d6d9a8929)](https://circleci.com/gh/bigcommerce/bonvoy/tree/main)

A simple go CLI tool to perform various tasks against Envoy sidecars in a Consul
Connect and Nomad environment.

## Usage

There are various commands you can run. Usually you are required to pass the
name of the service you want to query the sidecar for.

### List all Envoy Listeners

```bash
bonvoy listeners auth-grpc
```

### Show all Expired Certificates

```bash
bonvoy certs-expired auth-grpc
```