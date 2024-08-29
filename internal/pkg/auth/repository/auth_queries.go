package repository

const (
	GetUserByEmailQuery = `
		SELECT id, email, password_hash
		FROM public.user
		WHERE email = $1;
		`

	GetUserByIDQuery = `
		SELECT id, email, password_hash
		FROM public.user
		WHERE id = $1;
		`

	CreateUserQuery = `
		INSERT
		INTO public.user (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email;
		`
)
