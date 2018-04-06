package controller

type EtcdConfig struct {
}

func (c *EtcdConfig) New() (*EtcdController, error) {
	return &EtcdController{}, nil
}
