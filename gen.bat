cd protos && protoc --go_out=plugins=grpc:../services goods.proto
protoc --go_out=plugins=grpc:../services orders.proto
protoc --go_out=plugins=grpc:../services models.proto
cd ..