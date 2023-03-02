CREATE SCHEMA IF NOT EXISTS `stp_dao_v2_pre_bsc` DEFAULT CHARACTER SET utf8mb4 ;
USE `stp_dao_v2_pre_bsc`;

SET GLOBAL TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
SET SESSION TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;

CREATE TABLE `event_historical_data` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `message_sender` VARCHAR(44) NOT NULL,
    `event_type` VARCHAR(45) NOT NULL,
    `address` VARCHAR(44) NOT NULL,
    `topic0` VARCHAR(66) NOT NULL,
    `topic1` VARCHAR(66) NOT NULL,
    `topic2` VARCHAR(66) NOT NULL,
    `topic3` VARCHAR(66) NOT NULL,
    `data` VARCHAR(514) NOT NULL,
    `block_number` VARCHAR(66) NOT NULL,
    `time_stamp` VARCHAR(66) NOT NULL,
    `gas_price` VARCHAR(66) NOT NULL,
    `gas_used` VARCHAR(66) NOT NULL,
    `log_index` VARCHAR(66) NOT NULL,
    `transaction_hash` VARCHAR(66) NOT NULL,
    `transaction_index` VARCHAR(66) NOT NULL,
    `chain_id` INT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `index_message_sender` (`message_sender` ASC),
    INDEX `index_event_type` (`event_type` ASC),
    INDEX `index_address` (`address` ASC),
    INDEX `index_topic0` (`topic0` ASC),
    INDEX `index_topic1` (`topic1` ASC),
    INDEX `index_topic2` (`topic2` ASC),
    INDEX `index_topic3` (`topic3` ASC),
    INDEX `index_block_number` (`block_number` ASC),
    INDEX `index_chain_id` (`chain_id` ASC),
    UNIQUE INDEX `unique_index_chain_id_log_index_transaction_hash` (`chain_id` ASC, `log_index` ASC, `transaction_hash` ASC));

CREATE TABLE `scan_task` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `event_type` VARCHAR(45) NOT NULL,
    `address` VARCHAR(44) NOT NULL,
    `last_block_number` INT UNSIGNED NOT NULL,
    `rest_parameter` VARCHAR(514) NOT NULL,
    `chain_id` INT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `index_address` (`address` ASC),
    INDEX `index_event_type` (`event_type` ASC),
    INDEX `index_chain_id` (`chain_id` ASC));


CREATE TABLE `tb_nonce`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `chain_id` INT NOT NULL,
    `account` VARCHAR(128) NOT NULL,
    `nonce` INT NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `chain_id` (`chain_id` ASC),
    INDEX `account` (`account` ASC),
    UNIQUE INDEX `unique_index_chain_id_account` (`chain_id` ASC, `account` ASC)
);

