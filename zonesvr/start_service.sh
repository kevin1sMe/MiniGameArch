#!/bin/sh
python3 /code/service.py &
nohup envoy -c /etc/service-envoy.yaml --service-cluster service${SERVICE_NAME}