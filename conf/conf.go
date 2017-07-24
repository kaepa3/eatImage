package conf

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

//Config 設定ファイル
type Config struct {
	TeamName   string
	APIKey     string
	ExportRoot string
}
type Configer interface {
	// メソッドリスト
	ReadConfig(path string) error
}

func (c *Config) ReadConfig(path string) error {
	if !exists(path) {
		return errors.New("no file")
	}
	_, err := toml.DecodeFile(path, &c)
	return err

}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
