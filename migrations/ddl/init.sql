DROP TABLE IF EXISTS eg_sw_connection CASCADE;
DROP TABLE IF EXISTS eg_sw_service CASCADE;
DROP TABLE IF EXISTS eg_sw_connectionholder CASCADE;
DROP TABLE IF EXISTS eg_sw_plumberinfo CASCADE;
DROP TABLE IF EXISTS eg_sw_roadcuttinginfo CASCADE;

CREATE TABLE eg_sw_connection (
    id character varying(256) PRIMARY KEY,
    tenantid character varying(256) NOT NULL,
    applicationno character varying(256) NOT NULL UNIQUE,
    connectionno character varying(256),
    oldconnectionno character varying(256),
    property_id character varying(256) NOT NULL,
    connectiontype character varying(256),
    roadtype character varying(256),
    roadcuttingarea INTEGER,
    applicationstatus character varying(256) NOT NULL,
    status character varying(256) NOT NULL,
    action character varying(256),
    channel character varying(256),
    createdby character varying(256),
    lastmodifiedby character varying(256),
    createdtime BIGINT,
    lastmodifiedtime BIGINT,
    applicationtype character varying(256),
    dateEffectiveFrom BIGINT,
    locality character varying(256),
    isoldapplication BOOLEAN,
    additionaldetails JSONB,
    isDisconnectionTemporary BOOLEAN,
    disconnectionReason character varying(256),

    -- Go service compatibility columns
    connection_holders JSONB DEFAULT '[]'::jsonb,
    plumber_info JSONB DEFAULT '[]'::jsonb,
    road_cutting_info JSONB DEFAULT '[]'::jsonb
);

CREATE TABLE eg_sw_service (
    connection_id character varying(256) REFERENCES eg_sw_connection(id) ON DELETE CASCADE,
    connectionExecutionDate BIGINT,
    noOfWaterClosets INTEGER,
    noOfToilets INTEGER,
    connectiontype character varying(256),
    proposedWaterClosets INTEGER,
    proposedToilets INTEGER,
    appCreatedDate BIGINT,
    createdby character varying(256),
    lastmodifiedby character varying(256),
    createdtime BIGINT,
    lastmodifiedtime BIGINT,
    disconnectionExecutionDate BIGINT
);

CREATE TABLE eg_sw_connectionholder (
    tenantid character varying(256),
    connectionid character varying(256) UNIQUE REFERENCES eg_sw_connection(id) ON DELETE CASCADE,
    userid character varying(256),
    status character varying(256),
    isprimaryholder BOOLEAN,
    connectionholdertype character varying(256),
    holdershippercentage NUMERIC,
    relationship character varying(256),
    createdby character varying(256),
    createdtime BIGINT,
    lastmodifiedby character varying(256),
    lastmodifiedtime BIGINT
);

CREATE TABLE eg_sw_plumberinfo (
    id character varying(256) PRIMARY KEY,
    tenantId character varying(256),
    name character varying(256),
    licenseno character varying(256),
    mobilenumber character varying(256),
    gender character varying(256),
    fatherorhusbandname character varying(256),
    correspondenceaddress character varying(256),
    relationship character varying(256),
    createdBy character varying(256),
    lastModifiedBy character varying(256),
    createdTime BIGINT,
    lastModifiedTime BIGINT,
    swid character varying(256) REFERENCES eg_sw_connection(id) ON DELETE CASCADE
);

CREATE TABLE eg_sw_roadcuttinginfo (
    id character varying(256) PRIMARY KEY,
    tenantId character varying(256),
    swid character varying(256) REFERENCES eg_sw_connection(id) ON DELETE CASCADE,
    active BOOLEAN,
    roadtype character varying(256),
    roadcuttingarea INTEGER,
    createdBy character varying(256),
    lastModifiedBy character varying(256),
    createdTime BIGINT,
    lastModifiedTime BIGINT
);
