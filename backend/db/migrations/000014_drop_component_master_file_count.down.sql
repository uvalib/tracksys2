START TRANSACTION;

ALTER table components ADD COLUMN master_files_count int NOT NULL DEFAULT 0;

COMMIT;