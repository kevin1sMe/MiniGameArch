#coding=utf8
from concurrent import futures
import time
import grpc

import sys
sys.path.append("../proto/rpc/src")

from room_pb2 import Room, RoomList, RoomRequest
import room_pb2_grpc as room_pb2_grpc

# GetRoom 的具体实现
class RoomServiceImpl(room_pb2_grpc.RoomServiceServicer):
    def GetRoom(self, request, context):
        result = Room(id = request.id)
        result.name = "Room #%d" % result.id
        result.size = 100
        result.type = 0
        return result

def start_service():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    room_pb2_grpc.add_RoomServiceServicer_to_server(RoomServiceImpl(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    try:
        while True:
            time.sleep(24 * 3600)
    except KeyboardInterrupt:
        server.stop(1)

if __name__ == '__main__':
    start_service()
