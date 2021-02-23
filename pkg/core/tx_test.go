package core

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ledgerhq/bitcoin-lib-grpc/pkg/chaincfg"
)

func TestCreateTransaction(t *testing.T) {
	tests := []struct {
		name                    string
		tx                      *Tx
		chainParams             chaincfg.ChainParams
		wantErr                 error
		wantNotEnoughUtxoAmount *NotEnoughUtxo
	}{
		{
			name: "mainnet P2WPKH",
			tx: &Tx{
				LockTime: 0,
				Inputs: []Input{
					{
						OutputHash:  "2f5dae23c2e18588c86cfc4e154f3b68bd8eb4265fe0b4b1341ad5aa40422f66",
						OutputIndex: 0,
						Script:      []byte("76a914e18c90d108c3509e952c1d79121f1776facf1c6788ac"),
						Value:       110000,
					},
				},
				Outputs: []Output{
					{
						Address: "1MZbRqZGpiSWGRLg8DUdVrDKHwNe1oesUZ",
						Value:   100000,
					},
				},
				ChangeAddress: "1GgX4cGLiqF9p4Sd1XcPQhEAAhNDA4wLYS",
				FeeSatPerKb:   1234,
			},
			chainParams: chaincfg.BitcoinMainNetParams,
		},
		{
			name: "mainnet P2WPKH with not enough utxo",
			tx: &Tx{
				LockTime: 0,
				Inputs: []Input{
					{
						OutputHash:  "2f5dae23c2e18588c86cfc4e154f3b68bd8eb4265fe0b4b1341ad5aa40422f66",
						OutputIndex: 0,
						Script:      []byte("76a914e18c90d108c3509e952c1d79121f1776facf1c6788ac"),
						Value:       100000,
					},
				},
				Outputs: []Output{
					{
						Address: "1MZbRqZGpiSWGRLg8DUdVrDKHwNe1oesUZ",
						Value:   100000,
					},
				},
				ChangeAddress: "1GgX4cGLiqF9p4Sd1XcPQhEAAhNDA4wLYS",
				FeeSatPerKb:   1234,
			},
			chainParams:             chaincfg.BitcoinMainNetParams,
			wantNotEnoughUtxoAmount: &NotEnoughUtxo{MissingAmount: 134},
		},
	}

	s := &Service{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawTx, err := s.CreateTransaction(tt.tx, tt.chainParams)
			if err != nil && tt.wantErr == nil {
				t.Fatalf("CreateTransaction() got error '%v'", err)
			}

			if rawTx.NotEnoughUtxo != nil && tt.wantNotEnoughUtxoAmount == nil {
				t.Fatalf("CreateTransaction() got NotEnoughUtxo '%v'", rawTx.NotEnoughUtxo)
			}

			if rawTx.NotEnoughUtxo == nil && tt.wantNotEnoughUtxoAmount != nil {
				t.Fatalf("CreateTransaction() should have '%v'", tt.wantNotEnoughUtxoAmount)
			}

			if rawTx.NotEnoughUtxo != nil && tt.wantNotEnoughUtxoAmount != nil && rawTx.NotEnoughUtxo.MissingAmount != tt.wantNotEnoughUtxoAmount.MissingAmount {
				t.Fatalf("CreateTransaction() got NotEnoughUtxo '%v'", rawTx.NotEnoughUtxo)
			}

			if tt.wantNotEnoughUtxoAmount == nil {
				if rawTx == nil {
					t.Fatalf("CreateTransaction() got nil response")
				}

				if len(rawTx.Hex) == 0 {
					t.Fatalf("CreateTransaction() got empty raw hex")
				}
			}
		})
	}
}

