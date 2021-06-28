# bitcoin-lib-grpc

bitcoin-lib-grpc is a modular service that exposes a gRPC interface to wrap
protocol-centric logic related to Bitcoin, and related forks.

It is based on btcd - a Bitcoin SDK and full node implementation written in Go.

### Supported chains

* Bitcoin
  * mainnet
  * testnet3
  * regtest
* Litecoin
  * mainnet

### Summary of gRPC methods

Full details available [here](https://github.com/LedgerHQ/bitcoin-lib-grpc/blob/master/pb/bitcoin/service.proto).

* [`ValidateAddress`](https://github.com/LedgerHQ/bitcoin-lib-grpc/blob/e515b9797f25565955207594448664e32a5e35b0/pb/bitcoin/service.proto#L13-L16)
  ➔ check whether an address is valid for the given chain parameters.

* [`DeriveExtendedKey`](https://github.com/LedgerHQ/bitcoin-lib-grpc/blob/e515b9797f25565955207594448664e32a5e35b0/pb/bitcoin/service.proto#L18-L21)
  ➔ derive a child extended key based on BIP0032 derivation rules.

* [`EncodeAddress`](https://github.com/LedgerHQ/bitcoin-lib-grpc/blob/e515b9797f25565955207594448664e32a5e35b0/pb/bitcoin/service.proto#L23-L32)
  ➔ serialize public key into a network-specific address.

* [`GetAccountExtendedKey`](https://github.com/LedgerHQ/bitcoin-lib-grpc/blob/e515b9797f25565955207594448664e32a5e35b0/pb/bitcoin/service.proto#L34-L36)
  ➔ serialize public key-material into xPub.

* [`CreateTransaction`](https://github.com/LedgerHQ/bitcoin-lib-grpc/blob/e515b9797f25565955207594448664e32a5e35b0/pb/bitcoin/service.proto#L38-L39)
  ➔ prepare a raw TX based on TX params.

* [`SignTransaction`](https://github.com/LedgerHQ/bitcoin-lib-grpc/blob/e515b9797f25565955207594448664e32a5e35b0/pb/bitcoin/service.proto#L48-L51)
  ➔ combine a raw TX with DER signatures to produce the signed TX.

* [`GetKeypair`](https://github.com/LedgerHQ/bitcoin-lib-grpc/blob/e515b9797f25565955207594448664e32a5e35b0/pb/bitcoin/service.proto#L41-L43)
  ➔ derive an extended public-private keypair from a seed on a given derivation.
  
  ⚠️ **Note:** for use in tests only.
* [`GenerateDerSignatures`](https://github.com/LedgerHQ/bitcoin-lib-grpc/blob/e515b9797f25565955207594448664e32a5e35b0/pb/bitcoin/service.proto#L45-L46)
  ➔ sign a raw TX with a private key to produce DER signatures.

  ⚠️ **Note:** for use in tests only.

### Development

1. Install [mage](https://magefile.org).
2. Build the project to produce the executable.
    ```
    $ mage build
    ```
3. Run the binary.
  ```
  $ ./lbs
  ```
