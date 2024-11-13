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
    "mappingId" BIGINT NOT NULL,
    "codeSystemRoleId" INTEGER NOT NULL,
    "conceptId" BIGINT NULL
);
ALTER TABLE
    "element" ADD PRIMARY KEY("mappingId");
ALTER TABLE
    "element" ADD PRIMARY KEY("codeSystemRoleId");
CREATE TABLE "CodeSystem"(
    "id" SERIAL NOT NULL,
    "uri" TEXT NOT NULL,
    "version" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "title" TEXT NULL,
    "description" TEXT NULL,
    "author" TEXT NULL,
    "created" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "modified" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
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
    "codeSystemId" INTEGER NOT NULL,
    "code" TEXT NOT NULL,
    "display" TEXT NOT NULL,
    "displaySearchVector" TEXT NOT NULL
);
ALTER TABLE
    "Concept" ADD PRIMARY KEY("id");
CREATE INDEX "concept_codesystemid_index" ON
    "Concept"("codeSystemId");
CREATE TABLE "Mapping"(
    "id" bigserial NOT NULL,
    "projectId" INTEGER NOT NULL,
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
CREATE INDEX "mapping_projectid_index" ON
    "Mapping"("projectId");
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
    "project_permission" ADD PRIMARY KEY("user");
ALTER TABLE
    "project_permission" ADD PRIMARY KEY("project");
CREATE TABLE "User"(
    "id" UUID NOT NULL,
    "userName" TEXT NOT NULL,
    "fullName" TEXT NOT NULL,
    "email" TEXT NOT NULL
);
ALTER TABLE
    "User" ADD PRIMARY KEY("id");
CREATE TABLE "code_system_role"(
    "id" SERIAL NOT NULL,
    "type" VARCHAR(255) CHECK
        ("type" IN('source', 'target')) NOT NULL,
        "projectId" INTEGER NOT NULL,
        "codeSystemId" INTEGER NOT NULL,
        "name" TEXT NOT NULL,
        "position" INTEGER NOT NULL
);
ALTER TABLE
    "code_system_role" ADD PRIMARY KEY("id");
CREATE INDEX "code_system_role_projectid_index" ON
    "code_system_role"("projectId");
COMMENT
ON COLUMN
    "code_system_role"."id" IS 'id braucht man, um von element aus einfacher zu referenzierten';
COMMENT
ON COLUMN
    "code_system_role"."type" IS 'soll Type auch primary Schlüssel sein? Falls ja, kann concept sowohl als source als auch als target im gleichen mapping vorkommen kann';
ALTER TABLE
    "Mapping" ADD CONSTRAINT "mapping_projectid_foreign" FOREIGN KEY("projectId") REFERENCES "Project"("id");
ALTER TABLE
    "code_system_role" ADD CONSTRAINT "code_system_role_id_foreign" FOREIGN KEY("id") REFERENCES "element"("codeSystemRoleId");
ALTER TABLE
    "code_system_role" ADD CONSTRAINT "code_system_role_projectid_foreign" FOREIGN KEY("projectId") REFERENCES "Project"("id");
ALTER TABLE
    "User" ADD CONSTRAINT "user_id_foreign" FOREIGN KEY("id") REFERENCES "project_permission"("user");
ALTER TABLE
    "Mapping" ADD CONSTRAINT "mapping_id_foreign" FOREIGN KEY("id") REFERENCES "element"("mappingId");
ALTER TABLE
    "Concept" ADD CONSTRAINT "concept_codesystemid_foreign" FOREIGN KEY("codeSystemId") REFERENCES "CodeSystem"("id");
ALTER TABLE
    "code_system_role" ADD CONSTRAINT "code_system_role_codesystemid_foreign" FOREIGN KEY("codeSystemId") REFERENCES "CodeSystem"("id");
ALTER TABLE
    "element" ADD CONSTRAINT "element_conceptid_foreign" FOREIGN KEY("conceptId") REFERENCES "Concept"("id");
ALTER TABLE
    "Project" ADD CONSTRAINT "project_id_foreign" FOREIGN KEY("id") REFERENCES "project_permission"("project");