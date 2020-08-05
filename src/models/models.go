package models

type Configuration struct {
	ServicePort int `json:"ServicePort"`

	CSVSourceFile string `json:"CSVSourceFile"`

	Host     string `json:"Host"`
	DbPort   int    `json:"DbPort"`
	User     string `json:"User"`
	Password string `json:"Password"`
	Dbname   string `json:"Dbname"`
}

type Restaurant struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
}
