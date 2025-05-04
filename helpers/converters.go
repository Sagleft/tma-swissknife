package helpers

import "github.com/mitchellh/mapstructure"

func ConvertToStruct(data map[string]any, destinationPointer any) error {
	decoderConfig := &mapstructure.DecoderConfig{
		Result:  destinationPointer,
		TagName: "json", // учитываем теги json в структуре
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return err
	}
	return decoder.Decode(data)
}
