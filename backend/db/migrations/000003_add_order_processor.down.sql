START TRANSACTION;

ALTER TABLE orders DROP  FOREIGN KEY fk_processor_id;
ALTER TABLE orders DROP processor_id;

COMMIT;