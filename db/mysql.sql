CREATE SCHEMA IF NOT EXISTS `stp_dao_v2` DEFAULT CHARACTER SET utf8mb4 ;
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
INSERT INTO tb_category (category_name) VALUES ('Social'),('Protocol'),('NFT'),('Metaverse'),('Gaming'),('Other');

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
`account_level` VARCHAR(128),   #superAdmin;admin
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
    `push_switch` INT NOT NULL DEFAULT 0,
    `fans_num` INT NOT NULL DEFAULT 0,
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
    `dao_name` VARCHAR(200) NOT NULL,
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

CREATE TABLE `tb_swap` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `chain_id` INT NOT NULL,
    `title` VARCHAR(200) NOT NULL,
    `creator` VARCHAR(44) NOT NULL,
    `sale_way` VARCHAR(30) NOT NULL COMMENT 'general;discount',
    `sale_token` VARCHAR(44) NOT NULL,
    `sale_token_img` VARCHAR(200) NOT NULL,
    `sale_amount` DECIMAL(65,0) UNSIGNED NOT NULL,
    `sale_price` DECIMAL(65,0) UNSIGNED NOT NULL,
    `original_discount` VARCHAR(44) NOT NULL,
    `receive_token` VARCHAR(44) NOT NULL,
    `receive_token_img` VARCHAR(200) NOT NULL,
    `limit_min` DECIMAL(65,0) UNSIGNED NOT NULL,
    `limit_max` DECIMAL(65,0) UNSIGNED NOT NULL,
    `start_time` INT NOT NULL,
    `end_time` INT NOT NULL,
    `white_list` TEXT NOT NULL,
    `about` TEXT NOT NULL,
    `sold_amount` DECIMAL(65,0) UNSIGNED NOT NULL DEFAULT '0',
    `status` VARCHAR(30) NOT NULL DEFAULT 'pending' COMMENT 'pending;soon;normal;ended;cancel',
    PRIMARY KEY (`id`),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_creator` (`creator` ASC)
);

CREATE TABLE `tb_swap_token` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `chain_id` INT NOT NULL,
    `token_address` VARCHAR(44) NOT NULL,
    `token_name` VARCHAR(40) NOT NULL DEFAULT '',
    `symbol` VARCHAR(40) NOT NULL DEFAULT '',
    `decimals` INT NOT NULL,
    `coin_ids` VARCHAR(30) NOT NULL DEFAULT '',
    `price` FLOAT NOT NULL DEFAULT -1,
    `img` VARCHAR(200) NOT NULL DEFAULT '',
    `url_coingecko` VARCHAR(200) NOT NULL DEFAULT '',
    `url_coinmarketcap` VARCHAR(200) NOT NULL DEFAULT '',
    `isSync` BOOL NOT NULL DEFAULT false,
    PRIMARY KEY (`id`),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_token_address` (`token_address` ASC)
);

# token list main
INSERT INTO tb_swap_token (chain_id,token_address,token_name,symbol,decimals,coin_ids,price,img,url_coingecko,url_coinmarketcap,isSync) VALUES
(1,'0xde7d85157d9714eadf595045cc12ca4a5f3e2adb','STP','STPT',18,'stp-network',0,'https://s2.coinmarketcap.com/static/img/coins/128x128/4006.png','https://www.coingecko.com/en/coins/stp-network','https://coinmarketcap.com/currencies/standard-tokenization-protocol',1),
(1,'0xb49fa25978abf9a248b8212ab4b87277682301c0','RAI Finance','SOFI',18,'rai-finance',0,'https://s2.coinmarketcap.com/static/img/coins/128x128/16552.png','https://www.coingecko.com/en/coins/rai-finance','https://coinmarketcap.com/currencies/rai-finance-sofi',1),
(1,'0x62959c699a52ec647622c91e79ce73344e4099f5','DeFine','DFA',18,'define',0,'https://s2.coinmarketcap.com/static/img/coins/128x128/11150.png','https://www.coingecko.com/en/coins/define','https://coinmarketcap.com/currencies/define',1),
(1,'0x0000000000000000000000000000000000000000','ETH','ETH',18,'ethereum',0,'https://etherscan.io/images/svg/brands/ethereum-original.svg','https://www.coingecko.com/en/coins/ethereum','https://coinmarketcap.com/currencies/ethereum',1);

