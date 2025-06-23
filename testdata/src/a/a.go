package a

type A struct {
	UserName string   `db:"user_name"`
	IsAdmin  bool     `json:"is_admin"`
	Roles    []string `yaml:"roles"`
	Age      int
}
