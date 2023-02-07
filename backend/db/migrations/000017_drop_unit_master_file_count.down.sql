START TRANSACTION;

ALTER table units ADD COLUMN master_files_count int DEFAULT 0;

COMMIT;