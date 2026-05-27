CREATE TABLE IF NOT EXISTS `ap_trust_submissions` (
   `id` int NOT NULL AUTO_INCREMENT,
   `metadata_id` bigint NOT NULL,
   `bag` varchar(255) NOT NULL,
   `requested_at` datetime NOT NULL,
   `submitted_at` datetime DEFAULT NULL,
   PRIMARY KEY (`id`),
   KEY `ap_trust_submissions_on_metadata_id` (`metadata_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;