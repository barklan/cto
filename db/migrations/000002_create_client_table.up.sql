BEGIN;

CREATE TABLE IF NOT EXISTS client(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    tg_nick VARCHAR (100) NOT NULL
);

CREATE TABLE IF NOT EXISTS project (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    client_id UUID NOT NULL REFERENCES client(id),
    pretty_title VARCHAR (100),
    secret_key VARCHAR (100) NOT NULL
);

CREATE TABLE IF NOT EXISTS chat (
    id INTEGER PRIMARY KEY,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    project_id UUID,
    UNIQUE(project_id),
    FOREIGN KEY(project_id) REFERENCES project(id)
);

COMMIT;
