package config

type Middleware struct {
	CoreData CoreData `yaml:"CoreData"`
}

type CoreData struct {
	Port                  string `yaml:"Port"`
	CessRpcAddr           string `yaml:"CessRpcAddr"`
	IdAccountPhraseOrSeed string `yaml:"IdAccountPhraseOrSeed"`
}

var Data = new(Middleware)