CREATE TABLE `tb_dao` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,	
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,	
  `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,	
    `dao_logo` VARCHAR(500) NOT NULL,
    `dao_name` VARCHAR(30) NOT NULL,
    `dao_address`  VARCHAR(128) NOT NULL,
    `creator` VARCHAR(128) NOT NULL,
    `handle` VARCHAR(30) NOT NULL,
    `description` VARCHAR(500) NOT NULL,
    `chain_id` INT NOT NULL,
    `token_chain_id` INT NOT NULL,
    `token_address` VARCHAR(128) NOT NULL,
    `proposal_threshold` DECIMAL(65,0) UNSIGNED NOT NULL,
    `voting_quorum` DECIMAL(65,0) UNSIGNED NOT NULL,
    `voting_period` INT NOT NULL,
    `voting_type` VARCHAR(128) NOT NULL,
	`twitter` VARCHAR(256),
	`github` VARCHAR(256),
	`discord` VARCHAR(256),
	`website` VARCHAR(256),
	`update_bool` bool NOT NULL,
	`weight` INT,
	`approve` bool NOT NULL,
	`deprecated` bool NOT NULL DEFAULT false,
	`members` INT NOT NULL DEFAULT 0,
	`total_proposals` INT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `index_dao_address` (`dao_address` ASC),
  INDEX `index_dao_name` (`dao_name` ASC),
  INDEX `index_token_address` (`token_address` ASC),
  INDEX `index_creator` (`creator` ASC),
  INDEX `index_deprecated` (`deprecated` ASC)
);

CREATE TABLE `tb_category`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `category_name` VARCHAR(30) NOT NULL UNIQUE,
    PRIMARY KEY (`id`)
);

CREATE TABLE `tb_dao_category`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `dao_id` INT UNSIGNED NOT NULL,
    `category_id` INT UNSIGNED NOT NULL,
    FOREIGN KEY (dao_id) REFERENCES tb_dao(id) ON UPDATE CASCADE,
    FOREIGN KEY (category_id) REFERENCES tb_category(id) ON UPDATE CASCADE,
    PRIMARY KEY (`id`)
);

CREATE TABLE `tb_member` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `dao_address` VARCHAR(128) NOT NULL,
  `chain_id` INT NOT NULL,
  `account` VARCHAR(128) NOT NULL,
  `join_switch` INT NOT NULL,
  INDEX `dao_address` (`dao_address` ASC),
  INDEX `account` (`account` ASC),
  INDEX `chain_id` (`chain_id` ASC),
  INDEX `join_switch` (`join_switch` ASC),
  UNIQUE INDEX `unique_index_chain_id_dao_address_account` (`chain_id` ASC, `dao_address` ASC, `account` ASC),
  PRIMARY KEY (`id`)
);

CREATE TABLE `tb_admin` (
`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
`create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
`update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
`dao_address` VARCHAR(128) NOT NULL,
`chain_id` INT NOT NULL,
`account` VARCHAR(128) NOT NULL,
`account_level` VARCHAR(128),
INDEX `dao_address` (`dao_address` ASC),
INDEX `account` (`account` ASC),
INDEX `chain_id` (`chain_id` ASC),
UNIQUE INDEX `unique_index_chain_id_dao_address_account_account_level` (`chain_id` ASC, `dao_address` ASC, `account` ASC, account_level ASC),
PRIMARY KEY (`id`)
);

CREATE TABLE `tb_holder_data` (
 `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
 `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 `token_address` VARCHAR(44) NOT NULL,
 `holder_address` VARCHAR(66) NOT NULL,
 `balance` DECIMAL(65,0) UNSIGNED NOT NULL,
 `chain_id` INT UNSIGNED NOT NULL,
 PRIMARY KEY (`id`),
 INDEX `index_token_address` (`token_address` ASC),
 INDEX `index_holder_address` (`holder_address` ASC),
 UNIQUE INDEX `unique_index_chain_id_token_address_holder_address` (`chain_id` ASC, `token_address` ASC, `holder_address` ASC),
 INDEX `index_chain_id` (`chain_id` ASC),
 INDEX `index_balance` (`balance` DESC)
);

CREATE TABLE `error_info` (
`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
`create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
`update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
`title` VARCHAR(512) NOT NULL,
`content` TEXT NOT NULL,
`func` VARCHAR(128) NOT NULL,
`params` VARCHAR(512) NOT NULL,
PRIMARY KEY (`id`)
);

CREATE TABLE `tb_account` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `account` VARCHAR(128) NOT NULL UNIQUE,
    `account_logo` VARCHAR(128),
    `nickname` VARCHAR(128),
    `introduction` VARCHAR(200),
    `twitter` VARCHAR(128),
    `github` VARCHAR(128),
    `discord` VARCHAR(128),
    `email` VARCHAR(128),
    `country` VARCHAR(128),
    `youtube` VARCHAR(128),
    `opensea` VARCHAR(128),
    PRIMARY KEY (`id`),
    INDEX `index_account` (`account` ASC)
);

CREATE TABLE `tb_account_follow` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `account` VARCHAR(66) NOT NULL,
    `followed` VARCHAR(66) NOT NULL,
    `status` BOOL NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `index_account` (`account` ASC),
    INDEX `index_followed` (`followed` ASC)
);

CREATE TABLE `tb_account_record` (
    `id`          INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP    NOT NULL                             DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP    NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `creator` VARCHAR(128) NOT NULL,
    `types` VARCHAR(30) NOT NULL,
    `chain_id` INT NOT NULL,
    `address` VARCHAR(128) NOT NULL,#(dao address / token address)
    `activity_id` INT NOT NULL,#(proposal / airdrop)
    `avatar` VARCHAR(500) NOT NULL,#(dao / token)
    `dao_name` VARCHAR(30) NOT NULL,#(dao name)
    `titles` VARCHAR(500) NOT NULL,#(title of proposal / token name='' / airdrop name)
    `time` INT NOT NULL,
    `update_bool` bool NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `index_creator` (`creator` ASC),
    INDEX `index_types` (`types` ASC),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_address` (`address` ASC),
    INDEX `index_activity_id` (`activity_id` ASC),
    INDEX `index_update_bool` (`update_bool` ASC)
);

