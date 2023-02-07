START TRANSACTION;

ALTER table metadata DROP COLUMN units_count;

COMMIT;