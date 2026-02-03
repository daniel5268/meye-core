CREATE TYPE invitation_state AS ENUM ('pending', 'accepted');

CREATE TABLE campaign_invitations (
    id VARCHAR(255) PRIMARY KEY,
    campaign_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    state invitation_state NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_campaign
        FOREIGN KEY (campaign_id)
        REFERENCES campaigns(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_campaign_invitations_campaign_id ON campaign_invitations(campaign_id);
CREATE INDEX idx_campaign_invitations_user_id ON campaign_invitations(user_id);
CREATE INDEX idx_campaign_invitations_state ON campaign_invitations(state);