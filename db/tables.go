package db

import (
	"database/sql"
	oo "github.com/Anna2024/liboo"
	"stp_dao_v2/consts"
	"stp_dao_v2/db/o"
)

type TbEventHistoricalModel struct {
	Id               uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime       string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime       string `db:"update_time,omitempty" sqler:"skips"`
	MessageSender    string `db:"message_sender"`
	EventType        string `db:"event_type"`
	Address          string `db:"address"`
	Topic0           string `db:"topic0"`
	Topic1           string `db:"topic1"`
	Topic2           string `db:"topic2"`
	Topic3           string `db:"topic3"`
	Data             string `db:"data"`
	BlockNumber      string `db:"block_number"`
	TimeStamp        string `db:"time_stamp"`
	GasPrice         string `db:"gas_price"`
	GasUsed          string `db:"gas_used"`
	LogIndex         string `db:"log_index"`
	TransactionHash  string `db:"transaction_hash"`
	TransactionIndex string `db:"transaction_index"`
	ChainId          int    `db:"chain_id"`
}

type TbScanTaskModel struct {
	Id              uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime      string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime      string `db:"update_time,omitempty" sqler:"skips"`
	EventType       string `db:"event_type"`
	Address         string `db:"address"`
	LastBlockNumber int64  `db:"last_block_number"`
	RestParameter   string `db:"rest_parameter"`
	ChainId         uint64 `db:"chain_id"`
}

func GetTbScanTaskModel(w ...[][]interface{}) (data TbScanTaskModel, err error) {
	err = oo.SqlGet(o.Pre(consts.TbNameScanTask, w).Select(), &data)
	return data, err
}

type TbAccountModel struct {
	Id           int64          `db:"id,omitempty" sqler:"skips"`
	CreateTime   string         `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime   string         `db:"update_time,omitempty" sqler:"skips"`
	Account      string         `db:"account"`
	AccountLogo  sql.NullString `db:"account_logo"`
	Nickname     sql.NullString `db:"nickname"`
	Introduction sql.NullString `db:"introduction"`
	Twitter      sql.NullString `db:"twitter"`
	Github       sql.NullString `db:"github"`
	Discord      sql.NullString `db:"discord"`
	Email        sql.NullString `db:"email"`
	Country      sql.NullString `db:"country"`
	Youtube      sql.NullString `db:"youtube"`
	Opensea      sql.NullString `db:"opensea"`
	PushSwitch   int            `db:"push_switch"`
	FansNum      int64          `db:"fans_num"`
}

func GetTbAccountModel(w ...[][]interface{}) (data TbAccountModel, err error) {
	err = oo.SqlGet(o.Pre(consts.TbNameAccount, w).Select(), &data)
	return data, err
}

type TbAccountFollowModel struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	Account    string `db:"account"`
	Followed   string `db:"followed"`
	Status     bool   `db:"status"`
}

type TbAccountSignModel struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	ChainId    int    `db:"chain_id"`
	DaoAddress string `db:"dao_address"`
	Account    string `db:"account"`
	Operate    string `db:"operate"`
	Signature  string `db:"signature"`
	Message    string `db:"message"`
	Timestamp  int64  `db:"timestamp"`
}

type TbAccountRecordModel struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	Creator    string `db:"creator"`
	Types      string `db:"types"`
	ChainId    int    `db:"chain_id"`
	Address    string `db:"address"`
	ActivityId int    `db:"activity_id"`
	Avatar     string `db:"avatar"`
	DaoName    string `db:"dao_name"`
	Titles     string `db:"titles"`
	Time       int64  `db:"time"`
	UpdateBool bool   `db:"update_bool"`
}

type NonceModel struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	ChainId    int    `db:"chain_id"`
	Account    string `db:"account"`
	Nonce      uint64 `db:"nonce"`
}

type TbMemberModel struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	DaoAddress string `db:"dao_address"`
	ChainId    int    `db:"chain_id"`
	Account    string `db:"account"`
	JoinSwitch int    `db:"join_switch"`
}

type TbDaoModel struct {
	Id                uint64        `db:"id,omitempty" sqler:"skips"`
	CreateTime        string        `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime        string        `db:"update_time,omitempty" sqler:"skips"`
	DaoLogo           string        `db:"dao_logo"`
	DaoName           string        `db:"dao_name"`
	DaoAddress        string        `db:"dao_address"`
	Creator           string        `db:"creator"`
	Handle            string        `db:"handle"`
	Description       string        `db:"description"`
	ChainId           int           `db:"chain_id"`
	TokenChainId      int           `db:"token_chain_id"`
	TokenAddress      string        `db:"token_address"`
	ProposalThreshold string        `db:"proposal_threshold"`
	VotingQuorum      string        `db:"voting_quorum"`
	VotingPeriod      int           `db:"voting_period"`
	VotingType        string        `db:"voting_type"`
	Twitter           string        `db:"twitter"`
	Github            string        `db:"github"`
	Discord           string        `db:"discord"`
	Website           string        `db:"website"`
	UpdateBool        bool          `db:"update_bool"`
	Weight            sql.NullInt64 `db:"weight"`
	Approve           bool          `db:"approve"`
	Deprecated        bool          `db:"deprecated"`
	Members           int64         `db:"members"`
	TotalProposals    int64         `db:"total_proposals"`
}

