package userrepository

const createNewUserQuery = `INSERT INTO users(username, email, password) VALUES($1,$2,$3) RETURNING id, created_at, updated_at`

const findOneUserByEmailQuery = `SELECT users.id, username, email, password, users.created_at, users.updated_at, roles.id, roles.role_name
	FROM users JOIN users_roles ON users.id=users_roles.user_id
	JOIN roles on users_roles.role_id=roles.id WHERE email=$1`

const findOneUserByUsernameQuery = `SELECT users.id, username, email, password, users.created_at, users.updated_at, roles.id, roles.role_name
	FROM users JOIN users_roles ON users.id=users_roles.user_id
	JOIN roles on users_roles.role_id=roles.id WHERE username=$1`

const findOneUserByIdQuery = `SELECT users.id, username, email, password, users.created_at, users.updated_at, roles.id, roles.role_name
	FROM users JOIN users_roles ON users.id=users_roles.user_id
	JOIN roles on users_roles.role_id=roles.id WHERE users.id=$1`

const findAllUserQuery = `SELECT users.id, username, email, password, users.created_at, users.updated_at, roles.id, roles.role_name
	FROM users JOIN users_roles ON users.id=users_roles.user_id
	JOIN roles on users_roles.role_id=roles.id ORDER BY users.id ASC`

const deleteUserQuery = `DELETE FROM users WHERE id=$1 RETURNING id`
