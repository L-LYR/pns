package settings

import jsoniter "github.com/json-iterator/go"

// frontend private settings

type Settings struct {
	EmbedPages map[string]*EmbedPageSettings `json:"embedPages"`
}

type EmbedPageSettings struct {
	URL string `json:"url"`
}

var (
	GlobalSettings = &Settings{}
)

func MustLoadConfig() {
	err := jsoniter.UnmarshalFromString(_RawSettings, GlobalSettings)
	if err != nil {
		panic(err)
	}
}

func (s *Settings) MustGetEmbedPageSetting(name string) *EmbedPageSettings {
	eps, ok := s.EmbedPages[name]
	if !ok {
		panic("unknown settings")
	}
	return eps
}
