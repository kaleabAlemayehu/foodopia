package utility

var CreateUserQueryStr string = `mutation MyMutation($username: String!, $password_hash: String!, $email: String!) {
	insert_users_one(object: {username: $username, email: $email, password_hash: $password_hash}) {
	  id
	}
  }`

var CheckUser string = `query GetUserByEmail($email: String!) {
	users(where: {email: {_eq: $email}}) {
	  username
	  password_hash
	  email
	  id
	}
  }`
