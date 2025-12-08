START TRANSACTION;

ALTER TABLE projects
   ADD column `priority` int DEFAULT '0',
   ADD column `viu_number` varchar(255) DEFAULT NULL;

COMMIT;