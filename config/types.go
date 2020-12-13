package config

type Config struct {
	DevEnv bool `json:"-"`

	DB     `json:"db"`
	Log    `json:"log"`
	Secret string `json:"secret"`
	//RedisURL string `json:"redisurl"`
}

type DB struct {
	Connection     string `json:"connection"`
	ConnectionProd string `json:"connectionProd"`
}

type Log struct {
	Out string `json:"out"`
}
