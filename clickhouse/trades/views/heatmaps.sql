CREATE OR REPLACE VIEW trades.heatmap_view AS
WITH params AS (
    SELECT
        {base:String} AS base,
        {quote:String} AS quote,
        {num_bins:UInt32} AS num_bins,
        {start_time:DateTime} AS start_time,
        {end_time:DateTime} AS end_time,
        {threshold:Float64} AS threshold
),
inliers AS (
    SELECT *
    FROM trades.inliers_view(
        threshold = (SELECT threshold FROM params),
        start_time = (SELECT start_time FROM params),
        end_time = (SELECT end_time FROM params)
    )
    WHERE base = (SELECT base FROM params) AND quote = (SELECT quote FROM params)
),
stats AS (
    SELECT
        min(volume_base) AS x_min,
        max(volume_base) AS x_max,
        min(price) AS y_min,
        max(price) AS y_max
    FROM inliers
),
heatmap AS (
    SELECT
        x_bin,
        y_bin,
        count() AS count,
        sum(volume_base) AS volume
    FROM (
        SELECT
            volume_base,
            price,
            floor((volume_base - x_min) / ((x_max - x_min) / num_bins)) + 1 AS x_bin,
            floor((price - y_min) / ((y_max - y_min) / num_bins)) + 1 AS y_bin
        FROM inliers, stats, params
    )
    GROUP BY x_bin, y_bin
    ORDER BY x_bin, y_bin
)
SELECT
    x_bin,
    y_bin,
    count,
    volume,
    log1p(count) AS count_log,
    log1p(volume) AS volume_log,
    x_min + (x_bin - 1) * (x_max - x_min) / num_bins AS x_start,
    x_min + x_bin * (x_max - x_min) / num_bins AS x_end,
    y_min + (y_bin - 1) * (y_max - y_min) / num_bins AS y_start,
    y_min + y_bin * (y_max - y_min) / num_bins AS y_end,
    x_min,
    x_max,
    y_min,
    y_max
FROM heatmap, stats, params;

-- Example query to use the heatmap view
SELECT *
FROM trades.heatmap_view(
    base = 'bitcoin',
    quote = 'tether',
    num_bins = 600,
    start_time = subtractHours(now(), 24),
    end_time = now(),
    threshold = 2
);