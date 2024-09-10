package repository

const (
	GetUserByEmailQuery = `
		SELECT u.id, u.email, u.password_hash, p.name, p.surname
		FROM public.user u
		JOIN public.profile p
		ON u.id = p.user_id
		WHERE u.email = $1;
		`

	GetUserByIDQuery = `
		SELECT u.id, u.email, u.password_hash, p.name, p.surname
		FROM public.user u
		JOIN public.profile p
		ON u.id = p.user_id
		WHERE u.id = $1;
		`

	CreateUserQuery = `
		INSERT
		INTO public.user (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email;
		`

	UpdateEmailQuery = `
		UPDATE public.user
		SET email = $1
		WHERE id = $2
		RETURNING id, email;
		`
)
