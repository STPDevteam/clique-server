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

type JsonRPCTokenHoldersModel struct {
	Id      uint64  `json:"id"`
	Jsonrpc string  `json:"jsonrpc"`
	Result  Holders `json:"result"`
}

type Holders struct {
	Total uint64 `json:"total"`
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

type JsonRPCReserveTokenModel struct {
	Id      uint64 `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  Input  `json:"result"`
}

type Input struct {
	Input string `json:"input"`
}

type JsonRPCBalanceModel struct {
	Id      uint64  `json:"id"`
	Jsonrpc string  `json:"jsonrpc"`
	Result  Balance `json:"result"`
}

type Balance struct {
	Value string `json:"Value"`
}

type JsonRPCGetBlockNumber struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

type JsonRPCAccountNFT struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  struct {
		Data []struct {
			ContractAddress string  `json:"contractAddress"`
			TokenId         string  `json:"tokenId"`
			Amount          float64 `json:"amount"`
			Quantity        int     `json:"quantity"`
			UsdAmount       float64 `json:"usdAmount"`
			Currency        string  `json:"currency"`
			Standard        string  `json:"standard"`
			Metadata        struct {
				ImageURL        string `json:"imageURL"`
				GatewayImageURL string `json:"gatewayImageURL"`
				Name            string `json:"name"`
				CollectionName  string `json:"collectionName"`
			} `json:"metadata"`
		} `json:"data"`
		NextPageIndex int `json:"nextPageIndex"`
		Total         int `json:"total"`
	} `json:"result"`
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
	Followers    int    `json:"followers"`
	Following    int    `json:"following"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Twitter      string `json:"twitter"`
	Github       string `json:"github"`
	Discord      string `json:"discord"`
	Email        string `json:"email"`
	Country      string `json:"country"`
	Youtube      string `json:"youtube"`
	Opensea      string `json:"opensea"`
	//MyTokens     []ResMyTokens `json:"myTokens"`
	AdminDao             []ResDao `json:"adminDao"`
	MemberDao            []ResDao `json:"memberDao"`
	AllDaosICreateOrJoin bool     `json:"allDaosICreateOrJoin"`
	NewDao               bool     `json:"newDao"`
	AllDaoAirdrop        bool     `json:"allDaoAirdrop"`
	AllDaoProposal       bool     `json:"allDaoProposal"`
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
	Param UpdateAccountParam `json:"param"`
}

type UpdateAccountParam struct {
	AccountLogo  string `json:"accountLogo"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Twitter      string `json:"twitter"`
	Github       string `json:"github"`
	Discord      string `json:"discord"`
	Email        string `json:"email"`
	Country      string `json:"country"`
	Youtube      string `json:"youtube"`
	Opensea      string `json:"opensea"`
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
	TotalProposals int64 `json:"totalProposals"`
	// activity proposals
	ActiveProposals uint64 `json:"activeProposals"`
	// soon proposals
	SoonProposals uint64 `json:"soonProposals"`
	// closed proposals
	ClosedProposals uint64 `json:"closedProposals"`
	// members total
	Members int64 `json:"members"`
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
	Deadline     int64  `json:"deadline"`
}

type JoinDaoWithSignParam struct {
	Sign   SignData     `json:"sign"`
	Params JoinDaoParam `json:"params"`
}

type JoinDaoParam struct {
	ChainId    int    `json:"chainId"`
	DaoAddress string `json:"daoAddress"`
	// join/quit
	JoinSwitch string `json:"joinSwitch"`
	Timestamp  int64  `json:"timestamp"`
}

type ResLeftDaoCreator struct {
	Account    string `json:"account"`
	DaoAddress string `json:"daoAddress"`
	ChainId    int    `json:"chainId"`
	Role       string `json:"role"`
}

type ResProposalsListPage struct {
	List  []ResProposalsList `json:"list"`
	Total uint64             `json:"total"`
}

type ResProposalsList struct {
	ChainId      int    `json:"chainId"`
	DaoAddress   string `json:"daoAddress"`
	DaoAddressV1 string `json:"daoAddressV1"`
	ProposalId   int    `json:"proposalId"`
	Proposer     string `json:"proposer"`
	Title        string `json:"title"`
	ContentV1    string `json:"contentV1"`
	StartTime    int64  `json:"startTime"`
	EndTime      int64  `json:"endTime"`
	Version      string `json:"version"`
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

type ReqAccountQuery struct {
	Account string `json:"account"`
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

type V1ProposalHistory struct {
	Code int              `json:"code"`
	Data []V1ProposalData `json:"data"`
	Msg  string           `json:"msg"`
}

type V1ProposalData struct {
	Id      int    `json:"id"`
	Address string `json:"address"`
	Topic1  string `json:"topic1"`
	Topic2  string `json:"topic2"`
	Data    string `json:"data"`
}

type KlaytnBlock struct {
	Code    int      `json:"code"`
	Data    []string `json:"data"`
	Message string   `json:"msg"`
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
	TokenAddress string `json:"tokenAddress"` // ReserveToken
	DaoLogo      string `json:"daoLogo"`
	DaoName      string `json:"daoName"`
	ProposalId   int    `json:"proposalId"`   // NewProposal
	ProposalName string `json:"proposalName"` // NewProposal
	ActivityId   int    `json:"activityId"`   // Airdrop
	ActivityName string `json:"activityName"` // Airdrop
	TokenLogo    string `json:"tokenLogo"`    //PublicSale
	//StartTime    int64  `json:"startTime"`
	Creator string `json:"creator"` //PublicSaleCreator
	Buyer   string `json:"buyer"`
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

type ResAccountRecordPage struct {
	List  []ResAccountRecord `json:"list"`
	Total uint64             `json:"total"`
}

type ResAccountRecord struct {
	Creator    string `json:"creator"`
	Types      string `json:"types"`
	ChainId    int    `json:"chainId"`
	Address    string `json:"address"`
	ActivityId int    `json:"activityId"`
	Avatar     string `json:"avatar"`
	DaoName    string `json:"daoName"`
	Titles     string `json:"titles"`
	Time       int64  `json:"time"`
}

type ResOverview struct {
	TotalDao        int    `json:"totalDao"`
	TotalApproveDao int    `json:"totalApproveDao"`
	TotalAccount    uint64 `json:"totalAccount"`
	TotalProposal   int    `json:"totalProposal"`
}

type ResAccountSignPage struct {
	List  []ResAccountSign `json:"list"`
	Total uint64           `json:"total"`
}

type ResAccountSign struct {
	ChainId     int    `json:"chainId"`
	DaoAddress  string `json:"daoAddress"`
	Account     string `json:"account"`
	Operate     string `json:"operate"`
	Signature   string `json:"signature"`
	Message     string `json:"message"`
	Timestamp   int64  `json:"timestamp"`
	AccountLogo string `json:"accountLogo"`
}

type FollowWithSignParam struct {
	Params FollowParam `json:"params"`
}

type FollowParam struct {
	FollowAccount string `json:"followAccount"`
	// status false: unfollow true: follow
	Status bool `json:"status"`
}

type ResAccountFollowPage struct {
	List  []ResAccountFollow `json:"list"`
	Total uint64             `json:"total"`
}

type ResAccountFollow struct {
	Account     string `json:"account"`
	FollowTime  string `json:"followTime"`
	Following   string `json:"following"`
	AccountLogo string `json:"accountLogo"`
	Nickname    string `json:"nickname"`
	// following or mutualFollowing
	Relation string `json:"relation"`
}

type ResAccountFollowersPage struct {
	List  []ResAccountFollowers `json:"list"`
	Total uint64                `json:"total"`
}

type ResAccountFollowers struct {
	Account     string `json:"account"`
	FollowTime  string `json:"followTime"`
	Followers   string `json:"followers"`
	AccountLogo string `json:"accountLogo"`
	Nickname    string `json:"nickname"`
	// following or mutualFollowing
	Relation string `json:"relation"`
}

type ResAccountTopList struct {
	Account  string `json:"account"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	FansNum  int64  `json:"fansNum"`
}

type ReqCreateSale struct {
	ChainId int    `json:"chainId"`
	Title   string `json:"title"`
	Creator string `json:"creator"`
	// SaleWay:general/discount
	SaleWay      string   `json:"saleWay"`
	SaleToken    string   `json:"saleToken"`
	SaleAmount   string   `json:"saleAmount"`
	SalePrice    string   `json:"salePrice"`
	ReceiveToken string   `json:"receiveToken"`
	LimitMin     string   `json:"limitMin"`
	LimitMax     string   `json:"limitMax"`
	StartTime    int64    `json:"startTime"`
	EndTime      int64    `json:"endTime"`
	WhiteList    []string `json:"whiteList"`
	About        string   `json:"about"`
}

type ResCreateSale struct {
	SaleId    int64  `json:"saleId"`
	Signature string `json:"signature"`
}

type ReqPurchased struct {
	SaleId    int64  `json:"saleId"`
	Account   string `json:"account"`
	BuyAmount string `json:"buyAmount"`
}

type ResPurchased struct {
	Signature string `json:"signature"`
}

type ResSwapListPage struct {
	List  []ResSwapList `json:"list"`
	Total int64         `json:"total"`
}

type ResSwapList struct {
	SaleId           int64  `json:"saleId"`
	SaleWay          string `json:"saleWay"`
	Title            string `json:"title"`
	CreateTime       int64  `json:"createTime"`
	ChainId          int    `json:"chainId"`
	Creator          string `json:"creator"`
	SaleToken        string `json:"saleToken"`
	SaleTokenImg     string `json:"saleTokenImg"`
	SaleAmount       string `json:"saleAmount"`
	SalePrice        string `json:"salePrice"`
	ReceiveToken     string `json:"receiveToken"`
	ReceiveTokenImg  string `json:"receiveTokenImg"`
	LimitMin         string `json:"limitMin"`
	LimitMax         string `json:"limitMax"`
	StartTime        int64  `json:"startTime"`
	EndTime          int64  `json:"endTime"`
	Status           string `json:"status"`
	About            string `json:"about"`
	OriginalDiscount string `json:"originalDiscount"`
	SoldAmount       string `json:"soldAmount"`
}

type ResSwapTransactionListPage struct {
	List  []ResSwapTransactionList `json:"list"`
	Total int64                    `json:"total"`
}

type ResSwapTransactionList struct {
	SaleId       int64  `json:"saleId"`
	Buyer        string `json:"buyer"`
	BuyAmount    string `json:"buy_amount"`
	PayAmount    string `json:"payAmount"`
	Time         int64  `json:"time"`
	BuyTokenName string `json:"buyTokenName"`
	PayTokenName string `json:"payTokenName"`
}

type ResSwapPrices struct {
	ChainId          int     `json:"chainId"`
	TokenAddress     string  `json:"tokenAddress"`
	Price            float64 `json:"price"`
	Img              string  `json:"img"`
	UrlCoingecko     string  `json:"urlCoingecko"`
	UrlCoinmarketcap string  `json:"urlCoinmarketcap"`
	TokenName        string  `json:"tokenName"`
	Symbol           string  `json:"symbol"`
	Decimals         int64   `json:"decimals"`
	UpdateAt         int64   `json:"updateAt"`
}

type ResIsWhite struct {
	IsWhite bool `json:"isWhite"`
}

type UpdateAccountPushSwitchParam struct {
	AllDaosICreateOrJoin bool `json:"allDaosICreateOrJoin"`
	NewDao               bool `json:"newDao"`
	AllDaoAirdrop        bool `json:"allDaoAirdrop"`
	AllDaoProposal       bool `json:"allDaoProposal"`
}

//type SignDataForTask struct {
//	Account   string `json:"account" validate:"eth_addr"`              // personal_sign address,0x
//	Signature string `json:"signature" validate:"len=130,hexadecimal"` // personal_sign sign result,no 0x
//	// Timestamp: expired timestamp, cannot exceed one day(86400)
//	Timestamp  int64  `json:"timestamp"`
//	ChainId    int64  `json:"chainId"`
//	DaoAddress string `json:"daoAddress"`
//}

type ReqJobsApply struct {
	ChainId    int64  `json:"chainId"`
	DaoAddress string `json:"daoAddress"`
	// ApplyRole: C_member/B_admin
	ApplyRole string `json:"applyRole"`
	Message   string `json:"message"`
}

type ReqCreateTask struct {
	SpacesId int64  `json:"spacesId"`
	TaskName string `json:"taskName"`
	Content  string `json:"content"`
	Deadline int64  `json:"deadline"`
	// A_low;B_medium;C_high
	Priority      string `json:"priority"`
	AssignAccount string `json:"assignAccount"`
	ProposalId    int    `json:"proposalId"`
	Reward        string `json:"reward"`
	// A_notStarted;B_inProgress;C_done;D_notStatus
	Status string `json:"status"`
}

type ReqUpdateTask struct {
	SpacesId int64  `json:"spacesId"`
	TaskId   int64  `json:"taskId"`
	TaskName string `json:"taskName"`
	Content  string `json:"content"`
	Deadline int64  `json:"deadline"`
	// A_low;B_medium;C_high
	Priority      string `json:"priority"`
	AssignAccount string `json:"assignAccount"`
	ProposalId    int    `json:"proposalId"`
	Reward        string `json:"reward"`
	// A_notStarted;B_inProgress;C_done;D_notStatus
	Status string  `json:"status"`
	Weight float64 `json:"weight"`
}

type ReqRemoveTask struct {
	SpacesId int64   `json:"spacesId"`
	TaskId   []int64 `json:"taskId"`
}

type ReqCreateTeamSpaces struct {
	ChainId    int64  `json:"chainId"`
	DaoAddress string `json:"daoAddress"`
	Title      string `json:"title"`
	// public;private
	Access string `json:"access"`
}

type ReqUpdateTeamSpaces struct {
	TeamSpacesId int64  `json:"teamSpacesId"`
	ChainId      int64  `json:"chainId"`
	DaoAddress   string `json:"daoAddress"`
	Title        string `json:"title"`
	Access       string `json:"access"`
	Url          string `json:"url"`
}

type ResTeamSpacesList struct {
	TeamSpacesId       int64  `json:"teamSpacesId"`
	ChainId            int64  `json:"chainId"`
	DaoAddress         string `json:"daoAddress"`
	Creator            string `json:"creator"`
	AvatarCreator      string `json:"avatarCreator"`
	NicknameCreator    string `json:"nicknameCreator"`
	Title              string `json:"title"`
	Url                string `json:"url"`
	LastEditTime       int64  `json:"lastEditTime"`
	LastEditBy         string `json:"lastEditBy"`
	AvatarLastEditBy   string `json:"avatarLastEditBy"`
	NicknameLastEditBy string `json:"nicknameLastEditBy"`
	Access             string `json:"access"`
}

type ReqRemoveTeamSpaces struct {
	ChainId      int64  `json:"chainId"`
	DaoAddress   string `json:"daoAddress"`
	TeamSpacesId int64  `json:"teamSpacesId"`
}

type ResTaskList struct {
	SpacesId       int64   `json:"spacesId"`
	TaskId         int64   `json:"taskId"`
	TaskName       string  `json:"taskName"`
	Deadline       int64   `json:"deadline"`
	Priority       string  `json:"priority"`
	AssignAccount  string  `json:"assignAccount"`
	AssignAvatar   string  `json:"assignAvatar"`
	AssignNickname string  `json:"assignNickname"`
	Status         string  `json:"status"`
	Weight         float64 `json:"weight"`
}

type ResTaskDetail struct {
	TaskId         int64   `json:"taskId"`
	SpacesId       int64   `json:"spacesId"`
	TaskName       string  `json:"taskName"`
	Content        string  `json:"content"`
	Deadline       int64   `json:"deadline"`
	Priority       string  `json:"priority"`
	AssignAccount  string  `json:"assignAccount"`
	AssignAvatar   string  `json:"assignAvatar"`
	AssignNickname string  `json:"assignNickname"`
	ProposalId     int     `json:"proposalId"`
	Reward         string  `json:"reward"`
	Status         string  `json:"status"`
	Weight         float64 `json:"weight"`
}

type ResJobsList struct {
	JobId      int64  `json:"jobId"`
	ChainId    int    `json:"chainId"`
	DaoAddress string `json:"daoAddress"`
	Account    string `json:"account"`
	Jobs       string `json:"jobs"`
	Avatar     string `json:"avatar"`
	Nickname   string `json:"nickname"`
	Twitter    string `json:"twitter"`
	Discord    string `json:"discord"`
	Youtube    string `json:"youtube"`
	Opensea    string `json:"opensea"`
}

type ResJobsApplyList struct {
	ApplyId    int64  `json:"applyId"`
	ChainId    int    `json:"chainId"`
	DaoAddress string `json:"daoAddress"`
	Account    string `json:"account"`
	Avatar     string `json:"avatar"`
	Nickname   string `json:"nickname"`
	ApplyRole  string `json:"applyRole"`
	ApplyTime  int64  `json:"applyTime"`
	Message    string `json:"message"`
}

type ReqJobsApplyReview struct {
	ChainId     int64  `json:"chainId"`
	DaoAddress  string `json:"daoAddress"`
	JobsApplyId int64  `json:"jobsApplyId"`
	IsPass      bool   `json:"isPass"`
}

type ReqJobsAlter struct {
	ChainId    int64  `json:"chainId"`
	DaoAddress string `json:"daoAddress"`
	JobId      int64  `json:"jobId"`
	// ChangeTo: B_admin/C_member/noRole
	ChangeTo string `json:"changeTo"`
}
