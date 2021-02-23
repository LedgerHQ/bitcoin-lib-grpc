package core

import (
	"reflect"
	"testing"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ledgerhq/bitcoin-lib-grpc/pkg/chaincfg"
	"github.com/pkg/errors"
)

// h indicates the BIP32 harden bit, equivalent to 2^31.
var h uint32 = 0x80000000

func TestDeriveExtendedKey(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		derivation []uint32
		want       PublicKeyMaterial
		wantErr    error
	}{
		{
			// BIP0032: Test Vector 1 (chain m/0H/1/2H)
			name:       "ErrDeriveHardFromPublic",
			key:        "xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5",
			derivation: []uint32{0 + h, 1, 2 + h},
			wantErr:    hdkeychain.ErrDeriveHardFromPublic,
		},
		{
			// BIP0032: Test Vector 2 (chain m)
			name:       "mainnet derive from root",
			key:        "xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB",
			derivation: []uint32{0, 2147483647, 1, 2147483646, 2}, // m/0/2147483647/1/2147483646/2
			want: PublicKeyMaterial{
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
			name:       "mainnet derive from account level",
			key:        "xpub6Cc939fyHvfB9pPLWd3bSyyQFvgKbwhidca49jGCM5Hz5ypEPGf9JVXB4NBuUfPgoHnMjN6oNgdC9KRqM11RZtL8QLW6rFKziNwHDYhZ6Kx",
			derivation: []uint32{0, 1}, // m/44'/0'/0'/0/1
			want: PublicKeyMaterial{
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
			name:       "testnet3 derive chain path",
			key:        "tpubDDRojdS4jYQXNugn4t2WLrZ7mjfAyoVQu7MLk4eurqFCbrc7cHLZX8W5YRS8ZskGR9k9t3PqVv68bVBjAyW4nWM9pTGRddt3GQftg6MVQsm",
			derivation: []uint32{2}, // m/0'/1/2'/2
			want: PublicKeyMaterial{
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
		{
			// BIP0032: Test Vector 1 (chain m/0H/1/2H)
			// https://github.com/btcsuite/btcutil/blob/4649e4b/hdkeychain/extendedkey_test.go#L583-L593
			name:       "no derivation",
			key:        "xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5",
			derivation: []uint32{},
			want: PublicKeyMaterial{
				ExtendedKey: "xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5",
				PublicKey: []byte{
					0x03, 0x57, 0xbf, 0xe1, 0xe3, 0x41, 0xd0, 0x1c,
					0x69, 0xfe, 0x56, 0x54, 0x30, 0x99, 0x56, 0xcb,
					0xea, 0x51, 0x68, 0x22, 0xfb, 0xa8, 0xa6, 0x01,
					0x74, 0x3a, 0x01, 0x2a, 0x78, 0x96, 0xee, 0x8d,
					0xc2,
				},
				ChainCode: []byte{
					0x04, 0x46, 0x6b, 0x9c, 0xc8, 0xe1, 0x61, 0xe9,
					0x66, 0x40, 0x9c, 0xa5, 0x29, 0x86, 0xc5, 0x84,
					0xf0, 0x7e, 0x9d, 0xc8, 0x1f, 0x73, 0x5d, 0xb6,
					0x83, 0xc3, 0xff, 0x6e, 0xc7, 0xb1, 0x50, 0x3f,
				},
			},
		},
		{
			// invalid extended key length
			name:       "ErrInvalidKeyLen",
			key:        "deadbeef",
			derivation: []uint32{},
			wantErr:    hdkeychain.ErrInvalidKeyLen,
		},
	}

	s := &Service{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.DeriveExtendedKey(tt.key, tt.derivation)

			if err != nil && tt.wantErr == nil {
				t.Fatalf("DeriveExtendedKey() unexpected error: %v", err)
			}

			if err == nil && tt.wantErr != nil {
				t.Fatalf("DeriveExtendedKey() got no error, want '%v'",
					tt.wantErr)
			}

			if err != nil && tt.wantErr.Error() != errors.Cause(err).Error() {
				t.Fatalf("DeriveExtendedKey() got error '%v', want '%v'",
					err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("DeriveExtendedKey() got error '%v', want '%v'",
					got, tt.want)
			}
		})
	}
}

func TestGetAccountExtendedKey(t *testing.T) {
	tests := []struct {
		name               string
		publicKey          []byte
		chainCode          []byte
		accountIndex       uint32
		chainParams        chaincfg.ChainParams
		want               string
		wantAddress        string
		encodingForAddress AddressEncoding
		wantErr            error
	}{
		{
			// https://github.com/LedgerHQ/lib-ledger-core/blob/54ddf50/core/test/bitcoin/address_test.cpp#L81
			name:         "mainnet legacy",
			accountIndex: 0,
			publicKey: []byte{
				0x02, 0xc3, 0x68, 0xbd, 0xec, 0x47, 0xa1, 0xb6,
				0xfa, 0xa7, 0x6d, 0x62, 0x4e, 0xad, 0x0c, 0xd2,
				0x78, 0x32, 0x34, 0x98, 0x3c, 0x46, 0x67, 0x67,
				0x21, 0x6e, 0xcd, 0xac, 0x8c, 0x47, 0x2d, 0xf3,
				0xa6,
			},
			chainCode: []byte{
				0xb6, 0xb8, 0xa4, 0x9c, 0x62, 0x34, 0xb2, 0x6c,
				0x91, 0xbf, 0xaf, 0xac, 0xd9, 0x05, 0x4c, 0x18,
				0x56, 0x21, 0x30, 0x23, 0x4d, 0xc3, 0x9e, 0x94,
				0x63, 0x56, 0x1c, 0xa6, 0x66, 0x7f, 0x40, 0xf8,
			},
			chainParams:        chaincfg.BitcoinMainNetParams,
			wantAddress:        "14QcTVDFpuGsmNSLeDexB1kWCdoBnTTtgr",
			encodingForAddress: Legacy,
			want:               "xpub6DVHQNhjvVchuKeMGnKbbNSdczQ4yMqEW1H1qhQzk1oPxkSqyHZR9Pn7zZ494sVhZqK2WD8kxo9rqiJFL41P67JCdNYka2W5LnANDVWSjzm",
		},
	}

	s := &Service{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetAccountExtendedKey(
				tt.publicKey, tt.chainCode, tt.accountIndex, tt.chainParams)

			if err != nil && tt.wantErr == nil {
				t.Fatalf("GetAccountExtendedKey() unexpected error: %v", err)
			}

			if err == nil && tt.wantErr != nil {
				t.Fatalf("GetAccountExtendedKey() got no error, want '%v'",
					tt.wantErr)
			}

			if err != nil && tt.wantErr.Error() != errors.Cause(err).Error() {
				t.Fatalf("GetAccountExtendedKey() got error '%v', want '%v'",
					err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GetAccountExtendedKey() got error '%v', want '%v'",
					got, tt.want)
			}

			deriveAddress := func(key string, derivation []uint32) (string, error) {
				keyMaterial, err := s.DeriveExtendedKey(key, derivation)
				if err != nil {
					return "", err
				}

				addr, err := s.EncodeAddress(
					keyMaterial.PublicKey, tt.encodingForAddress, tt.chainParams)
				if err != nil {
					return "", err
				}

				return addr, nil
			}

			gotAddress, err := deriveAddress(tt.want, []uint32{0, 0})
			if err != nil {
				t.Fatalf("error: '%v' - cannot derive for '%v' at path 0/0",
					err.Error(), tt.want)
			}

			if gotAddress != tt.wantAddress {
				t.Fatalf("GetAccountExtendedKey() got wrong address: '%v', want '%v'",
					gotAddress, tt.wantAddress)
			}
		})
	}
}

