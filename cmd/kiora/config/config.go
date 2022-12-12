package config

import (
	"io/ioutil"

	"github.com/awalterschulze/gographviz"
)

type ConfigFile struct {
}

func LoadConfigFile(path string) (ConfigFile, error) {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return ConfigFile{}, err
	}

	graphAst, err := gographviz.ParseString(string(body))
	if err != nil {
		return ConfigFile{}, err
	}

	configGraph := newConfigGraph()
	if err := gographviz.Analyse(graphAst, &configGraph); err != nil {
		return ConfigFile{}, err
	}

	return ConfigFile{}, nil
}
