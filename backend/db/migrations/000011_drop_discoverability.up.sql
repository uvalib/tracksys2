START TRANSACTION;

ALTER table components DROP COLUMN discoverability;

ALTER table metadata
   DROP COLUMN discoverability,
   DROP column qdc_generated_at;

COMMIT;