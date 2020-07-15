package bitcoin

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/btcsuite/btcd/chaincfg"
	proto "github.com/ledgerhq/lama-bitcoin-svc"
)

var (
	mainnetChainParamsProto = &proto.ChainParams{
		Network: &proto.ChainParams_BitcoinNetwork{
			BitcoinNetwork: proto.BitcoinNetwork_BITCOIN_NETWORK_MAINNET,
		},
	}

	testnet3ChainParamsProto = &proto.ChainParams{
		Network: &proto.ChainParams_BitcoinNetwork{
			BitcoinNetwork: proto.BitcoinNetwork_BITCOIN_NETWORK_TESTNET3,
		},
	}

	regtestChainParamsProto = &proto.ChainParams{
		Network: &proto.ChainParams_BitcoinNetwork{
			BitcoinNetwork: proto.BitcoinNetwork_BITCOIN_NETWORK_REGTEST,
		},
	}

	invalidChainParamsProto = &proto.ChainParams{
		Network: &proto.ChainParams_BitcoinNetwork{
			BitcoinNetwork: 99999,
		},
	}
)

func Test_getChainParams(t *testing.T) {
	tests := []struct {
		name             string
		chainParamsProto *proto.ChainParams
		want             *chaincfg.Params
		wantErr          error
	}{
		{
			name:             "get mainnet chain params",
			chainParamsProto: mainnetChainParamsProto,
			want:             &chaincfg.MainNetParams,
			wantErr:          nil,
		},
		{
			name:             "get testnet3 chain params",
			chainParamsProto: testnet3ChainParamsProto,
			want:             &chaincfg.TestNet3Params,
			wantErr:          nil,
		},
		{
			name:             "get regtest chain params",
			chainParamsProto: regtestChainParamsProto,
			want:             &chaincfg.RegressionNetParams,
			wantErr:          nil,
		},
		{
			name:             "get unknown chain params",
			chainParamsProto: invalidChainParamsProto,
			want:             nil,
			wantErr:          ErrUnknownNetwork("99999"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getChainParams(tt.chainParamsProto)
			if err != nil && err != tt.wantErr {
				t.Errorf("getChainParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getChainParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ValidateAddress(t *testing.T) {
	tests := []struct {
		name    string
		request *proto.ValidateAddressRequest
		want    *proto.ValidateAddressResponse
		wantErr *status.Status
	}{
		{
			name: "mainnet P2PKH valid",
			request: &proto.ValidateAddressRequest{
				Address:     "1MirQ9bwyQcGVJPwKUgapu5ouK2E2Ey4gX",
				ChainParams: mainnetChainParamsProto,
			},
			want: &proto.ValidateAddressResponse{
				Address: "1MirQ9bwyQcGVJPwKUgapu5ouK2E2Ey4gX",
				IsValid: true,
			},
			wantErr: nil,
		},
		{
			name: "mainnet P2WPKH invalid checksum",
			request: &proto.ValidateAddressRequest{
				Address:     "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t5",
				ChainParams: mainnetChainParamsProto,
			},
			want: &proto.ValidateAddressResponse{
				Address:       "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t5",
				IsValid:       false,
				InvalidReason: "checksum failed. Expected v8f3t4, got v8f3t5.",
			},
			wantErr: nil,
		},
		{
			name: "testnet3 P2WPKH invalid mixed case",
			request: &proto.ValidateAddressRequest{
				Address:     "tb1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3q0sL5k7",
				ChainParams: testnet3ChainParamsProto,
			},
			want: &proto.ValidateAddressResponse{
				Address:       "tb1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3q0sL5k7",
				IsValid:       false,
				InvalidReason: "string not all lowercase or all uppercase",
			},
			wantErr: nil,
		},
		{
			name: "get unknown chain params",
			request: &proto.ValidateAddressRequest{
				Address:     "1MirQ9bwyQcGVJPwKUgapu5ouK2E2Ey4gX",
				ChainParams: invalidChainParamsProto,
			},
			want:    nil,
			wantErr: status.New(codes.InvalidArgument, ErrUnknownNetwork("99999").Error()),
		},
	}

	ctx := context.Background()
	s := &Service{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ValidateAddress(ctx, tt.request)
			grpcErr := status.Convert(err)

			if tt.wantErr != nil {
				if grpcErr.Code() != tt.wantErr.Code() {
					t.Errorf("getChainParams() gRPC error code = %v, wantErr %v",
						grpcErr.Code(), tt.wantErr.Code())
					return
				}

				if grpcErr.Message() != tt.wantErr.Message() {
					t.Errorf("getChainParams() gRPC error msg = %v, wantErr %v",
						grpcErr.Message(), tt.wantErr.Message())
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