func TestGetKeypair(t *testing.T) {
	tests := []struct {
		name        string
		seed        string
		chainParams chaincfg.ChainParams
		derivation  []uint32
		want        Keypair
		wantErr     error
	}{
		{
			name:        "Get keypair from a seed",
			seed:        "I am the Lama from Lama land",
			chainParams: chaincfg.BitcoinMainNetParams,
			derivation:  []uint32{44 + h, 0 + h, 0 + h},
			want: Keypair{
				ExtendedPublicKey: "xpub6CuV4qnYG4mQb6Q4qHy4dnovUzrt9PXGzA9v7yPjYeTKuQjACFXCFQbkFfvCTz8WsR3ggq7MaNDkwLjvoy6FZby3rZ9PLNGy51rVFdmwhrZ",
				PrivateKey:        "xprv9yv8fLFeRhD7NcKbjGS4GesBvy2PjvoRcwEKKaz7zJvM2cQ1eiCwhcHGQNEBwsXthHbPtZNQg5SBBEWS1QH941SKitBdaUT7VDTxzdS8vu7",
			},
		},
	}

	s := &Service{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetKeypair(tt.seed, tt.chainParams, tt.derivation)

			if err != nil && tt.wantErr == nil {
				t.Fatalf("GetKeypair() unexpected error: %v", err)
			}

			if err == nil && tt.wantErr != nil {
				t.Fatalf("GetKeypair() got no error, want '%v'",
					tt.wantErr)
			}

			if err != nil && tt.wantErr.Error() != errors.Cause(err).Error() {
				t.Fatalf("GetKeypair() got error '%v', want '%v'",
					err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GetKeypair() got error '%v', want '%v'",
					got, tt.want)
			}
		})
	}
}
