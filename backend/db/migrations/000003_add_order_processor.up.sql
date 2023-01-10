START TRANSACTION;

ALTER TABLE orders ADD processor_id int;
ALTER TABLE orders ADD CONSTRAINT fk_processor_id FOREIGN KEY (processor_id) REFERENCES staff_members(id);

COMMIT;