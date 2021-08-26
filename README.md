# Bonvoy

[![CircleCI](https://circleci.com/gh/bigcommerce/bonvoy/tree/main.svg?style=svg&circle-token=df991e7eb7eb4c38c2ece44f81cc259d6d9a8929)](https://circleci.com/gh/bigcommerce/bonvoy/tree/main)

A simple go CLI tool to perform various tasks against Envoy sidecars in a Consul
Connect and Nomad environment.

Supports Consul 1.10+, Nomad 1.1.3+, and Envoy 1.18+. 

## Usage

There are various commands you can run. Usually you are required to pass the
name of the service you want to query the sidecar for.

All commands have the ability to output in JSON format with: `-o json`

### Listeners

List all listeners:
```bash
bonvoy listeners list auth-grpc
```

### Clusters

List all clusters for a sidecar:
```bash
bonvoy clusters list auth-grpc
```

### Certificates

List all certificates:
```bash
bonvoy certificates list auth-grpc
```

List all expired certificates as compared to the Consul Agent:
```bash
bonvoy certificates expired auth-grpc
```

Or show all certs expired on a host:
```bash
bonvoy certificates expired all
```

### Config

Dump the config:
```bash
bonvoy config dump auth-grpc
```

### Server

Output server information:
```bash
bonvoy server info auth-grpc
```

Show server memory statistics:
```bash
bonvoy server memory auth-grpc
```

### Logging

Set the log level for a sidecar:
```bash
bonvoy log level auth-grpc -l debug
```