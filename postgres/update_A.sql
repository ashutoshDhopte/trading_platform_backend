
ALTER TABLE stocks ADD current_price_cents BIGINT;

UPDATE stocks SET current_price_cents = opening_price_cents;