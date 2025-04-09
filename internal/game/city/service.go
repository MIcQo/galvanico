package city

type Service interface {
}

type ServiceImpl struct {
}

func NewService() *ServiceImpl {
	return &ServiceImpl{}
}
