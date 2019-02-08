package config

type Config struct {
	Test bool   `json:"test"`
	Port string `json:"port"`

	DB  `json:"db"`
	Log `json:"log"`
	//RedisURL string `json:"redisurl"`
}

type DB struct {
	Connection string `json:"connection"`
}

type Log struct {
	Out string `json:"out"`
}
