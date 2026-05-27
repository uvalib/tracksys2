START TRANSACTION;

CREATE TABLE IF NOT EXISTS `preservation_tiers` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `preservation_tiers` VALUES 
   (1,'Tier 1','Backed-up'),
   (2,'Tier 2','Duplicated in separate physical location'),
   (3,'Tier 3','Multiple storage technologies, multiple geographic regions');

ALTER TABLE metadata
  ADD COLUMN preservation_tier_id bigint DEFAULT NULL;

COMMIT;