# Containerlab virtual rack example with Netbox integration

This example illustrates integration with Netbox

## Lab setup

1. Start the devices
```
sudo containerlab deploy
```

2. Start Netbox
```
docker-compose up
```

# Use cases

## Build virtual replica of an existing physical rack

* Generate Containerlab YAML configuration file from data contained in Netbox
  - mgmt ip addresses
  - device types
  - ...
