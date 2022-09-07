package models

import "database/sql"

type EventHistoricalModel struct {
	Id               uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime       string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime       string `db:"update_time,omitempty" sqler:"skips"`
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
	Id                uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime        string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime        string `db:"update_time,omitempty" sqler:"skips"`
	DaoLogo           string `db:"dao_logo"`
	DaoName           string `db:"dao_name"`
	DaoAddress        string `db:"dao_address"`
	Creator           string `db:"creator"`
	Handle            string `db:"handle"`
	Description       string `db:"description"`
	ChainId           int    `db:"chain_id"`
	TokenChainId      int    `db:"token_chain_id"`
	TokenAddress      string `db:"token_address"`
	ProposalThreshold string `db:"proposal_threshold"`
	VotingQuorum      string `db:"voting_quorum"`
	VotingPeriod      int    `db:"voting_period"`
	VotingType        string `db:"voting_type"`
	Twitter           string `db:"twitter"`
	Github            string `db:"github"`
	Discord           string `db:"discord"`
	UpdateBool        bool   `db:"update_bool"`
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
}

type ProposalModel struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	ProposalId int    `db:"proposal_id"`
	ChainId    int    `db:"chain_id"`
	DaoAddress string `db:"dao_address"`
	Proposer   string `db:"proposer"`
	StartTime  int64  `db:"start_time"`
	EndTime    int64  `db:"end_time"`
}

type AirdropAddressModel struct {
	Id         uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
	Content    string `db:"content"`
}

type ActivityModel struct {
	Id           uint64 `db:"id,omitempty" sqler:"skips"`
	CreateTime   string `db:"create_time,omitempty" sqler:"skips"`
	UpdateTime   string `db:"update_time,omitempty" sqler:"skips"`
	ChainId      int    `db:"chain_id"`
	DaoAddress   string `db:"dao_address"`
	Creator      string `db:"creator"`
	AirdropId    int    `db:"airdrop_id"`
	TokenAddress string `db:"token_address"`
	Amount       string `db:"amount"`
	MerkleRoot   string `db:"merkle_root"`
	StartTime    int    `db:"start_time"`
	EndTime      int    `db:"end_time"`
	Price        string `db:"price"`
}

//type VoteModel struct {
//	Id         uint64 `db:"id,omitempty" sqler:"skips"`
//	CreateTime string `db:"create_time,omitempty" sqler:"skips"`
//	UpdateTime string `db:"update_time,omitempty" sqler:"skips"`
//	ProposalId int    `db:"proposal_id"`
//	ChainId    int    `db:"chain_id"`
//	DaoAddress string `db:"dao_address"`
//	Voter      string `db:"voter"`
//}
//
//type VoteVotesModel struct {
//	Id          uint64 `db:"id,omitempty" sqler:"skips"`
//	CreateTime  string `db:"create_time,omitempty" sqler:"skips"`
//	UpdateTime  string `db:"update_time,omitempty" sqler:"skips"`
//	VoteId      int    `db:"vote_id"`
//	OptionIndex int    `db:"option_index"`
//	Amount      string `db:"amount"`
//}
