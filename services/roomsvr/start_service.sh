/usr/bin/consul agent -retry-join consul-server-bootstrap -client 0.0.0.0 -data-dir /etc/consul.d & 
python3 /code/service.py &
envoy -c /etc/service-envoy.yaml --service-cluster roomsvr --service-node roomsvr --v2-config-only -l debug
