-- Recipe Variations Table
CREATE TABLE IF NOT EXISTS recipe_variations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES users(id),
    markdown_content TEXT NOT NULL,
    prep_time_minutes INT,
    cook_time_minutes INT,
    servings INT,
    difficulty VARCHAR(20),
    notes TEXT,
    is_published BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Ensure one variation per author per recipe
    CONSTRAINT unique_author_variation UNIQUE (recipe_id, author_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_variations_recipe ON recipe_variations(recipe_id);
CREATE INDEX IF NOT EXISTS idx_variations_author ON recipe_variations(author_id);
CREATE INDEX IF NOT EXISTS idx_variations_published ON recipe_variations(is_published, created_at DESC);
