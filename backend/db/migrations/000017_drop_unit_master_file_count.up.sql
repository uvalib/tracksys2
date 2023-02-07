START TRANSACTION;

ALTER table units DROP COLUMN master_files_count;

COMMIT;