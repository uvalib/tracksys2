CREATE TABLE IF NOT EXISTS `checkouts` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `metadata_id` bigint DEFAULT NULL,
  `checkout_at` datetime DEFAULT NULL,
  `return_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_checkouts_on_metadata_id` (`metadata_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;