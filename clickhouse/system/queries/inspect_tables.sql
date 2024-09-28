-- Display all Views from a database
SELECT name, engine, create_table_query
FROM system.tables
WHERE database = 'trades'
  AND engine = 'View';
ORDER BY database, name;

-- Display all Views from all databases
SELECT database, name, engine, create_table_query
FROM system.tables
WHERE engine = 'View'
ORDER BY database, name;
