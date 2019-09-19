package config

type SelecterConfig struct {
	Name    string // example img
	Attr    string // example src
	Pattern string // regex pattern
}

// get new selecterconfig
func NewSelecterConfig() *SelecterConfig {
	se := new(SelecterConfig)
	se.Name = "img"
	se.Attr = "src"
	se.Pattern = "\\.png"
	return se
}

// test this and init
func (this *SelecterConfig) TestConfig() error {

	if this.Name == "" {
		this.Name = "img"
	}
	if this.Attr == "" {
		this.Attr = "src"
	}
	if this.Pattern == "" {
		this.Pattern = "\\.png"
	}
	return nil
}
