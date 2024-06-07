package accountactivationrepository

const createAccountActivationDataQuery = `INSERT INTO accounts_activation(user_id, token, expiration_time, next_request_available_at) VALUES($1, $2, $3, $4)`
const findOneAccountActivationDataQuery = `SELECT user_id, token, request_time, expiration_time, next_request_available_at FROM accounts_activation WHERE token=$1`
const findOneAccountByUserIdActivationDataQuery = `SELECT user_id, token, request_time, expiration_time, next_request_available_at FROM accounts_activation WHERE user_id=$1 ORDER BY expiration_time DESC`
const deleteAccountActivationDataQuery = `DELETE FROM accounts_activation WHERE token=$1`
const updateRequestTimeAccountActivationDataQuery = `UPDATE accounts_activation SET request_time=$1, next_request_available_at=$2 WHERE token=$3`
