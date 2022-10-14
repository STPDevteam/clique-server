package models

import "database/sql"

type EventHistoricalModel struct {
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

type ScanTaskModel struct {
	Id              uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime      string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime      string `db:"update_time,omitempty" sqler:"skips"`
	EventType       string `db:"event_type"`
	Address         string `db:"address"`
	LastBlockNumber int    `db:"last_block_number"`
	RestParameter   string `db:"rest_parameter"`
	ChainId         uint64 `db:"chain_id"`
}

type AccountModel struct {
	Id           uint64         `db:"id,omitempty" sqler:"skips"`
	CreateTime   string         `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime   string         `db:"update_time,omitempty" sqler:"skips"`
	Account      string         `db:"account"`
	AccountLogo  sql.NullString `db:"account_logo"`
	Nickname     sql.NullString `db:"nickname"`
	Introduction sql.NullString `db:"introduction"`
	Twitter      sql.NullString `db:"twitter"`
	Github       sql.NullString `db:"github"`
	Discord      sql.NullString `db:"discord"`
}

type AccountRecordModel struct {
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

type MemberModel struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	DaoAddress string `db:"dao_address"`
	ChainId    int    `db:"chain_id"`
	Account    string `db:"account"`
	JoinSwitch int    `db:"join_switch"`
}

type DaoModel struct {
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

type HolderDataModel struct {
	Id            uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime    string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime    string `db:"update_time,omitempty" sqler:"skips"`
	TokenAddress  string `db:"token_address"`
	HolderAddress string `db:"holder_address"`
	Balance       string `db:"balance"`
	ChainId       int    `db:"chain_id"`
}

type AdminModel struct {
	Id           uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime   string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime   string `db:"update_time,omitempty" sqler:"skips"`
	DaoAddress   string `db:"dao_address"`
	ChainId      int    `db:"chain_id"`
	Account      string `db:"account"`
	AccountLevel string `db:"account_level"`
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

type TokensImgModel struct {
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

type ProposalModel struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	ProposalId int    `db:"proposal_id"`
	ChainId    int    `db:"chain_id"`
	DaoAddress string `db:"dao_address"`
	Title      string `db:"title"`
	IdV1       int    `db:"id_v1"`
	ContentV1  string `db:"content_v1"`
	Proposer   string `db:"proposer"`
	StartTime  int64  `db:"start_time"`
	EndTime    int64  `db:"end_time"`
	Version    string `db:"version"`
}

type ProposalV1Model struct {
	Id           uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime   string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime   string `db:"update_time,omitempty" sqler:"skips"`
	ChainId      int    `db:"chain_id"`
	DaoAddress   string `db:"dao_address"`
	DaoAddressV1 string `db:"dao_address_v1"`
	VotingV1     string `db:"voting_v1"`
	StartIdV1    int    `db:"start_id_v1"`
}

type AirdropModel struct {
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

type AirdropPrepareModel struct {
	Id             uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime     string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime     string `db:"update_time,omitempty" sqler:"skips"`
	AirdropId      int    `db:"airdrop_id"`
	Root           string `db:"root"`
	PrepareAddress string `db:"prepare_address"`
}

type AirdropUserSubmit struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	AirdropId  int    `db:"airdrop_id"`
	Account    string `db:"account"`
	SubmitInfo string `db:"submit_info"`
	Timestamp  int64  `db:"timestamp"`
}

type ActivityModel struct {
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

type VoteModel struct {
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

type NotificationModel struct {
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

type NotificationAccountModel struct {
	Id               uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime       string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime       string `db:"update_time,omitempty" sqler:"skips"`
	NotificationId   int    `db:"notification_id"`
	Account          string `db:"account"`
	AlreadyRead      bool   `db:"already_read"`
	NotificationTime int64  `db:"notification_time"`
}

type HandleLockModel struct {
	Id           uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime   string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime   string `db:"update_time,omitempty" sqler:"skips"`
	Handle       string `db:"handle"`
	HandleKeccak string `db:"handle_keccak"`
	LockBlock    int    `db:"lock_block"`
	ChainId      int    `db:"chain_id"`
	Account      string `db:"account"`
}
