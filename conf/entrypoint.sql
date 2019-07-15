CREATE TABLE customers (
  "id" integer NOT NULL,
  "name" character varying(200) NOT NULL,
  "email" character varying(200) NOT NULL,
  "status" character varying(30) NOT NULL,
   PRIMARY KEY (id)
);

CREATE INDEX customers_status ON "customers" ("status");
