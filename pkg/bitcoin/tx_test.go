package bitcoin

import (
	"bytes"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	tests := []struct {
		name    string
		tx      *Tx
		net     ChainParams
		wantErr error
	}{
		{
			name: "mainnet P2WPKH",
			tx: &Tx{
				LockTime: 0,
				Inputs: []Input{
					{
						OutputHash:  "2f5dae23c2e18588c86cfc4e154f3b68bd8eb4265fe0b4b1341ad5aa40422f66",
						OutputIndex: 1,
						Script:      "160014513f387619014109e25764de4df3467b786ad125",
						Sequence:    16777215,
					},
				},
				Outputs: []Output{
					{
						Address: "1MZbRqZGpiSWGRLg8DUdVrDKHwNe1oesUZ",
						Value:   100000,
					},
				},
			},
			net: Mainnet,
		},
	}

	s := &Service{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var buf bytes.Buffer

			err := s.CreateTransaction(&buf, tt.tx, tt.net)
			if err != nil && tt.wantErr == nil {
				t.Fatalf("CreateTransaction() got error '%v'", err)
			}

			if buf.Len() == 0 {
				t.Fatalf("Created Transaction is empty")
			}
		})
	}
}
