package helpers

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type BaseEnv struct {
	Host string
	Port string
}

type RedisEnv struct {
	SessionKey string
	MaxConn    int
	Secret     string
	BaseEnv
}

type DBEnv struct {
	User     string
	Password string
	Name     string
	BaseEnv
}

type GoogleCloudEnv struct {
	Filepath string
}

func RetrieveRedisEnv() *RedisEnv {
	sessionKey := os.Getenv("REDIS_SESSION_KEY")
	host := os.Getenv("REDISHOST")
	port := os.Getenv("REDISPORT")
	maxConn, err := strconv.Atoi(os.Getenv("REDIS_MAX_CONNECTIONS"))
	if err != nil {
		panic(err)
	}
	secret := os.Getenv("REDIS_SECRET_KEY")

	env := RedisEnv{
		SessionKey: sessionKey,
		MaxConn:    maxConn,
		Secret:     secret,
		BaseEnv: BaseEnv{
			Host: host,
			Port: port,
		},
	}
	return &env
}

func RetrieveGoogleCloudEnv() *GoogleCloudEnv {
	filepath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	env := GoogleCloudEnv{
		Filepath: filepath,
	}
	return &env
}

func RetrieveWebAppEnv() *BaseEnv {
	addr := os.Getenv("WEBAPP_ADDRESS")
	port := os.Getenv("WEBAPP_PORT")
	env := BaseEnv{
		Host: addr,
		Port: port,
	}
	return &env
}

func RetrieveProdDBEnv() *DBEnv {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	env := DBEnv{
		User:     user,
		Password: password,
		Name:     name,
		BaseEnv: BaseEnv{
			Host: host,
			Port: port,
		},
	}
	return &env
}

func RetrieveTestDBEnv() *DBEnv {
	host := os.Getenv("DB_TEST_HOST")
	port := os.Getenv("DB_TEST_PORT")
	name := os.Getenv("DB_TEST_NAME")
	user := os.Getenv("DB_TEST_USER")
	password := os.Getenv("DB_TEST_PASSWORD")
	env := DBEnv{
		User:     user,
		Password: password,
		Name:     name,
		BaseEnv: BaseEnv{
			Host: host,
			Port: port,
		},
	}
	return &env
}

func (env *BaseEnv) Address() string {
	return fmt.Sprintf("%s:%s", env.Host, env.Port)
}

func (env *DBEnv) DataSourceName() string {
	dataSourceName := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v",
		env.Host,
		env.User,
		env.Password,
		env.Name,
		env.Port,
	)
	return dataSourceName
}

func RetrieveBackendEnv() *BaseEnv {
	host := os.Getenv("BACKEND_HOST")
	port := os.Getenv("BACKEND_PORT")
	env := BaseEnv{
		Host: host,
		Port: port,
	}
	return &env
}

func RetrieveClientAddress() string {
	clientAddress := os.Getenv("FRONTEND_BASE_URL")
	log.Println("Client address is:", clientAddress)
	return clientAddress
}

func RetrieveBackendAddress() string {
	backendAddress := os.Getenv("BACKEND_BASE_URL")
	log.Println("Set backend address to:", backendAddress)
	return backendAddress
}

func getBoolEnvVar(key string) bool {
	envVar := os.Getenv(key)
	val, err := strconv.ParseBool(envVar)
	if err != nil {
		return false
	}
	return val
}

func RetrieveDevMode() bool {
	return getBoolEnvVar("IS_RELEASE")
}
