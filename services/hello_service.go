package services

// HelloService handles the business logic for hello operations
type HelloService struct{}

// NewHelloService creates a new instance of HelloService
func NewHelloService() *HelloService {
    return &HelloService{}
}

// GetHelloMessage returns a hello world message
func (s *HelloService) GetHelloMessage() string {
    return "Hello, World!"
}
