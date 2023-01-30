START TRANSACTION;

delete from events where job_status_id not in (select id from job_statuses js);

ALTER TABLE job_statuses MODIFY COLUMN id BIGINT AUTO_INCREMENT NOT NULL;

ALTER TABLE `events`
  ADD CONSTRAINT fk_job_status_id
  FOREIGN KEY (job_status_id)
  REFERENCES job_statuses(id)
  ON DELETE CASCADE;

COMMIT;