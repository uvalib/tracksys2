ALTER table projects
   ADD COLUMN `order_id` bigint DEFAULT NULL,
   ADD COLUMN `customer_id` bigint  DEFAULT NULL,
   ADD COLUMN `agency_id` bigint  DEFAULT NULL,
   ADD COLUMN `call_number` varchar(255)  DEFAULT NULL,
   ADD COLUMN `title` text DEFAULT NULL;
