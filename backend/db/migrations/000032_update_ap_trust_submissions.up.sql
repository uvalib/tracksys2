START TRANSACTION;

ALTER table ap_trust_submissions
   ADD COLUMN processed_at datetime DEFAULT NULL,
   ADD COLUMN success tinyint(1) DEFAULT 0;

update ap_trust_submissions apts2  set success=1, processed_at =
   (select finished_at from ap_trust_statuses apts1 where apts2.metadata_id = apts1.metadata_id);

COMMIT;