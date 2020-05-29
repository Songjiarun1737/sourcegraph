CREATE TABLE "meta" (
    "id" integer PRIMARY KEY NOT NULL,
    "lsifVersion" text NOT NULL,
    "sourcegraphVersion" text NOT NULL,
    "numResultChunks" integer NOT NULL
);

CREATE TABLE "documents" (
    "path" text PRIMARY KEY NOT NULL,
    "data" blob NOT NULL
);

CREATE TABLE "resultChunks" (
    "id" integer PRIMARY KEY NOT NULL,
    "data" blob NOT NULL
);

CREATE TABLE "definitions" (
    "id" integer PRIMARY KEY NOT NULL,
    "scheme" text NOT NULL,
    "identifier" text NOT NULL,
    "data" blob NOT NULL
);

CREATE TABLE "references" (
    "id" integer PRIMARY KEY NOT NULL,
    "scheme" text NOT NULL,
    "identifier" text NOT NULL,
    "data" blob NOT NULL
);
