module github.com/ledgerhq/bitcoin-lib-grpc

go 1.16

require (
	github.com/btcsuite/btcd v0.23.2
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.1
	github.com/btcsuite/btcutil v1.0.3-0.20201104004401-a21f014935da
	github.com/btcsuite/btcwallet/wallet/txauthor v1.0.0
	github.com/btcsuite/btcwallet/wallet/txrules v1.0.0
	github.com/btcsuite/btcwallet/wallet/txsizes v1.0.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/magefile/mage v1.11.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/viper v1.3.2
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20201005172224-997123666555 // indirect
	google.golang.org/genproto v0.0.0-20201006033701-bcad7cf615f2 // indirect
	google.golang.org/grpc v1.33.2
)

replace github.com/ledgerhq/bitcoin-lib-grpc/pb => ./pb
