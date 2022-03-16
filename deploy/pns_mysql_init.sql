-- TODO: add all table definition SQL
create database if not exists `pns` default character set utf8mb4 collate utf8mb4_unicode_ci;
use pns;
create table if not exists `target`(
    `device_id` varchar(256) not null comment 'device id',
    `os` char(8) default '' comment 'device os specification',
    `brand` varchar(32) default '' comment 'device brand',
    `model` varchar(128) default '' comment 'device model',
    `tz_name` varchar(32) default '' comment 'time zone name',
    `app_id` smallint not null comment 'app id',
    `app_version` varchar(64) default '' comment 'app version string',
    `push_sdk_version` varchar(64) default '' comment 'push sdk version string',
    `create_time` DATETIME not null comment 'create time of target',
    unique key `uniq_device_id_app_id` (`device_id`, `app_id`)
) engine = InnoDB charset = utf8mb4 collate utf8mb4_unicode_ci;