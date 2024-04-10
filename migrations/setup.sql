/* should be added manually at this point */
CREATE DATABASE dbasik;
\c dbasik

CREATE ROLE 'dbasik' WITH LOGIN PASSWORD 'dbasik';
CREATE EXTENSION IF NOT EXISTS citext;
