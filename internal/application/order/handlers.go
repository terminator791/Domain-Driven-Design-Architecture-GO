package order

import (
	"context"

	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/common"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/customer"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/order"
	"github.com/terminator791/Domain-Driven-Design-Architecture-GO/internal/domain/product"
)

// CreateOrderHandler handles the creation of orders
type CreateOrderHandler struct {
	orderRepo    order.Repository
	productRepo  product.Repository
	customerRepo customer.Repository
	domainSvc    *order.DomainService
	eventPub     order.EventPublisher
}

func NewCreateOrderHandler(
	orderRepo order.Repository,
	productRepo product.Repository,
	customerRepo customer.Repository,
	domainSvc *order.DomainService,
	eventPub order.EventPublisher,
) *CreateOrderHandler {
	return &CreateOrderHandler{
		orderRepo:    orderRepo,
		productRepo:  productRepo,
		customerRepo: customerRepo,
		domainSvc:    domainSvc,
		eventPub:     eventPub,
	}
}

func (h *CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrderRequest) (*CreateOrderResponse, error) {
	// Validate customer exists
	customerID, err := order.NewCustomerIDFromString(cmd.CustomerID)
	if err != nil {
		return nil, err
	}

	_, err = h.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return nil, common.ErrCustomerNotFound
	}

	// Create new order
	newOrder := order.NewOrder(customerID)

	// Add items to order
	for _, itemReq := range cmd.Items {
		productID, err := order.NewProductIDFromString(itemReq.ProductID)
		if err != nil {
			return nil, err
		}

		quantity, err := order.NewQuantity(itemReq.Quantity)
		if err != nil {
			return nil, err
		}

		unitPrice, err := order.NewMoney(itemReq.UnitPrice, itemReq.Currency)
		if err != nil {
			return nil, err
		}

		err = newOrder.AddItem(productID, quantity, unitPrice)
		if err != nil {
			return nil, err
		}
	}

	// Validate order using domain service
	err = h.domainSvc.ValidateOrderCreation(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	// Confirm the order
	err = newOrder.Confirm()
	if err != nil {
		return nil, err
	}

	// Save order
	err = h.orderRepo.Save(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	// Publish domain events
	for _, event := range newOrder.Events() {
		err = h.eventPub.Publish(ctx, event)
		if err != nil {
			// Log error but don't fail the operation
			// In a real system, you might want to use an outbox pattern
		}
	}

	newOrder.ClearEvents()

	return &CreateOrderResponse{
		OrderID: newOrder.ID().String(),
		Message: "Order created successfully",
	}, nil
}

// UpdateOrderStatusHandler handles order status updates
type UpdateOrderStatusHandler struct {
	orderRepo order.Repository
	eventPub  order.EventPublisher
}

func NewUpdateOrderStatusHandler(orderRepo order.Repository, eventPub order.EventPublisher) *UpdateOrderStatusHandler {
	return &UpdateOrderStatusHandler{
		orderRepo: orderRepo,
		eventPub:  eventPub,
	}
}

func (h *UpdateOrderStatusHandler) Handle(ctx context.Context, orderID string, cmd UpdateOrderStatusRequest) error {
	id, err := order.NewOrderIDFromString(orderID)
	if err != nil {
		return err
	}

	existingOrder, err := h.orderRepo.FindByID(ctx, id)
	if err != nil {
		return common.ErrOrderNotFound
	}

	status := order.OrderStatus(cmd.Status)
	err = existingOrder.UpdateStatus(status)
	if err != nil {
		return err
	}

	err = h.orderRepo.Save(ctx, existingOrder)
	if err != nil {
		return err
	}

	// Publish domain events
	for _, event := range existingOrder.Events() {
		err = h.eventPub.Publish(ctx, event)
		if err != nil {
			// Log error but don't fail the operation
		}
	}

	existingOrder.ClearEvents()

	return nil
}

// GetOrderHandler handles order retrieval
type GetOrderHandler struct {
	orderRepo order.Repository
}

func NewGetOrderHandler(orderRepo order.Repository) *GetOrderHandler {
	return &GetOrderHandler{
		orderRepo: orderRepo,
	}
}

func (h *GetOrderHandler) Handle(ctx context.Context, orderID string) (*OrderResponse, error) {
	id, err := order.NewOrderIDFromString(orderID)
	if err != nil {
		return nil, err
	}

	existingOrder, err := h.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, common.ErrOrderNotFound
	}

	return h.mapToResponse(existingOrder), nil
}

func (h *GetOrderHandler) HandleByCustomer(ctx context.Context, customerID string) (*OrderListResponse, error) {
	id, err := order.NewCustomerIDFromString(customerID)
	if err != nil {
		return nil, err
	}

	orders, err := h.orderRepo.FindByCustomerID(ctx, id)
	if err != nil {
		return nil, err
	}

	var responses []OrderResponse
	for _, o := range orders {
		responses = append(responses, *h.mapToResponse(o))
	}

	return &OrderListResponse{
		Orders: responses,
		Total:  len(responses),
	}, nil
}

func (h *GetOrderHandler) mapToResponse(o *order.Order) *OrderResponse {
	var items []OrderItemResponse
	for _, item := range o.Items() {
		items = append(items, OrderItemResponse{
			ProductID: item.ProductID().String(),
			Quantity:  item.Quantity().Value(),
			UnitPrice: MoneyResponse{
				Amount:   item.UnitPrice().Amount(),
				Currency: item.UnitPrice().Currency(),
			},
			Total: MoneyResponse{
				Amount:   item.TotalPrice().Amount(),
				Currency: item.TotalPrice().Currency(),
			},
		})
	}

	return &OrderResponse{
		ID:         o.ID().String(),
		CustomerID: o.CustomerID().String(),
		Items:      items,
		Status:     o.Status().String(),
		Total: MoneyResponse{
			Amount:   o.Total().Amount(),
			Currency: o.Total().Currency(),
		},
		CreatedAt: o.CreatedAt(),
		UpdatedAt: o.UpdatedAt(),
	}
}