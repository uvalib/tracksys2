START TRANSACTION;

ALTER table metadata ADD COLUMN units_count  int DEFAULT 0;

COMMIT;