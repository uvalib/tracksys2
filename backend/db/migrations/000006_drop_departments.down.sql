START TRANSACTION;
CREATE TABLE IF NOT EXISTS `departments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `customers_count` int DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_departments_on_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE customers
  ADD COLUMN department_id INT DEFAULT NULL;

ALTER TABLE customers
   ADD CONSTRAINT customers_department_id_fk
   FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE;

COMMIT;