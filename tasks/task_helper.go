package tasks

import (
	"stp_dao_v2/consts"
	"stp_dao_v2/utils"
)

func EventTypes(event string) string {
	var (
		//CreateDAO(uint256 indexed handler, address indexed creator, address indexed daoAddress, uint256 chainId, address tokenAddress)
		createDao = utils.Keccak256("CreateDAO(uint256,address,address,uint256,address)")
		//CreateProposal(uint256 indexed proposalId, address indexed proposer, uint256 startTime, uint256 endTime)
		createProposal = utils.Keccak256("CreateProposal(uint256,address,uint256,uint256)")
		//CancelProposal(uint256 indexed proposalId)
		cancelProposal = utils.Keccak256("CancelProposal(uint256)")
		//Vote(uint256 indexed proposalId, address indexed voter, uint256 indexed optionIndex, uint256 amount)
		vote = utils.Keccak256("Vote(uint256,address,uint256,uint256)")
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

		//CreatedSale(uint256 indexed saleId, address indexed saleToken, address indexed receiveToken, uint256 saleAmount, uint256 pricePer, uint256 limitMin, uint256 limitMax, uint256 startTime, uint256 endTime)
		createdSale = utils.Keccak256("CreatedSale(uint256,address,address,uint256,uint256,uint256,uint256,uint256,uint256)")
		//Purchased(uint256 indexed saleId, uint256 indexed buyAmount)
		purchased = utils.Keccak256("Purchased(uint256,uint256,uint256)")
		//CancelSale(uint256 indexed saleId)
		cancelSale = utils.Keccak256("CancelSale(uint256)")
	)
	switch event {
	case createDao:
		event = consts.EvCreateDao
		break
	case createProposal:
		event = consts.EvCreateProposal
		break
	case cancelProposal:
		event = consts.EvCancelProposal
		break
	case vote:
		event = consts.EvVote
		break
	case admin:
		event = consts.EvAdmin
		break
	case setting:
		event = consts.EvSetting
		break
	case ownershipTransferred:
		event = consts.EvOwnershipTransferred
		break
	case createERC20:
		event = consts.EvCreateERC20
		break
	case transfer:
		event = consts.EvTransfer
		break
	case claimReserve:
		event = consts.EvClaimReserve
		break
	case createAirdrop:
		event = consts.EvCreateAirdrop
		break
	case settleAirdrop:
		event = consts.EvSettleAirdrop
		break
	case claimed:
		event = consts.EvClaimed
		break
	case createdSale:
		event = consts.EvCreatedSale
		break
	case purchased:
		event = consts.EvPurchased
		break
	case cancelSale:
		event = consts.EvCancelSale
		break
	default:
		event = "Undefined"
		break
	}
	return event
}
