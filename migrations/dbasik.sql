CREATE TABLE "datamap" (
  "id" int PRIMARY KEY,
  "name" varchar,
  "description" varchar
);

CREATE TABLE "datamapline" (
  "id" uuid PRIMARY KEY,
  "datamap_id" uuid,
  "name" string,
  "cellref" string,
  "sheet" string,
  "data_type" string
);

CREATE TABLE "returnitem" (
  "id" uuid PRIMARY KEY,
  "datamapline_id" uuid,
  "return_id" int,
  "value" string
);

CREATE TABLE "return" (
  "id" int PRIMARY KEY,
  "name" string
);

ALTER TABLE "datamapline" ADD FOREIGN KEY ("datamap_id") REFERENCES "datamap" ("id");

ALTER TABLE "returnitem" ADD FOREIGN KEY ("datamapline_id") REFERENCES "datamapline" ("id");

ALTER TABLE "returnitem" ADD FOREIGN KEY ("return_id") REFERENCES "return" ("id");
