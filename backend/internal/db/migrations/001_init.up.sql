-- Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Categories
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    icon VARCHAR(50),
    description TEXT,
    order_index INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Tags
CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    color VARCHAR(7) DEFAULT '#6366f1',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Recipe Groups
CREATE TABLE IF NOT EXISTS recipe_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Recipes
CREATE TABLE IF NOT EXISTS recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    markdown_content TEXT NOT NULL,
    author_id UUID REFERENCES users(id),
    category_id UUID REFERENCES categories(id),
    description TEXT,
    prep_time_minutes INT,
    cook_time_minutes INT,
    servings INT,
    difficulty VARCHAR(20),
    featured_image_path VARCHAR(500),
    is_published BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    published_at TIMESTAMP
);

-- Recipe Tags (Many-to-Many)
CREATE TABLE IF NOT EXISTS recipe_tags (
    recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (recipe_id, tag_id)
);

-- Recipe Groups (Many-to-Many)
CREATE TABLE IF NOT EXISTS recipe_groupings (
    group_id UUID REFERENCES recipe_groups(id) ON DELETE CASCADE,
    recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
    order_index INT DEFAULT 0,
    PRIMARY KEY (group_id, recipe_id)
);

-- Recipe Images
CREATE TABLE IF NOT EXISTS recipe_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
    file_path VARCHAR(500) NOT NULL,
    webp_path VARCHAR(500),
    thumbnail_path VARCHAR(500),
    caption TEXT,
    order_index INT DEFAULT 0,
    uploaded_at TIMESTAMP DEFAULT NOW()
);

-- Share Codes
CREATE TABLE IF NOT EXISTS share_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
    code VARCHAR(20) UNIQUE NOT NULL,
    expires_at TIMESTAMP,
    max_uses INT DEFAULT NULL,
    use_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- User Invites
CREATE TABLE IF NOT EXISTS user_invites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(255),
    role VARCHAR(50) DEFAULT 'user',
    created_by UUID REFERENCES users(id),
    expires_at TIMESTAMP,
    used_at TIMESTAMP,
    used_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- App Settings
CREATE TABLE IF NOT EXISTS app_settings (
    key VARCHAR(100) PRIMARY KEY,
    value JSONB NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Instance Sharing
CREATE TABLE IF NOT EXISTS shared_recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
    source_instance_url VARCHAR(500),
    source_recipe_id UUID,
    signature VARCHAR(255) NOT NULL,
    imported_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_recipes_author ON recipes(author_id);
CREATE INDEX IF NOT EXISTS idx_recipes_category ON recipes(category_id);
CREATE INDEX IF NOT EXISTS idx_recipes_published ON recipes(is_published, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_recipe_groups_slug ON recipe_groups(slug);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