func TestGenerateDerSignatures(t *testing.T) {
	hashStrToHash := func(str string) *chainhash.Hash {
		hash, err := chainhash.NewHashFromStr("864608ddfcb050c8a9a0c275687186ee2957e0853bee198aa464de798b7696db")
		if err != nil {
			panic(err)
		}
		return hash
	}
	hexStrToBytes := func(hexStr string) []byte {
		script, err := hex.DecodeString(hexStr)
		if err != nil {
			panic(err)
		}
		return script
	}

	tests := []struct {
		name    string
		msgTx   *wire.MsgTx
		utxos   []Utxo
		privKey string
		wantErr error
	}{
		{
			name: "generate DER signatures",
			msgTx: &wire.MsgTx{
				Version: 1,
				TxIn: []*wire.TxIn{
					wire.NewTxIn(
						wire.NewOutPoint(hashStrToHash("864608ddfcb050c8a9a0c275687186ee2957e0853bee198aa464de798b7696db"), 0),
						nil,
						nil,
					),
				},
				LockTime: 0x0,
			},
			utxos: []Utxo{
				{
					Script:     hexStrToBytes("001457f683080ee4491f1979950333e3240a0a9695d5"),
					Value:      1000000,
					Derivation: []uint32{84 + h, 0 + h, 0 + h, 0, 2},
				},
			},
			privKey: "tprv8g5UXufoRYtBEU5g2ueXDG32joLBLEzsGnoTUBZayNm8cAFGS56CcJmGwuSNEBLguQ3ja5betvc6kas1BXPpVzwuh8MWKr2ijzXJWuoJBqL",
		},
	}

	s := &Service{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			derSignatures, err := s.GenerateDerSignatures(tt.msgTx, tt.utxos, tt.privKey)
			if err != nil && tt.wantErr == nil {
				t.Fatalf("GenerateDerSignatures() got error '%v'", err)
			}

			if tt.wantErr == nil {
				if derSignatures == nil {
					t.Fatal("GenerateDerSignatures() got nil response")
				}

				countDerSignatures := len(derSignatures)
				countUtxos := len(tt.utxos)

				if countDerSignatures != countUtxos {
					t.Fatalf(
						"GenerateDerSignatures has generated %d signatures instead of %d",
						countDerSignatures,
						countUtxos,
					)
				}
			}
		})
	}
}

func TestSignTransaction(t *testing.T) {
	// Helper to derive extended key and return the btcec public key.
	// Use this in unit-tests to get input public key.
	getPublicKey := func(extendedKey string, derivation []uint32, chainParams chaincfg.ChainParams) *btcec.PublicKey {
		s := &Service{}

		pubKeyMat, err := s.DeriveExtendedKey(extendedKey, derivation)
		if err != nil {
			panic(err)
		}

		addressPubKey, err := btcutil.NewAddressPubKey(pubKeyMat.PublicKey, chainParams)
		if err != nil {
			panic(err)
		}

		return addressPubKey.PubKey()
	}

	tests := []struct {
		name               string
		msgTx              *wire.MsgTx
		chainParams        chaincfg.ChainParams
		utxos              []Utxo
		privKey            string
		signaturesMetadata []SignatureMetadata
		wantErr            error
	}{
		{
			name: "sign transaction",
			msgTx: &wire.MsgTx{
				Version: 1,
				TxIn: []*wire.TxIn{
					wire.NewTxIn(
						wire.NewOutPoint(btcutil.NewTx(wire.NewMsgTx(1)).Hash(), 0),
						nil,
						nil,
					),
				},
				LockTime: 0x0,
			},
			chainParams: chaincfg.BitcoinMainNetParams,
			utxos: []Utxo{
				{
					Script:     nil,
					Value:      100000,
					Derivation: []uint32{44 + h, 0 + h, 0 + h, 0, 0},
				},
			},
			privKey: "xprv9yv8fLFeRhD7NcKbjGS4GesBvy2PjvoRcwEKKaz7zJvM2cQ1eiCwhcHGQNEBwsXthHbPtZNQg5SBBEWS1QH941SKitBdaUT7VDTxzdS8vu7",
			signaturesMetadata: []SignatureMetadata{
				{
					DerSig: nil,
					PubKey: getPublicKey(
						"xpub6Cc939fyHvfB9pPLWd3bSyyQFvgKbwhidca49jGCM5Hz5ypEPGf9JVXB4NBuUfPgoHnMjN6oNgdC9KRqM11RZtL8QLW6rFKziNwHDYhZ6Kx",
						[]uint32{1, 1},
						chaincfg.BitcoinMainNetParams,
					),
					AddrEncoding: Legacy,
				},
			},
		},
	}

	s := &Service{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			derSignatures, err := s.GenerateDerSignatures(tt.msgTx, tt.utxos, tt.privKey)
			if err != nil && tt.wantErr == nil {
				t.Fatalf("GenerateDerSignatures() got error '%v'", err)
			}

			// Put DER signatures
			for idx := range tt.signaturesMetadata {
				tt.signaturesMetadata[idx].DerSig = derSignatures[idx]
			}

			signedRawTx, err := s.SignTransaction(tt.msgTx, tt.chainParams, tt.signaturesMetadata)
			if err != nil && tt.wantErr == nil {
				t.Fatalf("SignTransaction() got error '%v'", err)
			}

			if tt.wantErr == nil {
				if signedRawTx == nil {
					t.Fatal("SignTransaction() got nil response")
				}

				if len(signedRawTx.Hex) == 0 {
					t.Fatal("SignTransaction() got empty raw hex")
				}
			}
		})
	}
}
