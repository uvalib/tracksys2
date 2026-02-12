START TRANSACTION;

CREATE TABLE IF NOT EXISTS `message_recipients` (
   `id` bigint NOT NULL AUTO_INCREMENT,
   `message_id` bigint DEFAULT NULL,
   `staff_id` bigint DEFAULT NULL,
   `read` tinyint(1) DEFAULT '0',
   `deleted` tinyint(1) DEFAULT '0',
   `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `inde_message_recipients_on_staff_id` (`staff_id`),
  KEY `index_message_recipients_on_message_id` (`message_id`),
  CONSTRAINT `message_recipients_message_id_fk` FOREIGN KEY (`message_id`) REFERENCES `messages` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

 insert into message_recipients (`message_id`, `staff_id`, `read`, `deleted`, `deleted_at`)
   select `id`,`to_id`, `read`, `deleted`, `deleted_at` from messages;

ALTER table messages
   DROP COLUMN `to_id`,
   DROP COLUMN `read`,
   DROP COLUMN `deleted`,
   DROP COLUMN `deleted_at`;

COMMIT;