package consts

const (
	TbNameEventHistorical = "event_historical_data"
	TbNameScanTask        = "scan_task"
	TbSysConfig           = "sysconfig"

	TbNameAccount             = "tb_account"
	TbNameAccountFollow       = "tb_account_follow"
	TbNameAccountRecord       = "tb_account_record"
	TbNameAccountSign         = "tb_account_sign"
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
	TbNameSwap                = "tb_swap"
	TbNameSwapToken           = "tb_swap_token"
	TbNameSwapTransaction     = "tb_swap_transaction"
	TbJobs                    = "tb_jobs"
	TbJobsApply               = "tb_jobs_apply"
	TbTask                    = "tb_task"

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
	EvCreatedSale          = "CreatedSale"
	EvPurchased            = "Purchased"
	EvCancelSale           = "CancelSale"

	LevelSuperAdmin = "superAdmin"
	LevelAdmin      = "admin"
	LevelMember     = "member"
	LevelNoRole     = "noRole"

	TypesNameAirdrop             = "Airdrop"
	TypesNamePublicSaleCreated   = "PublicSaleCreated"   // notification
	TypesNamePublicSalePurchased = "PublicSalePurchased" // notification
	TypesNamePublicSaleCanceled  = "PublicSaleCanceled"  // notification
	TypesNameNewProposal         = "NewProposal"
	TypesNameReserveToken        = "ReserveToken"

	MaxValue        = 0x7fffffff
	MaxIntUnsigned  = 4294967295
	ZeroAddress0x64 = "0x0000000000000000000000000000000000000000000000000000000000000000"
	ZeroAddress0x40 = "0x0000000000000000000000000000000000000000"

	CacheTokenHolders = "CacheTokenHolders"

	SignMessagePrefix = "Welcome come Clique"

	GoerliTestnet5      = 5
	EthMainnet1         = 1
	PolygonTestnet80001 = 80001
	PolygonMainnet137   = 137
	KlaytnTestnet1001   = 1001
	KlaytnMainnet8217   = 8217
	BSCTestnet97        = 97
	BSCMainnet56        = 56

	StatusPending = "pending"
	StatusSoon    = "soon"
	StatusNormal  = "normal"
	StatusEnded   = "ended"
	StatusCancel  = "cancel"

	Jobs_Status_InApplication = "inApplication"
	Jobs_Status_Agree         = "agree"
	Jobs_Status_Reject        = "reject"

	Jobs_A_superAdmin = "A_superAdmin"
	Jobs_B_admin      = "B_admin"
	Jobs_C_member     = "C_member"

	Task_status_A_notStarted = "A_notStarted"
	Task_status_B_inProgress = "B_inProgress"
	Task_status_C_done       = "C_done"
	Task_status_D_notStatus  = "D_notStatus"
)
