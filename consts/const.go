package consts

import (
	_ "crypto/cipher"
	"stp_dao_v2/utils"
)

const (
	TbNameEventHistorical = "event_historical_data"
	TbNameScanTask        = "scan_task"
	TbNameAccount         = "tb_account"
	TbNameNonce           = "tb_nonce"
	TbNameMember          = "tb_member"
	TbNameDao             = "tb_dao"
	TbNameCategory        = "tb_category"
	TbNameDaoCategory     = "tb_dao_category"
	TbNameHolderData      = "tb_holder_data"

	EvCreateDao      = "CreateDao"
	EvCreateProposal = "CreateProposal"
	EvVote           = "Vote"
	EvCancelProposal = "CancelProposal"
	EvAdmin          = "Admin"
	EvSetting        = "Setting"
	EvCreateERC20    = "CreateERC20"
	EvTransfer       = "Transfer"

	MaxValue    = 0x7fffffff
	ZeroAddress = "0x0000000000000000000000000000000000000000000000000000000000000000"

	SignMessagePrefix = "Welcome come Clique"
)

func NeedScanUrl() map[string]map[string]string {
	ScanMap := make(map[string]map[string]string)

	ScanMap["url01"] = make(map[string]string, 2)
	ScanMap["url01"]["chainId"] = "4"
	ScanMap["url01"]["url"] = "https://eth-rinkeby.blockvision.org/v1/2AruYS477FHohqhPsDiMNWxxiVL"

	ScanMap["url02"] = make(map[string]string, 2)
	ScanMap["url02"]["chainId"] = "1"
	ScanMap["url02"]["url"] = "https://eth-mainnet.blockvision.org/v1/29dy5nutBpNq2hRJVWt9xbjUKxC"

	return ScanMap
}

func EventTypes(event string) string {
	var (
		//CreateDAO(address indexed creator, address indexed daoAddress, uint256 chainId, address tokenAddress)
		createDao = utils.Keccak256("CreateDAO(address,address,uint256,address)")

		//CreateProposal(uint256 indexed proposalId, address indexed proposer, uint256 nonce, uint256 startTime, uint256 endTime)
		createProposal = utils.Keccak256("CreateProposal(uint256,address,uint256,uint256,uint256)")

		//CancelProposal(uint256 indexed proposalId)
		cancelProposal = utils.Keccak256("CancelProposal(uint256)")

		//Vote(uint256 indexed proposalId, address indexed voter, uint256 indexed optionIndex, uint256 amount, uint256 nonce)
		vote = utils.Keccak256("Vote(uint256,address,uint256,uint256,uint256)")

		//Admin(address indexed admin, bool enable)
		admin = utils.Keccak256("Admin(address,bool)")

		//Setting(uint256 indexed settingType)
		setting = utils.Keccak256("Setting(uint256)")

		//CreateERC20(address indexed creator, address token)
		createERC20 = utils.Keccak256("CreateERC20(address,address)")

		//Transfer(address indexed from, address indexed to, uint256 value);
		transfer = utils.Keccak256("Transfer(address,address,uint256)")
	)
	switch event {
	case createDao:
		event = EvCreateDao
		break
	case createProposal:
		event = EvCreateProposal
		break
	case cancelProposal:
		event = EvCancelProposal
		break
	case vote:
		event = EvVote
		break
	case admin:
		event = EvAdmin
		break
	case setting:
		event = EvSetting
		break
	case createERC20:
		event = EvCreateERC20
		break
	case transfer:
		event = EvTransfer
		break
	default:
		event = "Undefined"
		break
	}
	return event
}
