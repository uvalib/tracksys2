START TRANSACTION;

ALTER table customers
   ADD COLUMN master_files_count int DEFAULT 0,
   ADD COLUMN orders_count  int DEFAULT 0;

COMMIT;