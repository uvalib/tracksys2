CREATE TABLE IF NOT EXISTS `metadata_versions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `metadata_id` bigint NOT NULL,
  `staff_member_id` bigint NOT NULL,
  `desc_metadata` text,
  `version_tag` varchar(40) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `comment` text,
  PRIMARY KEY (`id`),
  KEY `index_metadata_versions_on_metadata_id` (`metadata_id`),
  KEY `index_metadata_versions_on_staff_member_id` (`staff_member_id`),
  KEY `index_metadata_versions_on_version_tag` (`version_tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;