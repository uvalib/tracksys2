START TRANSACTION;
ALTER table customers DROP FOREIGN KEY customers_department_id_fk;
ALTER table customers DROP COLUMN department_id;
DROP TABLE iF EXISTS `departments`;
COMMIT;