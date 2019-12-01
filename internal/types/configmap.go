package types

type configMap struct {
	configMap map[string]interface{}
}

type ConfigMap interface {
	GetString(key string) string
	SetValue(key string, value interface{})
	SetValueMap(mapSI map[string]interface{})
}

func NewConfigMap() ConfigMap {
	return &configMap{
		configMap: map[string]interface{}{},
	}
}

func NewConfigMapFromMapSI(mapSI map[string]interface{}) ConfigMap {
	return &configMap{
		configMap: mapSI,
	}
}

func (cm *configMap) SetValueMap(data map[string]interface{}) {
	for key, intf := range data {
		cm.configMap[key] = intf
	}
}

func (cm *configMap) SetValue(key string, value interface{}) {
	cm.configMap[key] = value
}

func (cm *configMap) GetString(key string) string {
	intf, ok := cm.configMap[key]
	if !ok {
		return ""
	}
	val, ok := intf.(string)
	if !ok {
		return ""
	}
	return val
}
