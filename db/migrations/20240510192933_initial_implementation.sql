-- +goose Up
-- +goose StatementBegin
CREATE TABLE "Project"(
    "id" SERIAL PRIMARY KEY NOT NULL,
    "name" TEXT NOT NULL,
    "version" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "equivalence_required" BOOLEAN NOT NULL,
    "status_required" BOOLEAN NOT NULL,
    "created" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "modified" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE "element"(
    "mapping" BIGINT NOT NULL,
    "code_system_role" INTEGER NOT NULL,
    "concept" BIGINT NULL,
    PRIMARY KEY ("mapping", "code_system_role")
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE "CodeSystem"(
    "id" SERIAL PRIMARY KEY NOT NULL,
    "uri" TEXT NOT NULL,
    "version" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "title" TEXT NULL,
    "description" TEXT NULL,
    "author" TEXT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE "Concept"(
    "id" bigserial PRIMARY KEY NOT NULL,
    "system" INTEGER NOT NULL,
    "code" TEXT NOT NULL,
    "display" TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE "Mapping"(
    "id" bigserial PRIMARY KEY NOT NULL,
    "project" INTEGER NOT NULL,
    "equivalence" VARCHAR(255) CHECK
        (
            "equivalence" IN(
                'equivalent',
                'not_equivalent',
                'partial'
            )
        ) NULL,
        "status" VARCHAR(255)
    CHECK
        (
            "status" IN('active', 'inactive', 'pending')
        ) NULL,
        "comment" TEXT NULL,
        "created" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
        "modified" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE "project_permission"(
    "user" UUID NOT NULL,
    "project" INTEGER NOT NULL,
    "role" VARCHAR(255) CHECK
        (
            "role" IN(
                'reviewer',
                'project_owner',
                'editor'
            )
        ) NOT NULL,
    PRIMARY KEY("user", "project")
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE "User"(
    "id" UUID PRIMARY KEY NOT NULL,
    "user_name" TEXT NOT NULL,
    "log_name" TEXT NOT NULL,
    "affiliation" TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE "code_system_role"(
    "id" SERIAL PRIMARY KEY NOT NULL,
    "type" VARCHAR(255) CHECK
        ("type" IN('source', 'target')) NOT NULL,
        "project" INTEGER NOT NULL,
        "system" INTEGER NOT NULL,
        "name" TEXT NOT NULL,
        "position" INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE
    "Mapping" ADD CONSTRAINT "mapping_project_foreign" FOREIGN KEY("project") REFERENCES "Project"("id");
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE
    "element" ADD CONSTRAINT "element_code_system_role_foreign" FOREIGN KEY("code_system_role") REFERENCES "code_system_role"("id");
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE
    "code_system_role" ADD CONSTRAINT "code_system_role_project_foreign" FOREIGN KEY("project") REFERENCES "Project"("id");
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE
    "project_permission" ADD CONSTRAINT "project_permission_user_foreign" FOREIGN KEY("user") REFERENCES "User"("id");
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE
    "element" ADD CONSTRAINT "element_mapping_foreign" FOREIGN KEY("mapping") REFERENCES "Mapping"("id");
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE
    "Concept" ADD CONSTRAINT "concept_system_foreign" FOREIGN KEY("system") REFERENCES "CodeSystem"("id");
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE
    "code_system_role" ADD CONSTRAINT "code_system_role_system_foreign" FOREIGN KEY("system") REFERENCES "CodeSystem"("id");
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE
    "element" ADD CONSTRAINT "element_concept_foreign" FOREIGN KEY("concept") REFERENCES "Concept"("id");
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE
    "project_permission" ADD CONSTRAINT "project_permission_project_foreign" FOREIGN KEY("project") REFERENCES "Project"("id");
-- +goose StatementEnd

-- create Test data

-- +goose StatementBegin
-- Insert a test user
INSERT INTO "User" ("id", "user_name", "log_name", "affiliation")
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'TestUser', 'TestLog', 'TestAffiliation');

-- Insert a test code system
INSERT INTO "CodeSystem" ("uri", "version", "name", "title", "description", "author")
VALUES ('http://test.com', '1.0', 'TestCodeSystem', 'TestTitle', 'TestDescription', 'TestAuthor');
-- +goose StatementEnd
