ALTER TABLE public.layer
ADD CONSTRAINT positive_num
    CHECK (number >= 0);

ALTER TABLE public.node
DROP COLUMN IF EXISTS gender;

ALTER TABLE public.node
DROP COLUMN IF EXISTS is_spouse;