func GetTbDao(w ...[][]interface{}) (data TbDaoModel, err error) {
	err = oo.SqlGet(o.Pre(consts.TbNameDao, w).Select(), &data)
	return data, err
}

func SelectTbDaoModel(w ...[][]interface{}) (arr []TbDaoModel, err error) {
	err = oo.SqlSelect(o.Pre(consts.TbNameDao, w).Select(), &arr)
	return arr, err
}

type CategoryModel struct {
	Id           uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime   string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime   string `db:"update_time,omitempty" sqler:"skips"`
	CategoryName string `db:"category_name"`
}

type DaoCategoryModel struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	DaoId      int    `db:"dao_id"`
	CategoryId int    `db:"category_id"`
}

type TbHolderDataModel struct {
	Id            uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime    string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime    string `db:"update_time,omitempty" sqler:"skips"`
	TokenAddress  string `db:"token_address"`
	HolderAddress string `db:"holder_address"`
	Balance       string `db:"balance"`
	ChainId       int    `db:"chain_id"`
}

type TbAdminModel struct {
	Id           uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime   string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime   string `db:"update_time,omitempty" sqler:"skips"`
	DaoAddress   string `db:"dao_address"`
	ChainId      int    `db:"chain_id"`
	Account      string `db:"account"`
	AccountLevel string `db:"account_level"`
}

func SelectTbAdmin(w ...[][]interface{}) (arr []TbAdminModel, err error) {
	err = oo.SqlSelect(o.Pre(consts.TbNameAdmin, w).Select(), &arr)
	return arr, err
}

func GetTbAdmin(w ...[][]interface{}) (data TbAdminModel, err error) {
	err = oo.SqlGet(o.Pre(consts.TbNameAdmin, w).Select(), &data)
	return data, err
}

type ErrorInfoModel struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	Title      string `db:"title"`
	Content    string `db:"content"`
	Func       string `db:"func"`
	Params     string `db:"params"`
}

type ProposalInfoModel struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	Uuid       string `db:"uuid"`
	Content    string `db:"content"`
}

type TbTokensImgModel struct {
	Id           uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime   string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime   string `db:"update_time,omitempty" sqler:"skips"`
	TokenChainId int    `db:"token_chain_id"`
	TokenAddress string `db:"token_address"`
	Thumb        string `db:"thumb"`
	Small        string `db:"small"`
	Large        string `db:"large"`
	OwnImg       string `db:"own_img"`
}

type TbProposalModel struct {
	Id          uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime  string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime  string `db:"update_time,omitempty" sqler:"skips"`
	ProposalId  int    `db:"proposal_id"`
	ChainId     int    `db:"chain_id"`
	DaoAddress  string `db:"dao_address"`
	Title       string `db:"title"`
	IdV1        int    `db:"id_v1"`
	ContentV1   string `db:"content_v1"`
	Proposer    string `db:"proposer"`
	StartTime   int64  `db:"start_time"`
	EndTime     int64  `db:"end_time"`
	Version     string `db:"version"`
	Deprecated  bool   `db:"deprecated"`
	BlockNumber string `db:"block_number"`
}

