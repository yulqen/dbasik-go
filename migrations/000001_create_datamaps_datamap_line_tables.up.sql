CREATE TABLE IF NOT EXISTS datamaps (
  id bigserial PRIMARY KEY,
  name text,
  description text,
  created timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS datamap_lines (
  datamap_line_id bigserial PRIMARY KEY,
  datamap_id bigserial REFERENCES datamaps ON DELETE CASCADE,
  key text,
  sheet text,
  data_type text,
  cellref text
);
