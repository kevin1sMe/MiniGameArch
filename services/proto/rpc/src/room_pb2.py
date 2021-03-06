# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: room.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='room.proto',
  package='naruto',
  syntax='proto3',
  serialized_pb=_b('\n\nroom.proto\x12\x06naruto\"*\n\x08RoomList\x12\x1e\n\x08roomList\x18\x01 \x03(\x0b\x32\x0c.naruto.Room\"\x19\n\x0bRoomRequest\x12\n\n\x02id\x18\x01 \x01(\x05\"<\n\x04Room\x12\n\n\x02id\x18\x01 \x01(\r\x12\x0c\n\x04name\x18\x02 \x01(\t\x12\x0c\n\x04size\x18\x03 \x01(\r\x12\x0c\n\x04type\x18\x04 \x01(\r2=\n\x0bRoomService\x12.\n\x07GetRoom\x12\x13.naruto.RoomRequest\x1a\x0c.naruto.Room\"\x00\x62\x06proto3')
)




_ROOMLIST = _descriptor.Descriptor(
  name='RoomList',
  full_name='naruto.RoomList',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='roomList', full_name='naruto.RoomList.roomList', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=22,
  serialized_end=64,
)


_ROOMREQUEST = _descriptor.Descriptor(
  name='RoomRequest',
  full_name='naruto.RoomRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='naruto.RoomRequest.id', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=66,
  serialized_end=91,
)


_ROOM = _descriptor.Descriptor(
  name='Room',
  full_name='naruto.Room',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='naruto.Room.id', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='name', full_name='naruto.Room.name', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='size', full_name='naruto.Room.size', index=2,
      number=3, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='type', full_name='naruto.Room.type', index=3,
      number=4, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=93,
  serialized_end=153,
)

_ROOMLIST.fields_by_name['roomList'].message_type = _ROOM
DESCRIPTOR.message_types_by_name['RoomList'] = _ROOMLIST
DESCRIPTOR.message_types_by_name['RoomRequest'] = _ROOMREQUEST
DESCRIPTOR.message_types_by_name['Room'] = _ROOM
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

RoomList = _reflection.GeneratedProtocolMessageType('RoomList', (_message.Message,), dict(
  DESCRIPTOR = _ROOMLIST,
  __module__ = 'room_pb2'
  # @@protoc_insertion_point(class_scope:naruto.RoomList)
  ))
_sym_db.RegisterMessage(RoomList)

RoomRequest = _reflection.GeneratedProtocolMessageType('RoomRequest', (_message.Message,), dict(
  DESCRIPTOR = _ROOMREQUEST,
  __module__ = 'room_pb2'
  # @@protoc_insertion_point(class_scope:naruto.RoomRequest)
  ))
_sym_db.RegisterMessage(RoomRequest)

Room = _reflection.GeneratedProtocolMessageType('Room', (_message.Message,), dict(
  DESCRIPTOR = _ROOM,
  __module__ = 'room_pb2'
  # @@protoc_insertion_point(class_scope:naruto.Room)
  ))
_sym_db.RegisterMessage(Room)



_ROOMSERVICE = _descriptor.ServiceDescriptor(
  name='RoomService',
  full_name='naruto.RoomService',
  file=DESCRIPTOR,
  index=0,
  options=None,
  serialized_start=155,
  serialized_end=216,
  methods=[
  _descriptor.MethodDescriptor(
    name='GetRoom',
    full_name='naruto.RoomService.GetRoom',
    index=0,
    containing_service=None,
    input_type=_ROOMREQUEST,
    output_type=_ROOM,
    options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_ROOMSERVICE)

DESCRIPTOR.services_by_name['RoomService'] = _ROOMSERVICE

# @@protoc_insertion_point(module_scope)
