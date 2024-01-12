package util

import "github.com/spf13/viper"

type Config struct {
	Environment          string `mapstructure:"ENVIRONMENT"`
	DBDriver             string `mapstructure:"DB_DRIVER"`
	DBSource             string `mapstructure:"DB_SOURCE"`
	TestDBDriver         string `mapstructure:"TEST_DB_DRIVER"`
	TestDBSource         string `mapstructure:"TEST_DB_SOURCE"`
	DBName               string `mapstructure:"DB_NAME"`
	TestDBName           string `mapstructure:"TEST_DB_NAME"`
	HotelsCollection     string `mapstructure:"HOTEL_COLLECTION"`
	RoomsCollection      string `mapstructure:"ROOM_COLLECTION"`
	UsersCollection      string `mapstructure:"USER_COLLECTION"`
	TestUsersCollection  string `mapstructure:"TEST_USERS_COLLECTION"`
	TestHotelsCollection string `mapstructure:"TEST_HOTELS_COLLECTION"`
	TestRoomsCollection  string `mapstructure:"TEST_ROOMS_COLLECTION"`
	HTTPServerAddress    string `mapstructure:"HTTP_SERVER_ADDRESS"`
	EmailSenderPassword  string `mapstructure:"EMAIL_SENDER_PASSWORD"`
	BcryptCost           int    `mapstructure:"BCRYPT_COST"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
