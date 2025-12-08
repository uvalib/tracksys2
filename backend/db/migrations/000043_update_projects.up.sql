START TRANSACTION;

ALTER TABLE projects
   DROP column priority,
   DROP column viu_number;

COMMIT;