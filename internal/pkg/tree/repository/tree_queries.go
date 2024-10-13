package repository

const (
	CheckPermissionForTreeQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM public.permission pm
			WHERE pm.tree_id = $1 AND pm.user_id = $2
		
			UNION ALL
		
			SELECT 1
			FROM public.tree tr
			WHERE tr.user_id = $2 AND tr.id = $1
		);
		`

	GetTreeQuery = `
		WITH relation_data AS (
			SELECT
				nd.id AS node_id,
				COALESCE(
					json_agg(DISTINCT rl.node_id) FILTER (WHERE rl.relation_type = 'Родитель'),
					'[]'::json
				) AS parent_ids,
				COALESCE(
					json_agg(DISTINCT rl.node_id) FILTER (WHERE rl.relation_type = 'Супруг'),
					'[]'::json
				) AS spouse_ids
			FROM public.node nd
				LEFT JOIN public.relation rl ON rl.relative_id = nd.id
			GROUP BY nd.id
		),
			 tree_data AS (
				 SELECT
					 tr.id AS tree_id,
					 tr.name AS tree_name,
					 lr.id AS layer_id,
					 lr.number AS number,
					 nd.id AS node_id,
					 nd.layer_id,
					 nd.name,
					 nd.birthdate,
					 nd.deathdate,
					 nd.preview_path
				 FROM public.tree tr
					  JOIN public.layer lr ON lr.tree_id = tr.id
					  JOIN public.node nd ON lr.id = nd.layer_id
				 WHERE NOT nd.is_deleted AND tr.id = $1
			 )
		SELECT
			td.*,
			COALESCE(jsonb_build_object(
				'children', relation_data.parent_ids,
				'spouses', relation_data.spouse_ids
			), '{}'::jsonb) AS relation_json
		FROM tree_data td
			LEFT JOIN relation_data ON td.node_id = relation_data.node_id
		ORDER BY td.number;
		`

	GetCreatedTreesListQuery = `
		SELECT id, user_id, name
		FROM public.tree
		WHERE user_id = $1;
		`

	GetAvailableTreesListQuery = `
		SELECT tr.id, tr.user_id, tr.name
		FROM public.tree tr
			JOIN public.permission pm ON tr.id = pm.tree_id
		WHERE pm.user_id = $1;
		`

	CreateTreeQuery = `
		INSERT
		INTO public.tree (user_id, name)
		VALUES ($1, $2)
		RETURNING id, user_id, name;
		`

	AddPermissionQuery = `
		INSERT 
		INTO public.permission (user_id, tree_id) 
		VALUES ($1, $2); 
		`
)
