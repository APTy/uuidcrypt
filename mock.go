package main

type mockConfig struct {
	config Config
}

func newMockConfig(config Config) RunConfig {
	return &mockConfig{config: config}
}

func (c *mockConfig) Load() error {
	return nil
}

func (c *mockConfig) Config() Config {
	return c.config
}

func (c *mockConfig) Done() error {
	return nil
}
