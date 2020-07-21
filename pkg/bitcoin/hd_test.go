package bitcoin

import (
	"context"
	"reflect"
	"testing"

	"github.com/btcsuite/btcutil/hdkeychain"
	pb "github.com/ledgerhq/lama-bitcoin-svc/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// h indicates the BIP32 harden bit, equivalent to 2^31.
var h uint32 = 0x80000000

func TestDeriveExtendedKey(t *testing.T) {
	tests := []struct {
		name    string
		request *pb.DeriveExtendedKeyRequest
		want    *pb.DeriveExtendedKeyResponse
		wantErr *status.Status
	}{
		{
			// BIP0032: Test Vector 1 (chain m/0H/1/2H)
			name: "ErrDeriveHardFromPublic",
			request: &pb.DeriveExtendedKeyRequest{
				ExtendedKey: "xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5",
				Derivation:  []uint32{0 + h, 1, 2 + h},
			},
			want:    nil,
			wantErr: status.New(codes.Internal, hdkeychain.ErrDeriveHardFromPublic.Error()),
		},
		{
			// BIP0032: Test Vector 2 (chain m)
			name: "mainnet derive from root",
			request: &pb.DeriveExtendedKeyRequest{
				ExtendedKey: "xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB",
				Derivation:  []uint32{0, 2147483647, 1, 2147483646, 2}, // m/0/2147483647/1/2147483646/2
			},
			want: &pb.DeriveExtendedKeyResponse{
				ExtendedKey: "xpub6H7WkJf547AiSwAbX6xsm8Bmq9M9P1Gjequ5SipsjipWmtXSyp4C3uwzewedGEgAMsDy4jEvNTWtxLyqqHY9C12gaBmgUdk2CGmwachwnWK",
				PublicKey: []byte{
					0x03, 0x54, 0xf9, 0x40, 0xcd, 0xd9, 0x6e, 0xeb,
					0x6b, 0xb5, 0xea, 0x60, 0x76, 0x95, 0x1e, 0x28,
					0xc2, 0x5f, 0xec, 0x76, 0xad, 0xf2, 0xc7, 0x8d,
					0x1e, 0xdc, 0xae, 0xc4, 0x20, 0xe7, 0xc8, 0xc7,
					0x0b,
				},
				ChainCode: []byte{
					0x1a, 0xc7, 0x97, 0x78, 0x3c, 0x3e, 0x3d, 0x79,
					0x83, 0xdc, 0x80, 0xb9, 0xa5, 0x80, 0xd4, 0xba,
					0x0d, 0x2d, 0x6b, 0xf9, 0x89, 0xfd, 0x35, 0xa2,
					0x1f, 0x41, 0x83, 0xd4, 0x33, 0x6f, 0xdf, 0x24,
				},
			},
		},
		{
			// https://github.com/LedgerHQ/lib-ledger-core/blob/54ddf50/core/test/bitcoin/address_test.cpp#L81
			name: "mainnet derive from account level",
			request: &pb.DeriveExtendedKeyRequest{
				ExtendedKey: "xpub6Cc939fyHvfB9pPLWd3bSyyQFvgKbwhidca49jGCM5Hz5ypEPGf9JVXB4NBuUfPgoHnMjN6oNgdC9KRqM11RZtL8QLW6rFKziNwHDYhZ6Kx",
				Derivation:  []uint32{0, 1}, // m/44'/0'/0'/0/1
			},
			want: &pb.DeriveExtendedKeyResponse{
				ExtendedKey: "xpub6HHu39JZziv1GHo1Yvm3DBa7Wztu93uMyrssG9DFgXsnmaRs7JmCtrcGRvJVd5gnvtRfDXW1CqfR7Q4CwCFsWWAYUHnWPEAKEdr35q51JY3",
				PublicKey: []byte{
					0x02, 0xc3, 0x68, 0xbd, 0xec, 0x47, 0xa1, 0xb6,
					0xfa, 0xa7, 0x6d, 0x62, 0x4e, 0xad, 0x0c, 0xd2,
					0x78, 0x32, 0x34, 0x98, 0x3c, 0x46, 0x67, 0x67,
					0x21, 0x6e, 0xcd, 0xac, 0x8c, 0x47, 0x2d, 0xf3,
					0xa6,
				},
				ChainCode: []byte{
					0xb6, 0xb8, 0xa4, 0x9c, 0x62, 0x34, 0xb2, 0x6c,
					0x91, 0xbf, 0xaf, 0xac, 0xd9, 0x05, 0x4c, 0x18,
					0x56, 0x21, 0x30, 0x23, 0x4d, 0xc3, 0x9e, 0x94,
					0x63, 0x56, 0x1c, 0xa6, 0x66, 0x7f, 0x40, 0xf8,
				},
			},
		},
		{
			// BIP0032: Test vector 1 (combining chains m/0H/1/2H and m/0H/1/2H/2 for testnet3)
			name: "testnet3 derive chain path",
			request: &pb.DeriveExtendedKeyRequest{
				ExtendedKey: "tpubDDRojdS4jYQXNugn4t2WLrZ7mjfAyoVQu7MLk4eurqFCbrc7cHLZX8W5YRS8ZskGR9k9t3PqVv68bVBjAyW4nWM9pTGRddt3GQftg6MVQsm",
				Derivation:  []uint32{2}, // m/0'/1/2'/2
			},
			want: &pb.DeriveExtendedKeyResponse{
				ExtendedKey: "tpubDFfCa4Z1v25WTPAVm9EbEMiRrYwucPocLbEe12BPBGooxxEUg42vihy1DkRWyftztTsL23snYezF9uXjGGwGW6pQjEpcTpmsH6ajpf4CVPn",
				PublicKey: []byte{
					0x02, 0xe8, 0x44, 0x50, 0x82, 0xa7, 0x2f, 0x29,
					0xb7, 0x5c, 0xa4, 0x87, 0x48, 0xa9, 0x14, 0xdf,
					0x60, 0x62, 0x2a, 0x60, 0x9c, 0xac, 0xfc, 0xe8,
					0xed, 0x0e, 0x35, 0x80, 0x45, 0x60, 0x74, 0x1d,
					0x29,
				},
				ChainCode: []byte{
					0xcf, 0xb7, 0x18, 0x83, 0xf0, 0x16, 0x76, 0xf5,
					0x87, 0xd0, 0x23, 0xcc, 0x53, 0xa3, 0x5b, 0xc7,
					0xf8, 0x8f, 0x72, 0x4b, 0x1f, 0x8c, 0x28, 0x92,
					0xac, 0x12, 0x75, 0xac, 0x82, 0x2a, 0x3e, 0xdd,
				},
			},
		},
	}

	ctx := context.Background()
	s := &Service{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.DeriveExtendedKey(ctx, tt.request)
			grpcErr := status.Convert(err)

			if grpcErr != nil && tt.wantErr == nil {
				t.Errorf("unexpected error in DeriveExtendedKey(): %v", grpcErr.Message())
				return
			}

			if tt.wantErr != nil {
				if grpcErr.Code() != tt.wantErr.Code() {
					t.Errorf("DeriveExtendedKey() gRPC error code = %v, want %v",
						grpcErr.Code(), tt.wantErr.Code())
					return
				}

				if grpcErr.Message() != tt.wantErr.Message() {
					t.Errorf("DeriveExtendedKey() gRPC error msg = %v, want %v",
						grpcErr.Message(), tt.wantErr.Message())
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeriveExtendedKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}
