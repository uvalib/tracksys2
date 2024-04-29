START TRANSACTION;

ALTER TABLE archivesspace_reviews DROP FOREIGN KEY `archivesspace_reviews_unit_id_fk`;
ALTER TABLE archivesspace_reviews DROP COLUMN `unit_id`;

COMMIT;