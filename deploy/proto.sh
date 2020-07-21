# PROTOBUF

# should run these commands inside proto dir (need to improve that)
protoc --go_out=plugins=grpc:../pkg/authentication  --go_opt=paths=source_relative --proto_path=../proto authentication.proto
protoc --go_out=plugins=grpc:../pkg/person  --go_opt=paths=source_relative --proto_path=../proto person.proto
protoc --go_out=plugins=grpc:../pkg/chat  --go_opt=paths=source_relative --proto_path=../proto chat.proto
