ALTER table hathitrust_statuses
   ADD COLUMN `requested_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
   ADD COLUMN `package_created_at` datetime DEFAULT NULL;