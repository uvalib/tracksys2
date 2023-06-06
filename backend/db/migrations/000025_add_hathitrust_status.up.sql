CREATE TABLE IF NOT EXISTS `hathitrust_statuses` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `metadata_id` bigint DEFAULT NULL,
  `metadata_submitted_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `metadata_status` varchar(255) DEFAULT NULL,
  `package_submitted_at` datetime DEFAULT NULL,
  `package_status` varchar(255) DEFAULT NULL,
  `notes` text,
  `finished_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_hathitrust_statuses_on_metadata_id` (`metadata_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;