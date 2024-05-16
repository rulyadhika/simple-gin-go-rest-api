package accountactivationrepository

const createAccountActivationDataQuery = `INSERT INTO accounts_activation(user_id, token, request_time, expiration_time) VALUES($1, $2, $3, $4)`
const findOneAccountActivationDataQuery = `SELECT user_id, token, request_time, expiration_time FROM accounts_activation WHERE token=$1`
const findOneAccountByUserIdActivationDataQuery = `SELECT user_id, token, request_time, expiration_time FROM accounts_activation WHERE user_id=$1 ORDER_BY expiration_time DESC`
const deleteAccountActivationDataQuery = `DELETE FROM accounts_activation WHERE token=$1`
