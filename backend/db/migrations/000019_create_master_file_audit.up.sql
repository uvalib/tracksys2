CREATE TABLE IF NOT EXISTS `master_file_audit` (
  `id` int NOT NULL AUTO_INCREMENT,
  `master_file_id` int NOT NULL,
  `audited_at` datetime DEFAULT NULL,
  `archive_exists` tinyint(1) DEFAULT '1',
  `checksum_match` tinyint(1) DEFAULT '1',
  `audit_checksum` varchar(40) DEFAULT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `audit_master_file_id_fk` FOREIGN KEY (`master_file_id`) REFERENCES `master_files` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;