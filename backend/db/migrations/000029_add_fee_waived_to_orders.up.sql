ALTER table orders
   ADD COLUMN fee_waived bool NOT NULL DEFAULT 0,
   ADD COLUMN date_fee_waived datetime DEFAULT NULL;
