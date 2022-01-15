begin;

ALTER TABLE client
ALTER COLUMN email SET NOT NULL;

commit;
