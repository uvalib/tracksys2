START TRANSACTION;

ALTER table customers
   DROP COLUMN master_files_count,
   DROP COLUMN orders_count;

COMMIT;