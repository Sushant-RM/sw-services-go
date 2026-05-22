CREATE TABLE IF NOT EXISTS eg_sw_connection (
    id UUID PRIMARY KEY,

    tenant_id TEXT NOT NULL,

    application_number TEXT NOT NULL UNIQUE,

    connection_no TEXT,

    property_id TEXT NOT NULL,

    connection_type TEXT NOT NULL,

    road_type TEXT NOT NULL,

    road_cutting_area INTEGER NOT NULL,

    application_status TEXT NOT NULL,

    channel TEXT NOT NULL,

    connection_holders JSONB DEFAULT '[]'::jsonb,

    plumber_info JSONB DEFAULT '[]'::jsonb,

    road_cutting_info JSONB DEFAULT '[]'::jsonb,

    created_by TEXT,

    last_modified_by TEXT,

    created_time BIGINT,

    last_modified_time BIGINT
);
