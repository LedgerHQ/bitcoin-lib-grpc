package bitcoin

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestService_ValidateAddress(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		chainParams ChainParams
		want        string
		wantErr     error
	}{
		{
			name:        "mainnet P2PKH valid",
			address:     "1MirQ9bwyQcGVJPwKUgapu5ouK2E2Ey4gX",
			chainParams: Mainnet,
			want:        "1MirQ9bwyQcGVJPwKUgapu5ouK2E2Ey4gX",
		},
		{
			name:        "mainnet P2WPKH invalid checksum",
			address:     "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t5",
			chainParams: Mainnet,
			wantErr:     errors.New("checksum failed. Expected v8f3t4, got v8f3t5."),
		},
		{
			name:        "testnet3 P2WPKH invalid mixed case",
			address:     "tb1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3q0sL5k7",
			chainParams: Testnet3,
			wantErr:     errors.New("string not all lowercase or all uppercase"),
		},
	}

	s := &Service{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ValidateAddress(tt.address, tt.chainParams)
			if err != nil && tt.wantErr == nil {
				t.Fatalf("ValidateAddress() unexpected error: %v", err)
			}

			if err == nil && tt.wantErr != nil {
				t.Fatalf("ValidateAddress() got no error, want '%v'",
					tt.wantErr)
			}

			if err != nil && tt.wantErr.Error() != errors.Cause(err).Error() {
				t.Fatalf("ValidateAddress() got error '%v', want '%v'",
					err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ValidateAddress() got error '%v', want '%v'",
					got, tt.want)
			}
		})
	}
}

