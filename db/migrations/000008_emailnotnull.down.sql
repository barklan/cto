begin;

ALTER TABLE client
ALTER COLUMN email DROP NOT NULL;

commit;
