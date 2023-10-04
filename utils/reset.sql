-- To drop the database
DROP DATABASE IF EXISTS pcbe;

-- To create the database again
CREATE DATABASE pcbe;

-- Create organization table
DROP TABLE IF EXISTS projects CASCADE;
DROP TABLE IF EXISTS toolkits CASCADE;
DROP TABLE IF EXISTS organizations CASCADE;

-- Create organization table
CREATE TABLE organizations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

-- Create toolkit table
CREATE TABLE toolkits (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    description TEXT
);

-- Create project table
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    organization_id INT REFERENCES organizations(id),
    toolkit_id INT REFERENCES toolkits(id)
);

-- Insert into organization
INSERT INTO organizations (name) VALUES ('Packlify');

-- Insert into toolkit
INSERT INTO toolkits (name, description) VALUES ('React', 'This is a React toolkit');

-- Insert into project, associating with the organization and toolkit
INSERT INTO projects (name, organization_id, toolkit_id) VALUES ('Test Project', 1, 1);

