package models

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type JsonRPCModel struct {
	Id      uint64      `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
}

type JsonRPCScanBlockModel struct {
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

type JsonRPCTransactionByHashModel struct {
	Id      uint64 `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  From   `json:"result"`
}

type From struct {
	From string `json:"from"`
}

type JsonRPCBalanceModel struct {
	Id      uint64  `json:"id"`
	Jsonrpc string  `json:"jsonrpc"`
	Result  Balance `json:"result"`
}

type Balance struct {
	Value string `json:"Value"`
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
	Account      string `json:"account"`
	AccountLogo  string `json:"accountLogo"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Twitter      string `json:"twitter"`
	Github       string `json:"github"`
	//MyTokens     []ResMyTokens `json:"myTokens"`
	AdminDao  []ResDao `json:"adminDao"`
	MemberDao []ResDao `json:"memberDao"`
	//Activity  []ResActivity `json:"activity"`
}

//type ResMyTokens struct {
//	TokenAddress string `json:"tokenAddress"`
//	ChainId      int    `json:"chainId"`
//	Balance      string `json:"balance"`
//}

type ResDao struct {
	DaoAddress   string `json:"daoAddress"`
	ChainId      int    `json:"chainId"`
	AccountLevel string `json:"accountLevel"`
	DaoName      string `json:"daoName"`
	DaoLogo      string `json:"daoLogo"`
}

type ResActivity struct {
	EventType   string `json:"eventType"`
	ChainId     int    `json:"chainId"`
	DaoAddress  string `json:"daoAddress"`
	ProposalId  int    `json:"proposalId"`
	OptionIndex int    `json:"optionIndex"`
	Amount      string `json:"amount"`
}

type UpdateAccountWithSignParam struct {
	Sign  SignData           `json:"sign"`
	Param UpdateAccountParam `json:"param"`
}

type UpdateAccountParam struct {
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
	// approve:true,not approve:false
	Approve bool `json:"approve"`
	// total proposals
	TotalProposals uint64 `json:"totalProposals"`
	// activity proposals
	ActiveProposals uint64 `json:"activeProposals"`
	// soon proposals
	SoonProposals uint64 `json:"soonProposals"`
	// closed proposals
	ClosedProposals uint64 `json:"closedProposals"`
	// members total
	Members uint64 `json:"members"`
	// 0:not joined Dao, 1:joined Dao,default:0
	JoinSwitch int `json:"joinSwitch"`
}

type SignCreateDataParam struct {
	// SignType:"0":CreateProposal,"1":Vote
	SignType string `json:"signType"`
	// if SignType:"1":Vote,need proposalId
	ProposalId int    `json:"proposalId"`
	Account    string `json:"account"`
	ChainId    int    `json:"chainId"`
	DaoAddress string `json:"daoAddress"`
}

type ResSignCreateData struct {
	Account      string `json:"account"`
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
	Account    string `json:"account"`
	DaoAddress string `json:"daoAddress"`
	ChainId    int    `json:"chainId"`
}

type ResProposalsListPage struct {
	List  []ResProposalsList `json:"list"`
	Total uint64             `json:"total"`
}

type ResProposalsList struct {
	ChainId    int    `json:"chainId"`
	DaoAddress string `json:"daoAddress"`
	ProposalId int    `json:"proposalId"`
	Proposer   string `json:"proposer"`
	StartTime  int64  `json:"startTime"`
	EndTime    int64  `json:"endTime"`
}

type ResVotesListPage struct {
	List  []ResVotesList `json:"list"`
	Total int            `json:"total"`
}

type ResVotesList struct {
	ProposalId  int    `json:"proposalId"`
	Voter       string `json:"voter"`
	OptionIndex int    `json:"optionIndex"`
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
	Account   string `json:"account" validate:"eth_addr"`              // personal_sign address,0x
	Signature string `json:"signature" validate:"len=130,hexadecimal"` // personal_sign sign result,no 0x
}

type AirdropAdminSignData struct {
	ChainId    int                 `json:"chainId"`                                  // airdrop1:need ChainId
	DaoAddress string              `json:"daoAddress"`                               // airdrop1:need DaoAddress
	AirdropId  int64               `json:"airdropId"`                                // airdrop2/airdropDownload:need AirdropId
	Account    string              `json:"account" validate:"eth_addr"`              // personal_sign address,0x
	Message    string              `json:"message"`                                  //{"expired":1244,"root": "","type":"airdrop1/airdrop2/airdropDownload"}
	Signature  string              `json:"signature" validate:"len=130,hexadecimal"` // personal_sign sign result,no 0x
	Array      AirdropAddressArray `json:"array"`                                    // airdrop2:need Array
}

type AirdropAddressArray struct {
	Address []string `json:"address"`
	Amount  []string `json:"amount"`
}

type AdminMessage struct {
	Expired int64  `json:"expired"`
	Root    string `json:"root"`
	Type    string `json:"type"`
}

type AccountParam struct {
	Account string `json:"account"`
}

type ResDaoInfo struct {
	Members uint64 `json:"members"`
	// 0:not joined Dao, 1:joined Dao,default:0
	JoinSwitch int `json:"joinSwitch"`
	// approve:true,not approve:false
	Approve bool `json:"approve"`
}

type ResAdminsList struct {
	Account string `json:"account"`
	//AccountLevel string `json:"accountLevel"`
}

type ErrorInfoParam struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Func    string `json:"func"`
	Params  string `json:"params"`
}

type ProposalInfoParam struct {
	Content string `json:"content"`
}

type ResProposalUuid struct {
	Uuid string `json:"uuid"`
}

type ResProposalContent struct {
	Content string `json:"content"`
}

type TokensInfo struct {
	Id        string                 `json:"id"`
	Symbol    string                 `json:"symbol"`
	Name      string                 `json:"name"`
	Platforms map[string]interface{} `json:"platforms"`
}

type TokenImg struct {
	Image struct {
		Thumb string `json:"thumb"`
		Small string `json:"small"`
		Large string `json:"large"`
	} `json:"image"`
}

type ResTokenImg struct {
	TokenChainId int    `json:"tokenChainId"`
	TokenAddress string `json:"tokenAddress"`
	Thumb        string `json:"thumb"`
	Small        string `json:"small"`
	Large        string `json:"large"`
	OwnImg       string `json:"ownImg"`
}

type ResSnapshot struct {
	ChainId    int    `json:"chainId"`
	DaoAddress string `json:"daoAddress"`
	ProposalId int    `json:"proposalId"`
	Snapshot   int64  `json:"snapshot"`
}

type ResAirdropId struct {
	AirdropId int    `json:"airdropId"`
	Signature string `json:"signature"`
}

type AddressData struct {
	Id      uint64
	Amount  *big.Int
	Address common.Address
}

type ClaimInfo struct {
	Index  uint64   `json:"index"`
	Amount string   `json:"amount"`
	Proof  []string `json:"proof"`
}

type ResProof struct {
	AirdropTotalAmount string   `json:"airdropTotalAmount"`
	AirdropNumber      int      `json:"airdropNumber"`
	Title              string   `json:"title"`
	Index              uint64   `json:"index"`
	Amount             string   `json:"amount"`
	Proof              []string `json:"proof"`
}

type ResActivityPage struct {
	List  []ResActivityList `json:"list"`
	Total uint64            `json:"total"`
}

type ResActivityList struct {
	Title             string  `json:"title"`
	Types             string  `json:"types"`
	ChainId           int     `json:"chainId"`
	DaoAddress        string  `json:"daoAddress"`
	Creator           string  `json:"creator"`
	ActivityId        int     `json:"activityId"`
	TokenChainId      int     `json:"tokenChainId"`
	TokenAddress      string  `json:"tokenAddress"`
	StakingAmount     string  `json:"stakingAmount"`
	StartTime         int64   `json:"startTime"`
	EndTime           int64   `json:"endTime"`
	AirdropStartTime  int64   `json:"airdropStartTime"`
	AirdropEndTime    int64   `json:"airdropEndTime"`
	PublishTime       int64   `json:"publishTime"`
	Price             string  `json:"price"`
	AirdropNumber     int     `json:"airdropNumber"`
	ClaimedPercentage float64 `json:"claimedPercentage"`
}

type ResNotificationPage struct {
	List        []ResNotification `json:"list"`
	Total       uint64            `json:"total"`
	UnreadTotal int               `json:"unreadTotal"`
}

type ResNotification struct {
	Account          string           `json:"account"`
	AlreadyRead      bool             `json:"alreadyRead"` // true:have read
	NotificationId   int              `json:"notificationId"`
	NotificationTime int64            `json:"notificationTime"`
	Types            string           `json:"types"` //Airdrop PublicSale NewProposal ReserveToken
	Info             NotificationInfo `json:"info"`
}

type NotificationInfo struct {
	ChainId      int    `json:"chainId"`
	DaoAddress   string `json:"daoAddress"`
	DaoLogo      string `json:"daoLogo"`
	DaoName      string `json:"daoName"`
	ProposalId   int    `json:"proposalId"`   // NewProposal
	ProposalName string `json:"proposalName"` // NewProposal
	ActivityId   int    `json:"activityId"`   // Airdrop
	ActivityName string `json:"activityName"` // Airdrop
	//StartTime    int64  `json:"startTime"`
}

type NotificationReadParam struct {
	NotificationId int    `json:"notificationId"`
	Account        string `json:"account"`
	// read all:true
	ReadAll bool `json:"readAll"`
}

type ResNotificationUnreadTotal struct {
	UnreadTotal int `json:"unreadTotal"`
}

type SignDaoHandleParam struct {
	Account string `json:"account"`
	ChainId int    `json:"chainId"`
	Handle  string `json:"handle"`
}

type ResSignDaoHandleData struct {
	Signature    string `json:"signature"`
	Account      string `json:"account"`
	ChainId      int    `json:"chainId"`
	LockBlockNum int    `json:"lockBlockNum"`
}

type CreateAirdropParam struct {
	Title              string               `json:"title"`
	Description        string               `json:"description"`
	CollectInformation []CollectInfo        `json:"collectInformation"`
	TokenChainId       int                  `json:"tokenChainId"`
	TokenAddress       string               `json:"tokenAddress"`
	MaxAirdropAmount   string               `json:"maxAirdropAmount"`
	StartTime          uint64               `json:"startTime"`
	EndTime            uint64               `json:"endTime"`
	AirdropStartTime   uint64               `json:"airdropStartTime"`
	AirdropEndTime     uint64               `json:"airdropEndTime"`
	Sign               AirdropAdminSignData `json:"sign"`
}

type CollectInfo struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

type ResAirdropInfo struct {
	Creator          string        `json:"creator"`
	ChainId          int           `json:"chainId"`
	DaoAddress       string        `json:"daoAddress"`
	Title            string        `json:"title"`
	Description      string        `json:"description"`
	TokenChainId     int           `json:"tokenChainId"`
	TokenAddress     string        `json:"tokenAddress"`
	StartTime        int64         `json:"startTime"`
	EndTime          int64         `json:"endTime"`
	AirdropStartTime int64         `json:"airdropStartTime"`
	AirdropEndTime   int64         `json:"airdropEndTime"`
	AddressNum       int           `json:"addressNum"`
	CollectCount     int           `json:"collectCount"`
	Collect          []CollectInfo `json:"collect"`
}

type UserInformationParam struct {
	AirdropId  int64  `json:"airdropId"`
	Account    string `json:"account"`
	UserSubmit string `json:"userSubmit"`
}

type ResTreeRoot struct {
	Root string `json:"root"`
}

type ResUploadAddressList struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}
