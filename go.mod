module github.com/datachainlab/cosmos-sdk-interchain-dns

go 1.14

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200622203133-4716260a6e2d
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/grpc-ecosystem/grpc-gateway v1.14.7
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/otiai10/copy v1.2.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/iavl v0.14.0 // indirect
	github.com/tendermint/tendermint v0.33.5
	github.com/tendermint/tm-db v0.5.1
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.31.0
)

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

replace github.com/cosmos/cosmos-sdk => github.com/datachainlab/cosmos-sdk v0.34.4-0.20200623040429-5c9fe2c0d7e5
