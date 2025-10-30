package environments

import (
	"github.com/joho/godotenv"
	"os"
)

func Init() {
	//Setup the environment
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	if env == "dev" {
		err := godotenv.Load("./environments/.env." + env)
		if err != nil {
			return
		}
	}
}
