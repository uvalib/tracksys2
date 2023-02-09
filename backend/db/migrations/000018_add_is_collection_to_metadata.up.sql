START TRANSACTION;
ALTER TABLE metadata ADD COLUMN is_collection tinyint(1) NOT NULL DEFAULT 0;

-- flag some known collections as such
UPDATE metadata set is_collection=1 where id in (3002, 3009, 3109, 6405, 15784, 16315, 16585);

-- 3002  : Negatives from the Charlottesville photographic studio
-- 3009  : University of Virginia Visual History Collection
-- 3109  : Papers and photographs of Jackson Davis
-- 6405  : Cecil Lang Collection of Vanity Fair Illustrations
-- 15784 : DPLA Collection Record for UVA Digital Library Text Collections
-- 16315 : Frances Benjamin Johnston Photographic Collection
-- 16585 : Online Artifacts

COMMIT;