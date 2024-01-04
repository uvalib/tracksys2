CREATE TABLE IF NOT EXISTS `ap_trust_statuses` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `metadata_id` bigint DEFAULT NULL,
  `etag` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `note` text,
  `object_id` varchar(255) DEFAULT NULL,
  `submitted_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `finished_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_ap_trust_statuses_on_metadata_id` (`metadata_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;