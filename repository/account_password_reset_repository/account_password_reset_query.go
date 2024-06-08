package accountpasswordresetrepository

const createAccountPasswordResetDataQuery = `INSERT INTO accounts_password_resets(user_id, token, expiration_time, next_request_available_at) VALUES($1, $2, $3, $4)`
const findOneAccountPasswordResetDataByUserIdQuery = `SELECT user_id, token, request_time, expiration_time, next_request_available_at FROM accounts_password_resets WHERE user_id=$1 ORDER BY expiration_time DESC LIMIT 1`
const updateRequestTimeAccountPasswordResetDataQuery = `UPDATE accounts_password_resets SET request_time=$1, next_request_available_at=$2 WHERE token=$3`
