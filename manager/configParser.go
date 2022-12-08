package manager

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"YandexEdaParser/structs"
)

func NewConfig(configPath string) (*structs.Config, error) {
	// Create config structure
	config := structs.Config{}
	//забираем конфиги из файлика
	var (
		file *os.File
		err  error
	)
	if file, err = os.Open(configPath); err != nil {
		return nil, errors.Wrap(err, "Error open file config")
	}
	defer file.Close()
	// Декодироем данные из файлика
	d := yaml.NewDecoder(file)
	// Отправляем данные в структуру
	if err := d.Decode(&config); err != nil {
		return nil, errors.Wrap(err, "Error Decode file in Config")
	}
	return &config, nil
}
