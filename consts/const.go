package consts

import (
	_ "crypto/cipher"
	"stp_dao_v2/utils"
)

const (
	TbNameEventHistorical     = "event_historical_data"
	TbNameScanTask            = "scan_task"
	TbNameAccount             = "tb_account"
	TbNameAccountRecord       = "tb_account_record"
	TbNameNonce               = "tb_nonce"
	TbNameMember              = "tb_member"
	TbNameDao                 = "tb_dao"
	TbNameCategory            = "tb_category"
	TbNameDaoCategory         = "tb_dao_category"
	TbNameHolderData          = "tb_holder_data"
	TbNameAdmin               = "tb_admin"
	TbNameErrorInfo           = "error_info"
	TbNameProposalInfo        = "tb_proposal_info"
	TbNameTokensImg           = "tb_tokens_img"
	TbNameProposal            = "tb_proposal"
	TbNameProposalV1          = "tb_proposal_v1"
	TbNameAirdrop             = "tb_airdrop"
	TbNameAirdropPrepare      = "tb_airdrop_prepare"
	TbNameAirdropUserSubmit   = "tb_airdrop_user_submit"
	TbNameActivity            = "tb_activity"
	TbNameClaimed             = "tb_claimed"
	TbNameVote                = "tb_vote"
	TbNameNotification        = "tb_notification"
	TbNameNotificationAccount = "tb_notification_account"
	TbNameHandleLock          = "tb_handle_lock"

	EvCreateDao            = "CreateDao"
	EvCreateProposal       = "CreateProposal"
	EvVote                 = "Vote"
	EvCancelProposal       = "CancelProposal"
	EvAdmin                = "Admin"
	EvSetting              = "Setting"
	EvOwnershipTransferred = "OwnershipTransferred"
	EvCreateERC20          = "CreateERC20"
	EvTransfer             = "Transfer"
	EvClaimReserve         = "ClaimReserve"
	EvCreateAirdrop        = "CreateAirdrop"
	EvSettleAirdrop        = "SettleAirdrop"
	EvClaimed              = "Claimed"

	LevelSuperAdmin = "superAdmin"
	LevelAdmin      = "admin"
	LevelMember     = "member"
	LevelNoRole     = "noRole"

	TypesNameAirdrop      = "Airdrop"
	TypesNamePublicSale   = "PublicSale" // notification
	TypesNameNewProposal  = "NewProposal"
	TypesNameReserveToken = "ReserveToken"

	MaxValue        = 0x7fffffff
	MaxIntUnsigned  = 4294967295
	ZeroAddress0x64 = "0x0000000000000000000000000000000000000000000000000000000000000000"
	ZeroAddress0x40 = "0x0000000000000000000000000000000000000000"

	SignMessagePrefix = "Welcome come Clique"
)

func EventTypes(event string) string {
	var (
		//CreateDAO(uint256 indexed handler, address indexed creator, address indexed daoAddress, uint256 chainId, address tokenAddress)
		createDao = utils.Keccak256("CreateDAO(uint256,address,address,uint256,address)")
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
		//OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
		ownershipTransferred = utils.Keccak256("OwnershipTransferred(address,address)")

		//CreateERC20(address indexed creator, address token)
		createERC20 = utils.Keccak256("CreateERC20(address,address)")
		//Transfer(address indexed from, address indexed to, uint256 value);
		transfer = utils.Keccak256("Transfer(address,address,uint256)")

		//ClaimReserve(address indexed account, address indexed token, uint256 amount)
		claimReserve = utils.Keccak256("ClaimReserve(address,address,uint256)")

		//CreateAirdrop(address indexed creator, uint256 indexed airdropId, address token, uint256 amount, uint256 startTime, uint256 endTime)
		createAirdrop = utils.Keccak256("CreateAirdrop(address,uint256,address,uint256,uint256,uint256)")
		//SettleAirdrop(uint256 indexed airdropId, uint256 amount, bytes32 merkleTreeRoot)
		settleAirdrop = utils.Keccak256("SettleAirdrop(uint256,uint256,bytes32)")
		//Claimed(uint256 indexed airdropId, uint256 index, address account, uint256 amount)
		claimed = utils.Keccak256("Claimed(uint256,uint256,address,uint256)")

		//event CreateProposal(uint indexed id, address indexed from, address indexed to, uint amount, uint startTime, uint endTime, address daoToken);
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
	case ownershipTransferred:
		event = EvOwnershipTransferred
		break
	case createERC20:
		event = EvCreateERC20
		break
	case transfer:
		event = EvTransfer
		break
	case claimReserve:
		event = EvClaimReserve
		break
	case createAirdrop:
		event = EvCreateAirdrop
		break
	case settleAirdrop:
		event = EvSettleAirdrop
		break
	case claimed:
		event = EvClaimed
		break
	default:
		event = "Undefined"
		break
	}
	return event
}

func GetAnkrArchive(chainId int) string {
	if chainId == 80001 {
		return "https://rpc.ankr.com/polygon_mumbai"
	}
	if chainId == 137 {
		return "https://rpc.ankr.com/polygon"
	}
	if chainId == 5 {
		return "https://rpc.ankr.com/eth_goerli"
	}
	return ""
}
