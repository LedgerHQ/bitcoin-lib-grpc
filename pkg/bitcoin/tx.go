package bitcoin

import (
	"bytes"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type Input struct {
	OutputHash  string
	OutputIndex uint32
	Script      string
	Sequence    uint32
}

type Output struct {
	Address string
	Value   int64
}

type Tx struct {
	Inputs   []Input
	Outputs  []Output
	LockTime uint32
}

func (s *Service) CreateTransaction(buf *bytes.Buffer, tx *Tx, chainParams ChainParams) error {
	// Create a new btcd transaction
	msgTx := wire.NewMsgTx(wire.TxVersion)

	for _, input := range tx.Inputs {
		// hash string to hash byteArray
		outputHash, err := chainhash.NewHashFromStr(input.OutputHash)
		if err != nil {
			return err
		}

		// Outpoint = hash + index
		prevOut := wire.NewOutPoint(outputHash, uint32(input.OutputIndex))

		// Create new Input from outpoint and input script
		txIn := wire.NewTxIn(prevOut, []byte(input.Script), nil)

		// Add sequence
		txIn.Sequence = input.Sequence

		// Add input to btcd transaction
		msgTx.AddTxIn(txIn)
	}

	for _, output := range tx.Outputs {
		// Guess Address Type from address string
		address, err := btcutil.DecodeAddress(output.Address, chainParams)
		if err != nil {
			return err
		}

		// Create a public key script that pays to the address depending on the address type.
		script, err := txscript.PayToAddrScript(address)
		if err != nil {
			return err
		}

		// Create Output from value and script
		txOut := wire.NewTxOut(output.Value, script)

		// Add Output to btcd Transaction
		msgTx.AddTxOut(txOut)
	}

	// Add LockTime
	msgTx.LockTime = tx.LockTime

	// Encode transaction in Hexadecimal
	return msgTx.BtcEncode(buf, wire.ProtocolVersion, wire.WitnessEncoding)
}