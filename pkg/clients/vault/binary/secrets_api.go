package binary

func (c *client) Put(_ string, _ map[string]string) error {
	return nil
}

func (c *client) Get(_ string, _ string) (string, error) {
	return "", nil
}