func TestEncodeAddress(t *testing.T) {
	// Helper to derive extended key and return the serialized public key.
	// Use this in unit-tests to ensure extended key derivation and address
	// encoding work together as expected.
	derivePublicKey := func(extendedKey string, derivation []uint32) []byte {
		s := &Service{}

		response, err := s.DeriveExtendedKey(extendedKey, derivation)
		if err != nil {
			panic(err)
		}

		return response.PublicKey
	}

	tests := []struct {
		name        string
		publicKey   []byte
		encoding    AddressEncoding
		chainParams ChainParams
		want        string
		wantErr     error
	}{
		{
			// https://github.com/LedgerHQ/lib-ledger-core/blob/978a496/core/test/bitcoin/address_test.cpp#L89
			name: "xpub P2PKH",
			publicKey: derivePublicKey(
				"xpub6Cc939fyHvfB9pPLWd3bSyyQFvgKbwhidca49jGCM5Hz5ypEPGf9JVXB4NBuUfPgoHnMjN6oNgdC9KRqM11RZtL8QLW6rFKziNwHDYhZ6Kx",
				[]uint32{1, 1},
			),
			encoding:    Legacy,
			chainParams: Mainnet,
			want:        "1AkRBkUZQe5Zqj5syxn1cHCvKUV6DjL9Po",
		},
		{
			// https://github.com/trezor/blockbook/blob/7919486/bchain/coins/btc/bitcoinparser_test.go#L537-L546
			name: "ypub P2SH-P2WPKH",
			publicKey: derivePublicKey(
				"ypub6Ww3ibxVfGzLrAH1PNcjyAWenMTbbAosGNB6VvmSEgytSER9azLDWCxoJwW7Ke7icmizBMXrzBx9979FfaHxHcrArf3zbeJJJUZPf663zsP",
				[]uint32{0, 0},
			),
			encoding:    WrappedSegwit,
			chainParams: Mainnet,
			want:        "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf",
		},
		{
			// https://github.com/trezor/blockbook/blob/7919486/bchain/coins/btc/bitcoinparser_test.go#L549-L558
			name: "zpub P2WPKH",
			publicKey: derivePublicKey(
				"zpub6rFR7y4Q2AijBEqTUquhVz398htDFrtymD9xYYfG1m4wAcvPhXNfE3EfH1r1ADqtfSdVCToUG868RvUUkgDKf31mGDtKsAYz2oz2AGutZYs",
				[]uint32{0, 0},
			),
			encoding:    NativeSegwit,
			chainParams: Mainnet,
			want:        "bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu",
		},
		{
			// generated using https://gist.github.com/onyb/c022bc1a35aae47a327ce5356f2c6a31
			name: "tpub P2PKH",
			publicKey: derivePublicKey(
				"tpubDC5FSnBiZDMmhiuCmWAYsLwgLYrrT9rAqvTySfuCCrgsWz8wxMXUS9Tb9iVMvcRbvFcAHGkMD5Kx8koh4GquNGNTfohfk7pgjhaPCdXpoba",
				[]uint32{0, 0}),
			encoding:    Legacy,
			chainParams: Testnet3,
			want:        "mkpZhYtJu2r87Js3pDiWJDmPte2NRZ8bJV",
		},
		{
			// https://github.com/trezor/blockbook/blob/7919486/bchain/coins/btc/bitcoinparser_test.go#L559-L569
			name: "upub P2SH-P2WPKH",
			publicKey: derivePublicKey(
				"upub5DR1Mg5nykixzYjFXWW5GghAU7dDqoPVJ2jrqFbL8sJ7Hs7jn69MP7KBnnmxn88GeZtnH8PRKV9w5MMSFX8AdEAoXY8Qd8BJPoXtpMeHMxJ",
				[]uint32{0, 0}),
			encoding:    WrappedSegwit,
			chainParams: Testnet3,
			want:        "2N4Q5FhU2497BryFfUgbqkAJE87aKHUhXMp",
		},
		{
			// generated using https://gist.github.com/onyb/c022bc1a35aae47a327ce5356f2c6a31
			name: "vpub P2WPKH",
			publicKey: derivePublicKey(
				"vpub5Y6cjg78GGuNLsaPhmYsiw4gYX3HoQiRBiSwDaBXKUafCt9bNwWQiitDk5VZ5BVxYnQdwoTyXSs2JHRPAgjAvtbBrf8ZhDYe2jWAqvZVnsc",
				[]uint32{1, 1}),
			encoding:    NativeSegwit,
			chainParams: Testnet3,
			want:        "tb1qkwgskuzmmwwvqajnyr7yp9hgvh5y45kg8wvdmd",
		},
		{
			// https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
			name: "pubkey P2PKH mainnet",
			publicKey: []byte{
				0x02, 0x50, 0x86, 0x3a, 0xd6, 0x4a, 0x87, 0xae,
				0x8a, 0x2f, 0xe8, 0x3c, 0x1a, 0xf1, 0xa8, 0x40,
				0x3c, 0xb5, 0x3f, 0x53, 0xe4, 0x86, 0xd8, 0x51,
				0x1d, 0xad, 0x8a, 0x04, 0x88, 0x7e, 0x5b, 0x23,
				0x52,
			},
			encoding:    Legacy,
			chainParams: Mainnet,
			want:        "1PMycacnJaSqwwJqjawXBErnLsZ7RkXUAs",
		},
		{
			// http://bitcoinscri.pt/pages/segwit_p2sh_p2wpkh_address
			name: "pubkey P2SH-P2WPKH mainnet",
			publicKey: []byte{
				0x02, 0xf1, 0x18, 0xcc, 0x40, 0x97, 0x75, 0x41,
				0x9a, 0x93, 0x1c, 0x57, 0x66, 0x4d, 0x0c, 0x19,
				0xc4, 0x05, 0xe8, 0x56, 0xac, 0x0e, 0xe2, 0xf0,
				0xe2, 0xa4, 0x13, 0x7d, 0x82, 0x50, 0x53, 0x11,
				0x28,
			},
			encoding:    WrappedSegwit,
			chainParams: Mainnet,
			want:        "3Mwz6cg8Fz81B7ukexK8u8EVAW2yymgWNd",
		},
		{
			// http://bitcoinscri.pt/pages/segwit_native_p2wpkh_address
			name: "pubkey P2WPKH mainnet",
			publicKey: []byte{
				0x02, 0x53, 0x0c, 0x54, 0x8d, 0x40, 0x26, 0x70,
				0xb1, 0x3a, 0xd8, 0x88, 0x7f, 0xf9, 0x9c, 0x29,
				0x4e, 0x67, 0xfc, 0x18, 0x09, 0x7d, 0x23, 0x6d,
				0x57, 0x88, 0x0c, 0x69, 0x26, 0x1b, 0x42, 0xde,
				0xf7,
			},
			encoding:    NativeSegwit,
			chainParams: Mainnet,
			want:        "bc1qg9stkxrszkdqsuj92lm4c7akvk36zvhqw7p6ck",
		},
		// TODO: add pubkey P2PKH testnet3
		// TODO: add pubkey P2SH-P2WPKH testnet3
		{
			// https://github.com/libbitcoin/libbitcoin-system/blob/4dda9d0/test/wallet/witness_address.cpp#L82-L93
			name: "pubkey P2WPKH testnet3",
			publicKey: []byte{
				0x03, 0x82, 0x62, 0xa6, 0xc6, 0xce, 0xc9, 0x3c,
				0x2d, 0x3e, 0xcd, 0x6c, 0x60, 0x72, 0xef, 0xea,
				0x86, 0xd0, 0x2f, 0xf8, 0xe3, 0x32, 0x8b, 0xbd,
				0x02, 0x42, 0xb2, 0x0a, 0xf3, 0x42, 0x59, 0x90,
				0xac,
			},
			encoding:    NativeSegwit,
			chainParams: Testnet3,
			want:        "tb1qr47dd36u96r0fjle36hdygdnp0v6pwfgqe6jxg",
		},
		{
			// https://github.com/LedgerHQ/lib-ledger-core/blob/8c068f/core/test/integration/BaseFixture.cpp#L75-L77
			// https://github.com/LedgerHQ/lib-ledger-core/blob/2e5500a/core/test/integration/keychains/p2sh_keychain_test.cpp#L43-L48
			name: "non-standard key P2SH-P2WPKH",
			publicKey: derivePublicKey(
				"tpubDCcvqEHx7prGddpWTfEviiew5YLMrrKy4oJbt14teJZenSi6AYMAs2SNXwYXFzkrNYwECSmobwxESxMCrpfqw4gsUt88bcr8iMrJmbb8P2q",
				[]uint32{0, 0}),
			encoding:    WrappedSegwit,
			chainParams: Testnet3,
			want:        "2MvuUMAG1NFQmmM69Writ6zTsYCnQHFG9BF",
		},
		{
			// https://github.com/LedgerHQ/lib-ledger-core/blob/8c068fc/core/test/integration/BaseFixture.cpp#L42-L44
			// Instances of the address spread across the integration tests.
			name: "non-standard key P2WPKH",
			publicKey: derivePublicKey(
				"xpub6CMeLkY9TzXyLYXPWMXB5LWtprVABb6HwPEPXnEgESMNrSUBsvhXNsA7zKS1ZRKhUyQG4HjZysEP8v7gDNU4J6PvN5yLx4meEm3mpEapLMN",
				[]uint32{0, 0}),
			encoding:    NativeSegwit,
			chainParams: Mainnet,
			want:        "bc1qh4kl0a0a3d7su8udc2rn62f8w939prqpl34z86",
		},
		{
			// https://github.com/LedgerHQ/lib-ledger-core/blob/8c068fc/core/test/integration/BaseFixture.cpp#L46-L50
			// Address verified on https://iancoleman.io/bitcoin-key-compression
			name: "uncompressed pubkey P2PKH",
			publicKey: []byte{
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
			encoding:    Legacy,
			chainParams: Mainnet,
			want:        "18iytmdAvJcQwCHWfWppDB5hR3YHNsYhRr",
			wantErr:     nil,
		},
		{
			name:        "invalid length public key",
			publicKey:   []byte{0x04},
			encoding:    NativeSegwit,
			chainParams: Mainnet,
			wantErr:     errors.New("invalid pub key length 1"),
		},
		{
			// https://github.com/btcsuite/btcd/blob/69773a7/btcec/pubkey_test.go#L162-L174
			name: "invalid public key (X > P)",
			publicKey: []byte{
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
			encoding:    NativeSegwit,
			chainParams: Mainnet,
			wantErr:     errors.New("pubkey X parameter is >= to P"),
		},
		{
			// https://github.com/decred/dcrd/blob/b60c60f/dcrec/secp256k1/pubkey_test.go#L105-L109
			name: "invalid public key (not on curve)",
			publicKey: []byte{
				0x03, 0xce, 0x0b, 0x14, 0xfb, 0x84, 0x2b, 0x1b,
				0xa5, 0x49, 0xfd, 0xd6, 0x75, 0xc9, 0x80, 0x75,
				0xf1, 0x2e, 0x9c, 0x51, 0x0f, 0x8e, 0xf5, 0x2b,
				0xd0, 0x21, 0xa9, 0xa1, 0xf4, 0x80, 0x9d, 0x3b,
				0x4c,
			},
			encoding:    NativeSegwit,
			chainParams: Mainnet,
			wantErr:     errors.New("invalid square root"), // FIXME: Improve error in btcd
		},
		{
			name: "ErrUnknownAddressType",
			publicKey: []byte{
				0x02, 0x53, 0x0c, 0x54, 0x8d, 0x40, 0x26, 0x70,
				0xb1, 0x3a, 0xd8, 0x88, 0x7f, 0xf9, 0x9c, 0x29,
				0x4e, 0x67, 0xfc, 0x18, 0x09, 0x7d, 0x23, 0x6d,
				0x57, 0x88, 0x0c, 0x69, 0x26, 0x1b, 0x42, 0xde,
				0xf7,
			},
			encoding:    9999,
			chainParams: Mainnet,
			wantErr:     ErrUnknownAddressType,
		},
	}

	s := &Service{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.EncodeAddress(tt.publicKey, tt.encoding, tt.chainParams)
			if err != nil && tt.wantErr == nil {
				t.Fatalf("EncodeAddress() unexpected error: %v", err)
			}

			if err == nil && tt.wantErr != nil {
				t.Fatalf("EncodeAddress() got no error, want '%v'",
					tt.wantErr)
			}

			if err != nil && tt.wantErr.Error() != errors.Cause(err).Error() {
				t.Fatalf("EncodeAddress() got error '%v', want '%v'",
					err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("EncodeAddress() got error '%v', want '%v'",
					got, tt.want)
			}
		})
	}
}