CREATE TABLE `tb_account_sign` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `chain_id` INT NOT NULL,
    `dao_address` VARCHAR(128) NOT NULL,
    `account` VARCHAR(128) NOT NULL,
    `operate` VARCHAR(30) NOT NULL,
    `signature` VARCHAR(132) NOT NULL,
    `message` VARCHAR(132) NOT NULL,
    `timestamp` INT NOT NULL,
     PRIMARY KEY (`id`),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_dao_address` (`dao_address` ASC),
    INDEX `index_account` (`account` ASC),
    INDEX `index_timestamp` (`timestamp` ASC)
);

CREATE TABLE `tb_proposal_info` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `uuid` VARCHAR(36) NOT NULL,
  `content` TEXT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `index_uuid` (`uuid` ASC)
);

CREATE TABLE `tb_tokens_img` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `token_chain_id` INT NOT NULL,
    `token_address` VARCHAR(128) NOT NULL,
    `thumb` VARCHAR(300) NOT NULL,
    `small` VARCHAR(300) NOT NULL,
    `large` VARCHAR(300) NOT NULL,
    `own_img` VARCHAR(300) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `index_chain_id` (`token_chain_id` ASC),
    INDEX `index_token_address` (`token_address` ASC),
    UNIQUE INDEX `unique_index_token_chain_id_token_address` (`token_chain_id` ASC, `token_address` ASC)
);

CREATE TABLE `tb_proposal` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `proposal_id` INT NOT NULL,
    `chain_id` INT NOT NULL,
    `dao_address` VARCHAR(128) NOT NULL,
    `title` VARCHAR(500) NOT NULL,
    `id_v1` INT NOT NULL,
    `content_v1` TEXT NOT NULL,
	`proposer` VARCHAR(128) NOT NULL,
	`start_time` INT NOT NULL,
	`end_time` INT NOT NULL,
	`version` VARCHAR(10) NOT NULL,
    `deprecated` bool NOT NULL DEFAULT false,
    `block_number` VARCHAR(66) NOT NULL DEFAULT '',
	PRIMARY KEY (`id`),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_dao_address` (`dao_address` ASC),
    INDEX `index_proposalId` (`proposal_id` ASC),
    INDEX `index_deprecated` (`deprecated` ASC),
    INDEX `index_block_number` (`block_number` ASC),
    INDEX `index_start_time` (`start_time` ASC),
    INDEX `index_end_time` (`end_time` ASC)
);

CREATE TABLE `tb_proposal_v1` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `chain_id` INT NOT NULL,
    `dao_address` VARCHAR(128) NOT NULL,
    `dao_address_v1` VARCHAR(128) NOT NULL,
    `voting_v1` VARCHAR(128) NOT NULL,
    `start_id_v1` INT NOT NULL DEFAULT -1,
    PRIMARY KEY (`id`),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_dao_address` (`dao_address` ASC),
    INDEX `index_dao_address_v1` (`dao_address_v1` ASC),
    UNIQUE INDEX `index_voting_v1` (`voting_v1` ASC),
    INDEX `index_start_id_v1` (`start_id_v1` ASC)
);

CREATE TABLE `tb_airdrop` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `creator` VARCHAR(128) NOT NULL,
    `chain_id` INT NOT NULL,
    `dao_address` VARCHAR(128) NOT NULL,
    `title` VARCHAR(500) NOT NULL,
	`airdrop_address` TEXT NOT NULL,
	`description` TEXT NOT NULL,
	`collect_information` TEXT NOT NULL,
	`token_chain_id` INT NOT NULL,
	`token_address` VARCHAR(128) NOT NULL,
	`max_airdrop_amount` DECIMAL(65,0) UNSIGNED NOT NULL,
    `start_time` INT NOT NULL,
    `end_time` INT NOT NULL,
    `airdrop_start_time` INT NOT NULL,
    `airdrop_end_time` INT NOT NULL,
	PRIMARY KEY (`id`),
    INDEX `index_id` (`id` ASC),
    INDEX `index_creator` (`creator` ASC),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_dao_address` (`dao_address` ASC),
    INDEX `index_token_chain_id` (`token_chain_id` ASC),
    INDEX `index_token_address` (`token_address` ASC)
)AUTO_INCREMENT 1000;

