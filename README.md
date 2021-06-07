# Bonvoy

A simple go CLI tool to perform various tasks against Envoy sidecars in a Consul
Connect and Nomad environment.

## Usage

There are various commands you can run. Usually you are required to pass the
name of the service you want to query the sidecar for.

### List all Envoy Listeners:

```bash
bonvoy listeners auth-grpc
```