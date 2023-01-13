package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	Db       string `mapstructure:"db" json:"db"`
}


type ServerConfig struct {
	Name            string      `mapstructure:"name" json:"name"`
	MysqlInfo       MysqlConfig `mapstructure:"mysql" json:"mysql"`
}