CREATE TABLE `tb_airdrop_prepare` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `airdrop_id` INT NOT NULL,
    `root` VARCHAR(128) NOT NULL,
    `prepare_address` TEXT NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `index_airdrop_id` (`airdrop_id` ASC),
    INDEX `index_root` (`root` ASC)
);

CREATE TABLE `tb_airdrop_user_submit`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `airdrop_id` INT NOT NULL,
    `account` VARCHAR(128) NOT NULL,
    `submit_info` TEXT NOT NULL,
    `timestamp` INT NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `index_id` (`id` ASC),
    INDEX `index_airdrop_id` (`airdrop_id` ASC),
    INDEX `index_account` (`account` ASC)
);

CREATE TABLE `tb_activity` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `types` VARCHAR(30) NOT NULL,
    `chain_id` INT NOT NULL,
    `dao_address` VARCHAR(128) NOT NULL,
    `creator` VARCHAR(128) NOT NULL,
    `activity_id` INT NOT NULL,
    `token_chain_id` INT NOT NULL,
    `token_address` VARCHAR(128) NOT NULL,
    `staking_amount` DECIMAL(65,0) UNSIGNED NOT NULL,
    `airdrop_amount` DECIMAL(65,0) UNSIGNED NOT NULL,
    `merkle_root` VARCHAR(128) NOT NULL,
    `start_time` INT NOT NULL,
    `end_time` INT NOT NULL,
    `airdrop_start_time` INT NOT NULL,
    `airdrop_end_time` INT NOT NULL,
    `publish_time` INT NOT NULL,
    `price` VARCHAR(128) NOT NULL,
    `weight` INT,
    `deprecated` bool NOT NULL DEFAULT false,
    PRIMARY KEY (`id`),
    INDEX `index_id` (`id` ASC),
    INDEX `index_airdropId` (`activity_id` ASC),
    INDEX `index_types` (`types` ASC),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_dao_address` (`dao_address` ASC),
    INDEX `index_start_time` (`start_time` ASC),
    INDEX `index_end_time` (`end_time` ASC),
    INDEX `index_deprecated` (`deprecated` ASC)
);

CREATE TABLE `tb_claimed` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `chain_id` INT NOT NULL,
    `dao_address` VARCHAR(128) NOT NULL,
    `airdrop_id` INT NOT NULL,
    `index_id` INT NOT NULL,
    `account` VARCHAR(128) NOT NULL,
    `amount` DECIMAL(65,0) UNSIGNED NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `index_account` (`account` ASC),
    INDEX `index_airdrop_id` (`airdrop_id` ASC),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_dao_address` (`dao_address` ASC)
);

CREATE TABLE `tb_vote` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `chain_id` INT NOT NULL,
    `dao_address` VARCHAR(128) NOT NULL,
    `proposal_id` INT NOT NULL,
    `voter` VARCHAR(128) NOT NULL,
    `option_index` INT NOT NULL,
    `amount` DECIMAL(65,0) UNSIGNED NOT NULL,
    `nonce` INT NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_dao_address` (`dao_address` ASC),
    INDEX `index_proposal_id` (`proposal_id` ASC),
    INDEX `index_voter` (`voter` ASC)
);

CREATE TABLE `tb_notification` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `chain_id` INT NOT NULL,
    `dao_address` VARCHAR(128) NOT NULL,# if types=ReserveToken -> token address
    `types` VARCHAR(30) NOT NULL,
    `activity_id` INT NOT NULL,
    `dao_logo` VARCHAR(500) NOT NULL,
    `dao_name` VARCHAR(30) NOT NULL,
    `activity_name` VARCHAR(500) NOT NULL,
    `start_time` INT NOT NULL,
    `update_bool` bool NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `index_id` (`id` ASC),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_dao_address` (`dao_address` ASC)
);

