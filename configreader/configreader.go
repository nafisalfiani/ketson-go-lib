package configreader

const (
	Viper  = "viper"
	DotEnv = "dotenv"
)

type Interface interface {
	ReadConfig(cfg interface{}) error
	AllSettings() map[string]interface{}
}

type Options struct {
	Type string
	// location of the main config file
	ConfigFile string
	// additional configuration to append to the main config
	AdditionalConfig []AdditionalConfigOptions
}
type AdditionalConfigOptions struct {
	// key is the location in the main config to append the additional config
	// with dot delimiter. e.g. 'Parser.ExcelOptions'
	ConfigKey string
	// location of the additional config
	ConfigFile string
}

// Init creates new instance of configreader.
// If you want to use DotEnv, make sure that the config struct has 'env' struct tag in each field
func Init(opt Options) Interface {
	switch opt.Type {
	case Viper:
		return initViper(opt)
	case DotEnv:
		return initDotEnv(opt)
	default:
		return initDotEnv(opt)
	}
}
