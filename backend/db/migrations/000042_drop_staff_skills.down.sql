CREATE TABLE IF NOT EXISTS `staff_skills` (
  `id` int NOT NULL AUTO_INCREMENT,
  `staff_member_id` int DEFAULT NULL,
  `category_id` int DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_staff_skills_on_staff_member_id` (`staff_member_id`),
  KEY `index_staff_skills_on_category_id` (`category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;