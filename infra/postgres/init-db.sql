-- shortener-service-db
CREATE TABLE short_links (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    short_code VARCHAR(20) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expire_at TIMESTAMP,
    click_limit INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true
);
CREATE TABLE tenants (
    id UUID PRIMARY KEY,
    name VARCHAR(100),
    owner_id UUID REFERENCES users(id),
    custom_domain VARCHAR(255),
    plan VARCHAR(50) DEFAULT 'free',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- auth-service -> user-service-db
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
PARTITION BY toYYYYMM(timestamp)
ORDER BY (short_code, timestamp);
--billing-service-db 
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY,
    tenant_id UUID REFERENCES tenants(id),
    stripe_customer_id TEXT,
    stripe_subscription_id TEXT,
    plan VARCHAR(50),
    is_active BOOLEAN,
    started_at TIMESTAMP,
    ended_at TIMESTAMP
);
