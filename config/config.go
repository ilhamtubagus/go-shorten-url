package config

type RedisConfig struct {
	Port     string `env:"REDIS_PORT" defaultEnv:"6379"`
	Host     string `env:"REDIS_HOST"`
	Password string `env:"REDIS_PASSWORD"`
	TTL      int    `env:"REDIS_TTL"`
}

type MongoConfig struct {
	Host     string `env:"MONGODB_HOST"`
	User     string `env:"MONGODB_USER"`
	Password string `env:"MONGODB_PASSWORD"`
	Database string `env:"MONGODB_DATABASE_NAME"`
	Options  string `env:"MONGODB_OPTIONS"`
}

type Config struct {
	Host     string `env:"SERVICE_HOST"`
	Port     string `env:"SERVICE_PORT"`
	Protocol string `env:"SERVICE_PROTOCOL" default:"http"`
	Redis    RedisConfig
	Mongo    MongoConfig
}