type TbProposalV1Model struct {
	Id           uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime   string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime   string `db:"update_time,omitempty" sqler:"skips"`
	ChainId      int    `db:"chain_id"`
	DaoAddress   string `db:"dao_address"`
	DaoAddressV1 string `db:"dao_address_v1"`
	VotingV1     string `db:"voting_v1"`
	StartIdV1    int    `db:"start_id_v1"`
}

type TbAirdropModel struct {
	Id                 uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime         string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime         string `db:"update_time,omitempty" sqler:"skips"`
	Creator            string `db:"creator"`
	ChainId            int    `db:"chain_id"`
	DaoAddress         string `db:"dao_address"`
	Title              string `db:"title"`
	AirdropAddress     string `db:"airdrop_address"`
	Description        string `db:"description"`
	CollectInformation string `db:"collect_information"`
	TokenChainId       int    `db:"token_chain_id"`
	TokenAddress       string `db:"token_address"`
	MaxAirdropAmount   string `db:"max_airdrop_amount"`
	StartTime          int64  `db:"start_time"`
	EndTime            int64  `db:"end_time"`
	AirdropStartTime   int64  `db:"airdrop_start_time"`
	AirdropEndTime     int64  `db:"airdrop_end_time"`
}

type TbAirdropPrepareModel struct {
	Id             uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime     string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime     string `db:"update_time,omitempty" sqler:"skips"`
	AirdropId      int    `db:"airdrop_id"`
	Root           string `db:"root"`
	PrepareAddress string `db:"prepare_address"`
}

type TbAirdropUserSubmit struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	AirdropId  int    `db:"airdrop_id"`
	Account    string `db:"account"`
	SubmitInfo string `db:"submit_info"`
	Timestamp  int64  `db:"timestamp"`
}

type TbActivityModel struct {
	Id               uint64        `db:"id,omitempty" sqler:"skips"`
	CreateTime       string        `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime       string        `db:"update_time,omitempty" sqler:"skips"`
	Types            string        `db:"types"`
	ChainId          int           `db:"chain_id"`
	DaoAddress       string        `db:"dao_address"`
	Creator          string        `db:"creator"`
	ActivityId       int           `db:"activity_id"`
	TokenChainId     int           `db:"token_chain_id"`
	TokenAddress     string        `db:"token_address"`
	StakingAmount    string        `db:"staking_amount"`
	AirdropAmount    string        `db:"airdrop_amount"`
	MerkleRoot       string        `db:"merkle_root"`
	StartTime        int64         `db:"start_time"`
	EndTime          int64         `db:"end_time"`
	AirdropStartTime int64         `db:"airdrop_start_time"`
	AirdropEndTime   int64         `db:"airdrop_end_time"`
	PublishTime      int64         `db:"publish_time"`
	Price            string        `db:"price"`
	Weight           sql.NullInt64 `db:"weight"`
	Deprecated       bool          `db:"deprecated"`
}

type ClaimedModel struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	ChainId    int    `db:"chain_id"`
	DaoAddress string `db:"dao_address"`
	AirdropId  int    `db:"airdrop_id"`
	IndexId    int    `db:"index_id"`
	Account    string `db:"account"`
	Amount     string `db:"amount"`
}

type TbVoteModel struct {
	Id          uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime  string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime  string `db:"update_time,omitempty" sqler:"skips"`
	ChainId     int    `db:"chain_id"`
	DaoAddress  string `db:"dao_address"`
	ProposalId  int    `db:"proposal_id"`
	Voter       string `db:"voter"`
	OptionIndex int    `db:"option_index"`
	Amount      string `db:"amount"`
	Nonce       int    `db:"nonce"`
}

type TbNotificationModel struct {
	Id           uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime   string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime   string `db:"update_time,omitempty" sqler:"skips"`
	ChainId      int    `db:"chain_id"`
	DaoAddress   string `db:"dao_address"`
	Types        string `db:"types"`
	DaoLogo      string `db:"dao_logo"`
	DaoName      string `db:"dao_name"`
	ActivityId   int    `db:"activity_id"`
	ActivityName string `db:"activity_name"`
	StartTime    int64  `db:"start_time"`
	UpdateBool   bool   `db:"update_bool"`
}

type TbNotificationAccountModel struct {
	Id               uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime       string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime       string `db:"update_time,omitempty" sqler:"skips"`
	NotificationId   int    `db:"notification_id"`
	Account          string `db:"account"`
	AlreadyRead      bool   `db:"already_read"`
	NotificationTime int64  `db:"notification_time"`
}

type TbHandleLockModel struct {
	Id           uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime   string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime   string `db:"update_time,omitempty" sqler:"skips"`
	Handle       string `db:"handle"`
	HandleKeccak string `db:"handle_keccak"`
	LockBlock    int    `db:"lock_block"`
	ChainId      int    `db:"chain_id"`
	Account      string `db:"account"`
}

type TbSwap struct {
	Id               int64  `db:"id,omitempty" sqler:"skips"`
	CreateTime       string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime       string `db:"update_time,omitempty" sqler:"skips"`
	ChainId          int    `db:"chain_id"`
	Title            string `db:"title"`
	Creator          string `db:"creator"`
	SaleWay          string `db:"sale_way"`
	SaleToken        string `db:"sale_token"`
	SaleTokenImg     string `db:"sale_token_img"`
	SaleAmount       string `db:"sale_amount"`
	SalePrice        string `db:"sale_price"`
	OriginalDiscount string `db:"original_discount"`
	ReceiveToken     string `db:"receive_token"`
	ReceiveTokenImg  string `db:"receive_token_img"`
	LimitMin         string `db:"limit_min"`
	LimitMax         string `db:"limit_max"`
	StartTime        int64  `db:"start_time"`
	EndTime          int64  `db:"end_time"`
	WhiteList        string `db:"white_list"`
	About            string `db:"about"`
	SoldAmount       string `db:"sold_amount"`
	Status           string `db:"status"`
}

type TbSwapToken struct {
	Id               int64   `db:"id,omitempty" sqler:"skips"`
	CreateTime       string  `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime       string  `db:"update_time,omitempty" sqler:"skips"`
	ChainId          int     `db:"chain_id"`
	TokenAddress     string  `db:"token_address"`
	TokenName        string  `db:"token_name"`
	Symbol           string  `db:"symbol"`
	Decimals         int64   `db:"decimals"`
	CoinIds          string  `db:"coin_ids"`
	Price            float64 `db:"price"`
	Img              string  `db:"img"`
	UrlCoingecko     string  `db:"url_coingecko"`
	UrlCoinmarketcap string  `db:"url_coinmarketcap"`
	IsSync           bool    `db:"isSync"`
}

