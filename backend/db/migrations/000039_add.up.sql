ALTER table archivesspace_reviews
   ADD COLUMN unit_id int DEFAULT null,
   ADD CONSTRAINT `archivesspace_reviews_unit_id_fk` FOREIGN KEY (`unit_id`) REFERENCES `units` (`id`);