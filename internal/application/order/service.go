package order

import (
	"context"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/customer"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/product"
)

// ApplicationService orchestrates the application use cases
type ApplicationService struct {
	createOrderHandler       *CreateOrderHandler
	updateOrderStatusHandler *UpdateOrderStatusHandler
	getOrderHandler          *GetOrderHandler
}

func NewApplicationService(
	orderRepo order.Repository,
	productRepo product.Repository,
	customerRepo customer.Repository,
	eventPub order.EventPublisher,
) *ApplicationService {
	// Create adapter for product repository
	productAdapter := &ProductRepositoryAdapter{productRepo: productRepo}
	domainSvc := order.NewDomainService(orderRepo, productAdapter)

	return &ApplicationService{
		createOrderHandler:       NewCreateOrderHandler(orderRepo, productRepo, customerRepo, domainSvc, eventPub),
		updateOrderStatusHandler: NewUpdateOrderStatusHandler(orderRepo, eventPub),
		getOrderHandler:          NewGetOrderHandler(orderRepo),
	}
}

func (s *ApplicationService) CreateOrder(ctx context.Context, cmd CreateOrderRequest) (*CreateOrderResponse, error) {
	return s.createOrderHandler.Handle(ctx, cmd)
}

func (s *ApplicationService) UpdateOrderStatus(ctx context.Context, orderID string, cmd UpdateOrderStatusRequest) error {
	return s.updateOrderStatusHandler.Handle(ctx, orderID, cmd)
}

func (s *ApplicationService) GetOrder(ctx context.Context, orderID string) (*OrderResponse, error) {
	return s.getOrderHandler.Handle(ctx, orderID)
}

func (s *ApplicationService) GetOrdersByCustomer(ctx context.Context, customerID string) (*OrderListResponse, error) {
	return s.getOrderHandler.HandleByCustomer(ctx, customerID)
}

// ProductRepositoryAdapter adapts product.Repository to order.ProductRepository
type ProductRepositoryAdapter struct {
	productRepo product.Repository
}

func (a *ProductRepositoryAdapter) FindByID(ctx context.Context, id order.ProductID) (order.ProductInfo, error) {
	return a.productRepo.FindByID(ctx, id)
}