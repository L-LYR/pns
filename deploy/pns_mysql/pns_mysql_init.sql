-- TODO: add all table definition SQL
create database if not exists `pns` default character set utf8mb4 collate utf8mb4_unicode_ci;
use pns;
-- create table if not exists `target`(
--     `device_id` varchar(256) not null comment 'device id',
--     `os` char(8) default '' comment 'device os specification',
--     `brand` varchar(32) default '' comment 'device brand',
--     `model` varchar(128) default '' comment 'device model',
--     `tz_name` varchar(32) default '' comment 'time zone name',
--     `app_id` smallint not null comment 'app id',
--     `app_version` varchar(64) default '' comment 'app version string',
--     `push_sdk_version` varchar(64) default '' comment 'push sdk version string',
--     `create_time` DATETIME not null comment 'create time of target',
--     unique key `uniq_device_id_app_id` (`device_id`, `app_id`)
-- ) engine = InnoDB charset = utf8mb4 collate utf8mb4_unicode_ci;
create table if not exists `app_pusher_config`(
    `appId` int not null comment 'app id',
    `pusherId` tinyint not null comment 'pusher id',
    `config` JSON not null comment 'app pusher config',
    unique key `uniq_app_id_pusher_id` (`appId`, `pusherId`)
) engine = InnoDB charset = utf8mb4 collate utf8mb4_unicode_ci;
create table if not exists `app_config`(
    `id` int unique not null comment 'app id',
    `name` varchar(256) unique not null comment 'app name'
) engine = InnoDB charset = utf8mb4 collate utf8mb4_unicode_ci;
create table if not exists `biz_rule`(
    `name` varchar(256) unique not null comment 'rule name',
    `description` varchar(1024) comment 'rule description',
    `salience` int not null comment 'rule salience',
    `content` varchar(4096) not null comment 'rule content',
    `status` tinyint not null comment 'rule status, 0 for disable, 1 for enable'
) engine = InnoDB charset = utf8mb4 collate utf8mb4_unicode_ci;
create table if not exists `message_template`(
    `appId` int not null comment 'app id',
    `templateId` bigint unique not null comment 'template id',
    `template` JSON not null comment 'message template'
) engine = InnoDB charset = utf8mb4 collate utf8mb4_unicode_ci;
-- Test Data
insert into `app_pusher_config`
values(
        12345,
        1,
        JSON_OBJECT(
            "pusherKey",
            "test_app_name",
            "pusherSecret",
            "test_app_name",
            "receiverKey",
            "test_app_name",
            "receiverSecret",
            "test_app_name"
        )
    );
insert into `app_config`
values(12345, "test_app_name");