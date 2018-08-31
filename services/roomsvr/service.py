#coding=utf8
from concurrent import futures
import time
import grpc
import consulate

import sys
sys.path.append("../proto/rpc/src")

from room_pb2 import Room, RoomList, RoomRequest
import room_pb2_grpc as room_pb2_grpc

PORT=50051
ENVOY_PORT=80
# GetRoom 的具体实现
class RoomServiceImpl(room_pb2_grpc.RoomServiceServicer):
    def GetRoom(self, request, context):
        result = Room(id = request.id)
        result.name = "Room #%d" % result.id
        result.size = 100
        result.type = 0
        return result

#register service to consul
def registerService():
    consul = consulate.Consul()
    consul.agent.service.register('roomsvr', 
            port=ENVOY_PORT, 
            tags=['master'], 
            #ttl = '30s', 
            #grpc的health check还需要研究
            #check = consul.agent.Check(
                #name="roomsvr health status",
                #grpc='127.0.0.1:%d'%(PORT),
                #grpc_use_tls=True,
                #interval='10s'
                #))
            #先用tcp检测
            #check = consul.agent.Check(name='roomsvr health status', tcp='localhost:50051', interval='30s')
            )
    #consul.agent.check.register('roomsvr1', script='/bin/tcping -t 1 localhost %d'%PORT, interval='30s')
    consul.agent.check.register('roomsvr1', script='nc -z -w5 localhost %d'%ENVOY_PORT, interval='30s')

def start_service():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    room_pb2_grpc.add_RoomServiceServicer_to_server(RoomServiceImpl(), server)
    server.add_insecure_port("[::]:%d"%PORT)
    server.start()
    try:
        while True:
            time.sleep(24 * 3600)
    except KeyboardInterrupt:
        server.stop(1)


if __name__ == '__main__':
    registerService()
    start_service()
