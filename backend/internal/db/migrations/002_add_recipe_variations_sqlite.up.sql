-- Recipe Variations Table (SQLite compatible)
CREATE TABLE IF NOT EXISTS recipe_variations (
    id TEXT PRIMARY KEY,
    recipe_id TEXT NOT NULL,
    author_id TEXT NOT NULL,
    markdown_content TEXT NOT NULL,
    prep_time_minutes INTEGER,
    cook_time_minutes INTEGER,
    servings INTEGER,
    difficulty TEXT,
    notes TEXT,
    is_published INTEGER DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id),

    -- Ensure one variation per author per recipe
    CONSTRAINT unique_author_variation UNIQUE (recipe_id, author_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_variations_recipe ON recipe_variations(recipe_id);
CREATE INDEX IF NOT EXISTS idx_variations_author ON recipe_variations(author_id);
CREATE INDEX IF NOT EXISTS idx_variations_published ON recipe_variations(is_published, created_at DESC);
