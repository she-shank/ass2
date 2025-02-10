package server

type RestApi struct {
	port int
}

func NewRestApi() (*RestApi, error) {
	return &RestApi{}, nil
}

func (ra *RestApi) Start() error {
	return nil
}

func (ra *RestApi) Close() error {
	return nil
}
