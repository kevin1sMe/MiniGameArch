/usr/bin/consul agent -retry-join consul-server-bootstrap -client 0.0.0.0 -data-dir /etc/consul.d & 
sleep 3
/usr/local/bin/envoy-xds -debug=true -xds=xds
