package backstage

type Config struct {
	Host   string `json:"host"`
	Port   string `json:"port"`
	ApiKey string `json:"apiKey"`
}

func (c *Config) Validate() {

}
