-- Add new talent columns as booleans
ALTER TABLE pjs
ADD COLUMN is_physical_talented BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN is_mental_talented BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN is_coordination_talented BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN is_physical_skills_talented BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN is_mental_skills_talented BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN is_energy_skills_talented BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN is_energy_talented BOOLEAN NOT NULL DEFAULT false;

-- Migrate existing data from basic_talent to new columns
UPDATE pjs SET is_physical_talented = true WHERE basic_talent = 'physical';
UPDATE pjs SET is_mental_talented = true WHERE basic_talent = 'mental';
UPDATE pjs SET is_coordination_talented = true WHERE basic_talent = 'coordination';

-- Migrate existing data from special_talent to new columns
UPDATE pjs SET is_physical_skills_talented = true WHERE special_talent = 'physical';
UPDATE pjs SET is_mental_skills_talented = true WHERE special_talent = 'mental';
UPDATE pjs SET is_energy_skills_talented = true WHERE special_talent = 'energy';

-- Note: For energy tank talent, basic_talent='energy' gives cheaper energy tank costs
-- This affects the energy_tank XP calculation, not the energy_skills
UPDATE pjs SET is_energy_talented = true WHERE basic_talent = 'energy';

-- Remove old talent columns
ALTER TABLE pjs
DROP COLUMN basic_talent,
DROP COLUMN special_talent;

-- Drop the unused ENUM types
DROP TYPE basic_talent_type;
DROP TYPE special_talent_type;