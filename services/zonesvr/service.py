from flask import Flask
from flask import request
import socket
import os
import sys
import json
sys.path.append("../proto/rpc/src")
import requests
import consulate

from google.protobuf.json_format import MessageToJson
import grpc
from room_pb2 import RoomRequest, Room
from room_pb2_grpc import RoomServiceStub

PORT = 8080

app = Flask(__name__)

TRACE_HEADERS_TO_PROPAGATE = [
    'X-Ot-Span-Context',
    'X-Request-Id',

    # Zipkin headers
    'X-B3-TraceId',
    'X-B3-SpanId',
    'X-B3-ParentSpanId',
    'X-B3-Sampled',
    'X-B3-Flags',

    # Jaeger header (for native client)
    "uber-trace-id"
]

@app.route('/zone/hello')
def homepage():
    return ('Hello from behind Envoy (zone)! hostname: {} resolved'
            'hostname: {}\n'.format(socket.gethostname(),
                                    socket.gethostbyname(socket.gethostname())))


@app.route('/zone/<service_number>')
def hello(service_number):
    return ('Hello from behind Envoy (zone {})! hostname: {} resolved'
            'hostname: {}\n'.format(os.environ['SERVICE_NAME'], 
                                    socket.gethostname(),
                                    socket.gethostbyname(socket.gethostname())))

@app.route("/zone/room/<room_id>")
def get_room(room_id):
    #get service's ip/port from consul
    #consul = consulate.Consul()
    #service = consul..services()

    with grpc.insecure_channel('localhost:50051') as channel:
        stub = RoomServiceStub(channel)
        return MessageToJson(stub.GetRoom(RoomRequest(id=1)))

#register service to consul
def registerService():
    consul = consulate.Consul()
    consul.agent.service.register('zonesvr', 
            port=PORT, 
            tags=['40001'], 
            )
    consul.agent.check.register('roomsvr1', script='nc -z -w5 localhost %d'%PORT, interval='30s')


if __name__ == "__main__":
    registerService()
    app.run(host='0.0.0.0', port=PORT, debug=True)
