package config

type Config struct {
	Dev  bool   `json:"dev"`
	Port string `json:"port"`

	DB     `json:"db"`
	Log    `json:"log"`
	Secret string `json:"secret"`
	//RedisURL string `json:"redisurl"`
	Email `json:"email"`
}

type DB struct {
	Connection     string `json:"connection"`
	ConnectionProd string `json:"connectionProd"`
}

type Log struct {
	Out string `json:"out"`
}

type Email struct {
	ApiKey string `json:"apikey"`
	Domain string `json:"domain"`
}
