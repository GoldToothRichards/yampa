-- Parts created
SELECT
    event_time,
    event_time_microseconds,
    rows
FROM system.part_log
WHERE (database = 'trades') AND (table = 'raw_trades') AND (event_type IN ['NewPart'])
ORDER BY event_time DESC
LIMIT 10;

-- Inspect Partitions
-- Name will be in the format:
-- <partition>_<min_data_block>_<max_data_block>_<chunk_level>_<mutation_version>
SELECT
    partition,
    name,
    active
FROM system.parts
WHERE table = 'raw_trades'

-- Inspect rows per part
SELECT 
    partition,
    count() AS parts,
    sum(rows) AS sum_rows,
    max(rows) AS max_rows,
    min(rows) AS min_rows,
    avg(rows) AS avg_rows,
    sum(bytes) AS total_bytes,
    formatReadableSize(sum(bytes)) AS size,
    formatReadableSize(sum(data_compressed_bytes)) AS compressed_size,
    formatReadableSize(sum(data_uncompressed_bytes)) AS uncompressed_size
FROM system.parts
WHERE table = 'raw_trades' AND active
GROUP BY partition
ORDER BY partition;

-- Inspect disk location for each part
SELECT
    name,
    disk_name
FROM system.parts
WHERE (table = 'raw_trades') AND (active = 1);

-- Disk Usage by Table
SELECT
    hostName(),
    database,
    table,
    sum(rows) AS rows,
    formatReadableSize(sum(bytes_on_disk)) AS total_bytes_on_disk,
    formatReadableSize(sum(data_compressed_bytes)) AS total_data_compressed_bytes,
    formatReadableSize(sum(data_uncompressed_bytes)) AS total_data_uncompressed_bytes,
    round(sum(data_compressed_bytes) / sum(data_uncompressed_bytes), 3) AS compression_ratio
FROM system.parts
WHERE database != 'system'
GROUP BY
    hostName(),
    database,
    table
ORDER BY sum(bytes_on_disk) DESC;

-- Parts moved
SELECT *
FROM system.moves;

-- Merge Errors
SELECT
    event_date,
    event_type,
    table,
    error AS error_code,
    errorCodeToName(error) AS error_code_name,
    count() as c
FROM system.part_log
WHERE (error_code != 0) AND (event_date > (now() - toIntervalMonth(1)))
GROUP BY
    event_date,
    event_type,
    error,
    table
ORDER BY
    event_date DESC,
    event_type ASC,
    table ASC,
    error ASC;
