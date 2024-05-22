package utility

var CreateUserQueryStr string = `mutation MyMutation($username: String!, $password_hash: String!, $email: String!) {
	insert_foodopia_users_one(object: {email: $email, password_hash: $password_hash, username: $username}) {
		id
		username
		password_hash
	}
}`

var CheckUser string = `query MyQuery($email: String!) {
	foodopia_users(where: {email: {_eq: $email}}) {
	  id
	  username
	  password_hash
	}
  }`
