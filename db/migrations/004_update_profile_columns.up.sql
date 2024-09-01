ALTER TABLE public.profile
ADD COLUMN IF NOT EXISTS avatar_path TEXT UNIQUE
    CONSTRAINT max_len_avatar_path CHECK (LENGTH(avatar_path) <= 256);
