package proto

//go:generate protoc --gogofaster_out=plugins=grpc:. traffic_quota.proto
//go:generate mockgen -source=traffic_quota.pb.go -package=${GOPACKAGE} -destination=mock_traffic_quota.pb.go
