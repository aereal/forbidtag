package a

type A struct {
	UserName string   `db:"user_name"`  // want "unexpected struct tag: db"
	IsAdmin  bool     `json:"is_admin"` // want "unexpected struct tag: json"
	Roles    []string `yaml:"roles"`    // want "unexpected struct tag: yaml"
	Age      int
}
