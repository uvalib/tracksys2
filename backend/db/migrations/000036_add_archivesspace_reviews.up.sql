CREATE TABLE IF NOT EXISTS `archivesspace_reviews` (
   `id` bigint NOT NULL AUTO_INCREMENT,
   `metadata_id` int DEFAULT NULL,
   `submit_staff_id` int DEFAULT NULL,
   `submitted_at` datetime DEFAULT CURRENT_TIMESTAMP,
   `review_staff_id` int DEFAULT NULL,
   `status` varchar(10) DEFAULT NULL,
   `notes` varchar(255) DEFAULT NULL,
   `published_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_archivesspace_reviews_on_metadata_id` (`metadata_id`),
  CONSTRAINT `archivesspace_reviews_metadata_id_fk` FOREIGN KEY (`metadata_id`) REFERENCES `metadata` (`id`),
  KEY `index_archivesspace_reviews_on_submit_staff_id` (`submit_staff_id`),
  CONSTRAINT `archivesspace_reviews_submit_staff_id_fk` FOREIGN KEY (`submit_staff_id`) REFERENCES `staff_members` (`id`),
  KEY `index_archivesspace_reviews_on_review_staff_id` (`review_staff_id`),
  CONSTRAINT `archivesspace_reviews_review_staff_id_fk` FOREIGN KEY (`review_staff_id`) REFERENCES `staff_members` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;