# token list test
# INSERT INTO tb_swap_token (chain_id,token_address,token_name,symbol,decimals,coin_ids,price,img) VALUES
# (11155111,'0x41526D8dE5ae045aCb88Eb0EedA752874B222ccD','18spt','18spt',18,'18spt',0.1,''),
# (11155111,'0x0090847C22856a346C6069B8d1ed08A4A1D18241','18RAI','18RAI',18,'18RAI',0.001,''),
# (11155111,'0x5c58eC0b4A18aFB85f9D6B02FE3e6454f988436E','6USDT','6USDT',6,'6USDT',0.1,''),
# (1,'0xde7d85157d9714eadf595045cc12ca4a5f3e2adb','STP','stpt',18,'',-1,''),
# (1,'0x006bea43baa3f7a6f765f14f10a1a1b08334ef45','Stox','stx',18,'',-1,''),
# (5,'0x3c0837064c3a440fe44c9002c743dcab94e16454','A','A',18,'a',1.5,''),
# (5,'0x2358fbd8a8e0470b593328503c0f9666540339a1','B','B',18,'b',0.5,''),
# (5,'0xe8a67c44933b8750204ca4ddd2307aab0547310d','C','C',18,'c',3,''),
# (5,'0x57F013F27360E62efc1904D8c4f4021648ABa7a9','D','D',6,'d',10,''),
# (5,'0x53C0475aa628D9C8C5724A2eb8B5Fd81c32a9267','E','E',18,'e',5,'');

CREATE TABLE `tb_swap_transaction` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `sale_id` INT NOT NULL,
    `buyer` VARCHAR(44) NOT NULL,
    `buy_amount` DECIMAL(65,0) UNSIGNED NOT NULL,
    `pay_amount` DECIMAL(65,0) UNSIGNED NOT NULL,
    `time` INT NOT NULL,
    `chain_id` INT NOT NULL,
    `buy_token` VARCHAR(44) NOT NULL,
    `pay_token` VARCHAR(44) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `index_sale_id` (`sale_id` ASC)
);

CREATE TABLE `sysconfig` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `cfg_name` VARCHAR(128) NULL DEFAULT '',#cfg_swap_creator_white_list
    `cfg_val` TEXT NULL,
    `cfg_type` VARCHAR(128) NULL DEFAULT '',
    `cfg_comment` VARCHAR(128) NULL DEFAULT '',
    `cfg_is_enabled` BOOL NOT NULL DEFAULT false,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `cfg_name_UNIQUE` (`cfg_name` ASC),
    INDEX `cfg_is_enabled` (`cfg_is_enabled` ASC)
);
INSERT INTO sysconfig (cfg_name,cfg_val,cfg_type,cfg_comment,cfg_is_enabled) VALUES
('cfg_swap_creator_white_list','0x5aEFAA34EaDaC483ea542077D30505eF2472cfe3','','',1);

CREATE TABLE `tb_jobs_apply` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `chain_id` INT NOT NULL,
    `dao_address` VARCHAR(44) NOT NULL,
    `account` VARCHAR(44) NOT NULL,
    `apply_role` VARCHAR(30) NOT NULL DEFAULT 'C_member' COMMENT 'B_admin;C_member',
    `message` TEXT NOT NULL,
    `status` VARCHAR(30) NOT NULL DEFAULT 'inApplication' COMMENT 'inApplication;agree;reject',
    PRIMARY KEY (`id`),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_dao_address` (`dao_address` ASC),
    INDEX `index_account` (`account` ASC)
);

CREATE TABLE `tb_jobs` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `chain_id` INT NOT NULL,
    `dao_address` VARCHAR(44) NOT NULL,
    `account` VARCHAR(44) NOT NULL,
    `job` VARCHAR(30) NOT NULL DEFAULT 'C_member' COMMENT 'A_superAdmin;B_admin;C_member',
    PRIMARY KEY (`id`),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_dao_address` (`dao_address` ASC),
    INDEX `index_account` (`account` ASC),
    INDEX `index_job` (`job` ASC)
);

