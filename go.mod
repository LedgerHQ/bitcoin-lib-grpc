module github.com/ledgerhq/bitcoin-lib-grpc

go 1.15

require (
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v1.0.3-0.20200713135911-4649e4b73b34
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/ledgerhq/bitcoin-lib-grpc/pb v0.1.0
	github.com/magefile/mage v1.10.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/viper v1.3.2
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899 // indirect
	golang.org/x/net v0.0.0-20201002202402-0a1ea396d57c // indirect
	golang.org/x/sys v0.0.0-20201005172224-997123666555 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20201006033701-bcad7cf615f2 // indirect
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0 // indirect
)

replace github.com/ledgerhq/bitcoin-lib-grpc/pb => ./pb
