package config

import (
	"flag"
	"log"
	"strconv"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// NewConfig creates a new configuration
func NewConfig() *viper.Viper {
	conf := viper.New()
	conf.AutomaticEnv()
	defineFlags(conf)
	bindFlags(conf)
	setValues(conf)
	return conf
}

func parseUint8(value string) uint8 {
	result, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		log.Panic("unable to parse uint8")
	}
	return uint8(result)
}

// Map sets the config for the Gateway. NB: Gateways start without a private key. Private keys are provided by a gateway admin client.
func Map(conf *viper.Viper) settings.AppSettings {
	registerRefreshDuration, err := time.ParseDuration(conf.GetString("REGISTER_REFRESH_DURATION"))
	if err != nil {
		registerRefreshDuration = settings.DefaultRegisterRefreshDuration
	}
	tcpInactivityTimeout, err := time.ParseDuration(conf.GetString("TCP_INACTIVITY_TIMEOUT"))
	if err != nil {
		tcpInactivityTimeout = settings.DefaultTCPInactivityTimeout
	}
	tcpLongInactivityTimeout, err := time.ParseDuration(conf.GetString("TCP_LONG_INACTIVITY_TIMEOUT"))
	if err != nil {
		tcpLongInactivityTimeout = settings.DefaultLongTCPInactivityTimeout
	}
	return settings.AppSettings{
		BindRestAPI:     conf.GetString("BIND_REST_API"),
		BindProviderAPI: conf.GetString("BIND_PROVIDER_API"),
		BindGatewayAPI:  conf.GetString("BIND_GATEWAY_API"),
		BindAdminAPI:    conf.GetString("BIND_ADMIN_API"),
		LogLevel:        conf.GetString("LOG_LEVEL"),
		LogTarget:       conf.GetString("LOG_TARGET"),
		LogDir:          conf.GetString("LOG_DIR"),
		LogFile:         conf.GetString("LOG_FILE"),
		LogMaxBackups:   conf.GetInt("LOG_MAX_BACKUPS"),
		LogMaxAge:       conf.GetInt("LOG_MAX_AGE"),
		LogMaxSize:      conf.GetInt("LOG_MAX_SIZE"),
		LogCompress:     conf.GetBool("LOG_COMPRESS"),
		GatewayID:       conf.GetString("GATEWAY_ID"),

		RegisterAPIURL:          conf.GetString("REGISTER_API_URL"),
		RegisterRefreshDuration: registerRefreshDuration,

		GatewayAddress:        conf.GetString("GATEWAY_ADDRESS"),
		NetworkInfoGateway:    conf.GetString("IP") + ":" + conf.GetString("BIND_GATEWAY_API"),
		GatewayRegionCode:     conf.GetString("GATEWAY_REGION_CODE"),
		GatewayRootSigningKey: conf.GetString("GATEWAY_ROOT_SIGNING_KEY"),
		GatewaySigningKey:     conf.GetString("GATEWAY_SIGNING_KEY"),

		NetworkInfoClient:   conf.GetString("IP") + ":" + conf.GetString("BIND_REST_API"),
		NetworkInfoProvider: conf.GetString("IP") + ":" + conf.GetString("BIND_PROVIDER_API"),
		NetworkInfoAdmin:    conf.GetString("IP") + ":" + conf.GetString("BIND_ADMIN_API"),

		TCPInactivityTimeout:     tcpInactivityTimeout,
		TCPLongInactivityTimeout: tcpLongInactivityTimeout,
	}
}

func defineFlags(conf *viper.Viper) {
	flag.String("host", "0.0.0.0", "help message for host")
	flag.String("ip", "127.0.0.1", "help message for ip")
}

func bindFlags(conf *viper.Viper) {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	conf.BindPFlags(pflag.CommandLine)
}

func setValues(conf *viper.Viper) {
	conf.Set("IP", conf.GetString("ip"))
}
