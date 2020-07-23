package bitcoin

import (
	"context"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	pb "github.com/ledgerhq/lama-bitcoin-svc/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct{}

func (s *Service) ValidateAddress(
	ctx context.Context, request *pb.ValidateAddressRequest,
) (*pb.ValidateAddressResponse, error) {
	params, err := getChainParams(request.ChainParams)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	addr, err := btcutil.DecodeAddress(request.Address, params)
	if err != nil {
		return &pb.ValidateAddressResponse{
			Address:       request.Address,
			IsValid:       false,
			InvalidReason: err.Error(),
		}, nil
	}

	return &pb.ValidateAddressResponse{
		Address: addr.EncodeAddress(), // Normalize the original address
		IsValid: true,
	}, nil
}

func (s *Service) EncodeAddress(
	ctx context.Context, request *pb.EncodeAddressRequest,
) (*pb.EncodeAddressResponse, error) {
	chainParams, err := getChainParams(request.ChainParams)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Load the serialized public key to a btcec.PublicKey type, in order to
	// ensure that the:
	//   * public point is on the secp256k1 elliptic curve.
	//   * public point coordinates belong to the finite field of secp256k1.
	//   * public key is well formed (valid magic, length, etc).
	//
	// Both compressed and uncompressed public keys are accepted.
	//
	// Using addresses encoded from incorrect public keys may lead to
	// irrevocable fund loss.
	publicKey, err := btcec.ParsePubKey(request.PublicKey, btcec.S256())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Calculate the RIPEMD160 of the SHA256 of a public key, aka HASH160.
	//
	// As per Bitcoin protocol, the serialized public key MUST be compressed
	// for P2SH-P2WPKH and P2WPKH addresses, whereas P2PKH addresses could be
	// either compressed or uncompressed. Regarding P2PKH addresses, the
	// convention at Ledger and in Bitcoin Wiki examples is to use compressed
	// public keys.
	publicKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())

	address, err := func() (btcutil.Address, error) {
		switch request.Encoding {
		case pb.AddressEncoding_ADDRESS_ENCODING_P2PKH:
			// Ref: https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
			return btcutil.NewAddressPubKeyHash(publicKeyHash, chainParams)
		case pb.AddressEncoding_ADDRESS_ENCODING_P2SH_P2WPKH:
			// Ref: https://bitcoincore.org/en/segwit_wallet_dev/#creation-of-p2sh-p2wpkh-address

			// Create a P2WPKH native-segwit address
			p2wpkhAddress, err := btcutil.NewAddressWitnessPubKeyHash(publicKeyHash, chainParams)
			if err != nil {
				return nil, err
			}

			// Create a P2SH redeemScript that pays to the P2WPKH address.
			//
			// The redeemScript is 22 bytes long, and starts with a OP_0,
			// followed by a canonical push of the keyhash. The keyhash
			// is HASH160 of the 33-byte compressed public key.
			//
			// scriptSig: OP_0 <hash160(compressed public key)>
			redeemScript, err := txscript.PayToAddrScript(p2wpkhAddress)
			if err != nil {
				return nil, err
			}

			return btcutil.NewAddressScriptHash(redeemScript, chainParams)
		case pb.AddressEncoding_ADDRESS_ENCODING_P2WPKH:
			// Ref: https://bitcoincore.org/en/segwit_wallet_dev/#native-pay-to-witness-public-key-hash-p2wpkh
			return btcutil.NewAddressWitnessPubKeyHash(publicKeyHash, chainParams)
		default:
			return nil, btcutil.ErrUnknownAddressType
		}
	}()

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.EncodeAddressResponse{Address: address.EncodeAddress()}, nil
}

func getChainParams(params *pb.ChainParams) (*chaincfg.Params, error) {
	switch network := params.GetBitcoinNetwork(); network {
	case pb.BitcoinNetwork_BITCOIN_NETWORK_MAINNET:
		return &chaincfg.MainNetParams, nil
	case pb.BitcoinNetwork_BITCOIN_NETWORK_TESTNET3:
		return &chaincfg.TestNet3Params, nil
	case pb.BitcoinNetwork_BITCOIN_NETWORK_REGTEST:
		return &chaincfg.RegressionNetParams, nil
	default:
		return nil, ErrUnknownNetwork(network.String())
	}
}
