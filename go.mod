module github.com/ledgerhq/bitcoin-lib-grpc

go 1.16

require (
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v1.0.3-0.20201104004401-a21f014935da
	github.com/btcsuite/btcwallet/wallet/txauthor v1.0.0
	github.com/btcsuite/btcwallet/wallet/txrules v1.0.0
	github.com/btcsuite/btcwallet/wallet/txsizes v1.0.0
	github.com/magefile/mage v1.11.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/viper v1.3.2
	google.golang.org/grpc v1.53.0
)

replace github.com/ledgerhq/bitcoin-lib-grpc/pb => ./pb
