package config

type AppConfig struct {
	Daemon                     bool     `json:"daemon"`
	BasePath                   string   `mapstructure:"base_path" toml:"base_path"`
	ServerPort                 int      `mapstructure:"server_port" toml:"server_port"`
	MinIdleAddressCount        uint64   `mapstructure:"min_idle_address_count" toml:"min_idle_address_count"`
	MaxIdleAddressCount        uint64   `mapstructure:"max_idle_address_count" toml:"max_idle_address_count"`
	DurationIdleAddressMonitor int64    `mapstructure:"duration_idle_address_monitor" toml:"duration_idle_address_monitor"`
	DurationTransactionMonitor int64    `mapstructure:"duration_transaction_monitor" toml:"duration_transaction_monitor"`
	BlockNumberPerReq          int      `mapstructure:"block_number_per_req" toml:"block_number_per_req"`
	DelayedBlockNumber         int      `mapstructure:"delayed_block_number" toml:"delayed_block_number"`
	MaxUpdateImgSize           int64    `mapstructure:"max_upload_img_size" toml:"max_upload_img_size"`
	SignMessagePriKey          string   `mapstructure:"sign_message_pri_key" toml:"sign_message_pri_key"`
	SwaggerUrl                 string   `mapstructure:"swagger_url" toml:"swagger_url"`
	CallbackUrl                string   `mapstructure:"callback_url" toml:"callback_url"`
	WithdrawCallbackUrl        string   `mapstructure:"withdraw_callback_url" toml:"withdraw_callback_url"`
	CallbackInterval           []int64  `mapstructure:"callback_interval" toml:"callback_interval"`
	MysqlConfName              string   `mapstructure:"mysql_conf_name" toml:"mysql_conf_name"`
	ScanInfoConfName           []string `mapstructure:"scan_info_conf_name" toml:"scan_info_conf_name"`
}

type MysqlConfig struct {
	User     string `mapstructure:"user" toml:"user"`
	Password string `mapstructure:"password" toml:"password"`
	Host     string `mapstructure:"host" toml:"host"`
	Port     int    `mapstructure:"port" toml:"port"`
	Name     string `mapstructure:"name" toml:"name"`
}

type ScanInfoConfig struct {
	ChainId               int64    `mapstructure:"chain_id" toml:"chain_id"`
	ScanUrl               string   `mapstructure:"scan_url" toml:"scan_url"`
	KeyLastBlock          string   `mapstructure:"key_last_block" toml:"key_last_block"`
	MinAmount             []string `mapstructure:"min_amount" toml:"min_amount"`
	SupportedCoin         []string `mapstructure:"supported_coin" toml:"supported_coin"`
	SupportedCoinDecimals []uint64 `mapstructure:"supported_coin_decimals" toml:"supported_coin_decimals"`
	SupportedCoinAddress  []string `mapstructure:"supported_coin_address" toml:"supported_coin_address"`
}
