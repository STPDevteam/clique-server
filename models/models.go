package models

type JsonRPCModel struct {
	Id      uint64   `json:"id"`
	Jsonrpc string   `json:"jsonrpc"`
	Result  []Result `json:"result"`
}

type Result struct {
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	Data             string   `json:"data"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
	Topics           []string `json:"topics"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}

type JsonRPCTimesTampModel struct {
	Id      uint64           `json:"id"`
	Jsonrpc string           `json:"jsonrpc"`
	Result  GetBlockByNumber `json:"result"`
}

type GetBlockByNumber struct {
	Timestamp string `json:"timestamp"`
	GasUsed   string `json:"gasUsed"`
}

type JsonRPCBalanceModel struct {
	Id      uint64  `json:"id"`
	Jsonrpc string  `json:"jsonrpc"`
	Result  Balance `json:"result"`
}

type Balance struct {
	Value string `json:"Value"`
}

type JsonRPCInfoModel struct {
	Id      uint64 `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
}

type BlockNumber struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"msg"`
}

type ResResult struct {
	Success bool `json:"success"`
}

type ResUploadImgPath struct {
	Path string `json:"path"`
}

type ResQueryAccount struct {
	Account      string        `json:"account"`
	AccountLogo  string        `json:"accountLogo"`
	Nickname     string        `json:"nickname"`
	Introduction string        `json:"introduction"`
	Twitter      string        `json:"twitter"`
	Github       string        `json:"github"`
	MyTokens     []ResMyTokens `json:"myTokens"`
	Daos         []ResDaos     `json:"daos"`
}

type ResMyTokens struct {
	TokenAddress string `json:"tokenAddress"`
	ChainId      int    `json:"chainId"`
	Balance      string `json:"balance"`
}

type ResDaos struct {
	DaoAddress   string `json:"daoAddress"`
	ChainId      int    `json:"chainId"`
	AccountLevel string `json:"accountLevel"`
}

type UpdateAccountWithSignParam struct {
	Sign  SignData           `json:"sign"`
	Param UpdateAccountParam `json:"param"`
}

type UpdateAccountParam struct {
	Account      string `json:"account"`
	AccountLogo  string `json:"accountLogo"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Twitter      string `json:"twitter"`
	Github       string `json:"github"`
}

type ResDaoListPage struct {
	List  []ResDaoList `json:"list"`
	Total uint64       `json:"total"`
}

type ResDaoList struct {
	DaoLogo    string `json:"daoLogo"`
	DaoName    string `json:"daoName"`
	DaoAddress string `json:"daoAddress"`
	ChainId    int    `json:"chainId"`
	// proposals total
	Proposals uint64 `json:"proposals"`
	// members total
	Members uint64 `json:"members"`
	// false:not joined Dao,true:joined Dao,default:false
	JoinSwitch bool `json:"joinSwitch"`
	// Account Level:noRole,member,admin,superAdmin
	AccountLevel string `json:"accountLevel"`
}

type SignCreateDataParam struct {
	// Sign Type:"0":CreateProposal,"1":Vote
	SignType   string `json:"signType"`
	Account    string `json:"account"`
	DaoAddress string `json:"daoAddress"`
}

type ResSignCreateData struct {
	TokenChainId int64  `json:"tokenChainId"`
	TokenAddress string `json:"tokenAddress"`
	Signature    string `json:"signature"`
	Balance      string `json:"balance"`
}

type JoinDaoWithSignParam struct {
	Sign   SignData     `json:"sign"`
	Params JoinDaoParam `json:"params"`
}

type JoinDaoParam struct {
	// 0:quit Dao,1:join Dao
	JoinSwitch int    `json:"joinSwitch"`
	DaoAddress string `json:"daoAddress"`
	ChainId    int    `json:"chainId"`
	Account    string `json:"account"`
}

type ResLeftDaoCreator struct {
	Account      string `json:"account"`
	DaoName      string `json:"daoName"`
	DaoAddress   string `json:"daoAddress"`
	TokenAddress string `json:"tokenAddress"`
	ChainId      int    `json:"chainId"`
}

type ResProposalsListPage struct {
	List  []ResProposalsList `json:"list"`
	Total uint64             `json:"total"`
}

type ResProposalsList struct {
	DaoAddress string `json:"daoAddress"`
	ProposalId string `json:"proposalId"`
	Proposer   string `json:"proposer"`
	StartTime  int64  `json:"startTime"`
	EndTime    int64  `json:"endTime"`
	Status     string `json:"status"`
}

type ResVotesListPage struct {
	List  []ResVotesList `json:"list"`
	Total uint64         `json:"total"`
}

type ResVotesList struct {
	OptionIndex string `json:"optionIndex"`
	Voter       string `json:"voter"`
	Amount      string `json:"amount"`
}

type ResTokenListPage struct {
	List  []ResTokenList `json:"list"`
	Total uint64         `json:"total"`
}

type ResTokenList struct {
	TokenAddress    string   `json:"tokenAddress"`
	ContractAddress string   `json:"contractAddress"`
	ChainId         int      `json:"chainId"`
	DaoName         []string `json:"daoName"`
	TotalSupply     string   `json:"totalSupply"`
}

type SignData struct {
	Account   string `json:"account" validate:"eth_addr"`              // personal_sign签名地址,0x开头
	Signature string `json:"signature" validate:"len=130,hexadecimal"` // personal_sign签名结果,不带0x开头
}
