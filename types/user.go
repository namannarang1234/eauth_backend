package types

type User struct {
	Name     string `bson:"name"`
	Email    string `bson:"email"`
	Phone    string `bson:"phone"`
	Password string `bson:"password"`
	Token    string `bson:"token"`
}

type FEUser struct {
	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Phone string `bson:"phone" json:"phone"`
}
