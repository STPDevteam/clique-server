CREATE SCHEMA IF NOT EXISTS `stp_dao_v2_pre` DEFAULT CHARACTER SET utf8 ;
USE `stp_dao_v2_pre`;

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
  PRIMARY KEY (`id`),
  INDEX `dao_address` (`dao_address` ASC),
  INDEX `dao_name` (`dao_name` ASC),
  INDEX `token_address` (`token_address` ASC),
  INDEX `creator` (`creator` ASC)
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
PRIMARY KEY (`id`)
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
	PRIMARY KEY (`id`),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_dao_address` (`dao_address` ASC),
    INDEX `index_proposalId` (`proposal_id` ASC)
);

CREATE TABLE `tb_proposal_v1` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `chain_id` INT NOT NULL,
    `dao_address` VARCHAR(128) NOT NULL,
    `dao_address_v1` VARCHAR(128) NOT NULL,
    `voting_v1` VARCHAR(128) NOT NULL,
    `start_id_v1` INT NOT NULL,#new data must -1
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
    PRIMARY KEY (`id`),
    INDEX `index_id` (`id` ASC),
    INDEX `index_airdropId` (`activity_id` ASC),
    INDEX `index_types` (`types` ASC),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_dao_address` (`dao_address` ASC),
    INDEX `index_start_time` (`start_time` ASC),
    INDEX `index_end_time` (`end_time` ASC)
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
INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
('CreateDao','0xa2d34aA709De897Ef62ee08274EC6e2c451a1CdC',34248089,'0x',137),
('CreateERC20','0xa2d34aA709De897Ef62ee08274EC6e2c451a1CdC',34248089,'0x',137),
('ClaimReserve','0xa2d34aA709De897Ef62ee08274EC6e2c451a1CdC',34248089,'0x',137),
('CreateAirdrop','0x9fBa77AA2957b2C47c0B80e14fdf7e7d28eDd127',34248107,'0x',137),
('SettleAirdrop','0x9fBa77AA2957b2C47c0B80e14fdf7e7d28eDd127',34248107,'0x',137),
('Claimed','0x9fBa77AA2957b2C47c0B80e14fdf7e7d28eDd127',34248107,'0x',137),
('CreateDao','0xD7a52a2Fe72A588351600Fa2feDD6132381f065d',15731643,'0x',1),
('CreateERC20','0xD7a52a2Fe72A588351600Fa2feDD6132381f065d',15731643,'0x',1),
('ClaimReserve','0xD7a52a2Fe72A588351600Fa2feDD6132381f065d',15731643,'0x',1),
('CreateAirdrop','0x1EFB2Cb5015FDd13120dF72BB152c8Ec91bCD68e',15731651,'0x',1),
('SettleAirdrop','0x1EFB2Cb5015FDd13120dF72BB152c8Ec91bCD68e',15731651,'0x',1),
('Claimed','0x1EFB2Cb5015FDd13120dF72BB152c8Ec91bCD68e',15731651,'0x',1);

INSERT INTO tb_category (category_name) VALUES ('Social'),('Protocol'),('NFT'),('Metaverse'),('Gaming'),('Dapp'),('Other');

#example
INSERT into tb_proposal_v1 (chain_id,dao_address,dao_address_v1,voting_v1,start_id_v1) VALUES
(5,'0xb6dbd00a199b3a616be3d38c621b337f48a065ce','0xbc61E252c79D76D9Eb23DAE0E524E80dBA6E54B4','0x6ada02cb261f864646a6fc2466a9350336fda5ad',-1),
(5,'0xf8c3b39b2533cb853620c5ccf580ad5cb2f744cd','0x53760E38B28d6882Ccf21151417Bc942E2300D00','0xef6d5b23a69b622851cca5bc2202e257021e4f7d',-1),
(5,'0xb61d2ab83f9c976bb28f8343e304d733b38832d0','0x9a151fAaca125f344E30BE6c9deF867a53a1e824','0x4ea954b5523226c671b767c3b8dfb05df8ae1561',-1),
(5,'0xc01123105f8478b56cf0a2ee67fac13d9f58e65d','0x31e7B9aF1643e96437d9DC49d3c546620A063FEC','0x29f3f68ffeff164e2d06558fb5760e1429073bd0',-1);