CREATE TYPE pj_type AS ENUM ('human', 'supernatural');
CREATE TYPE basic_talent_type AS ENUM ('physical', 'mental', 'coordination', 'energy');
CREATE TYPE special_talent_type AS ENUM ('physical', 'mental', 'energy');

CREATE TABLE pjs (
    id VARCHAR(255) PRIMARY KEY,
    campaign_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    weight INTEGER NOT NULL,
    height INTEGER NOT NULL,
    age INTEGER NOT NULL,
    look INTEGER NOT NULL,
    charisma INTEGER NOT NULL,
    villainy INTEGER NOT NULL,
    heroism INTEGER NOT NULL,
    pj_type pj_type NOT NULL,
    basic_talent basic_talent_type NOT NULL,
    special_talent special_talent_type NOT NULL,

    -- Basic Stats - Physical
    strength INTEGER NOT NULL,
    agility INTEGER NOT NULL,
    speed INTEGER NOT NULL,
    resistance INTEGER NOT NULL,

    -- Basic Stats - Mental
    inteligence INTEGER NOT NULL,
    wisdom INTEGER NOT NULL,
    concentration INTEGER NOT NULL,
    will INTEGER NOT NULL,

    -- Basic Stats - Coordination
    precision INTEGER NOT NULL,
    calculation INTEGER NOT NULL,
    range INTEGER NOT NULL,
    reflexes INTEGER NOT NULL,

    -- Basic Stats - Life
    life INTEGER NOT NULL,

    -- Special Stats - Physical
    empowerment INTEGER NOT NULL,
    vital_control INTEGER NOT NULL,

    -- Special Stats - Mental
    ilusion INTEGER NOT NULL,
    mental_control INTEGER NOT NULL,

    -- Special Stats - Energy
    object_handling INTEGER NOT NULL,
    energy_handling INTEGER NOT NULL,
    energy_tank INTEGER NOT NULL,

    -- Supernatural Stats (nullable for human PJs)
    supernatural_stats JSONB,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_pj_campaign
        FOREIGN KEY (campaign_id)
        REFERENCES campaigns(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_pj_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_pjs_campaign_id ON pjs(campaign_id);
CREATE INDEX idx_pjs_user_id ON pjs(user_id);
CREATE INDEX idx_pjs_pj_type ON pjs(pj_type);
