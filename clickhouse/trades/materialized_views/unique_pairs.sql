CREATE MATERIALIZED VIEW IF NOT EXISTS trades.unique_pairs
ENGINE = ReplacingMergeTree
ORDER BY (base, quote)
POPULATE
AS SELECT
    base,
    quote
FROM trades.raw_trades
GROUP BY base, quote;