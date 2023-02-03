START TRANSACTION;

ALTER table components ADD column discoverability tinyint(1) DEFAULT NULL;

ALTER table metadata
   ADD column discoverability tinyint(1) DEFAULT NULL,
   ADD column qdc_generated_at datetime DEFAULT NULL;

COMMIT;