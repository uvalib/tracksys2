START TRANSACTION;

ALTER TABLE agencies
   DROP column ancestry,
   DROP column names_depth_cache,
   DROP column orders_count;

COMMIT;