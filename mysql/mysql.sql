CREATE SCHEMA IF NOT EXISTS `stp_dao_v2` DEFAULT CHARACTER SET utf8 ;
USE `stp_dao_v2`;

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
    INDEX `index_chain_id` (`chain_id` ASC));

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
PRIMARY KEY (`id`)
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
    `title` VARCHAR(300) NOT NULL,
	`proposer` VARCHAR(128) NOT NULL,
	`start_time` INT NOT NULL,
	`end_time` INT NOT NULL,
	PRIMARY KEY (`id`),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_dao_address` (`dao_address` ASC),
    INDEX `index_proposalId` (`proposal_id` ASC)
);

CREATE TABLE `tb_airdrop_address` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `title` VARCHAR(500) NOT NULL,
	`content` TEXT NOT NULL,
	PRIMARY KEY (`id`),
    INDEX `index_id` (`id` ASC)
)AUTO_INCREMENT 1000;

CREATE TABLE `tb_activity` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `types` VARCHAR(30) NOT NULL,
    `chain_id` INT NOT NULL,
    `dao_address` VARCHAR(128) NOT NULL,
    `creator` VARCHAR(128) NOT NULL,
    `activity_id` INT NOT NULL,
    `token_address` VARCHAR(128) NOT NULL,
    `amount` DECIMAL(65,0) UNSIGNED NOT NULL,
    `merkle_root` VARCHAR(128) NOT NULL,
    `start_time` INT NOT NULL,
    `end_time` INT NOT NULL,
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
    `dao_address` VARCHAR(128) NOT NULL,
    `types` VARCHAR(30) NOT NULL,
    `activity_id` INT NOT NULL,
    `dao_logo` VARCHAR(500) NOT NULL,
    `dao_name` VARCHAR(30) NOT NULL,
    `activity_name` VARCHAR(300) NOT NULL,
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

INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
('CreateDao','0x0ac366978B0560bf12824d8A3a4B34B2C87ab385',28246055,'0x',80001),
('CreateERC20','0x0ac366978B0560bf12824d8A3a4B34B2C87ab385',28246055,'0x',80001),
('CreateDao','0x8E4BF15cc6fC3901aed9fFD07fb6A1211e3593ef',7642600,'0x',5),
('CreateERC20','0x8E4BF15cc6fC3901aed9fFD07fb6A1211e3593ef',7642600,'0x',5);

INSERT INTO tb_category (category_name) VALUES ('Social'),('Protocol'),('NFT'),('Metaverse'),('Gaming'),('Dapp'),('Other');