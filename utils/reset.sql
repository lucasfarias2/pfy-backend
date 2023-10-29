-- To drop the database
DROP DATABASE IF EXISTS pcbe;

-- To create the database again
CREATE DATABASE pcbe;

-- Create organization table
DROP TABLE IF EXISTS projects CASCADE;
DROP TABLE IF EXISTS organizations CASCADE;

-- Types first
CREATE TYPE task_status AS ENUM ('Pending', 'Running', 'Success', 'Failed');
CREATE TYPE toolkit AS ENUM ('react');

-- Create organization table
CREATE TABLE organizations
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE
);

-- Create project table
CREATE TABLE projects
(
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) UNIQUE,
    github_repo     VARCHAR(255) UNIQUE,
    organization_id INT REFERENCES organizations (id),
    toolkit         toolkit NOT NULL
);

CREATE TABLE tasks
(
    id         SERIAL PRIMARY KEY,
    project_id INT         NOT NULL REFERENCES projects (id),
    task_name  TEXT,
    status     task_status NOT NULL,
    message    TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert into organization
INSERT INTO organizations (name)
VALUES ('Packlify');

-- Insert into project, associating with the organization and toolkit
INSERT INTO projects (name, organization_id, toolkit, github_repo)
VALUES ('test-project-1', 1, 'react', 'test-repo-1');