type TbSwapTransaction struct {
	Id         int64  `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	SaleId     int64  `db:"sale_id"`
	Buyer      string `db:"buyer"`
	BuyAmount  string `db:"buy_amount"`
	PayAmount  string `db:"pay_amount"`
	Time       int64  `db:"time"`
	ChainId    int64  `db:"chain_id"`
	BuyToken   string `db:"buy_token"`
	PayToken   string `db:"pay_token"`
}

type TbSysConfig struct {
	Id           int64  `db:"id,omitempty" sqler:"skips"`
	CreateTime   string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime   string `db:"update_time,omitempty" sqler:"skips"`
	CfgName      string `db:"cfg_name"`
	CfgVal       string `db:"cfg_val"`
	CfgType      string `db:"cfg_type"`
	CfgComment   string `db:"cfg_comment"`
	CfgIsEnabled bool   `db:"cfg_is_enabled"`
}

type TbTask struct {
	Id            int64   `db:"id,omitempty" sqler:"skips"`
	CreateTime    string  `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime    string  `db:"update_time,omitempty" sqler:"skips"`
	SpacesId      int64   `db:"spaces_id"`
	TaskName      string  `db:"task_name"`
	Content       string  `db:"content"`
	Deadline      int64   `db:"deadline"`
	Priority      string  `db:"priority"`
	AssignAccount string  `db:"assign_account"`
	ProposalId    int     `db:"proposal_id"`
	Reward        string  `db:"reward"`
	Status        string  `db:"status"`
	Weight        float64 `db:"weight"`
	IsTrash       bool    `db:"is_trash"`
}

func SelectTbTask(w ...[][]interface{}) (arr []TbTask, err error) {
	err = oo.SqlSelect(o.Pre(consts.TbTask, w).Select(), &arr)
	return arr, err
}

func GetTbTask(w ...[][]interface{}) (data TbTask, err error) {
	err = oo.SqlGet(o.Pre(consts.TbTask, w).Select(), &data)
	return data, err
}

