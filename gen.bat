cd protos && protoc --go_out=plugins=grpc:../services goods.proto
protoc --go_out=plugins=grpc:../services orders.proto
protoc --go_out=plugins=grpc:../services users.proto

protoc --go_out=plugins=grpc:../services models.proto

protoc --grpc-gateway_out=logtostderr=true:../services goods.proto
protoc --grpc-gateway_out=logtostderr=true:../services orders.proto
protoc --grpc-gateway_out=logtostderr=true:../services users.proto
cd ..