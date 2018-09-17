#!/bin/sh
python3 /code/service.py &
consul agent -retry-join consul-server-bootstrap -client 0.0.0.0 -data-dir /etc/consul.d & 
envoy -c /etc/service-envoy.yaml --service-cluster zonesvr --service-node zonesvr -l debug
