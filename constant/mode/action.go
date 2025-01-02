package mode

type Mode string

const (
	Chat           Mode = "chat"
	ChatStream     Mode = "stream"
	Embedding      Mode = "embedding"
	Completions    Mode = "completions"
	Models         Mode = "models"
	Audio          Mode = "audio"
	Image          Mode = "image"
	Translate      Mode = "translate"
	Transcriptions Mode = "transcriptions"
	ImageEdit      Mode = "imageEdit"
	ImageVariation Mode = "imageVariation"
)
