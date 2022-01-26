# Bonvoy

[![CircleCI](https://circleci.com/gh/bigcommerce/bonvoy/tree/main.svg?style=svg&circle-token=df991e7eb7eb4c38c2ece44f81cc259d6d9a8929)](https://circleci.com/gh/bigcommerce/bonvoy/tree/main) [![Maintainability](https://api.codeclimate.com/v1/badges/63e2c2eb1a4e00d269e2/maintainability)](https://codeclimate.com/repos/61f18cc93e634f01a300d78c/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/63e2c2eb1a4e00d269e2/test_coverage)](https://codeclimate.com/repos/61f18cc93e634f01a300d78c/test_coverage)

A simple go CLI tool to perform various tasks against Envoy sidecars in a Consul
Connect and Nomad environment.

Supports Consul 1.10+, Nomad 1.1.3+, and Envoy 1.18+. Requires DOCKER_API_VERSION 1.39.

## Usage

There are various commands you can run. Usually you are required to pass the
name of the service you want to query the sidecar for.

All commands have the ability to output in JSON format with: `-o json`

### Listeners

List all listeners:
```bash
$ bonvoy listeners list products-grpc

Listeners for products-grpc Envoy (PID 14013)
--------------------------------------------------------------------------
              NAME                 |    ADDRESS
-----------------------------------+-----------------
  envoy_metrics_listener           | 0.0.0.0:1239
  public_listener:0.0.0.0:22904    | 0.0.0.0:22904
  product-variants:127.0.0.1:8003  | 127.0.0.1:8003
  product-options:127.0.0.1:8004   | 127.0.0.1:8004
```

### Clusters

List all clusters for a sidecar:
```bash
$ bonvoy clusters list products-grpc

products-grpc.default.int-us-east1.internal.asdf1234.consul
------------------------------------------------------------------------------------------------------------
  Outlier: Success Rate          | -1   | Outlier: Success Rate Ejection | -1
                                 |      | Threshold                      |
  Outlier: Local Origin -        | -1   | Outlier: Local Origin          | -1
  Success Rate                   |      | - Success Rate Ejection        |
                                 |      | Threshold                      |
  Default Priority - Max         | 2048 | Default Priority - Max Retries | 3
  Connections                    |      |                                |
  Default Priority - Max Pending | 1024 | Default Priority - Max         | 1024
  Requests                       |      | Requests                       |
  High Priority - Max            | 1024 | High Priority - Max Retries    | 3
  Connections                    |      |                                |
  High Priority - Max Pending    | 1024 | High Priority - Max Requests   | 1024
  Requests                       |      |                                |

---------------------
- Cluster Instances -
---------------------
         HOST   | CX ACTIVE | CX FAILED | CX TOTAL | REQ ACTIVE | REQ TIMEOUT | REQ SUCCESS | REQ ERROR | REQ TOTAL | SUCCESS RATE | LOCAL SUCCESS RATE | HEALTH FLAGS | REGION | ZONE | SUBZONE | CANARY
----------------+-----------+-----------+----------+------------+-------------+-------------+-----------+-----------+--------------+--------------------+--------------+--------+------+---------+---------
  1.2.3.4:23837 |         1 |         0 |        1 |          0 |           0 |    23436077 |         0 |  23436077 |         -1.0 |               -1.0 |      healthy |        |      |         |  false
  1.2.3.5:23910 |         1 |         0 |        1 |          0 |           0 |    23613685 |         0 |  23613685 |         -1.0 |               -1.0 |      healthy |        |      |         |  false
  1.2.3.6:23011 |         1 |         0 |        1 |          0 |           0 |    54314043 |         0 |  54314043 |         -1.0 |               -1.0 |      healthy |        |      |         |  false
```

### Certificates

List all certificates:

```bash
$ bonvoy certificates list products-grpc

----------------------------------------------------------------------------
products-grpc Envoy (PID 14013)
----------------------------------------------------------------------------
Certificate Chains:
  SAN                                                                   | SERIAL # |   PATH   |      VALID FROM      |   EXPIRATION TIME    | DAYS UNTIL EXPIRATION
------------------------------------------------------------------------+----------+----------+----------------------+----------------------+------------------------
  spiffe://asdf1234.consul/ns/default/dc/int-us-east1/svc/products-grpc | 17d220e0 | <inline> | 2022-01-24T22:56:46Z | 2022-01-27T22:56:46Z | 1
  spiffe://asdf1234.consul/ns/default/dc/int-us-east1/svc/products-grpc | 17d220e0 | <inline> | 2022-01-24T22:56:46Z | 2022-01-27T22:56:46Z | 1
  spiffe://asdf1234.consul/ns/default/dc/int-us-east1/svc/products-grpc | 17d220e0 | <inline> | 2022-01-24T22:56:46Z | 2022-01-27T22:56:46Z | 1

CA Certificates:
  SAN                      | SERIAL # |   PATH   |      VALID FROM      |   EXPIRATION TIME    | DAYS UNTIL EXPIRATION
---------------------------+----------+----------+----------------------+----------------------+------------------------
  spiffe://asdf1234.consul | 17ce3503 | <inline> | 2020-07-20T16:05:55Z | 2030-07-20T16:05:55Z | 3096
  spiffe://asdf1234.consul | 17ce3503 | <inline> | 2020-07-20T16:05:55Z | 2030-07-20T16:05:55Z | 3096
  spiffe://asdf1234.consul | 17ce3503 | <inline> | 2020-07-20T16:05:55Z | 2030-07-20T16:05:55Z | 3096

```

List all expired certificates as compared to the Consul Agent:
```bash
$ bonvoy certificates expired products-grpc

  SERVICE | PID | ENVOY EXPIRY | DAYS LEFT | CONSUL LEAF EXPIRY | RESTARTED
----------+-----+--------------+-----------+--------------------+------------
```

Or show all certs expired on a host:
```bash
$ bonvoy certificates expired all

  SERVICE | PID | ENVOY EXPIRY | DAYS LEFT | CONSUL LEAF EXPIRY | RESTARTED
----------+-----+--------------+-----------+--------------------+------------
```

### Config

Dump the config:
```bash
$ bonvoy config dump products-grpc

# configuration dumped here in JSON format
```

### Server

Output server information:
```bash
$ bonvoy server info products-grpc

----------------------
- Server Information -
----------------------
  Service             | products-grpc
  Envoy Pid           | 14013
  Version             | asdf/1.22.0/Clean/RELEASE/BoringSSL
  Hot Restart Version | disabled
  State               | LIVE
  Uptime              | 1115062s

--------------------
- Node Information -
--------------------
  Node ID       | _nomad-task-asdf-group-products-grpc-sidecar-proxy
  Node Cluster  | products-grpc
  User Agent    | envoy
  Envoy Version | 1.22.0
  Namespace     | default

------------------------
- Command Line Options -
------------------------
  Concurrency          | 24
  Mode                 | Serve
  Log Level            | info
  Component Log Level  |
  Log Format           | [%Y-%m-%d %T.%e][%t][%l][%n]
                       | [%g:%#] %v
  Drain Strategy       | Gradual
  Drain Time           | 100s
  Config Path          | /path/to/envoy_bootstrap.json
  Parent Shutdown Time | 200s
```

Show server memory statistics:
```bash
$ bonvoy server memory products-grpc

----------------------
- Server Memory Info -
----------------------
  Service              | products-grpc
  Envoy Pid            | 14013
  Allocated            | 16349720
  Heap Size            | 41943040
  Page Heap (Unmapped) | 0
  Page Heap (Free)     | 3317760
  Total Physical Bytes | 43894086
  Total Thread Cache   | 20066256
```

### Logging

Set the log level for a sidecar:
```bash
$ bonvoy log level products-grpc -l debug
```

## License

Copyright (c) 2021-present, BigCommerce Pty. Ltd. All rights reserved

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit
persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.