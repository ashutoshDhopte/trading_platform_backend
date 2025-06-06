-- RESET SCRIPT
TRUNCATE TABLE orders;
TRUNCATE TABLE holdings;
UPDATE stocks SET current_price_cents = opening_price_cents;
UPDATE users SET cash_balance_cents = 10000000;