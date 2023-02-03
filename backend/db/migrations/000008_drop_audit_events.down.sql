CREATE TABLE IF NOT EXISTS `audit_events` (
  `id` int NOT NULL AUTO_INCREMENT,
  `staff_member_id` int DEFAULT NULL,
  `auditable_id` int DEFAULT NULL,
  `auditable_type` varchar(255) DEFAULT NULL,
  `event` int DEFAULT NULL,
  `details` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_audit_events_on_staff_member_id` (`staff_member_id`),
  KEY `index_audit_events_on_auditable_type_and_auditable_id` (`auditable_type`,`auditable_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;