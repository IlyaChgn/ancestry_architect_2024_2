ALTER TABLE public.description
ADD CONSTRAINT max_len_description
    CHECK (LENGTH(description) <= 500);
