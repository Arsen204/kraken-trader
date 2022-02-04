package service

type Service struct {
	Algorithm AlgorithmService
	Exchange  ExchangeService
}

func NewService(algorithmService AlgorithmService, exchangeService ExchangeService) *Service {
	return &Service{
		Algorithm: algorithmService,
		Exchange:  exchangeService,
	}
}