CREATE TABLE `tb_team_spaces` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `chain_id` INT NOT NULL,
    `dao_address` VARCHAR(44) NOT NULL,
    `creator` VARCHAR(44) NOT NULL,
    `title` VARCHAR(20) NOT NULL,
    `url` VARCHAR(255) NOT NULL DEFAULT '',
    `last_edit_time` INT NOT NULL DEFAULT 0,
    `last_edit_by` VARCHAR(44) NOT NULL DEFAULT '',
    `access` VARCHAR(20) NOT NULL DEFAULT 'public' COMMENT 'public;private',
    `is_trash` BOOL NOT NULL DEFAULT false,
    PRIMARY KEY (`id`),
    INDEX `index_chain_id` (`chain_id` ASC),
    INDEX `index_dao_address` (`dao_address` ASC),
    INDEX `index_access` (`access` ASC),
    INDEX `index_is_trash` (`is_trash` ASC)
);

CREATE TABLE `tb_task` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `spaces_id` INT NOT NULL,
    `task_name` VARCHAR(255) NULL DEFAULT '',
    `content` TEXT NULL,
    `deadline` INT NOT NULL DEFAULT 0,
    `priority` VARCHAR(20) NULL DEFAULT '' COMMENT 'A_low;B_medium;C_high',
    `assign_account` VARCHAR(128) NULL DEFAULT '',
    `proposal_id` INT NOT NULL DEFAULT 0,
    `reward` DECIMAL(65,0) UNSIGNED NOT NULL DEFAULT '0',
    `status` VARCHAR(50) NULL DEFAULT '' COMMENT 'A_notStarted;B_inProgress;C_done;D_notStatus',
    `weight` FLOAT NOT NULL DEFAULT 0,
    `is_trash` BOOL NOT NULL DEFAULT false,
    PRIMARY KEY (`id`),
    INDEX `index_spaces_id` (`spaces_id` ASC),
    INDEX `index_status` (`status` ASC),
    INDEX `index_is_trash` (`is_trash` ASC)
);

# CREATE TABLE `tb_task_types` (
#     `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
#     `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
#     `update_time` TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
#     `task_id` INT NOT NULL,
#     `task_name` VARCHAR(255) NULL DEFAULT '',
#     PRIMARY KEY (`id`),
#     INDEX `index_task_id` (`task_id` ASC)
# );

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
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreateDao','0xD27879D1C09c4ded7d7860f22835De13cAA5885f',140075,'0x',1442),
# ('CreateERC20','0xD27879D1C09c4ded7d7860f22835De13cAA5885f',140075,'0x',1442),
# ('ClaimReserve','0xD27879D1C09c4ded7d7860f22835De13cAA5885f',140075,'0x',1442),
# ('CreateAirdrop','0x6d6aFa2C67BE77d440f5cabce62a9AB093B6085A',140075,'0x',1442),
# ('SettleAirdrop','0x6d6aFa2C67BE77d440f5cabce62a9AB093B6085A',140075,'0x',1442),
# ('Claimed','0x6d6aFa2C67BE77d440f5cabce62a9AB093B6085A',140075,'0x',1442);

#Polygon_zkEVM test proxy
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreateDao','0x837614a67877FE0C011d78740febA5b9b3f3B603',140325,'0x',1442),
# ('CreateERC20','0x837614a67877FE0C011d78740febA5b9b3f3B603',140325,'0x',1442),
# ('ClaimReserve','0x837614a67877FE0C011d78740febA5b9b3f3B603',140325,'0x',1442),
# ('CreateAirdrop','0x41526D8dE5ae045aCb88Eb0EedA752874B222ccD',140325,'0x',1442),
# ('SettleAirdrop','0x41526D8dE5ae045aCb88Eb0EedA752874B222ccD',140325,'0x',1442),
# ('Claimed','0x41526D8dE5ae045aCb88Eb0EedA752874B222ccD',140325,'0x',1442);

