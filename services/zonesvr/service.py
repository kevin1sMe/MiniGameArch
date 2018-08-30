from flask import Flask
from flask import request
import socket
import os
import sys
import json
sys.path.append("../proto/rpc/src")
import requests

from google.protobuf.json_format import MessageToJson
import grpc
from room_pb2 import RoomRequest, Room
from room_pb2_grpc import RoomServiceStub


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
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = RoomServiceStub(channel)
        return MessageToJson(stub.GetRoom(RoomRequest(id=1)))

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=8080, debug=True)
