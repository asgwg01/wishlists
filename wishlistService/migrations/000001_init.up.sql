CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS wishlists(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL,
    title VARCHAR(80) NOT NULL,
    description TEXT,
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS items(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wishlist_id UUID NOT NULL REFERENCES wishlists(id) ON DELETE CASCADE,
    name VARCHAR(80) NOT NULL,
    description TEXT,
    image_url TEXT,
    product_url TEXT,
    price INTEGER CHECK (price IS NULL OR price >= 0),
    booked_by UUID,
    booked_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT fk_wishlist FOREIGN KEY (wishlist_id) REFERENCES wishlists(id) ON DELETE CASCADE
);

CREATE INDEX wishlist_items_idx ON items(wishlist_id);
CREATE INDEX booked_items_idx ON items(booked_by);

CREATE INDEX owner_wishlists_idx ON wishlists(owner_id);
CREATE INDEX public_wishlists_idx ON wishlists(is_public);

-- Индекс для поиска забронированных предметов пользователя
CREATE INDEX booked_wishlist_items_idx ON items(booked_by, wishlist_id);
