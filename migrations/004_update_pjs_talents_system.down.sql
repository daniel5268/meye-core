-- Recreate the ENUM types
CREATE TYPE basic_talent_type AS ENUM ('physical', 'mental', 'coordination', 'energy');
CREATE TYPE special_talent_type AS ENUM ('physical', 'mental', 'energy');

-- Add back the old talent columns
ALTER TABLE pjs
ADD COLUMN basic_talent basic_talent_type,
ADD COLUMN special_talent special_talent_type;

-- Attempt to migrate data back (best effort - data loss may occur for multiple talents)
-- For basic_talent, prioritize in order: physical, mental, coordination, energy
UPDATE pjs SET basic_talent = 'physical' WHERE is_physical_talented = true;
UPDATE pjs SET basic_talent = 'mental' WHERE is_mental_talented = true AND basic_talent IS NULL;
UPDATE pjs SET basic_talent = 'coordination' WHERE is_coordination_talented = true AND basic_talent IS NULL;
UPDATE pjs SET basic_talent = 'energy' WHERE is_energy_skills_talented = true AND basic_talent IS NULL;

-- For special_talent, prioritize in order: physical, mental, energy
UPDATE pjs SET special_talent = 'physical' WHERE is_physical_skills_talented = true;
UPDATE pjs SET special_talent = 'mental' WHERE is_mental_skills_talented = true AND special_talent IS NULL;
UPDATE pjs SET special_talent = 'energy' WHERE is_energy_skills_talented = true AND special_talent IS NULL;

-- Set defaults for any remaining NULL values (shouldn't happen in practice)
UPDATE pjs SET basic_talent = 'physical' WHERE basic_talent IS NULL;
UPDATE pjs SET special_talent = 'physical' WHERE special_talent IS NULL;

-- Make columns NOT NULL
ALTER TABLE pjs
ALTER COLUMN basic_talent SET NOT NULL,
ALTER COLUMN special_talent SET NOT NULL;

-- Remove the new talent columns
ALTER TABLE pjs
DROP COLUMN is_physical_talented,
DROP COLUMN is_mental_talented,
DROP COLUMN is_coordination_talented,
DROP COLUMN is_physical_skills_talented,
DROP COLUMN is_mental_skills_talented,
DROP COLUMN is_energy_skills_talented,
DROP COLUMN is_energy_talented;