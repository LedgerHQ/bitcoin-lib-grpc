package bitcoin

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/btcsuite/btcd/chaincfg"
	pb "github.com/ledgerhq/lama-bitcoin-svc/pb/v1"
)

var (
	mainnetChainParamsProto = &pb.ChainParams{
		Network: &pb.ChainParams_BitcoinNetwork{
			BitcoinNetwork: pb.BitcoinNetwork_BITCOIN_NETWORK_MAINNET,
		},
	}

	testnet3ChainParamsProto = &pb.ChainParams{
		Network: &pb.ChainParams_BitcoinNetwork{
			BitcoinNetwork: pb.BitcoinNetwork_BITCOIN_NETWORK_TESTNET3,
		},
	}

	regtestChainParamsProto = &pb.ChainParams{
		Network: &pb.ChainParams_BitcoinNetwork{
			BitcoinNetwork: pb.BitcoinNetwork_BITCOIN_NETWORK_REGTEST,
		},
	}

	invalidChainParamsProto = &pb.ChainParams{
		Network: &pb.ChainParams_BitcoinNetwork{
			BitcoinNetwork: 99999,
		},
	}
)

func Test_getChainParams(t *testing.T) {
	tests := []struct {
		name             string
		chainParamsProto *pb.ChainParams
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
		request *pb.ValidateAddressRequest
		want    *pb.ValidateAddressResponse
		wantErr *status.Status
	}{
		{
			name: "mainnet P2PKH valid",
			request: &pb.ValidateAddressRequest{
				Address:     "1MirQ9bwyQcGVJPwKUgapu5ouK2E2Ey4gX",
				ChainParams: mainnetChainParamsProto,
			},
			want: &pb.ValidateAddressResponse{
				Address: "1MirQ9bwyQcGVJPwKUgapu5ouK2E2Ey4gX",
				IsValid: true,
			},
			wantErr: nil,
		},
		{
			name: "mainnet P2WPKH invalid checksum",
			request: &pb.ValidateAddressRequest{
				Address:     "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t5",
				ChainParams: mainnetChainParamsProto,
			},
			want: &pb.ValidateAddressResponse{
				Address:       "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t5",
				IsValid:       false,
				InvalidReason: "checksum failed. Expected v8f3t4, got v8f3t5.",
			},
			wantErr: nil,
		},
		{
			name: "testnet3 P2WPKH invalid mixed case",
			request: &pb.ValidateAddressRequest{
				Address:     "tb1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3q0sL5k7",
				ChainParams: testnet3ChainParamsProto,
			},
			want: &pb.ValidateAddressResponse{
				Address:       "tb1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3q0sL5k7",
				IsValid:       false,
				InvalidReason: "string not all lowercase or all uppercase",
			},
			wantErr: nil,
		},
		{
			name: "get unknown chain params",
			request: &pb.ValidateAddressRequest{
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

func TestEncodeAddress(t *testing.T) {
	// Helper to derive extended key and return the serialized public key.
	// Use this in unit-tests to ensure extended key derivation and address
	// encoding work together as expected.
	derivePublicKey := func(extendedKey string, derivation []uint32) []byte {
		ctx := context.Background()
		s := &Service{}

		request := &pb.DeriveExtendedKeyRequest{
			ExtendedKey: extendedKey,
			Derivation:  derivation,
		}
		response, err := s.DeriveExtendedKey(ctx, request)
		if err != nil {
			panic(err)
		}

		return response.PublicKey
	}

	tests := []struct {
		name    string
		request *pb.EncodeAddressRequest
		want    *pb.EncodeAddressResponse
		wantErr *status.Status
	}{
		{
			// https://github.com/LedgerHQ/lib-ledger-core/blob/978a496/core/test/bitcoin/address_test.cpp#L89
			name: "xpub P2PKH",
			request: &pb.EncodeAddressRequest{
				PublicKey: derivePublicKey(
					"xpub6Cc939fyHvfB9pPLWd3bSyyQFvgKbwhidca49jGCM5Hz5ypEPGf9JVXB4NBuUfPgoHnMjN6oNgdC9KRqM11RZtL8QLW6rFKziNwHDYhZ6Kx",
					[]uint32{1, 1},
				),
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2PKH,
				ChainParams: mainnetChainParamsProto,
			},
			want: &pb.EncodeAddressResponse{
				Address: "1AkRBkUZQe5Zqj5syxn1cHCvKUV6DjL9Po",
			},
			wantErr: nil,
		},
		{
			// https://github.com/trezor/blockbook/blob/7919486/bchain/coins/btc/bitcoinparser_test.go#L537-L546
			name: "ypub P2SH-P2WPKH",
			request: &pb.EncodeAddressRequest{
				PublicKey: derivePublicKey(
					"ypub6Ww3ibxVfGzLrAH1PNcjyAWenMTbbAosGNB6VvmSEgytSER9azLDWCxoJwW7Ke7icmizBMXrzBx9979FfaHxHcrArf3zbeJJJUZPf663zsP",
					[]uint32{0, 0}),
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2SH_P2WPKH,
				ChainParams: mainnetChainParamsProto,
			},
			want: &pb.EncodeAddressResponse{
				Address: "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf",
			},
			wantErr: nil,
		},
		{
			// https://github.com/trezor/blockbook/blob/7919486/bchain/coins/btc/bitcoinparser_test.go#L549-L558
			name: "zpub P2WPKH",
			request: &pb.EncodeAddressRequest{
				PublicKey: derivePublicKey(
					"zpub6rFR7y4Q2AijBEqTUquhVz398htDFrtymD9xYYfG1m4wAcvPhXNfE3EfH1r1ADqtfSdVCToUG868RvUUkgDKf31mGDtKsAYz2oz2AGutZYs",
					[]uint32{0, 0}),
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2WPKH,
				ChainParams: mainnetChainParamsProto,
			},
			want: &pb.EncodeAddressResponse{
				Address: "bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu",
			},
			wantErr: nil,
		},
		{
			// generated using https://gist.github.com/onyb/c022bc1a35aae47a327ce5356f2c6a31
			name: "tpub P2PKH",
			request: &pb.EncodeAddressRequest{
				PublicKey: derivePublicKey(
					"tpubDC5FSnBiZDMmhiuCmWAYsLwgLYrrT9rAqvTySfuCCrgsWz8wxMXUS9Tb9iVMvcRbvFcAHGkMD5Kx8koh4GquNGNTfohfk7pgjhaPCdXpoba",
					[]uint32{0, 0}),
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2PKH,
				ChainParams: testnet3ChainParamsProto,
			},
			want: &pb.EncodeAddressResponse{
				Address: "mkpZhYtJu2r87Js3pDiWJDmPte2NRZ8bJV",
			},
			wantErr: nil,
		},
		{
			// https://github.com/trezor/blockbook/blob/7919486/bchain/coins/btc/bitcoinparser_test.go#L559-L569
			name: "upub P2SH-P2WPKH",
			request: &pb.EncodeAddressRequest{
				PublicKey: derivePublicKey(
					"upub5DR1Mg5nykixzYjFXWW5GghAU7dDqoPVJ2jrqFbL8sJ7Hs7jn69MP7KBnnmxn88GeZtnH8PRKV9w5MMSFX8AdEAoXY8Qd8BJPoXtpMeHMxJ",
					[]uint32{0, 0}),
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2SH_P2WPKH,
				ChainParams: testnet3ChainParamsProto,
			},
			want: &pb.EncodeAddressResponse{
				Address: "2N4Q5FhU2497BryFfUgbqkAJE87aKHUhXMp",
			},
			wantErr: nil,
		},
		{
			// generated using https://gist.github.com/onyb/c022bc1a35aae47a327ce5356f2c6a31
			name: "vpub P2WPKH",
			request: &pb.EncodeAddressRequest{
				PublicKey: derivePublicKey(
					"vpub5Y6cjg78GGuNLsaPhmYsiw4gYX3HoQiRBiSwDaBXKUafCt9bNwWQiitDk5VZ5BVxYnQdwoTyXSs2JHRPAgjAvtbBrf8ZhDYe2jWAqvZVnsc",
					[]uint32{1, 1}),
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2WPKH,
				ChainParams: testnet3ChainParamsProto,
			},
			want: &pb.EncodeAddressResponse{
				Address: "tb1qkwgskuzmmwwvqajnyr7yp9hgvh5y45kg8wvdmd",
			},
			wantErr: nil,
		},
		{
			// https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
			name: "pubkey P2PKH mainnet",
			request: &pb.EncodeAddressRequest{
				PublicKey: []byte{
					0x02, 0x50, 0x86, 0x3a, 0xd6, 0x4a, 0x87, 0xae,
					0x8a, 0x2f, 0xe8, 0x3c, 0x1a, 0xf1, 0xa8, 0x40,
					0x3c, 0xb5, 0x3f, 0x53, 0xe4, 0x86, 0xd8, 0x51,
					0x1d, 0xad, 0x8a, 0x04, 0x88, 0x7e, 0x5b, 0x23,
					0x52,
				},
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2PKH,
				ChainParams: mainnetChainParamsProto,
			},
			want: &pb.EncodeAddressResponse{
				Address: "1PMycacnJaSqwwJqjawXBErnLsZ7RkXUAs",
			},
			wantErr: nil,
		},
		{
			// http://bitcoinscri.pt/pages/segwit_p2sh_p2wpkh_address
			name: "pubkey P2SH-P2WPKH mainnet",
			request: &pb.EncodeAddressRequest{
				PublicKey: []byte{
					0x02, 0xf1, 0x18, 0xcc, 0x40, 0x97, 0x75, 0x41,
					0x9a, 0x93, 0x1c, 0x57, 0x66, 0x4d, 0x0c, 0x19,
					0xc4, 0x05, 0xe8, 0x56, 0xac, 0x0e, 0xe2, 0xf0,
					0xe2, 0xa4, 0x13, 0x7d, 0x82, 0x50, 0x53, 0x11,
					0x28,
				},
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2SH_P2WPKH,
				ChainParams: mainnetChainParamsProto,
			},
			want: &pb.EncodeAddressResponse{
				Address: "3Mwz6cg8Fz81B7ukexK8u8EVAW2yymgWNd",
			},
			wantErr: nil,
		},
		{
			// http://bitcoinscri.pt/pages/segwit_native_p2wpkh_address
			name: "pubkey P2WPKH mainnet",
			request: &pb.EncodeAddressRequest{
				PublicKey: []byte{
					0x02, 0x53, 0x0c, 0x54, 0x8d, 0x40, 0x26, 0x70,
					0xb1, 0x3a, 0xd8, 0x88, 0x7f, 0xf9, 0x9c, 0x29,
					0x4e, 0x67, 0xfc, 0x18, 0x09, 0x7d, 0x23, 0x6d,
					0x57, 0x88, 0x0c, 0x69, 0x26, 0x1b, 0x42, 0xde,
					0xf7,
				},
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2WPKH,
				ChainParams: mainnetChainParamsProto,
			},
			want: &pb.EncodeAddressResponse{
				Address: "bc1qg9stkxrszkdqsuj92lm4c7akvk36zvhqw7p6ck",
			},
			wantErr: nil,
		},
		// TODO: add pubkey P2PKH testnet3
		// TODO: add pubkey P2SH-P2WPKH testnet3
		{
			// https://github.com/libbitcoin/libbitcoin-system/blob/4dda9d0/test/wallet/witness_address.cpp#L82-L93
			name: "pubkey P2WPKH testnet3",
			request: &pb.EncodeAddressRequest{
				PublicKey: []byte{
					0x03, 0x82, 0x62, 0xa6, 0xc6, 0xce, 0xc9, 0x3c,
					0x2d, 0x3e, 0xcd, 0x6c, 0x60, 0x72, 0xef, 0xea,
					0x86, 0xd0, 0x2f, 0xf8, 0xe3, 0x32, 0x8b, 0xbd,
					0x02, 0x42, 0xb2, 0x0a, 0xf3, 0x42, 0x59, 0x90,
					0xac,
				},
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2WPKH,
				ChainParams: testnet3ChainParamsProto,
			},
			want: &pb.EncodeAddressResponse{
				Address: "tb1qr47dd36u96r0fjle36hdygdnp0v6pwfgqe6jxg",
			},
			wantErr: nil,
		},
		{
			// https://github.com/LedgerHQ/lib-ledger-core/blob/8c068f/core/test/integration/BaseFixture.cpp#L75-L77
			// https://github.com/LedgerHQ/lib-ledger-core/blob/2e5500a/core/test/integration/keychains/p2sh_keychain_test.cpp#L43-L48
			name: "non-standard key P2SH-P2WPKH",
			request: &pb.EncodeAddressRequest{
				PublicKey: derivePublicKey(
					"tpubDCcvqEHx7prGddpWTfEviiew5YLMrrKy4oJbt14teJZenSi6AYMAs2SNXwYXFzkrNYwECSmobwxESxMCrpfqw4gsUt88bcr8iMrJmbb8P2q",
					[]uint32{0, 0}),
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2SH_P2WPKH,
				ChainParams: testnet3ChainParamsProto,
			},
			want: &pb.EncodeAddressResponse{
				Address: "2MvuUMAG1NFQmmM69Writ6zTsYCnQHFG9BF",
			},
			wantErr: nil,
		},
		{
			// https://github.com/LedgerHQ/lib-ledger-core/blob/8c068fc/core/test/integration/BaseFixture.cpp#L42-L44
			// Instances of the address spread across the integration tests.
			name: "non-standard key P2WPKH",
			request: &pb.EncodeAddressRequest{
				PublicKey: derivePublicKey(
					"xpub6CMeLkY9TzXyLYXPWMXB5LWtprVABb6HwPEPXnEgESMNrSUBsvhXNsA7zKS1ZRKhUyQG4HjZysEP8v7gDNU4J6PvN5yLx4meEm3mpEapLMN",
					[]uint32{0, 0}),
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2WPKH,
				ChainParams: mainnetChainParamsProto,
			},
			want: &pb.EncodeAddressResponse{
				Address: "bc1qh4kl0a0a3d7su8udc2rn62f8w939prqpl34z86",
			},
			wantErr: nil,
		},
		{
			// https://github.com/LedgerHQ/lib-ledger-core/blob/8c068fc/core/test/integration/BaseFixture.cpp#L46-L50
			// Address verified on https://iancoleman.io/bitcoin-key-compression
			name: "uncompressed pubkey P2PKH",
			request: &pb.EncodeAddressRequest{
				PublicKey: []byte{
					0x04, 0x37, 0xbc, 0x83, 0xa3, 0x77, 0xea, 0x02,
					0x5e, 0x53, 0xea, 0xfc, 0xd1, 0x8f, 0x29, 0x92,
					0x68, 0xd1, 0xce, 0xca, 0xe8, 0x9b, 0x4f, 0x15,
					0x40, 0x19, 0x26, 0xa0, 0xf8, 0xb0, 0x06, 0xc0,
					0xf7, 0xee, 0x1b, 0x99, 0x50, 0x47, 0xb3, 0xe1,
					0x59, 0x59, 0xc5, 0xd1, 0x0d, 0xd1, 0x56, 0x3e,
					0x22, 0xa2, 0xe6, 0xe4, 0xbe, 0x95, 0x72, 0xaa,
					0x70, 0x78, 0xe3, 0x2f, 0x31, 0x76, 0x77, 0xa9,
					0x01,
				},
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2PKH,
				ChainParams: mainnetChainParamsProto,
			},
			want: &pb.EncodeAddressResponse{
				Address: "18iytmdAvJcQwCHWfWppDB5hR3YHNsYhRr",
			},
			wantErr: nil,
		},
		{
			name: "xpub P2PKH invalid chain params",
			request: &pb.EncodeAddressRequest{
				PublicKey: derivePublicKey(
					"xpub6Cc939fyHvfB9pPLWd3bSyyQFvgKbwhidca49jGCM5Hz5ypEPGf9JVXB4NBuUfPgoHnMjN6oNgdC9KRqM11RZtL8QLW6rFKziNwHDYhZ6Kx",
					[]uint32{1, 1},
				),
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2PKH,
				ChainParams: invalidChainParamsProto,
			},
			want: nil,
			wantErr: status.New(codes.InvalidArgument,
				ErrUnknownNetwork("99999").Error()),
		},
		{
			name: "invalid length public key",
			request: &pb.EncodeAddressRequest{
				PublicKey:   []byte{0x04},
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2WPKH,
				ChainParams: mainnetChainParamsProto,
			},
			want: nil,
			wantErr: status.New(codes.InvalidArgument,
				"invalid pub key length 1"),
		},
		{
			// https://github.com/btcsuite/btcd/blob/69773a7/btcec/pubkey_test.go#L162-L174
			name: "invalid public key (X > P)",
			request: &pb.EncodeAddressRequest{
				PublicKey: []byte{
					0x04, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xFD,
					0x2F, 0xb2, 0xe0, 0xea, 0xdd, 0xfb, 0x84, 0xcc,
					0xf9, 0x74, 0x44, 0x64, 0xf8, 0x2e, 0x16, 0x0b,
					0xfa, 0x9b, 0x8b, 0x64, 0xf9, 0xd4, 0xc0, 0x3f,
					0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56, 0xb4, 0x12,
					0xa3,
				},
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2WPKH,
				ChainParams: mainnetChainParamsProto,
			},
			want: nil,
			wantErr: status.New(codes.InvalidArgument,
				"pubkey X parameter is >= to P"),
		},
		{
			// https://github.com/decred/dcrd/blob/b60c60f/dcrec/secp256k1/pubkey_test.go#L105-L109
			name: "invalid public key (not on curve)",
			request: &pb.EncodeAddressRequest{
				PublicKey: []byte{
					0x03, 0xce, 0x0b, 0x14, 0xfb, 0x84, 0x2b, 0x1b,
					0xa5, 0x49, 0xfd, 0xd6, 0x75, 0xc9, 0x80, 0x75,
					0xf1, 0x2e, 0x9c, 0x51, 0x0f, 0x8e, 0xf5, 0x2b,
					0xd0, 0x21, 0xa9, 0xa1, 0xf4, 0x80, 0x9d, 0x3b,
					0x4c,
				},
				Encoding:    pb.AddressEncoding_ADDRESS_ENCODING_P2WPKH,
				ChainParams: mainnetChainParamsProto,
			},
			want: nil,
			wantErr: status.New(codes.InvalidArgument,
				"invalid square root"), // FIXME: Improve error in btcd
		},
	}

	ctx := context.Background()
	s := &Service{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.EncodeAddress(ctx, tt.request)
			grpcErr := status.Convert(err)

			if grpcErr != nil && tt.wantErr == nil {
				t.Errorf("unexpected error in EncodeAddress(): %v", grpcErr.Message())
				return
			}

			if tt.wantErr != nil {
				if grpcErr.Code() != tt.wantErr.Code() {
					t.Errorf("EncodeAddress() gRPC error code = %v, want %v",
						grpcErr.Code(), tt.wantErr.Code())
					return
				}

				if grpcErr.Message() != tt.wantErr.Message() {
					t.Errorf("EncodeAddress() gRPC error msg = %v, want %v",
						grpcErr.Message(), tt.wantErr.Message())
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
