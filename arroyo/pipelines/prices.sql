CREATE VIEW trades as (
    SELECT
        extract_json_string(value, '$.source') as source,
        extract_json_string(value, '$.base') as base,
        extract_json_string(value, '$.quote') as quote,
        extract_json_string(value, '$.exchange') as exchange,
        CAST(array_element(extract_json(value, '$.volume_base'), 1) AS DOUBLE) as volume_base,
        CAST(array_element(extract_json(value, '$.volume_quote'), 1) AS DOUBLE) as volume_quote,
        CAST(array_element(extract_json(value, '$.price'), 1) AS DOUBLE) as price,
        extract_json_string(value, '$.timestamp') as 'timestamp'
    FROM trades_source
);

CREATE VIEW prices as (
    SELECT
        HOP(interval '1 second', interval '24 hour') as window,
        trades.source as source,
        trades.base as base,
        trades.quote as quote,
        trades.exchange as exchange,
        SUM(trades.volume_base) AS total_volume_base,
        SUM(trades.volume_quote) AS total_volume_quote
    FROM trades
    GROUP BY
        window,
        trades.source,
        trades.base,
        trades.quote,
        trades.exchange
);

INSERT INTO prices_sink
SELECT
    prices.source as source,
    prices.base as base,
    prices.quote as quote,
    prices.exchange as exchange,
    prices.total_volume_base as total_volume_base,
    prices.total_volume_quote as total_volume_quote,
    prices.total_volume_quote / prices.total_volume_base as vwap,
    prices.window as window
FROM prices;
