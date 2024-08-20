package app

import (
	"context"
	"time"

	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/orders/internal/database"
	"github.com/matthxwpavin/ticketing/orderstatus"
	"github.com/matthxwpavin/ticketing/serviceutil"
	"github.com/matthxwpavin/ticketing/streaming"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service is a type that represent domain/business logic
// of the application. It is high level language to communicate
// what exactly the application do.
type Service struct {
	subscriberCtx     context.Context
	or                database.OrderRepository
	tr                database.TicketRepository
	orderCreatedPub   streaming.OrderCreatedPublisher
	orderCancelledPub streaming.OrderCancelledPublisher
	ticketCreatedSub  streaming.TicketCreatedConsumer
	ticketUpdatedSub  streaming.TicketUpdateConsumer
}

func NewService(ctx context.Context, db database.Database, s streaming.OrderStreamer) (*Service, error) {
	logger := sugar.FromContext(ctx)

	svc := &Service{
		subscriberCtx: ctx,
		or:            db.OrderRepository(),
		tr:            db.TicketRepository(),
	}

	var err error

	svc.orderCreatedPub, err = s.OrderCreatedPublisher(ctx)
	if err != nil {
		logger.Errorw("unable to get ticket created publisher", "error", err)
		return nil, err
	}
	svc.orderCancelledPub, err = s.OrderCancelledPublisher(ctx)
	if err != nil {
		logger.Errorw("unable to get ticket created publisher", "error", err)
		return nil, err
	}
	svc.ticketCreatedSub, err = s.TicketCreatedConsumer(ctx, streaming.DefaultConsumeErrorHandler(ctx), "")
	if err != nil {
		logger.Errorw("unable to get ticket created consumer", "error", err)
		return nil, err
	}
	svc.ticketUpdatedSub, err = s.TicketUpdatedConsumer(ctx, streaming.DefaultConsumeErrorHandler(ctx), "")
	if err != nil {
		logger.Errorw("unable to get ticket updated consumer", "error", err)
		return nil, err
	}
	expirationSub, err := s.ExpirationCompletedConsumer(ctx, streaming.DefaultConsumeErrorHandler(ctx), "")
	if err != nil {
		logger.Errorw("could not get the expiration completed consumer", "error", err)
		return nil, err
	}
	if _, err := expirationSub.Consume(ctx, svc.handleOrderExpiration); err != nil {
		logger.Errorw("could not consume the order expiration", "error", err)
		return nil, err
	}
	paymentCreatedSub, err := s.PaymentCreatedConsumer(
		ctx,
		streaming.DefaultConsumeErrorHandler(ctx),
		streaming.PaymentCreatedStreamConfig.Subjects[0],
	)
	if err != nil {
		logger.Errorw("could not get the payment created consumer", "error", err)
		return nil, err
	}
	if _, err := paymentCreatedSub.Consume(ctx, svc.handlePaymentCreated); err != nil {
		logger.Errorw("could not consume a payment created message", "error", err)
		return nil, err
	}

	ticketCreatedSub, err := s.TicketCreatedConsumer(
		ctx,
		streaming.DefaultConsumeErrorHandler(ctx),
		streaming.TicketCreatedStreamConfig.Subjects[0],
	)
	if err != nil {
		logger.Errorw("Could not get Ticket created consumer", "error", err)
		return nil, err
	}
	ticketUpdatedSub, err := s.TicketUpdatedConsumer(
		ctx,
		streaming.DefaultConsumeErrorHandler(ctx),
		streaming.TicketUpdatedStreamConfig.Subjects[0],
	)
	if err != nil {
		logger.Errorw("Could not get Ticket updated consumer", "error", err)
		return nil, err
	}

	if _, err := ticketCreatedSub.Consume(ctx, svc.handleTicketCreated(ctx)); err != nil {
		logger.Errorw("Could not cosume Ticket created message", "error", err)
		return nil, err
	}
	if _, err := ticketUpdatedSub.Consume(ctx, svc.handleTicketUpdated(ctx)); err != nil {
		logger.Errorw("Could not consume Ticket updated message", "error", err)
		return nil, err
	}
	return svc, nil
}

type OrderCreate struct {
	TicketID string `json:"ticketID" validate:"required"`
}

type Order struct {
	ID        string    `json:"orderID"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expiresAt"`
	UserID    string    `json:"userID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Ticket struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Price int32  `json:"price"`
}

type TicketOrder struct {
	Order
	Ticket Ticket `json:"ticket"`
}

const expirationWindow = 15 * time.Minute

func (s *Service) CreateOrder(ctx context.Context, input OrderCreate) (*TicketOrder, error) {
	logger := sugar.FromContext(ctx)
	claims, err := serviceutil.AuthenticateClaims(ctx)
	if err != nil {
		logger.Errorln("the request is unauthorized")
		return nil, err
	}

	logger = logger.With("user_id", claims.Metadata.UserID)

	if err := serviceutil.ValidateStruct(input); err != nil {
		logger.Errorw("invalid order to be created", "error", err)
		return nil, err
	}

	logger = logger.With("ticket_id", input.TicketID)

	ticket, err := s.tr.FindByID(ctx, input.TicketID)
	if err != nil {
		logger.Errorw("could not find the ticket", "error", err)
		return nil, err
	}
	if ticket == nil {
		logger.Errorln("the ticket not found")
		return nil, serviceutil.NewServiceFailureError("the ticket not found")
	}

	reserved, err := s.or.IsTicketReserved(ctx, ticket.ID)
	if err != nil {
		logger.Errorw("could not check the ticket is reserved", "error", err)
		return nil, serviceutil.NewServiceFailureError("could not check the ticket is reserved")
	}
	if reserved {
		logger.Errorw("the ticket has already reserved")
		return nil, serviceutil.NewServiceFailureError("the ticket has already reserved")
	}

	order := &database.TicketIdOrder{
		Order: database.Order{
			ID:        primitive.NewObjectID().Hex(),
			Status:    orderstatus.Created,
			ExpiresAt: time.Now().Add(expirationWindow),
			UserID:    claims.Metadata.UserID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   1,
		},
		TicketId: ticket.ID,
	}

	orderId, err := s.or.Insert(ctx, order)
	if err != nil {
		logger.Errorw("could not insert an order", "error", err)
		return nil, serviceutil.NewServiceFailureError("could not insert an order")
	}

	message := &streaming.OrderCreatedMessage{
		OrderId:        orderId,
		OrderStatus:    order.Status,
		OrderExpiresAt: order.ExpiresAt,
		OrderVersion:   order.Version,
		OrderUserId:    order.UserID,
	}
	message.Ticket.Id = ticket.ID
	message.Ticket.Price = ticket.Price

	if err := s.orderCreatedPub.Publish(ctx, message); err != nil {
		logger.Errorw("could not publish the event", "error", err)
		return nil, serviceutil.NewServiceFailureError("could not publish an order created event")
	}

	return &TicketOrder{
		Order: Order{
			ID:        orderId,
			Status:    order.Status,
			ExpiresAt: order.ExpiresAt,
			UserID:    order.UserID,
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
		},
		Ticket: Ticket{
			ID:    ticket.ID,
			Title: ticket.Title,
			Price: ticket.Price,
		},
	}, nil
}

func (s *Service) ListAllOrders(ctx context.Context) ([]*TicketOrder, error) {
	logger := sugar.FromContext(ctx)
	claims, err := serviceutil.AuthenticateClaims(ctx)
	if err != nil {
		return nil, err
	}

	logger = logger.With("user_id", claims.Metadata.UserID)
	orders, err := s.or.FindTicketOrdersByUserId(ctx, claims.Metadata.UserID)
	if err != nil {
		logger.Errorw("cound not find user's ticket orders", "error", err)
		return nil, err
	}
	res := make([]*TicketOrder, len(orders))
	for i, order := range orders {
		res[i] = &TicketOrder{
			Order: Order{
				ID:        order.ID,
				Status:    order.Status,
				ExpiresAt: order.ExpiresAt,
				UserID:    order.UserID,
				CreatedAt: order.CreatedAt,
				UpdatedAt: order.UpdatedAt,
			},
			Ticket: Ticket{
				ID:    order.Ticket.ID,
				Title: order.Ticket.Title,
				Price: order.Ticket.Price,
			},
		}
	}
	return res, nil
}

type OrderGet struct {
	OrderId string `json:"orderId" validate:"required"`
}

func (s *Service) GetOrder(ctx context.Context, input *OrderGet) (*TicketOrder, error) {
	claims, err := serviceutil.AuthenticateClaims(ctx)
	if err != nil {
		return nil, err
	}

	logger := sugar.FromContext(ctx).With("user_id", claims.Metadata.UserID)
	if err := serviceutil.ValidateStruct(input); err != nil {
		logger.Errorw("could not validate input", "error", err)
		return nil, err
	}

	result, err := s.or.FindTicketOrderByOrderID(ctx, input.OrderId)
	if err != nil {
		logger.Errorw("could not find the ticket order", "error", err)
		return nil, serviceutil.NewServiceFailureError("could not find the ticket order")
	}
	if result == nil {
		logger.Errorw("the order not found", "order_id", input.OrderId)
		return nil, serviceutil.NewServiceFailureError("the order not found")
	}

	if result.UserID != claims.Metadata.UserID {
		logger.Errorw("the order is not owned by the user", "order_user_id", result.UserID)
		return nil, serviceutil.NewServiceFailureError("the order is not owned by the user")
	}

	order := result.Order
	ticket := result.Ticket
	return &TicketOrder{
		Order: Order{
			ID:        order.ID,
			Status:    order.Status,
			ExpiresAt: order.ExpiresAt,
			UserID:    order.UserID,
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
		},
		Ticket: Ticket{
			ID:    ticket.ID,
			Title: ticket.Title,
			Price: ticket.Price,
		},
	}, nil
}

type UpdateOrderInput struct {
	OrderId     string `validate:"required"`
	OrderStatus string `json:"orderStatus"`
}

func (s *Service) UpdateOrder(ctx context.Context, input *UpdateOrderInput) error {
	user, err := serviceutil.Authenticate(ctx)
	if err != nil {
		return err
	}

	logger := sugar.FromContext(ctx).With("user_id", user.UserID, "order_id", input.OrderId)
	order, err := s.or.FindByID(ctx, input.OrderId)
	if err != nil {
		logger.Errorw("could not find an order", "error", err)
		return serviceutil.NewServiceFailureError("could not find an error")
	}
	if order == nil {
		logger.Errorw("order not found")
		return serviceutil.NewServiceFailureError("order not found")
	}
	if order.UserID != user.UserID {
		logger.Errorw("the user has no permission")
		return serviceutil.NewServiceFailureError("no permission")
	}

	order.Version += 1
	if input.OrderStatus != "" {
		order.Status = input.OrderStatus
		if err := s.or.UpdateByID(ctx, order.ID, order); err != nil {
			logger.Errorw("could not update the order", "error", err)
			return serviceutil.NewServiceFailureError("could not update the order")
		}
		message := &streaming.OrderCancelledMessage{
			OrderId:      order.ID,
			OrderVersion: order.Version,
		}
		message.Ticket.Id = order.TicketId
		if err := s.orderCancelledPub.Publish(ctx, message); err != nil {
			logger.Errorw("could not publish the event", "error", err)
			return serviceutil.NewServiceFailureError("could not emit the order cancelled event")
		}
	}

	return nil
}
