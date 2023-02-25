module github.com/ledgerhq/bitcoin-lib-grpc

go 1.16

require (
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v1.0.3-0.20201104004401-a21f014935da
	github.com/btcsuite/btcwallet/wallet/txauthor v1.0.0
	github.com/btcsuite/btcwallet/wallet/txrules v1.0.0
	github.com/btcsuite/btcwallet/wallet/txsizes v1.0.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/magefile/mage v1.11.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/viper v1.3.2
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20201006033701-bcad7cf615f2 // indirect
	google.golang.org/grpc v1.33.2
)

replace github.com/ledgerhq/bitcoin-lib-grpc/pb => ./pb
