CREATE TABLE sessions (
    id VARCHAR(255) PRIMARY KEY,
    campaign_id VARCHAR(255) NOT NULL,
    summary TEXT NOT NULL,
    xp_assignations JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL,
    
    CONSTRAINT fk_sessions_campaign 
        FOREIGN KEY (campaign_id) 
        REFERENCES campaigns(id) 
        ON DELETE CASCADE
);

CREATE INDEX idx_sessions_campaign_id ON sessions(campaign_id);
CREATE INDEX idx_sessions_created_at ON sessions(created_at DESC);

-- Índice GIN para búsquedas en JSONB (buscar por pj_id)
CREATE INDEX idx_sessions_xp_assignations_gin ON sessions USING GIN (xp_assignations);
