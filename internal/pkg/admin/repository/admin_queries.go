package repository

const (
	GetAdminByEmailQuery = `
		SELECT u.id, u.email, u.password_hash
		FROM public.user u
		WHERE u.email = $1 AND u.is_admin;
		`

	EditUserPasswordQuery = `
		UPDATE public.user
		SET password_hash = $2
		WHERE id = $1
		RETURNING id, email, password_hash;
		`

	GetUsersListQuery = `
		SELECT u.id, u.email, u.password_hash
		FROM public.user u;
		`

	GetNodesListQuery = `
		SELECT nd.id, nd.name, nd.birthdate, nd.deathdate, nd.gender, nd.preview_path, nd.is_deleted,
			lr.id, lr.number,
			tr.id, tr.user_id
		FROM public.node nd
		JOIN public.layer lr ON lr.id = nd.layer_id
		JOIN public.tree tr ON tr.id = lr.tree_id
		WHERE tr.id = $1;
		`

	GetTreesListQuery = `
		SELECT tr.id, tr.user_id, tr.name
		FROM public.tree tr;
		`

	GetTreesListByUserIDQuery = `
		SELECT tr.id, tr.user_id, tr.name
		FROM public.tree tr
		WHERE tr.user_id = $1;
		`

	EditTreeNameQuery = `
		UPDATE public.tree
		SET name = $2
		WHERE id = $1
		RETURNING id, user_id, name;
		`
)
