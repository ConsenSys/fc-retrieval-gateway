module github.com/ConsenSys/fc-retrieval-gateway

go 1.14

require (
	github.com/ConsenSys/fc-retrieval-common v0.0.0-20210331012500-bf32c4270d66
	github.com/ConsenSys/fc-retrieval-register v0.0.0-20210402135918-b9f0b9869362
	github.com/ant0ine/go-json-rest v3.3.2+incompatible
	github.com/joho/godotenv v1.3.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
)

replace github.com/ConsenSys/fc-retrieval-common => ../fc-retrieval-common
