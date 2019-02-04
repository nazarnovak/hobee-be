package config

type Config struct {
	Test bool   `json:"test"`
	Port string `json:"port"`

	DB  `json:"db"`
	Log `json:"log"`
	//RedisURL string `json:"redisurl"`
}

type DB struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Log struct {
	Out string `json:"out"`
}
