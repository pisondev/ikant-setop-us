DROP INDEX IF EXISTS idx_stock_movements_created_at;
DROP INDEX IF EXISTS idx_stock_movements_stock_batch_id;
DROP INDEX IF EXISTS idx_stock_out_items_stock_batch_id;
DROP INDEX IF EXISTS idx_stock_out_items_stock_out_id;
DROP INDEX IF EXISTS idx_stock_outs_out_at;
DROP INDEX IF EXISTS idx_stock_batches_cold_storage_id;
DROP INDEX IF EXISTS idx_stock_batches_overall_fifo;
DROP INDEX IF EXISTS idx_stock_batches_fifo;

DROP TABLE IF EXISTS stock_movements;
DROP TABLE IF EXISTS stock_out_items;
DROP TABLE IF EXISTS stock_outs;
DROP TABLE IF EXISTS stock_batches;
DROP TABLE IF EXISTS cold_storages;
DROP TABLE IF EXISTS fish_types;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS stock_movement_type;
DROP TYPE IF EXISTS stock_status;
DROP TYPE IF EXISTS fish_quality;
DROP TYPE IF EXISTS user_role;