type TbJobs struct {
	Id         int64  `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	ChainId    int    `db:"chain_id"`
	DaoAddress string `db:"dao_address"`
	Account    string `db:"account"`
	Job        string `db:"job"`
}

func GetTbJobs(w ...[][]interface{}) (data TbJobs, err error) {
	err = oo.SqlGet(o.Pre(consts.TbJobs, w).Select(), &data)
	return data, err
}

func SelectTbJobs(w ...[][]interface{}) (arr []TbJobs, err error) {
	err = oo.SqlSelect(o.Pre(consts.TbJobs, w).Select(), &arr)
	return arr, err
}

type TbJobsPublish struct {
	Id         int64  `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	ChainId    int64  `db:"chain_id"`
	DaoAddress string `db:"dao_address"`
	Title      string `db:"title"`
	JobBio     string `db:"job_bio"`
	Access     string `db:"access"`
}

func GetTbJobsPublish(w ...[][]interface{}) (data TbJobsPublish, err error) {
	err = oo.SqlGet(o.Pre(consts.TbJobsPublish, w).Select(), &data)
	return data, err
}

func SelectTbJobsPublish(w ...[][]interface{}) (arr []TbJobsPublish, err error) {
	err = oo.SqlSelect(o.Pre(consts.TbJobsPublish, w).Select(), &arr)
	return arr, err
}

type TbJobsApply struct {
	Id         int64  `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	ChainId    int    `db:"chain_id"`
	DaoAddress string `db:"dao_address"`
	Account    string `db:"account"`
	ApplyRole  string `db:"apply_role"`
	Message    string `db:"message"`
	Status     string `db:"status"`
}

func GetTbJobsApply(w ...[][]interface{}) (data TbJobsApply, err error) {
	err = oo.SqlGet(o.Pre(consts.TbJobsApply, w).Select(), &data)
	return data, err
}

type TbTeamSpaces struct {
	Id           int64  `db:"id,omitempty" sqler:"skips"`
	CreateTime   string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime   string `db:"update_time,omitempty" sqler:"skips"`
	ChainId      int64  `db:"chain_id"`
	DaoAddress   string `db:"dao_address"`
	Creator      string `db:"creator"`
	Title        string `db:"title"`
	Url          string `db:"url"`
	LastEditTime int64  `db:"last_edit_time"`
	LastEditBy   string `db:"last_edit_by"`
	Access       string `db:"access"`
	IsTrash      bool   `db:"is_trash"`
}

func SelectTbTeamSpaces(w ...[][]interface{}) (arr []TbTeamSpaces, err error) {
	err = oo.SqlSelect(o.Pre(consts.TbTeamSpaces, w).Select(), &arr)
	return arr, err
}

func GetTbTeamSpaces(w ...[][]interface{}) (data TbTeamSpaces, err error) {
	err = oo.SqlGet(o.Pre(consts.TbTeamSpaces, w).Select(), &data)
	return data, err
}

type TbSBT struct {
	Id           uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime   string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime   string `db:"update_time,omitempty" sqler:"skips"`
	ChainId      int64  `db:"chain_id"`
	DaoAddress   string `db:"dao_address"`
	TokenChainId int64  `db:"token_chain_id"`
	TokenAddress string `db:"token_address"`
	FileUrl      string `db:"file_url"`
	ItemName     string `db:"item_name"`
	Symbol       string `db:"symbol"`
	Introduction string `db:"introduction"`
	TotalSupply  uint64 `db:"total_supply"`
	StartTime    int64  `db:"start_time"`
	EndTime      int64  `db:"end_time"`
	Way          string `db:"way"`
	WhiteList    string `db:"whitelist"`
	Status       string `db:"status"`
}

func GetTbSBT(w ...[][]interface{}) (data TbSBT, err error) {
	err = oo.SqlGet(o.Pre(consts.TbSBT, w).Select(), &data)
	return data, err
}
func SelectTbSBT(w ...[][]interface{}) (arr []TbSBT, err error) {
	err = oo.SqlSelect(o.Pre(consts.TbSBT, w).Select(), &arr)
	return arr, err
}

type TbSBTClaim struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	SBTId      uint64 `db:"sbt_id"`
	Account    string `db:"account"`
}
