package env

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	populateTaggedStructures(&DB_CONFIG, &AWS_CONFIG)
}

func populateTaggedStructures(configs ...interface{}) error {
	var missingVars []string

	for _, cfg := range configs {
		missing, err := populateStruct(cfg)
		if err != nil {
			return err
		}
		missingVars = append(missingVars, missing...)
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("environment variables not set: %v", missingVars)
	}

	return nil
}

func populateStruct(cfg interface{}) ([]string, error) {
	const ENV_TAG_NAME string = "env"

	var val = reflect.ValueOf(cfg).Elem()
	var typ = val.Type()

	var missingVars []string


	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		envTag := fieldType.Tag.Get(ENV_TAG_NAME)
		envValue, exists := os.LookupEnv(envTag)
		if !exists {
			missingVars = append(missingVars, envTag)
			continue
		}

		field.SetString(envValue)
	}

	return missingVars, nil
}
