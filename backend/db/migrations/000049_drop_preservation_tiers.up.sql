START TRANSACTION;

ALTER table metadata DROP COLUMN preservation_tier_id;
DROP TABLE iF EXISTS preservation_tiers;

COMMIT;