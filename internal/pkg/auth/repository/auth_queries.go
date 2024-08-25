package repository

const (
	GerUserByEmailQuery = `
		SELECT *
		FROM public.user
		WHERE email = $1;
		`

	GerUserByIDQuery = `
		SELECT *
		FROM public.user
		WHERE id = $1;
		`

	CreateUserQuery = `
		INSERT
		INTO public.user (email, password_hash)
		VALUES ($1, $2)
		`
)
