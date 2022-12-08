package structs

type Config struct {
	DB       DataBase `yaml:"DB"`
	Location Location `yaml:"location"`
	Rating   float64  `yaml:"rating"`
}

type DataBase struct {
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Location struct {
	Latitude  float64 `yaml:"latitude"`
	Longitude float64 `yaml:"longitude"`
}
