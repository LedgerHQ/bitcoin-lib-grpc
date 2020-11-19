package bitcoin

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/pkg/errors"
)

type Input struct {
	OutputHash  string
	OutputIndex uint32
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

// RawTx represents the serialized transaction encoded using legacy encoding
// of BIP144 Segregated Witness encoding.
//
// Hash and WitnessHash are the same if transaction has no witness data.
type RawTx struct {
	Hex         string
	Hash        string
	WitnessHash string
}

type DerSignature = []byte

type Utxo struct {
	Script     []byte
	Value      int64
	Derivation []uint32
}

type SignatureMetadata struct {
	DerSig       DerSignature
	PubKey       *btcec.PublicKey
	AddrEncoding AddressEncoding
}

func (s *Service) CreateTransaction(tx *Tx, chainParams ChainParams) (*RawTx, error) {
	// Create a new btcd transaction
	msgTx := wire.NewMsgTx(wire.TxVersion)

	// For each input to spend, add a TxIn
	for _, input := range tx.Inputs {
		// hash string to hash byteArray
		outputHash, err := chainhash.NewHashFromStr(input.OutputHash)
		if err != nil {
			return nil, errors.Wrapf(err,
				"failed to get hash string from output hash %s",
				input.OutputHash,
			)
		}

		// Previous outpoint = hash + index
		prevOut := wire.NewOutPoint(outputHash, uint32(input.OutputIndex))

		// Create new Input from previous outpoint
		txIn := wire.NewTxIn(prevOut, nil, nil)

		// Add TxIn to MsgTx
		msgTx.AddTxIn(txIn)
	}

	// For each output to send, add a TxOut
	for _, output := range tx.Outputs {
		// Decode address from string
		address, err := btcutil.DecodeAddress(output.Address, chainParams)

		if err != nil {
			return nil, errors.Wrapf(err,
				"failed to decode address from output address %s",
				output.Address,
			)
		}

		// Create a 'pay to' script that pays to the address depending on the address type.
		outputScript, err := txscript.PayToAddrScript(address)
		if err != nil {
			return nil, errors.Wrapf(err,
				"failed to build 'pay to' script from address %v",
				address,
			)
		}

		// Create Output from value and script
		txOut := wire.NewTxOut(output.Value, outputScript)

		// Add TxOut to MsgTx
		msgTx.AddTxOut(txOut)
	}

	// Add LockTime
	msgTx.LockTime = tx.LockTime

	// Encode MsgTx to RawTx
	rawTx, err := encodeMsgTx(msgTx)
	if err != nil {
		return nil, err
	}

	return rawTx, nil
}

func (s *Service) GenerateDerSignatures(msgTx *wire.MsgTx, utxos []Utxo, privKey string) ([]DerSignature, error) {
	// Validation
	if len(msgTx.TxIn) != len(utxos) {
		return nil, errors.New("inputs length != utxos length")
	}

	// Get extended key from private key
	extendedKey, err := hdkeychain.NewKeyFromString(privKey)
	if err != nil {
		return nil, errors.Wrapf(err,
			"failed to get extended key from private key %s",
			privKey,
		)
	}

	// Build sig hash type and sig hashes
	sigHashType := txscript.SigHashAll
	sigHashes := txscript.NewTxSigHashes(msgTx)

	derSignatures := make([]DerSignature, len(msgTx.TxIn))

	// Generate a valid der signature for each input
	for idx, input := range msgTx.TxIn {
		// Get the utxo assuming inputs and utxos are in the same order
		utxo := utxos[idx]

		script := utxo.Script

		amount := utxo.Value

		derivation := utxo.Derivation

		// Derive the extended key for given derivation path
		for _, childIndex := range derivation {
			extendedKey, err = extendedKey.Derive(childIndex)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to derive extendedKey %s at index %d",
					extendedKey, childIndex)
			}
		}

		// Get the private key for given derivation path
		ecPrivKey, err := extendedKey.ECPrivKey()
		if err != nil {
			return nil, err
		}

		derSig, err := txscript.RawTxInWitnessSignature(msgTx, sigHashes, idx, amount, script, sigHashType, ecPrivKey)

		if err != nil {
			return nil, errors.Wrapf(err,
				"failed to generate der signature for input %v",
				input,
			)
		}

		derSignatures[idx] = derSig
	}

	return derSignatures, nil
}

func (s *Service) SignTransaction(msgTx *wire.MsgTx, chainParams ChainParams, signatures []SignatureMetadata) (*RawTx, error) {
	// Validation
	if len(msgTx.TxIn) != len(signatures) {
		return nil, errors.New("inputs length != signatures length")
	}

	for inputIdx, input := range msgTx.TxIn {

		// Get the signature struct assuming inputs and signatures are in the same order
		signature := signatures[inputIdx]

		// Get the DER signature
		derSig := signature.DerSig

		// Get input address public key
		pubKey := signature.PubKey

		// Get address type for the input
		inputAddrEncoding := signature.AddrEncoding

		// Serialize input public key data
		pubKeyData := pubKey.SerializeCompressed()

		var sigScript []byte

		// If we're spending p2wkh output nested within a p2sh output, then
		// we'll need to attach a sigScript in addition to witness data.
		// Otherwise sigScript is an empty byte array
		if inputAddrEncoding == WrappedSegwit {
			pubKeyHash := btcutil.Hash160(pubKeyData)

			// Next, we'll generate a valid sigScript that will allow us to
			// spend the p2sh output. The sigScript will contain only a
			// single push of the p2wkh witness program corresponding to
			// the matching public key of this address.
			p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(
				pubKeyHash, chainParams,
			)

			if err != nil {
				return nil, err
			}

			redeemScript, err := txscript.PayToAddrScript(p2wkhAddr)
			if err != nil {
				return nil, err
			}

			bldr := txscript.NewScriptBuilder()
			bldr.AddData(redeemScript)
			sigScript, err = bldr.Script()

			if err != nil {
				return nil, err
			}
		}

		// Put signature script and witness data to the input
		input.SignatureScript = sigScript
		input.Witness = wire.TxWitness{derSig, pubKeyData}
	}

	// Encode signed MsgTx to RawTx
	signedRawTx, err := encodeMsgTx(msgTx)
	if err != nil {
		return nil, err
	}

	return signedRawTx, nil
}

// Encode MsgTx to RawTx
func encodeMsgTx(msgTx *wire.MsgTx) (*RawTx, error) {
	var buf bytes.Buffer
	if err := msgTx.Serialize(&buf); err != nil {
		return nil, errors.Wrap(err, "failed to encode transaction in hex")
	}

	rawTx := &RawTx{
		Hex:         hex.EncodeToString(buf.Bytes()),
		Hash:        msgTx.TxHash().String(),
		WitnessHash: msgTx.WitnessHash().String(),
	}

	return rawTx, nil
}

// Deserialize MsgTx from RawTx
func (s *Service) DeserializeMsgTx(rawTx *RawTx) (*wire.MsgTx, error) {
	// Instantiate a MsgTx
	msgTx := wire.NewMsgTx(wire.TxVersion)

	// Get bytes from hex string
	hexBytes, err := hex.DecodeString(rawTx.Hex)
	if err != nil {
		return nil, errors.Wrapf(err,
			"failed to derialize raw tx %v",
			rawTx,
		)
	}
	reader := bytes.NewReader(hexBytes)

	// Derialize into MsgTx
	msgTx.Deserialize(reader)

	return msgTx, nil
}
