package config

import (
	"dating/internal/pkg/utils"
	"fmt"
	"log"
	"reflect"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var EM ErrorMessage

type ErrorCode struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorMessage struct {
	vn         *viper.Viper
	ConfigPath string
	Success    ErrorCode
	Internal   struct {
		Database     ErrorCode
		DataNotFound ErrorCode
	}
	Invalid struct {
		Request ErrorCode
	}
}

// Initialization error message
func (em *ErrorMessage) Init() error {
	log.Println("initialzing error messages")
	vn := viper.New()
	vn.AddConfigPath(em.ConfigPath)
	vn.SetConfigName("errors")

	if err := vn.ReadInConfig(); err != nil {
		return err
	}
	em.vn = vn

	em.mapping("", reflect.ValueOf(em).Elem())

	vn.WatchConfig()
	vn.OnConfigChange(func(e fsnotify.Event) {
		log.Println("error messages change: %s", e.Name)
		em.vn = vn
		em.mapping("", reflect.ValueOf(em).Elem())
	})

	return nil
}

// mapping method is used to map the field name
func (em ErrorMessage) mapping(name string, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		fi := v.Field(i)
		if fi.Kind() != reflect.Struct {
			continue
		}

		fn := utils.Underscore(v.Type().Field(i).Name)
		if name != "" {
			fn = fmt.Sprint(name, ".", fn)
			fmt.Println(fn)
		}

		if fi.Type().Name() == "ErrorCode" {
			fi.Set(reflect.ValueOf(em.ErrorCode(fn)))
			continue
		}
		em.mapping(fn, fi)
	}
}

// ErrorCode method helps to get the value of error
func (em ErrorMessage) ErrorCode(name string) ErrorCode {
	rtn := ErrorCode{
		Code:    em.vn.GetString(fmt.Sprintf("error.%s.code", name)),
		Message: em.vn.GetString(fmt.Sprintf("error.%s.message", name)),
	}
	return rtn
}

//HasError method helps to verify error exists
func (ec ErrorCode) HasError() bool {
	if ec.Code != "" {
		return true
	}
	return false
}

// NewError method helps to create a new error. Avoid to use if dont have any special reason
func NewError(desc string) ErrorCode {
	return ErrorCode{
		Code:    "999",
		Message: desc,
	}
}
