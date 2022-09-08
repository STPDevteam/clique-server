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
	TbNameAdmin           = "tb_admin"
	TbNameErrorInfo       = "error_info"
	TbNameProposalInfo    = "tb_proposal_info"
	TbNameTokensImg       = "tb_tokens_img"
	TbNameProposal        = "tb_proposal"
	TbNameAirdropAddress  = "tb_airdrop_address"
	TbNameActivity        = "tb_activity"
	TbNameClaimed         = "tb_claimed"
	TbNameVote            = "tb_vote"

	EvCreateDao            = "CreateDao"
	EvCreateProposal       = "CreateProposal"
	EvVote                 = "Vote"
	EvCancelProposal       = "CancelProposal"
	EvAdmin                = "Admin"
	EvSetting              = "Setting"
	EvOwnershipTransferred = "OwnershipTransferred"
	EvCreateERC20          = "CreateERC20"
	EvTransfer             = "Transfer"
	EvCreateAirdrop        = "CreateAirdrop"
	EvClaimed              = "Claimed"

	LevelSuperAdmin = "superAdmin"
	LevelAdmin      = "admin"
	LevelNoRole     = "noRole"

	TypesNameAirdrop    = "Airdrop"
	TypesNamePublicSale = "PublicSale"

	MaxValue        = 0x7fffffff
	ZeroAddress0x64 = "0x0000000000000000000000000000000000000000000000000000000000000000"
	ZeroAddress0x40 = "0x0000000000000000000000000000000000000000"

	SignMessagePrefix = "Welcome come Clique"
)

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
		//OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
		ownershipTransferred = utils.Keccak256("OwnershipTransferred(address,address)")
		//CreateAirdrop(address indexed creator, uint256 indexed airdropId, address token, uint256 amount, bytes32 merkleRoot, uint256 startTime, uint256 endTime)
		createAirdrop = utils.Keccak256("CreateAirdrop(address,uint256,address,uint256,bytes32,uint256,uint256)")
		//Claimed(uint256 indexed airdropId, uint256 index, address account, uint256 amount)
		claimed = utils.Keccak256("Claimed(uint256,uint256,address,uint256)")

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
	case ownershipTransferred:
		event = EvOwnershipTransferred
		break
	case createERC20:
		event = EvCreateERC20
		break
	case transfer:
		event = EvTransfer
		break
	case createAirdrop:
		event = EvCreateAirdrop
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
