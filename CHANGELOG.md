The changelog for Bonvoy

## Pending Release

## 0.3.1

* Adjust envoy default address to 127.0.0.2; Nomad changed this in Nomad 1.1.3

## 0.3.0

* Add `clusters compare` to compare clusters between Consul and Envoy

## 0.2.0

* Add `-o json` to output data into a JSON output format
* Adjust `server restart` to restart the Nomad allocation rather than killing Envoy.
  This allows for a more graceful degradation of traffic to the service.
* Add cluster filter to `clusters list` to only show a specific cluster

## 0.1.0

* Add `server restart` to restart a sidecar
* Add `statistics dump` command to dump statistics
* Better resiliency around restarting sidecars
* Add `certificates list` to show all registered certificates on a sidecar
* Add `config dump` to dump the entire Envoy config

## 0.0.4

* Add `clusters list` to display all clusters and relevant statistics around them
* Add `-r/--restart` flag for `certificates expired` to restart all expired sidecars
* Add `certificates expired all` to show all expired certs on a host

## 0.0.3

* Ensure query for containers does not include dead containers

## 0.0.2

* Add `server memory` command to display server memory information
* Restructured to use cobra cli library rather than flag
* Add `server info` command to display information about the Envoy sidecar
* Add `log level` command to set Envoy's log level
* Restructure for better memory usage and a more closed API

## 0.0.1

Initial release