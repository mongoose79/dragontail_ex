package models

type Configuration struct {
	ServicePort int `json:"ServicePort"`

	CSVSourceFile string `json:"CSVSourceFile"`

	Host     string `json:"Host"`
	DbPort   int    `json:"DbPort"`
	User     string `json:"User"`
	Password string `json:"Password"`
	Dbname   string `json:"Dbname"`

	GoogleMapsAPIKey string `json:"GoogleMapsAPIKey"`
}

type Restaurant struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
}

type ByName []Restaurant

func (a ByName) Len() int           { return len(a) }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
