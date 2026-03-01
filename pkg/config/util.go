package config

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

func mapTo(vp *viper.Viper, section string, cfg interface{}) {
	if err := vp.UnmarshalKey(section, cfg); err != nil {
		log.Printf("Config parsing err: %v\n", err)
	}
}

func mapToList(vp *viper.Viper, section string, cfg interface{}, d interface{}) {
	if err := vp.UnmarshalKey(section, cfg,
		func(dc *mapstructure.DecoderConfig) {
			dc.TagName = "mapstructure"
			dc.DecodeHook = structDefaultHook(d)
		}); err != nil {
		log.Printf("Config parsing err: %v\n", err)
	}
}

func structDefaultHook(defaultStruct interface{}) mapstructure.DecodeHookFunc {
	defType := reflect.TypeOf(defaultStruct)
	return func(
		from reflect.Type, to reflect.Type, data interface{},
	) (interface{}, error) {
		if from.Kind() == reflect.Map && to == defType {
			rawMap, ok := data.(map[string]interface{})
			if !ok {
				return data, nil
			}
			defMap, err := structToMap(defaultStruct)
			if err != nil {
				return data, err
			}
			merged := mergeMaps(defMap, rawMap)
			return merged, nil
		}
		return data, nil
	}
}

func structToMap(s interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func mergeMaps(defaultMap, override map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(defaultMap))
	for k, v := range defaultMap {
		out[strings.ToLower(k)] = v
	}
	for k, v := range override {
		out[strings.ToLower(k)] = v
	}
	return out
}
