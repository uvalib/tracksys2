START TRANSACTION;

ALTER table components DROP COLUMN master_files_count;

COMMIT;