START TRANSACTION;

ALTER table messages
   ADD COLUMN `to_id` bigint DEFAULT NULL,
   ADD COLUMN `read` tinyint(1) DEFAULT '0',
   ADD COLUMN `deleted` tinyint(1) DEFAULT '0',
   ADD COLUMN `deleted_at` timestamp NULL DEFAULT NULL;


DROP TABLE iF EXISTS `message_recipients`;

COMMIT;