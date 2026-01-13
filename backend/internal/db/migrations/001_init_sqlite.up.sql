-- Users (SQLite-compatible)
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Categories (SQLite-compatible)
CREATE TABLE IF NOT EXISTS categories (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    icon TEXT,
    description TEXT,
    order_index INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tags (SQLite-compatible)
CREATE TABLE IF NOT EXISTS tags (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    color TEXT DEFAULT '#6366f1',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Recipe Groups (SQLite-compatible)
CREATE TABLE IF NOT EXISTS recipe_groups (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    description TEXT,
    icon TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Recipes (SQLite-compatible)
CREATE TABLE IF NOT EXISTS recipes (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    markdown_content TEXT NOT NULL,
    author_id TEXT REFERENCES users(id),
    category_id TEXT REFERENCES categories(id),
    description TEXT,
    prep_time_minutes INTEGER,
    cook_time_minutes INTEGER,
    servings INTEGER,
    difficulty TEXT,
    featured_image_path TEXT,
    is_published INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP
);

-- Recipe Tags (Many-to-Many)
CREATE TABLE IF NOT EXISTS recipe_tags (
    recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
    tag_id TEXT REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (recipe_id, tag_id)
);

-- Recipe Groups (Many-to-Many)
CREATE TABLE IF NOT EXISTS recipe_groupings (
    group_id TEXT REFERENCES recipe_groups(id) ON DELETE CASCADE,
    recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
    order_index INTEGER DEFAULT 0,
    PRIMARY KEY (group_id, recipe_id)
);

-- Recipe Images (SQLite-compatible)
CREATE TABLE IF NOT EXISTS recipe_images (
    id TEXT PRIMARY KEY,
    recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    webp_path TEXT,
    thumbnail_path TEXT,
    caption TEXT,
    order_index INTEGER DEFAULT 0,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Share Codes (SQLite-compatible)
CREATE TABLE IF NOT EXISTS share_codes (
    id TEXT PRIMARY KEY,
    recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
    code TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP,
    max_uses INTEGER,
    use_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User Invites (SQLite-compatible)
CREATE TABLE IF NOT EXISTS user_invites (
    id TEXT PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    email TEXT,
    role TEXT DEFAULT 'user',
    created_by TEXT REFERENCES users(id),
    expires_at TIMESTAMP,
    used_at TIMESTAMP,
    used_by TEXT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- App Settings (SQLite-compatible - JSONB to TEXT)
CREATE TABLE IF NOT EXISTS app_settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Instance Sharing (SQLite-compatible)
CREATE TABLE IF NOT EXISTS shared_recipes (
    id TEXT PRIMARY KEY,
    recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
    source_instance_url TEXT,
    source_recipe_id TEXT,
    signature TEXT NOT NULL,
    imported_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Users (SQLite-compatible)
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY DEFAULT (gen_random_uuid()),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Categories (SQLite-compatible)
CREATE TABLE IF NOT EXISTS categories (
    id TEXT PRIMARY KEY DEFAULT (gen_random_uuid()),
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    icon TEXT,
    description TEXT,
    order_index INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tags (SQLite-compatible)
CREATE TABLE IF NOT EXISTS tags (
    id TEXT PRIMARY KEY DEFAULT (gen_random_uuid()),
    name TEXT UNIQUE NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    color TEXT DEFAULT '#6366f1',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Recipe Groups (SQLite-compatible)
CREATE TABLE IF NOT EXISTS recipe_groups (
    id TEXT PRIMARY KEY DEFAULT (gen_random_uuid()),
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    description TEXT,
    icon TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Recipes (SQLite-compatible)
CREATE TABLE IF NOT EXISTS recipes (
    id TEXT PRIMARY KEY DEFAULT (gen_random_uuid()),
    title TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    markdown_content TEXT NOT NULL,
    author_id TEXT REFERENCES users(id),
    category_id TEXT REFERENCES categories(id),
    description TEXT,
    prep_time_minutes INTEGER,
    cook_time_minutes INTEGER,
    servings INTEGER,
    difficulty TEXT,
    featured_image_path TEXT,
    is_published INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP
);

-- Recipe Tags (Many-to-Many)
CREATE TABLE IF NOT EXISTS recipe_tags (
    recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
    tag_id TEXT REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (recipe_id, tag_id)
);

-- Recipe Groups (Many-to-Many)
CREATE TABLE IF NOT EXISTS recipe_groupings (
    group_id TEXT REFERENCES recipe_groups(id) ON DELETE CASCADE,
    recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
    order_index INTEGER DEFAULT 0,
    PRIMARY KEY (group_id, recipe_id)
);

-- Recipe Images (SQLite-compatible)
CREATE TABLE IF NOT EXISTS recipe_images (
    id TEXT PRIMARY KEY DEFAULT (gen_random_uuid()),
    recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    webp_path TEXT,
    thumbnail_path TEXT,
    caption TEXT,
    order_index INTEGER DEFAULT 0,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Share Codes (SQLite-compatible)
CREATE TABLE IF NOT EXISTS share_codes (
    id TEXT PRIMARY KEY DEFAULT (gen_random_uuid()),
    recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
    code TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP,
    max_uses INTEGER,
    use_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User Invites (SQLite-compatible)
CREATE TABLE IF NOT EXISTS user_invites (
    id TEXT PRIMARY KEY DEFAULT (gen_random_uuid()),
    code TEXT UNIQUE NOT NULL,
    email TEXT,
    role TEXT DEFAULT 'user',
    created_by TEXT REFERENCES users(id),
    expires_at TIMESTAMP,
    used_at TIMESTAMP,
    used_by TEXT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- App Settings (SQLite-compatible - JSONB to TEXT)
CREATE TABLE IF NOT EXISTS app_settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Instance Sharing (SQLite-compatible)
CREATE TABLE IF NOT EXISTS shared_recipes (
    id TEXT PRIMARY KEY DEFAULT (gen_random_uuid()),
    recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
    source_instance_url TEXT,
    source_recipe_id TEXT,
    signature TEXT NOT NULL,
    imported_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_recipes_author ON recipes(author_id);
CREATE INDEX IF NOT EXISTS idx_recipes_category ON recipes(category_id);
CREATE INDEX IF NOT EXISTS idx_recipes_published ON recipes(is_published, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_recipe_groups_slug ON recipe_groups(slug);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
