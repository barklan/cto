begin;

ALTER TABLE client
ALTER COLUMN personal_chat SET NOT NULL;

commit;
