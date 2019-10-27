package structs

//UserConnection usada para guardar os dados da conexão para uso posterior no PHP
type UserConnection struct {
	Dsn      string `json:"dsn"`
	User     string `json:"user"`
	Password string `json:"password"`
}
