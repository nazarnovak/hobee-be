package config

type Config struct {
	Test bool   `json:"test"`
	Port string `json:"port"`

	DB     `json:"db"`
	Log    `json:"log"`
	Secret string `json:"secret"`
	//RedisURL string `json:"redisurl"`
	Email `json:"email"`
}

type DB struct {
	Connection string `json:"connection"`
}

type Log struct {
	Out string `json:"out"`
}

type Email struct {
	ApiKey string `json:"apikey"`
	Domain string `json:"domain"`
}
