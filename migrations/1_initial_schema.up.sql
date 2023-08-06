CREATE TABLE projects
(
    id          bigint NOT NULL PRIMARY KEY,
    name        text,
    description text,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone
);

CREATE TABLE items
(
    id          bigint NOT NULL PRIMARY KEY,
    project_id  bigint NOT NULL,
    name        text,
    description text,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    CONSTRAINT fk_project FOREIGN KEY(project_id) REFERENCES projects(id),
    CONSTRAINT id_project_id_unique UNIQUE (id,project_id)
);
