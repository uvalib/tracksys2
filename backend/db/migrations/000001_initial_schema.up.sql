START TRANSACTION;

CREATE TABLE IF NOT EXISTS `academic_statuses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `customers_count` int DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_academic_statuses_on_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `active_admin_comments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `resource_id` int NOT NULL,
  `resource_type` varchar(255) NOT NULL,
  `author_id` int DEFAULT NULL,
  `author_type` varchar(255) DEFAULT NULL,
  `body` text,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `namespace` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_admin_notes_on_resource_type_and_resource_id` (`resource_type`,`resource_id`),
  KEY `index_active_admin_comments_on_namespace` (`namespace`),
  KEY `index_active_admin_comments_on_author_type_and_author_id` (`author_type`,`author_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `addresses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `addressable_id` int NOT NULL,
  `addressable_type` varchar(20) NOT NULL,
  `address_type` varchar(20) NOT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `first_name` varchar(255) DEFAULT NULL,
  `address_1` varchar(255) DEFAULT NULL,
  `address_2` varchar(255) DEFAULT NULL,
  `city` varchar(255) DEFAULT NULL,
  `state` varchar(255) DEFAULT NULL,
  `country` varchar(255) DEFAULT NULL,
  `post_code` varchar(255) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `organization` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `agencies` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `ancestry` varchar(255) DEFAULT NULL,
  `names_depth_cache` varchar(255) DEFAULT NULL,
  `orders_count` int DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_agencies_on_name` (`name`),
  KEY `index_agencies_on_ancestry` (`ancestry`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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

CREATE TABLE IF NOT EXISTS `ar_internal_metadata` (
  `key` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `assignments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `project_id` int DEFAULT NULL,
  `step_id` int DEFAULT NULL,
  `staff_member_id` int DEFAULT NULL,
  `assigned_at` datetime DEFAULT NULL,
  `started_at` datetime DEFAULT NULL,
  `finished_at` datetime DEFAULT NULL,
  `status` int DEFAULT '0',
  `duration_minutes` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_assignments_on_step_id` (`step_id`),
  KEY `index_assignments_on_staff_member_id` (`staff_member_id`),
  KEY `index_assignments_on_project_id` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `attachments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `unit_id` int DEFAULT NULL,
  `filename` varchar(255) DEFAULT NULL,
  `md5` varchar(255) DEFAULT NULL,
  `description` text,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `index_attachments_on_unit_id` (`unit_id`),
  CONSTRAINT `fk_rails_c777e65020` FOREIGN KEY (`unit_id`) REFERENCES `units` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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

CREATE TABLE IF NOT EXISTS `availability_policies` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `metadata_count` int DEFAULT '0',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `pid` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `categories` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `projects_count` int DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `checkouts` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `metadata_id` bigint DEFAULT NULL,
  `checkout_at` datetime DEFAULT NULL,
  `return_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_checkouts_on_metadata_id` (`metadata_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `collection_facets` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `component_types` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `components_count` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_component_types_on_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `components` (
  `id` int NOT NULL AUTO_INCREMENT,
  `component_type_id` int NOT NULL DEFAULT '0',
  `parent_component_id` int NOT NULL DEFAULT '0',
  `title` varchar(255) DEFAULT NULL,
  `label` varchar(255) DEFAULT NULL,
  `date` varchar(255) DEFAULT NULL,
  `content_desc` text,
  `idno` varchar(255) DEFAULT NULL,
  `barcode` varchar(255) DEFAULT NULL,
  `seq_number` int DEFAULT NULL,
  `pid` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `desc_metadata` text,
  `discoverability` tinyint(1) DEFAULT '1',
  `level` text,
  `ead_id_att` varchar(255) DEFAULT NULL,
  `date_dl_ingest` datetime DEFAULT NULL,
  `date_dl_update` datetime DEFAULT NULL,
  `master_files_count` int NOT NULL DEFAULT '0',
  `ancestry` varchar(255) DEFAULT NULL,
  `pids_depth_cache` varchar(255) DEFAULT NULL,
  `ead_id_atts_depth_cache` varchar(255) DEFAULT NULL,
  `followed_by_id` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_components_on_component_type_id` (`component_type_id`),
  KEY `index_components_on_ancestry` (`ancestry`),
  KEY `index_components_on_followed_by_id` (`followed_by_id`),
  KEY `index_components_on_pid` (`pid`),
  CONSTRAINT `components_component_type_id_fk` FOREIGN KEY (`component_type_id`) REFERENCES `component_types` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `container_types` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `has_folders` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `customers` (
  `id` int NOT NULL AUTO_INCREMENT,
  `department_id` int DEFAULT NULL,
  `academic_status_id` int NOT NULL DEFAULT '0',
  `last_name` varchar(255) DEFAULT NULL,
  `first_name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `master_files_count` int DEFAULT '0',
  `orders_count` int DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `index_customers_on_last_name` (`last_name`),
  KEY `index_customers_on_first_name` (`first_name`),
  KEY `index_customers_on_email` (`email`),
  KEY `index_customers_on_academic_status_id` (`academic_status_id`),
  KEY `index_customers_on_department_id` (`department_id`),
  CONSTRAINT `customers_academic_status_id_fk` FOREIGN KEY (`academic_status_id`) REFERENCES `academic_statuses` (`id`),
  CONSTRAINT `customers_department_id_fk` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `delayed_jobs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `priority` int NOT NULL DEFAULT '0',
  `attempts` int NOT NULL DEFAULT '0',
  `handler` text NOT NULL,
  `last_error` text,
  `run_at` datetime DEFAULT NULL,
  `locked_at` datetime DEFAULT NULL,
  `failed_at` datetime DEFAULT NULL,
  `locked_by` varchar(255) DEFAULT NULL,
  `queue` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `delayed_jobs_priority` (`priority`,`run_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `departments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `customers_count` int DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_departments_on_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `equipment` (
  `id` int NOT NULL AUTO_INCREMENT,
  `type` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `serial_number` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `status` int DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `events` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `job_status_id` bigint DEFAULT NULL,
  `level` int DEFAULT NULL,
  `text` text,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_events_on_job_status_id` (`job_status_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `external_systems` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `public_url` varchar(255) DEFAULT NULL,
  `api_url` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_external_systems_on_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `image_tech_meta` (
  `id` int NOT NULL AUTO_INCREMENT,
  `master_file_id` int NOT NULL DEFAULT '0',
  `image_format` varchar(255) DEFAULT NULL,
  `width` int DEFAULT NULL,
  `height` int DEFAULT NULL,
  `resolution` int DEFAULT NULL,
  `color_space` varchar(255) DEFAULT NULL,
  `depth` int DEFAULT NULL,
  `compression` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `color_profile` varchar(255) DEFAULT NULL,
  `equipment` varchar(255) DEFAULT NULL,
  `software` varchar(255) DEFAULT NULL,
  `model` varchar(255) DEFAULT NULL,
  `exif_version` varchar(255) DEFAULT NULL,
  `capture_date` datetime DEFAULT NULL,
  `iso` int DEFAULT NULL,
  `exposure_bias` varchar(255) DEFAULT NULL,
  `exposure_time` varchar(255) DEFAULT NULL,
  `aperture` varchar(255) DEFAULT NULL,
  `focal_length` decimal(10,0) DEFAULT NULL,
  `orientation` int DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `index_image_tech_meta_on_master_file_id` (`master_file_id`),
  CONSTRAINT `image_tech_meta_master_file_id_fk` FOREIGN KEY (`master_file_id`) REFERENCES `master_files` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `intended_uses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `description` varchar(255) DEFAULT NULL,
  `is_internal_use_only` tinyint(1) NOT NULL DEFAULT '0',
  `is_approved` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `units_count` int DEFAULT '0',
  `deliverable_format` varchar(255) DEFAULT NULL,
  `deliverable_resolution` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_intended_uses_on_description` (`description`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `invoices` (
  `id` int NOT NULL AUTO_INCREMENT,
  `order_id` int NOT NULL DEFAULT '0',
  `date_invoice` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `fee_amount_paid` int DEFAULT NULL,
  `date_fee_paid` datetime DEFAULT NULL,
  `date_second_notice_sent` datetime DEFAULT NULL,
  `transmittal_number` text,
  `notes` text,
  `permanent_nonpayment` tinyint(1) DEFAULT '0',
  `date_fee_declined` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_invoices_on_order_id` (`order_id`),
  CONSTRAINT `invoices_order_id_fk` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `job_statuses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `status` varchar(255) NOT NULL DEFAULT 'pending',
  `started_at` datetime DEFAULT NULL,
  `ended_at` datetime DEFAULT NULL,
  `failures` int NOT NULL DEFAULT '0',
  `error` text,
  `originator_id` int DEFAULT NULL,
  `originator_type` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_job_statuses_on_originator_type_and_originator_id` (`originator_type`,`originator_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `locations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `container_type_id` bigint DEFAULT NULL,
  `container_id` varchar(255) NOT NULL,
  `folder_id` varchar(255) DEFAULT NULL,
  `notes` text,
  `metadata_id` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_locations_on_container_type_id` (`container_type_id`),
  KEY `index_locations_on_metadata_id` (`metadata_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `master_file_locations` (
  `location_id` bigint DEFAULT NULL,
  `master_file_id` bigint DEFAULT NULL,
  KEY `index_master_file_locations_on_location_id` (`location_id`),
  KEY `index_master_file_locations_on_master_file_id` (`master_file_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `master_file_tags` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `master_file_id` bigint DEFAULT NULL,
  `tag_id` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_master_file_tags_on_master_file_id` (`master_file_id`),
  KEY `index_master_file_tags_on_tag_id` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `master_files` (
  `id` int NOT NULL AUTO_INCREMENT,
  `unit_id` int NOT NULL DEFAULT '0',
  `component_id` int DEFAULT NULL,
  `filename` varchar(255) DEFAULT NULL,
  `filesize` int DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `date_archived` datetime DEFAULT NULL,
  `description` text,
  `pid` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `transcription_text` text,
  `md5` varchar(255) DEFAULT NULL,
  `date_dl_ingest` datetime DEFAULT NULL,
  `date_dl_update` datetime DEFAULT NULL,
  `creation_date` varchar(255) DEFAULT NULL,
  `primary_author` varchar(255) DEFAULT NULL,
  `metadata_id` int DEFAULT NULL,
  `original_mf_id` int DEFAULT NULL,
  `deaccessioned_at` datetime DEFAULT NULL,
  `deaccession_note` text,
  `deaccessioned_by_id` int DEFAULT NULL,
  `text_source` int DEFAULT NULL,
  `exemplar` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `index_master_files_on_unit_id` (`unit_id`),
  KEY `index_master_files_on_filename` (`filename`),
  KEY `index_master_files_on_pid` (`pid`),
  KEY `index_master_files_on_component_id` (`component_id`),
  KEY `index_master_files_on_title` (`title`),
  KEY `index_master_files_on_date_dl_ingest` (`date_dl_ingest`),
  KEY `index_master_files_on_date_dl_update` (`date_dl_update`),
  KEY `index_master_files_on_metadata_id` (`metadata_id`),
  KEY `index_master_files_on_original_mf_id` (`original_mf_id`),
  CONSTRAINT `master_files_component_id_fk` FOREIGN KEY (`component_id`) REFERENCES `components` (`id`),
  CONSTRAINT `master_files_unit_id_fk` FOREIGN KEY (`unit_id`) REFERENCES `units` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `messages` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `subject` varchar(255) NOT NULL,
  `message` text,
  `read` tinyint(1) DEFAULT '0',
  `from_id` bigint DEFAULT NULL,
  `to_id` bigint DEFAULT NULL,
  `sent_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted` tinyint(1) DEFAULT '0',
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_messages_on_from_id` (`from_id`),
  KEY `index_messages_on_to_id` (`to_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `metadata` (
  `id` int NOT NULL AUTO_INCREMENT,
  `is_personal_item` tinyint(1) NOT NULL DEFAULT '0',
  `is_manuscript` tinyint(1) NOT NULL DEFAULT '0',
  `title` text,
  `creator_name` varchar(255) DEFAULT NULL,
  `catalog_key` varchar(255) DEFAULT NULL,
  `barcode` varchar(255) DEFAULT NULL,
  `call_number` varchar(255) DEFAULT NULL,
  `pid` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `parent_metadata_id` int NOT NULL DEFAULT '0',
  `desc_metadata` text,
  `units_count` int DEFAULT '0',
  `type` varchar(255) DEFAULT 'SirsiMetadata',
  `external_uri` varchar(255) DEFAULT NULL,
  `supplemental_uri` varchar(255) DEFAULT NULL,
  `collection_id` varchar(255) DEFAULT NULL,
  `ocr_hint_id` int DEFAULT NULL,
  `ocr_language_hint` varchar(255) DEFAULT NULL,
  `creator_death_date` int DEFAULT NULL,
  `preservation_tier_id` bigint DEFAULT NULL,
  `external_system_id` bigint DEFAULT NULL,
  `supplemental_system_id` bigint DEFAULT NULL,
  `use_right_id` bigint DEFAULT NULL,
  `availability_policy_id` bigint DEFAULT NULL,
  `discoverability` tinyint(1) DEFAULT NULL,
  `dpla` tinyint(1) DEFAULT NULL,
  `use_right_rationale` varchar(255) DEFAULT NULL,
  `collection_facet` varchar(255) DEFAULT NULL,
  `date_dl_ingest` datetime DEFAULT NULL,
  `date_dl_update` datetime DEFAULT NULL,
  `qdc_generated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_metadata_on_barcode` (`barcode`),
  KEY `index_metadata_on_call_number` (`call_number`),
  KEY `index_metadata_on_catalog_key` (`catalog_key`),
  KEY `index_metadata_on_pid` (`pid`),
  KEY `index_metadata_on_parent_metadata_id` (`parent_metadata_id`),
  KEY `index_metadata_on_ocr_hint_id` (`ocr_hint_id`),
  KEY `index_metadata_on_preservation_tier_id` (`preservation_tier_id`),
  KEY `index_metadata_on_external_system_id` (`external_system_id`),
  KEY `index_metadata_on_supplemental_system_id` (`supplemental_system_id`),
  KEY `index_metadata_on_use_right_id` (`use_right_id`),
  KEY `index_metadata_on_availability_policy_id` (`availability_policy_id`),
  CONSTRAINT `fk_rails_4e58402857` FOREIGN KEY (`ocr_hint_id`) REFERENCES `ocr_hints` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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

CREATE TABLE IF NOT EXISTS `notes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `staff_member_id` int DEFAULT NULL,
  `project_id` int DEFAULT NULL,
  `note` text,
  `note_type` int DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `step_id` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_notes_on_staff_member_id` (`staff_member_id`),
  KEY `index_notes_on_project_id` (`project_id`),
  KEY `index_notes_on_step_id` (`step_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `notes_problems` (
  `note_id` bigint DEFAULT NULL,
  `problem_id` bigint DEFAULT NULL,
  KEY `index_notes_problems_on_note_id` (`note_id`),
  KEY `index_notes_problems_on_problem_id` (`problem_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `ocr_hints` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `ocr_candidate` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `order_items` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `order_id` bigint DEFAULT NULL,
  `intended_use_id` bigint DEFAULT NULL,
  `title` text,
  `pages` text,
  `call_number` varchar(255) DEFAULT NULL,
  `author` varchar(255) DEFAULT NULL,
  `year` varchar(255) DEFAULT NULL,
  `location` varchar(255) DEFAULT NULL,
  `source_url` varchar(255) DEFAULT NULL,
  `description` text,
  `converted` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `index_order_items_on_order_id` (`order_id`),
  KEY `index_order_items_on_intended_use_id` (`intended_use_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `orders` (
  `id` int NOT NULL AUTO_INCREMENT,
  `customer_id` int NOT NULL DEFAULT '0',
  `agency_id` int DEFAULT NULL,
  `order_status` varchar(255) DEFAULT NULL,
  `is_approved` tinyint(1) NOT NULL DEFAULT '0',
  `order_title` varchar(255) DEFAULT NULL,
  `date_request_submitted` datetime DEFAULT NULL,
  `date_order_approved` datetime DEFAULT NULL,
  `date_deferred` datetime DEFAULT NULL,
  `date_canceled` datetime DEFAULT NULL,
  `date_due` date DEFAULT NULL,
  `date_customer_notified` datetime DEFAULT NULL,
  `fee` decimal(7,2) DEFAULT NULL,
  `special_instructions` text,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `staff_notes` text,
  `email` text,
  `date_patron_deliverables_complete` datetime DEFAULT NULL,
  `date_archiving_complete` datetime DEFAULT NULL,
  `date_finalization_begun` datetime DEFAULT NULL,
  `date_fee_estimate_sent_to_customer` datetime DEFAULT NULL,
  `units_count` int DEFAULT '0',
  `invoices_count` int DEFAULT '0',
  `master_files_count` int DEFAULT '0',
  `date_completed` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_orders_on_customer_id` (`customer_id`),
  KEY `index_orders_on_agency_id` (`agency_id`),
  KEY `index_orders_on_order_status` (`order_status`),
  KEY `index_orders_on_date_request_submitted` (`date_request_submitted`),
  KEY `index_orders_on_date_due` (`date_due`),
  KEY `index_orders_on_date_archiving_complete` (`date_archiving_complete`),
  KEY `index_orders_on_date_order_approved` (`date_order_approved`),
  CONSTRAINT `orders_agency_id_fk` FOREIGN KEY (`agency_id`) REFERENCES `agencies` (`id`),
  CONSTRAINT `orders_customer_id_fk` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `preservation_tiers` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `problems` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `label` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `project_equipment` (
  `id` int NOT NULL AUTO_INCREMENT,
  `project_id` int DEFAULT NULL,
  `equipment_id` int DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_project_equipment_on_equipment_id` (`equipment_id`),
  KEY `index_project_equipment_on_project_id` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `projects` (
  `id` int NOT NULL AUTO_INCREMENT,
  `workflow_id` int DEFAULT NULL,
  `unit_id` int DEFAULT NULL,
  `owner_id` int DEFAULT NULL,
  `current_step_id` int DEFAULT NULL,
  `priority` int DEFAULT '0',
  `due_on` date DEFAULT NULL,
  `item_condition` int DEFAULT NULL,
  `added_at` datetime DEFAULT NULL,
  `started_at` datetime DEFAULT NULL,
  `finished_at` datetime DEFAULT NULL,
  `category_id` int DEFAULT NULL,
  `viu_number` varchar(255) DEFAULT NULL,
  `capture_resolution` int DEFAULT NULL,
  `resized_resolution` int DEFAULT NULL,
  `resolution_note` varchar(255) DEFAULT NULL,
  `workstation_id` int DEFAULT NULL,
  `condition_note` text,
  `container_type_id` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_projects_on_workflow_id` (`workflow_id`),
  KEY `index_projects_on_unit_id` (`unit_id`),
  KEY `index_projects_on_category_id` (`category_id`),
  KEY `index_projects_on_workstation_id` (`workstation_id`),
  KEY `index_projects_on_container_type_id` (`container_type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `sirsi_metadata_components` (
  `sirsi_metadata_id` int DEFAULT NULL,
  `component_id` int DEFAULT NULL,
  KEY `bibl_id` (`sirsi_metadata_id`),
  KEY `component_id` (`component_id`),
  CONSTRAINT `sirsi_metadata_components_ibfk_2` FOREIGN KEY (`component_id`) REFERENCES `components` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `staff_members` (
  `id` int NOT NULL AUTO_INCREMENT,
  `computing_id` varchar(255) DEFAULT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `first_name` varchar(255) DEFAULT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `role` int DEFAULT '0',
  `notes` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_staff_members_on_computing_id` (`computing_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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

CREATE TABLE IF NOT EXISTS `statistics` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `value` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `group` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `steps` (
  `id` int NOT NULL AUTO_INCREMENT,
  `step_type` int DEFAULT '3',
  `name` varchar(255) DEFAULT NULL,
  `description` text,
  `workflow_id` int DEFAULT NULL,
  `next_step_id` int DEFAULT NULL,
  `fail_step_id` int DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `owner_type` int DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `index_steps_on_workflow_id` (`workflow_id`),
  KEY `index_steps_on_next_step_id` (`next_step_id`),
  KEY `index_steps_on_fail_step_id` (`fail_step_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `tags` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `tag` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `units` (
  `id` int NOT NULL AUTO_INCREMENT,
  `order_id` int NOT NULL DEFAULT '0',
  `metadata_id` int DEFAULT NULL,
  `unit_status` varchar(255) DEFAULT NULL,
  `unit_extent_estimated` int DEFAULT NULL,
  `unit_extent_actual` int DEFAULT NULL,
  `patron_source_url` text,
  `special_instructions` text,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `intended_use_id` int DEFAULT NULL,
  `staff_notes` text,
  `date_archived` datetime DEFAULT NULL,
  `date_patron_deliverables_ready` datetime DEFAULT NULL,
  `include_in_dl` tinyint(1) DEFAULT '0',
  `date_dl_deliverables_ready` datetime DEFAULT NULL,
  `remove_watermark` tinyint(1) DEFAULT '0',
  `master_files_count` int DEFAULT '0',
  `complete_scan` tinyint(1) DEFAULT '0',
  `reorder` tinyint(1) DEFAULT '0',
  `throw_away` tinyint(1) DEFAULT '0',
  `ocr_master_files` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `index_units_on_order_id` (`order_id`),
  KEY `index_units_on_date_archived` (`date_archived`),
  KEY `index_units_on_intended_use_id` (`intended_use_id`),
  KEY `index_units_on_date_dl_deliverables_ready` (`date_dl_deliverables_ready`),
  KEY `index_units_on_metadata_id` (`metadata_id`),
  CONSTRAINT `units_intended_use_id_fk` FOREIGN KEY (`intended_use_id`) REFERENCES `intended_uses` (`id`),
  CONSTRAINT `units_order_id_fk` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `use_rights` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `metadata_count` int DEFAULT '0',
  `uri` varchar(255) DEFAULT NULL,
  `statement` text,
  `commercial_use` tinyint(1) DEFAULT '0',
  `educational_use` tinyint(1) DEFAULT '0',
  `modifications` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_use_rights_on_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `workflows` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `description` text,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `active` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `workstation_equipment` (
  `id` int NOT NULL AUTO_INCREMENT,
  `workstation_id` int DEFAULT NULL,
  `equipment_id` int DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_workstation_equipment_on_workstation_id` (`workstation_id`),
  KEY `index_workstation_equipment_on_equipment_id` (`equipment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `workstations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `status` int DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

COMMIT;