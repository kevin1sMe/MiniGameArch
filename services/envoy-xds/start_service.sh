/usr/bin/consul agent -retry-join consul-server-bootstrap -client 0.0.0.0 -data-dir /etc/consul.d & 
sleep 3
/usr/local/bin/envoy-xds -debug=true -nodeID=front-proxy -xds=xds