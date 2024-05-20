CREATE TABLE "Project"(
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "version" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "equivalence_required" BOOLEAN NOT NULL,
    "status_required" BOOLEAN NOT NULL,
    "created" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "modified" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "Project" ADD PRIMARY KEY("id");
CREATE TABLE "element"(
    "mapping" BIGINT NOT NULL,
    "code_system_role" INTEGER NOT NULL,
    "concept" BIGINT NULL
);
CREATE TABLE "CodeSystem"(
    "id" SERIAL NOT NULL,
    "uri" TEXT NOT NULL,
    "version" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "title" TEXT NULL,
    "description" TEXT NULL,
    "author" TEXT NULL
);
ALTER TABLE
    "CodeSystem" ADD PRIMARY KEY("id");
COMMENT
ON COLUMN
    "CodeSystem"."name" IS 'computer friendly';
COMMENT
ON COLUMN
    "CodeSystem"."title" IS 'human friendly';
CREATE TABLE "Concept"(
    "id" bigserial NOT NULL,
    "system" INTEGER NOT NULL,
    "code" TEXT NOT NULL,
    "display" TEXT NOT NULL
);
ALTER TABLE
    "Concept" ADD PRIMARY KEY("id");
CREATE TABLE "Mapping"(
    "id" bigserial NOT NULL,
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
ALTER TABLE
    "Mapping" ADD PRIMARY KEY("id");
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
        ) NOT NULL
);
ALTER TABLE
    "project_permission" ADD CONSTRAINT "project_permission_user_unique" UNIQUE("user");
ALTER TABLE
    "project_permission" ADD CONSTRAINT "project_permission_project_unique" UNIQUE("project");
CREATE TABLE "User"(
    "id" UUID NOT NULL,
    "user_name" TEXT NOT NULL,
    "log_name" TEXT NOT NULL,
    "affiliation" TEXT NOT NULL
);
ALTER TABLE
    "User" ADD PRIMARY KEY("id");
CREATE TABLE "code_system_role"(
    "id" SERIAL NOT NULL,
    "type" VARCHAR(255) CHECK
        ("type" IN('source', 'target')) NOT NULL,
        "project" INTEGER NOT NULL,
        "system" INTEGER NOT NULL,
        "name" TEXT NOT NULL,
        "position" INTEGER NOT NULL
);
ALTER TABLE
    "code_system_role" ADD PRIMARY KEY("id");
COMMENT
ON COLUMN
    "code_system_role"."id" IS 'id braucht man, um von element aus einfacher zu referenzierten';
COMMENT
ON COLUMN
    "code_system_role"."type" IS 'soll Type auch primary Schlüssel sein? Falls ja, kann concept sowohl als source als auch als target im gleichen mapping vorkommen kann';
ALTER TABLE
    "Mapping" ADD CONSTRAINT "mapping_project_foreign" FOREIGN KEY("project") REFERENCES "Project"("id");
ALTER TABLE
    "element" ADD CONSTRAINT "element_code_system_role_foreign" FOREIGN KEY("code_system_role") REFERENCES "code_system_role"("id");
ALTER TABLE
    "code_system_role" ADD CONSTRAINT "code_system_role_project_foreign" FOREIGN KEY("project") REFERENCES "Project"("id");
ALTER TABLE
    "project_permission" ADD CONSTRAINT "project_permission_user_foreign" FOREIGN KEY("user") REFERENCES "User"("id");
ALTER TABLE
    "element" ADD CONSTRAINT "element_mapping_foreign" FOREIGN KEY("mapping") REFERENCES "Mapping"("id");
ALTER TABLE
    "Concept" ADD CONSTRAINT "concept_system_foreign" FOREIGN KEY("system") REFERENCES "CodeSystem"("id");
ALTER TABLE
    "code_system_role" ADD CONSTRAINT "code_system_role_system_foreign" FOREIGN KEY("system") REFERENCES "CodeSystem"("id");
ALTER TABLE
    "element" ADD CONSTRAINT "element_concept_foreign" FOREIGN KEY("concept") REFERENCES "Concept"("id");
ALTER TABLE
    "project_permission" ADD CONSTRAINT "project_permission_project_foreign" FOREIGN KEY("project") REFERENCES "Project"("id");