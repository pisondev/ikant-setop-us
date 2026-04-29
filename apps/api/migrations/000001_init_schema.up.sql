CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_role AS ENUM ('owner', 'warehouse_admin');
CREATE TYPE fish_quality AS ENUM ('baik', 'sedang', 'buruk');
CREATE TYPE stock_status AS ENUM ('available', 'depleted');
CREATE TYPE stock_movement_type AS ENUM ('in', 'out', 'quality_update', 'location_update', 'adjustment');

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(120) NOT NULL,
    role user_role NOT NULL,
    email varchar(120) UNIQUE,
    password_hash text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE fish_types (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(100) NOT NULL UNIQUE,
    image_url text,
    description text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE cold_storages (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(100) NOT NULL,
    location_label varchar(150),
    description text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE stock_batches (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    fish_type_id uuid NOT NULL REFERENCES fish_types(id),
    cold_storage_id uuid NOT NULL REFERENCES cold_storages(id),
    quality fish_quality NOT NULL,
    initial_weight_kg numeric(10,2) NOT NULL CHECK (initial_weight_kg > 0),
    remaining_weight_kg numeric(10,2) NOT NULL CHECK (remaining_weight_kg >= 0),
    entered_at timestamptz NOT NULL,
    status stock_status NOT NULL DEFAULT 'available',
    notes text,
    created_by uuid REFERENCES users(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CHECK (remaining_weight_kg <= initial_weight_kg)
);

CREATE TABLE stock_outs (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    destination varchar(150) NOT NULL,
    total_weight_kg numeric(10,2) NOT NULL CHECK (total_weight_kg > 0),
    out_at timestamptz NOT NULL,
    notes text,
    created_by uuid REFERENCES users(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE stock_out_items (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    stock_out_id uuid NOT NULL REFERENCES stock_outs(id) ON DELETE CASCADE,
    stock_batch_id uuid NOT NULL REFERENCES stock_batches(id),
    weight_kg numeric(10,2) NOT NULL CHECK (weight_kg > 0),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE stock_movements (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    stock_batch_id uuid NOT NULL REFERENCES stock_batches(id),
    movement_type stock_movement_type NOT NULL,
    weight_kg numeric(10,2) CHECK (weight_kg IS NULL OR weight_kg > 0),
    previous_quality fish_quality,
    new_quality fish_quality,
    previous_cold_storage_id uuid REFERENCES cold_storages(id),
    new_cold_storage_id uuid REFERENCES cold_storages(id),
    description text,
    created_by uuid REFERENCES users(id),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_stock_batches_fifo ON stock_batches (fish_type_id, status, entered_at);
CREATE INDEX idx_stock_batches_overall_fifo ON stock_batches (status, entered_at);
CREATE INDEX idx_stock_batches_cold_storage_id ON stock_batches (cold_storage_id);
CREATE INDEX idx_stock_outs_out_at ON stock_outs (out_at);
CREATE INDEX idx_stock_out_items_stock_out_id ON stock_out_items (stock_out_id);
CREATE INDEX idx_stock_out_items_stock_batch_id ON stock_out_items (stock_batch_id);
CREATE INDEX idx_stock_movements_stock_batch_id ON stock_movements (stock_batch_id);
CREATE INDEX idx_stock_movements_created_at ON stock_movements (created_at DESC);