#Base dev
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreateDao','0x25B084FC1de433D2EA72d8F0E7949f4ea040a69f',1500223,'0x',84531),
# ('CreateERC20','0x25B084FC1de433D2EA72d8F0E7949f4ea040a69f',1500223,'0x',84531),
# ('ClaimReserve','0x25B084FC1de433D2EA72d8F0E7949f4ea040a69f',1500223,'0x',84531),
# ('CreateAirdrop','0xE663f23F7326C5fdc884613FC53bC94c65F6C856',1500223,'0x',84531),
# ('SettleAirdrop','0xE663f23F7326C5fdc884613FC53bC94c65F6C856',1500223,'0x',84531),
# ('Claimed','0xE663f23F7326C5fdc884613FC53bC94c65F6C856',1500223,'0x',84531);

#Zetachain dev
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreateDao','0x626f936D28D758c9566d3EBC3A79491C23EB1015',2307566,'0x',7001),
# ('CreateERC20','0x626f936D28D758c9566d3EBC3A79491C23EB1015',2307566,'0x',7001),
# ('ClaimReserve','0x626f936D28D758c9566d3EBC3A79491C23EB1015',2307566,'0x',7001),
# ('CreateAirdrop','0xA7eFe998463f65A49080c848510698158C64500d',2307572,'0x',7001),
# ('SettleAirdrop','0xA7eFe998463f65A49080c848510698158C64500d',2307572,'0x',7001),
# ('Claimed','0xA7eFe998463f65A49080c848510698158C64500d',2307572,'0x',7001);

#Zetachain test
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreateDao','0xe3F5935fD2265fDC8EC1091e6A5269a4515B9989',2393902,'0x',7001),
# ('CreateERC20','0xe3F5935fD2265fDC8EC1091e6A5269a4515B9989',2393902,'0x',7001),
# ('ClaimReserve','0xe3F5935fD2265fDC8EC1091e6A5269a4515B9989',2393902,'0x',7001),
# ('CreateAirdrop','0xFB00facc857A05BdC216B791404d47432C612374',2393909,'0x',7001),
# ('SettleAirdrop','0xFB00facc857A05BdC216B791404d47432C612374',2393909,'0x',7001),
# ('Claimed','0xFB00facc857A05BdC216B791404d47432C612374',2393909,'0x',7001);

# goeril dev
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreatedSale','0x626f936D28D758c9566d3EBC3A79491C23EB1015',8669039,'0x',5),
# ('Purchased','0x626f936D28D758c9566d3EBC3A79491C23EB1015',8669039,'0x',5),
# ('CancelSale','0x626f936D28D758c9566d3EBC3A79491C23EB1015',8669039,'0x',5);

# sep dev
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreatedSale','0x8c4591ca2EaeC3698200C76d242782E1aC286c1E',3128425,'0x',11155111),
# ('Purchased','0x8c4591ca2EaeC3698200C76d242782E1aC286c1E',3128425,'0x',11155111),
# ('CancelSale','0x8c4591ca2EaeC3698200C76d242782E1aC286c1E',3128425,'0x',11155111);

# main pre
# INSERT INTO scan_task (event_type,address,last_block_number,rest_parameter,chain_id) VALUES
# ('CreatedSale','0xf161dF89C31c63f3a8DC60cAcceFC78FD53f1AFA',16972769,'0x',1),
# ('Purchased','0xf161dF89C31c63f3a8DC60cAcceFC78FD53f1AFA',16972769,'0x',1),
# ('CancelSale','0xf161dF89C31c63f3a8DC60cAcceFC78FD53f1AFA',16972769,'0x',1);