The changelog for Bonvoy

## Pending Release

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