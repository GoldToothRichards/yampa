-- Valid trades
CREATE OR REPLACE VIEW trades.clean_trades_view AS
SELECT *
FROM trades.raw_trades
WHERE source != ''
  AND base != ''
  AND quote != ''
  AND exchange != ''
  AND price > 0
  AND volume_base > 0
  AND volume_quote > 0
  AND base != quote;

-- Invalid trades
CREATE OR REPLACE VIEW trades.rejected_trades_view AS
SELECT *
FROM trades.raw_trades
WHERE source = ''
   OR base = ''
   OR quote = ''
   OR exchange = ''
   OR price <= 0
   OR volume_base <= 0
   OR volume_quote <= 0
   OR base = quote;

-- Compare counts
SELECT 
    (SELECT COUNT(*) FROM trades.raw_trades) AS raw_trades_count,
    (SELECT COUNT(*) FROM trades.clean_trades_view) AS clean_trades_count,
    (SELECT COUNT(*) FROM trades.rejected_trades_view) as rejected_trades_count