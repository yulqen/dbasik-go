/* should be added manually at this point */
CREATE DATABASE dbasik;
\c dbasik

CREATE ROLE 'dbasik' WITH LOGIN PASSWORD 'dbasik';
GRANT ALL PRIVILEGES ON DATABASE 'dbasik' TO 'dbasik';
ALTER DATABASE 'dbasik' OWNER TO 'dbasik';
CREATE EXTENSION IF NOT EXISTS citext;
