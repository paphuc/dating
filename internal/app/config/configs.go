package config

import (
	"dating/internal/pkg/glog"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/globalsign/mgo"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Stage string

const (
	StageLocal Stage = "dev"
	StageDEV   Stage = "dev"
)

func ParseStage(s string) Stage {
	switch s {
	case "local", "localhost", "l":
		return StageLocal
	case "dev", "develop", "development", "d":
		return StageDEV
	}
	return StageLocal
}

func New(path, state string) (*Configs, error) {
	conf := Configs{}
	stage := ParseStage(state)
	name := fmt.Sprintf("config.%s", stage)

	vn := viper.New()
	vn.AddConfigPath(path)
	vn.SetConfigName(name)
	vn.AutomaticEnv()
	vn.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := vn.ReadInConfig(); err != nil {
		errors.Wrap(err, "failed to read config")
		return nil, err
	}

	conf.binding(vn)

	vn.WatchConfig()
	vn.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("config file changed: %v", e.Name)
		if err := conf.binding(vn); err != nil {
			log.Printf("binding error:", err)
		}
		log.Printf("config: %+v", conf)
	})
	return &conf, nil
}

func (c *Configs) binding(v *viper.Viper) error {
	if err := v.Unmarshal(&c); err != nil {
		log.Printf("failed to unmarshal config: ", err)
		return err
	}
	return nil
}

type (
	Configs struct {
		Stage      Stage
		HTTPServer HTTPServer `mapstructure:"http_server"`
		Database   struct {
			Type  string  `mapstructure:"type"`
			Mongo MongoDB `mapstructure:"mongo"`
		} `mapstructure:"database"`
		Jwt struct {
			Duration time.Duration `mapstructure:"duration"`
		} `mapstructure:"jwt"`
	}

	// Config hold MongoDB configuration information
	MongoDB struct {
		Addresses []string      `envconfig:"MONGODB_ADDRS" default:"127.0.0.1:27017" mapstructure:"addresses"`
		Database  string        `envconfig:"MONGODB_DATABASE" default:"dating" mapstructure:"database"`
		Username  string        `mapstructure:"username"`
		Password  string        `mapstructure:"password"`
		Timeout   time.Duration `mapstructure:"timout"`
	}

	HTTPServer struct {
		Address           string        `mapstructure:"address"`
		Port              int           `mapstructure:"port"`
		ReadTimeout       time.Duration `mapstructure:"read_timeout"`
		WriteTimeout      time.Duration `mapstructure:"write_timeout"`
		ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
		ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
	}
)

// Dial dial to target server with Monotonic mode
func Dial(conf *MongoDB, logger glog.Logger) (*mgo.Session, error) {
	logger.Infof("dialing to target MongoDB at: %v, database: %v", conf.Addresses, conf.Database)
	ms, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    conf.Addresses,
		Database: conf.Database,
		Username: conf.Username,
		Password: conf.Password,
		Timeout:  conf.Timeout,
	})
	if err != nil {
		return nil, err
	}
	ms.SetMode(mgo.Monotonic, true)
	logger.Infof("successfully dialing to MongoDB at %v", conf.Addresses)
	return ms, nil
}
