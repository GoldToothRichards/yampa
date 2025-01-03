-- Display all Views from a database
SELECT name, engine, create_table_query
FROM system.tables
WHERE database = 'trades'
  AND engine = 'View';
ORDER BY database, name;

-- Display all Views/Materialized Views from trades database
SELECT database, name, engine
FROM system.tables
WHERE engine LIKE '%View'
AND database = 'trades'
ORDER BY database, name;
