#!/bin/sh
python3 /code/service.py &
#consul agent -retry-join consul-server-bootstrap -client 0.0.0.0 & 
envoy -c /etc/service-envoy.yaml --service-cluster zonesvr_${SERVICE_NAME}