CREATE TABLE `tb_notification_account` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `notification_id` INT NOT NULL,
    `account` VARCHAR(128) NOT NULL,
    `already_read` bool NOT NULL,
    `notification_time` INT NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `index_notification_id` (`notification_id` ASC),
    INDEX `index_account` (`account` ASC)
);

CREATE TABLE `tb_handle_lock` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `handle` VARCHAR(30) NOT NULL,
    `handle_keccak` VARCHAR(66) NOT NULL,
    `lock_block` INT UNSIGNED NOT NULL,
    `chain_id` INT NOT NULL,
    `account` VARCHAR(66) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `index_handle` (`handle` ASC),
    INDEX `index_lock_block` (`lock_block` ASC),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_account` (`account` ASC)
);

# dev
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreateDao','0x18Be998c31815d1C3d1dde881801112D9ee81532',28315453,'0x',80001),
# ('CreateERC20','0x18Be998c31815d1C3d1dde881801112D9ee81532',28315453,'0x',80001),
# ('ClaimReserve','0x18Be998c31815d1C3d1dde881801112D9ee81532',28315453,'0x',80001),
# ('CreateAirdrop','0x83D32D4618A8798112B0F6390b558761B8881348',28350844,'0x',80001),
# ('SettleAirdrop','0x83D32D4618A8798112B0F6390b558761B8881348',28350844,'0x',80001),
# ('Claimed','0x83D32D4618A8798112B0F6390b558761B8881348',28350844,'0x',80001),
# ('CreateDao','0x93e7A03239d62CC24D84A7A216E81FB2aDbC7D9b',7666402,'0x',5),
# ('CreateERC20','0x93e7A03239d62CC24D84A7A216E81FB2aDbC7D9b',7666402,'0x',5),
# ('ClaimReserve','0x93e7A03239d62CC24D84A7A216E81FB2aDbC7D9b',7666402,'0x',5),
# ('CreateAirdrop','0x4a9e8EeBd7e3d928E494A3ef43baD56838FB2Bf3',7678702,'0x',5),
# ('SettleAirdrop','0x4a9e8EeBd7e3d928E494A3ef43baD56838FB2Bf3',7678702,'0x',5),
# ('Claimed','0x4a9e8EeBd7e3d928E494A3ef43baD56838FB2Bf3',7678702,'0x',5);

# test
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreateDao','0x03919E8A7db18B89aC287ddb8ad5DE34F44E1E11',28561789,'0x',80001),
# ('CreateERC20','0x03919E8A7db18B89aC287ddb8ad5DE34F44E1E11',28561789,'0x',80001),
# ('ClaimReserve','0x03919E8A7db18B89aC287ddb8ad5DE34F44E1E11',28561789,'0x',80001),
# ('CreateAirdrop','0x4a9e8EeBd7e3d928E494A3ef43baD56838FB2Bf3',28561881,'0x',80001),
# ('SettleAirdrop','0x4a9e8EeBd7e3d928E494A3ef43baD56838FB2Bf3',28561881,'0x',80001),
# ('Claimed','0x4a9e8EeBd7e3d928E494A3ef43baD56838FB2Bf3',28561881,'0x',80001),
# ('CreateDao','0x8ab3364443396D81703444F2c29E86A2809c5eC1',7750148,'0x',5),
# ('CreateERC20','0x8ab3364443396D81703444F2c29E86A2809c5eC1',7750148,'0x',5),
# ('ClaimReserve','0x8ab3364443396D81703444F2c29E86A2809c5eC1',7750148,'0x',5),
# ('CreateAirdrop','0x8FC198E84e43474F95468300593539B8cb3bEe8f',7750200,'0x',5),
# ('SettleAirdrop','0x8FC198E84e43474F95468300593539B8cb3bEe8f',7750200,'0x',5),
# ('Claimed','0x8FC198E84e43474F95468300593539B8cb3bEe8f',7750200,'0x',5);

# pre
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreateDao','0xa2d34aA709De897Ef62ee08274EC6e2c451a1CdC',34248089,'0x',137),
# ('CreateERC20','0xa2d34aA709De897Ef62ee08274EC6e2c451a1CdC',34248089,'0x',137),
# ('ClaimReserve','0xa2d34aA709De897Ef62ee08274EC6e2c451a1CdC',34248089,'0x',137),
# ('CreateAirdrop','0x9fBa77AA2957b2C47c0B80e14fdf7e7d28eDd127',34248107,'0x',137),
# ('SettleAirdrop','0x9fBa77AA2957b2C47c0B80e14fdf7e7d28eDd127',34248107,'0x',137),
# ('Claimed','0x9fBa77AA2957b2C47c0B80e14fdf7e7d28eDd127',34248107,'0x',137),
# ('CreateDao','0xD7a52a2Fe72A588351600Fa2feDD6132381f065d',15731643,'0x',1),
# ('CreateERC20','0xD7a52a2Fe72A588351600Fa2feDD6132381f065d',15731643,'0x',1),
# ('ClaimReserve','0xD7a52a2Fe72A588351600Fa2feDD6132381f065d',15731643,'0x',1),
# ('CreateAirdrop','0x1EFB2Cb5015FDd13120dF72BB152c8Ec91bCD68e',15731651,'0x',1),
# ('SettleAirdrop','0x1EFB2Cb5015FDd13120dF72BB152c8Ec91bCD68e',15731651,'0x',1),
# ('Claimed','0x1EFB2Cb5015FDd13120dF72BB152c8Ec91bCD68e',15731651,'0x',1);

INSERT INTO tb_category (category_name) VALUES ('Social'),('Protocol'),('NFT'),('Metaverse'),('Gaming'),('Other');

#klaytn testnet dev
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreateDao','0x9fE8096d3C4D3cCb5E5537fa7761CdBb377d45bC',105807119,'0x',1001),
# ('CreateERC20','0x9fE8096d3C4D3cCb5E5537fa7761CdBb377d45bC',105807119,'0x',1001),
# ('ClaimReserve','0x9fE8096d3C4D3cCb5E5537fa7761CdBb377d45bC',105807119,'0x',1001),
# ('CreateAirdrop','0x8b09B9008f3bF5D4F48e0A251a872EB45a4372Dd',106211720,'0x',1001),
# ('SettleAirdrop','0x8b09B9008f3bF5D4F48e0A251a872EB45a4372Dd',106211720,'0x',1001),
# ('Claimed','0x8b09B9008f3bF5D4F48e0A251a872EB45a4372Dd',106211720,'0x',1001);

# klaytn mainnet pre
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreateDao','0x5bf53eAFd960AE3Cced46D2B7B1b8555334dBeF0',106593563,'0x',8217),
# ('CreateERC20','0x5bf53eAFd960AE3Cced46D2B7B1b8555334dBeF0',106593563,'0x',8217),
# ('ClaimReserve','0x5bf53eAFd960AE3Cced46D2B7B1b8555334dBeF0',106593563,'0x',8217),
# ('CreateAirdrop','0x8f606118B151A9235868DF966bf1604d24A1909B',106593575,'0x',8217),
# ('SettleAirdrop','0x8f606118B151A9235868DF966bf1604d24A1909B',106593575,'0x',8217),
# ('Claimed','0x8f606118B151A9235868DF966bf1604d24A1909B',106593575,'0x',8217);

#klaytn mainnet main
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreateDao','0x67dd666da5c03eC1Cf2faACa45064e793648ecA3',106681729,'0x',8217),
# ('CreateERC20','0x67dd666da5c03eC1Cf2faACa45064e793648ecA3',106681729,'0x',8217),
# ('ClaimReserve','0x67dd666da5c03eC1Cf2faACa45064e793648ecA3',106681729,'0x',8217),
# ('CreateAirdrop','0x3ba0116d115CC419B5cda1c2b59D402AF8D4056b',106681741,'0x',8217),
# ('SettleAirdrop','0x3ba0116d115CC419B5cda1c2b59D402AF8D4056b',106681741,'0x',8217),
# ('Claimed','0x3ba0116d115CC419B5cda1c2b59D402AF8D4056b',106681741,'0x',8217);

#BSC testnet dev
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreateDao','0xb6AE25DD2A30E4670077f0e177530bBC0921d7BA',24603051,'0x',97),
# ('CreateERC20','0xb6AE25DD2A30E4670077f0e177530bBC0921d7BA',24603051,'0x',97),
# ('ClaimReserve','0xb6AE25DD2A30E4670077f0e177530bBC0921d7BA',24603051,'0x',97),
# ('CreateAirdrop','0xae96637920430e4D05cCf11Db47cfa5cfC3224B7',24603072,'0x',97),
# ('SettleAirdrop','0xae96637920430e4D05cCf11Db47cfa5cfC3224B7',24603072,'0x',97),
# ('Claimed','0xae96637920430e4D05cCf11Db47cfa5cfC3224B7',24603072,'0x',97);

#BSC mainnet pre
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreateDao','0x5bf53eAFd960AE3Cced46D2B7B1b8555334dBeF0',23278269,'0x',56),
# ('CreateERC20','0x5bf53eAFd960AE3Cced46D2B7B1b8555334dBeF0',23278269,'0x',56),
# ('ClaimReserve','0x5bf53eAFd960AE3Cced46D2B7B1b8555334dBeF0',23278269,'0x',56),
# ('CreateAirdrop','0x2AC73343B61ec8C0301aebB39514d1cD12f9013A',23291445,'0x',56),
# ('SettleAirdrop','0x2AC73343B61ec8C0301aebB39514d1cD12f9013A',23291445,'0x',56),
# ('Claimed','0x2AC73343B61ec8C0301aebB39514d1cD12f9013A',23291445,'0x',56);

#BSC mainnet main
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreateDao','0xd56a88fDE1861297A32416e86Fc6f3347A1677bc',23319449,'0x',56),
# ('CreateERC20','0xd56a88fDE1861297A32416e86Fc6f3347A1677bc',23319449,'0x',56),
# ('ClaimReserve','0xd56a88fDE1861297A32416e86Fc6f3347A1677bc',23319449,'0x',56),
# ('CreateAirdrop','0x26E3a7841682D65e7a11e3C82067CeA0BbFC6aB4',23319449,'0x',56),
# ('SettleAirdrop','0x26E3a7841682D65e7a11e3C82067CeA0BbFC6aB4',23319449,'0x',56),
# ('Claimed','0x26E3a7841682D65e7a11e3C82067CeA0BbFC6aB4',23319449,'0x',56);

#Polygon_zkEVM dev
INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
('CreateDao','0x7D1d608630F311ebBa6bb56AA3Ac4ef8e0F52dD8',500,'0x',1442),
('CreateERC20','0x7D1d608630F311ebBa6bb56AA3Ac4ef8e0F52dD8',500,'0x',1442),
('ClaimReserve','0x7D1d608630F311ebBa6bb56AA3Ac4ef8e0F52dD8',500,'0x',1442),
('CreateAirdrop','0x210B84c4a10EA59EdF05964613161A6cbb7F3837',500,'0x',1442),
('SettleAirdrop','0x210B84c4a10EA59EdF05964613161A6cbb7F3837',500,'0x',1442),
('Claimed','0x210B84c4a10EA59EdF05964613161A6cbb7F3837',500,'0x',1442);

#Base dev
INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
('CreateDao','0x25B084FC1de433D2EA72d8F0E7949f4ea040a69f',1267171,'0x',84531),
('CreateERC20','0x25B084FC1de433D2EA72d8F0E7949f4ea040a69f',1267171,'0x',84531),
('ClaimReserve','0x25B084FC1de433D2EA72d8F0E7949f4ea040a69f',1267171,'0x',84531),
('CreateAirdrop','0xE663f23F7326C5fdc884613FC53bC94c65F6C856',1267171,'0x',84531),
('SettleAirdrop','0xE663f23F7326C5fdc884613FC53bC94c65F6C856',1267171,'0x',84531),
('Claimed','0xE663f23F7326C5fdc884613FC53bC94c65F6C856',1267171,'0x',84531);
