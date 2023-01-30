START TRANSACTION;

ALTER TABLE agencies
   ADD column ancestry varchar(255) DEFAULT NULL,
   ADD column names_depth_cache varchar(255) DEFAULT NULL,
   ADD column orders_count int DEFAULT 0;

COMMIT;