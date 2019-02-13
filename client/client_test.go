package client

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/orgplace/trafficquota/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

//go:generate mockgen -destination=mock_grpc_health_v1_test.go -package=${GOPACKAGE} google.golang.org/grpc/health/grpc_health_v1 HealthClient

func Test_newClient(t *testing.T) {
	cc := &grpc.ClientConn{}
	want := &client{
		clientConn:   cc,
		trafficQuota: proto.NewTrafficQuotaClient(cc),
		health:       grpc_health_v1.NewHealthClient(cc),
	}
	if got := newClient(cc); !reflect.DeepEqual(got, want) {
		t.Errorf("newClient() = %v, want %v", got, want)
	}
}

func Test_parseAddr(t *testing.T) {
	tests := []struct {
		addr        string
		wantAddr    string
		wantOptions []grpc.DialOption
	}{
		{
			addr:        "localhost:3895",
			wantAddr:    "localhost:3895",
			wantOptions: []grpc.DialOption{},
		},
		{
			addr:        "unix:/tmp/example.sock",
			wantAddr:    "/tmp/example.sock",
			wantOptions: []grpc.DialOption{unixDomainSocketDialer},
		},
	}
	for _, tt := range tests {
		t.Run(tt.addr, func(t *testing.T) {
			gotAddr, gotOptions := parseAddr(tt.addr)
			if gotAddr != tt.wantAddr {
				t.Errorf("parseAddr() gotAddr = %v, wantAddr %v", gotAddr, tt.wantAddr)
			}
			if !reflect.DeepEqual(gotOptions, tt.wantOptions) {
				t.Errorf("parseAddr() gotOptions = %v, wantOptions %v", gotOptions, tt.wantOptions)
			}
		})
	}
}

func Test_client_Take(t *testing.T) {
	type args struct {
		partitionKey   string
		clusteringKeys []string
	}
	type takeCall struct {
		req *proto.TakeRequest
		res *proto.TakeResponse
		err error
	}
	tests := []struct {
		name     string
		args     args
		takeCall takeCall
		want     bool
		wantErr  bool
	}{
		{
			name: "ok",
			args: args{
				partitionKey:   "partition key",
				clusteringKeys: []string{""},
			},
			takeCall: takeCall{
				req: &proto.TakeRequest{
					PartitionKey:   "partition key",
					ClusteringKeys: []string{""},
				},
				res: &proto.TakeResponse{
					Allowed: true,
				},
			},
			want: true,
		},
		{
			name: "error",
			args: args{
				partitionKey:   "partition key",
				clusteringKeys: []string{""},
			},
			takeCall: takeCall{
				req: &proto.TakeRequest{
					PartitionKey:   "partition key",
					ClusteringKeys: []string{""},
				},
				err: errors.New("error for test"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockTrafficQuota := proto.NewMockTrafficQuotaClient(ctrl)

			mockTrafficQuota.EXPECT().Take(
				context.Background(),
				tt.takeCall.req,
			).Return(tt.takeCall.res, tt.takeCall.err)

			c := &client{
				trafficQuota: mockTrafficQuota,
			}

			got, err := c.Take(tt.args.partitionKey, tt.args.clusteringKeys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.Take() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("client.Take() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_Ping(t *testing.T) {
	type checkCall struct {
		res *grpc_health_v1.HealthCheckResponse
		err error
	}
	tests := []struct {
		name      string
		checkCall checkCall
		wantErr   bool
	}{
		{
			name: "ok",
			checkCall: checkCall{
				res: &grpc_health_v1.HealthCheckResponse{
					Status: grpc_health_v1.HealthCheckResponse_SERVING,
				},
			},
		},
		{
			name: "not serving",
			checkCall: checkCall{
				res: &grpc_health_v1.HealthCheckResponse{
					Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
				},
			},
			wantErr: true,
		},
		{
			name: "error",
			checkCall: checkCall{
				err: errors.New("error for test"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockHealth := NewMockHealthClient(ctrl)
			mockHealth.EXPECT().Check(
				context.Background(),
				&grpc_health_v1.HealthCheckRequest{},
			).Return(tt.checkCall.res, tt.checkCall.err)

			c := &client{
				health: mockHealth,
			}

			if err := c.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("client.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
