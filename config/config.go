package config

type AppConfig struct {
	BasePath                string   `mapstructure:"base_path" toml:"base_path"`
	ServerPort              int      `mapstructure:"server_port" toml:"server_port"`
	BlockNumberPerReq       int      `mapstructure:"block_number_per_req" toml:"block_number_per_req"`
	MaxUpdateImgSize        int64    `mapstructure:"max_upload_img_size" toml:"max_upload_img_size"`
	SignMessagePriKey       string   `mapstructure:"sign_message_pri_key" toml:"sign_message_pri_key"`
	SwaggerUrl              string   `mapstructure:"swagger_url" toml:"swagger_url"`
	MysqlConfName           string   `mapstructure:"mysql_conf_name" toml:"mysql_conf_name"`
	ScanInfoConfName        []string `mapstructure:"scan_info_conf_name" toml:"scan_info_conf_name"`
	ArchiveBalanceSign      []int    `mapstructure:"archive_balance_sign" toml:"archive_balance_sign"`
	TestnetBalanceSign      []int    `mapstructure:"testnet_balance_sign" toml:"testnet_balance_sign"`
	ApiV1ProposalUrl        string   `mapstructure:"api_v1_proposal_url" toml:"api_v1_proposal_url"`
	ApiV1ProposalContentUrl string   `mapstructure:"api_v1_proposal_content_url" toml:"api_v1_proposal_content_url"`
	PolygonQuickNodeRPC     string   `mapstructure:"polygon_quick_node_rpc" toml:"polygon_quick_node_rpc"`
	MainnetChainstackRPC    string   `mapstructure:"mainnet_chainstack_rpc" toml:"mainnet_chainstack_rpc"`
}

type MysqlConfig struct {
	User     string `mapstructure:"user" toml:"user"`
	Password string `mapstructure:"password" toml:"password"`
	Host     string `mapstructure:"host" toml:"host"`
	Port     int    `mapstructure:"port" toml:"port"`
	Name     string `mapstructure:"name" toml:"name"`
}

type ScanInfoConfig struct {
	ChainId             []int    `mapstructure:"chain_id" toml:"chain_id"`
	ScanUrl             []string `mapstructure:"scan_url" toml:"scan_url"`
	HandleLockBlock     []int    `mapstructure:"handle_lock_block" toml:"handle_lock_block"`
	DelayedBlockNumber  []int    `mapstructure:"delayed_block_number" toml:"delayed_block_number"`
	QueryBlockNumberUrl []string `mapstructure:"query_block_number_url" toml:"query_block_number_url"`
}
