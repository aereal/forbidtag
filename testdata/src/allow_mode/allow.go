package allow_mode

type C struct {
	UserName string   `db:"user_name"` // want "unexpected struct tag: db"
	IsAdmin  bool     `json:"is_admin"`
	Roles    []string `yaml:"roles"` // want "unexpected struct tag: yaml"
	Age      int
}
