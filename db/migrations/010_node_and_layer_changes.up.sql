ALTER TABLE public.layer
DROP CONSTRAINT IF EXISTS positive_num;

ALTER TABLE public.node
ADD COLUMN IF NOT EXISTS gender genders
    NOT NULL
    DEFAULT 'Мужской';

ALTER TABLE public.node
ADD COLUMN IF NOT EXISTS is_spouse BOOLEAN
    NOT NULL
    DEFAULT FALSE;