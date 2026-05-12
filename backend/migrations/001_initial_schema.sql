-- Rygell Dashboard - Initial Schema Reference
-- Runtime schema is managed by GORM AutoMigrate in backend/internal/database.
-- Keep this file aligned with backend/internal/models; do not treat it as the migration runner.

-- ============================================================
-- AUTH
-- ============================================================

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(250) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- ============================================================
-- MASTER DATA
-- ============================================================

CREATE TABLE IF NOT EXISTS mills (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_mills_code ON mills(code);
CREATE INDEX IF NOT EXISTS idx_mills_deleted_at ON mills(deleted_at);

CREATE TABLE IF NOT EXISTS vendors (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    tax_id VARCHAR(50),
    status VARCHAR(50) DEFAULT 'active',
    address TEXT,
    contact_person VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_vendors_code ON vendors(code);
CREATE INDEX IF NOT EXISTS idx_vendors_deleted_at ON vendors(deleted_at);

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at);

CREATE TABLE IF NOT EXISTS zones (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_zones_deleted_at ON zones(deleted_at);

CREATE TABLE IF NOT EXISTS mots (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_mots_deleted_at ON mots(deleted_at);

CREATE TABLE IF NOT EXISTS uoms (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_uoms_deleted_at ON uoms(deleted_at);

-- ============================================================
-- CONTRACTS
-- ============================================================
-- GORM maps CostPerKGKM to cost_per_kgkm, not cost_per_kg_km.
-- GORM maps CostIDR and RunningCostIDR to cost_id_r and running_cost_id_r.

CREATE TABLE IF NOT EXISTS contract_dedicated_fixes (
    id BIGSERIAL PRIMARY KEY,
    mill_id BIGINT NOT NULL REFERENCES mills(id),
    vendor_id BIGINT NOT NULL REFERENCES vendors(id),
    product_id BIGINT REFERENCES products(id),
    license_plate VARCHAR(50),
    spk_number VARCHAR(100),
    area_category VARCHAR(100),
    proposal_cfas VARCHAR(100),
    fa_number VARCHAR(100),
    mot_id BIGINT REFERENCES mots(id),
    uom_id BIGINT REFERENCES uoms(id),
    validity_start TIMESTAMPTZ,
    validity_end TIMESTAMPTZ,
    cost_jan NUMERIC(15,2) DEFAULT 0,
    cost_feb NUMERIC(15,2) DEFAULT 0,
    cost_mar NUMERIC(15,2) DEFAULT 0,
    cost_apr NUMERIC(15,2) DEFAULT 0,
    cost_may NUMERIC(15,2) DEFAULT 0,
    cost_jun NUMERIC(15,2) DEFAULT 0,
    fix_cost NUMERIC(15,2) DEFAULT 0,
    cargo_carried NUMERIC(15,2) DEFAULT 0,
    unit_cost NUMERIC(15,2) DEFAULT 0,
    cost_per_kg NUMERIC(15,4) DEFAULT 0,
    cost_per_kgkm NUMERIC(15,4) DEFAULT 0,
    distributed_cost NUMERIC(15,2) DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_fixes_mill_id ON contract_dedicated_fixes(mill_id);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_fixes_vendor_id ON contract_dedicated_fixes(vendor_id);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_fixes_product_id ON contract_dedicated_fixes(product_id);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_fixes_mot_id ON contract_dedicated_fixes(mot_id);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_fixes_uom_id ON contract_dedicated_fixes(uom_id);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_fixes_spk_number ON contract_dedicated_fixes(spk_number);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_fixes_deleted_at ON contract_dedicated_fixes(deleted_at);

CREATE TABLE IF NOT EXISTS contract_dedicated_vars (
    id BIGSERIAL PRIMARY KEY,
    mill_id BIGINT NOT NULL REFERENCES mills(id),
    vendor_id BIGINT NOT NULL REFERENCES vendors(id),
    product_id BIGINT REFERENCES products(id),
    origin_zone_id BIGINT REFERENCES zones(id),
    dest_zone_id BIGINT REFERENCES zones(id),
    mot_id BIGINT REFERENCES mots(id),
    uom_id BIGINT REFERENCES uoms(id),
    spk_number VARCHAR(100),
    area_category VARCHAR(100),
    proposal_cfas VARCHAR(100),
    fa_number VARCHAR(100),
    distance NUMERIC(10,2) DEFAULT 0,
    validity_start TIMESTAMPTZ,
    validity_end TIMESTAMPTZ,
    payload NUMERIC(15,2) DEFAULT 0,
    cost_id_r NUMERIC(15,2) DEFAULT 0,
    cost_per_kg NUMERIC(15,4) DEFAULT 0,
    cost_per_kgkm NUMERIC(15,4) DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_vars_mill_id ON contract_dedicated_vars(mill_id);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_vars_vendor_id ON contract_dedicated_vars(vendor_id);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_vars_product_id ON contract_dedicated_vars(product_id);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_vars_origin_zone_id ON contract_dedicated_vars(origin_zone_id);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_vars_dest_zone_id ON contract_dedicated_vars(dest_zone_id);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_vars_mot_id ON contract_dedicated_vars(mot_id);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_vars_uom_id ON contract_dedicated_vars(uom_id);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_vars_spk_number ON contract_dedicated_vars(spk_number);
CREATE INDEX IF NOT EXISTS idx_contract_dedicated_vars_deleted_at ON contract_dedicated_vars(deleted_at);

CREATE TABLE IF NOT EXISTS contract_oncalls (
    id BIGSERIAL PRIMARY KEY,
    mill_id BIGINT NOT NULL REFERENCES mills(id),
    vendor_id BIGINT NOT NULL REFERENCES vendors(id),
    product_id BIGINT REFERENCES products(id),
    origin_zone_id BIGINT REFERENCES zones(id),
    dest_zone_id BIGINT REFERENCES zones(id),
    mot_id BIGINT REFERENCES mots(id),
    uom_id BIGINT REFERENCES uoms(id),
    spk_number VARCHAR(100),
    area_category VARCHAR(100),
    proposal_cfas VARCHAR(100),
    fa_number VARCHAR(100),
    validity_start TIMESTAMPTZ,
    validity_end TIMESTAMPTZ,
    distance NUMERIC(10,2) DEFAULT 0,
    payload NUMERIC(15,2) DEFAULT 0,
    loading_cost NUMERIC(15,2) DEFAULT 0,
    unloading_cost NUMERIC(15,2) DEFAULT 0,
    cost_id_r NUMERIC(15,2) DEFAULT 0,
    cost_per_kg NUMERIC(15,4) DEFAULT 0,
    cost_per_ton NUMERIC(15,4) DEFAULT 0,
    running_cost_id_r NUMERIC(15,4) DEFAULT 0,
    running_cost_usd NUMERIC(15,4) DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_contract_oncalls_mill_id ON contract_oncalls(mill_id);
CREATE INDEX IF NOT EXISTS idx_contract_oncalls_vendor_id ON contract_oncalls(vendor_id);
CREATE INDEX IF NOT EXISTS idx_contract_oncalls_product_id ON contract_oncalls(product_id);
CREATE INDEX IF NOT EXISTS idx_contract_oncalls_origin_zone_id ON contract_oncalls(origin_zone_id);
CREATE INDEX IF NOT EXISTS idx_contract_oncalls_dest_zone_id ON contract_oncalls(dest_zone_id);
CREATE INDEX IF NOT EXISTS idx_contract_oncalls_mot_id ON contract_oncalls(mot_id);
CREATE INDEX IF NOT EXISTS idx_contract_oncalls_uom_id ON contract_oncalls(uom_id);
CREATE INDEX IF NOT EXISTS idx_contract_oncalls_spk_number ON contract_oncalls(spk_number);
CREATE INDEX IF NOT EXISTS idx_contract_oncalls_deleted_at ON contract_oncalls(deleted_at);

-- ============================================================
-- AUDIT LOGS
-- ============================================================

CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    entity_type VARCHAR(100) NOT NULL,
    entity_id BIGINT NOT NULL,
    action VARCHAR(50) NOT NULL,
    changed_by VARCHAR(255),
    agreement_note TEXT,
    old_data JSONB,
    new_data JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_entity ON audit_logs(entity_type, entity_id);
