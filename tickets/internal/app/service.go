package app

import (
	"context"
	"time"

	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/serviceutil"
	"github.com/matthxwpavin/ticketing/streaming"
	"github.com/matthxwpavin/ticketing/tickets/internal/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service is a type that represent domain/business logic
// of the application. It is high level language to communicate
// what exactly the application do.
type Service struct {
	tr                database.TicketRepository
	ticketCreatedPub  streaming.TicketCreatedPublisher
	ticketUpdatedPub  streaming.TicketUpdatedPublisher
	orderCreatedSub   streaming.OrderCreatedConsumer
	orderCancelledSub streaming.OrderCancelledConsumer
}

func NewService(
	ctx context.Context,
	db database.Database,
	s streaming.TicketStreamer,
) (*Service, error) {
	logger := sugar.FromContext(ctx)

	ticketCreatedPub, err := s.TicketCreatedPublisher(ctx)
	if err != nil {
		logger.Errorw("could not get ticketed created publisher", "error", err)
		return nil, err
	}
	ticketUpdatedPub, err := s.TicketUpdatedPublisher(ctx)
	if err != nil {
		logger.Errorw("could not get ticketed updated publisher", "error", err)
		return nil, err
	}
	orderCreatedSub, err := s.OrderCreatedConsumer(
		ctx,
		streaming.DefaultConsumeErrorHandler(ctx),
		streaming.OrderCreatedStreamConfig.Subjects[1],
	)
	if err != nil {
		logger.Errorw("could not get order created consumer", "error", err)
		return nil, err
	}
	orderCancelledSub, err := s.OrderCancelledConsumer(
		ctx,
		streaming.DefaultConsumeErrorHandler(ctx),
		streaming.OrderCanceledStreamConfig.Subjects[1],
	)
	if err != nil {
		logger.Errorw("could not get order cancelled consumer", "error", err)
		return nil, err
	}
	svc := &Service{
		tr:                db.TicketRepository(),
		ticketCreatedPub:  ticketCreatedPub,
		ticketUpdatedPub:  ticketUpdatedPub,
		orderCreatedSub:   orderCreatedSub,
		orderCancelledSub: orderCancelledSub,
	}
	if _, err := orderCancelledSub.Consume(ctx, svc.handleOrderCanceled(ctx)); err != nil {
		logger.Errorw("could not subscribe order canceled", "error", err)
		return nil, err
	}

	if err := svc.Subscribe(ctx); err != nil {
		return nil, err
	}
	return svc, nil
}

func (s *Service) Subscribe(ctx context.Context) error {
	subscriptions := []func(context.Context) error{
		s.subscribeOrderCreated,
	}
	for _, subscribe := range subscriptions {
		if err := subscribe(ctx); err != nil {
			return err
		}
	}
	return nil
}

type Ticket struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Price     int32     `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TicketCreate struct {
	Title string `json:"title" validate:"required"`
	Price int32  `json:"price" validate:"gt=0"`
}

func (s *Service) CreateTicket(ctx context.Context, input *TicketCreate) (*Ticket, error) {
	claims, err := serviceutil.AuthenticateClaims(ctx)
	if err != nil {
		return nil, err
	}

	logger := sugar.FromContext(ctx)
	if err := serviceutil.ValidateStruct(input); err != nil {
		logger.Errorw("Ticket is invalid", "error", err)
		return nil, err
	}

	doc := &database.Ticket{
		ID:        primitive.NewObjectID().Hex(),
		Title:     input.Title,
		Price:     input.Price,
		UserID:    claims.Metadata.UserID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Version:   1,
	}
	if _, err := s.tr.Insert(ctx, doc); err != nil {
		logger.Errorw("Could not insert a ticket", "error", err)
		return nil, err
	}
	if err := s.ticketCreatedPub.Publish(ctx, &streaming.TicketCreatedMessage{
		TicketID:      doc.ID,
		TicketTitle:   doc.Title,
		TicketPrice:   doc.Price,
		UserID:        claims.Metadata.UserID,
		TicketVersion: doc.Version,
	}); err != nil {
		logger.Errorw("Could not publish the ticket created event", "error", err)
		return nil, serviceutil.NewServiceFailureError("Could not publish the ticket created event")
	}
	return &Ticket{
		ID:        doc.ID,
		Title:     doc.Title,
		Price:     doc.Price,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}, nil
}

func (s *Service) GetTicket(ctx context.Context, id string) (*Ticket, error) {
	if _, err := serviceutil.AuthenticateClaims(ctx); err != nil {
		return nil, err
	}
	logger := sugar.FromContext(ctx)

	doc, err := s.tr.FindByID(ctx, id)
	if err != nil {
		logger.Errorw("Could not find the ticket", "error", err, "ticket id", id)
		return nil, serviceutil.NewServiceFailureError("Could not find the ticket")
	}
	if doc == nil {
		logger.Errorw("The ticket not found", "id", id)
		return nil, serviceutil.NewServiceFailureError("The ticket not found")
	}

	return &Ticket{
		ID:        doc.ID,
		Title:     doc.Title,
		Price:     doc.Price,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}, nil
}

func (s *Service) ListAllTickets(ctx context.Context) ([]*Ticket, error) {
	logger := sugar.FromContext(ctx)
	if _, err := serviceutil.AuthenticateClaims(ctx); err != nil {
		logger.Errorw("The user is unauthorized")
		return nil, err
	}

	docs, err := s.tr.FindAvailableTickets(ctx)
	if err != nil {
		logger.Errorw("Could not find tickets", "error", err)
		return nil, serviceutil.NewServiceFailureError("Could not find tickets")
	}

	res := []*Ticket{}
	for _, doc := range docs {
		res = append(res, &Ticket{
			ID:        doc.ID,
			Title:     doc.Title,
			Price:     doc.Price,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
		})
	}

	return res, nil
}

type TicketUpdate struct {
	Title string `json:"title" validate:"required"`
	Price int32  `json:"price" validate:"gt=0"`
}

func (s *Service) UpdateTicket(ctx context.Context, id string, input *TicketUpdate) error {
	logger := sugar.FromContext(ctx)
	claims, err := serviceutil.AuthenticateClaims(ctx)
	if err != nil {
		logger.Errorw("The user is unauthorized")
		return err
	}

	doc, err := s.mustFindTicketByID(ctx, id)
	if err != nil {
		return err
	}
	if claims.Metadata.UserID != doc.UserID {
		logger.Errorw("The user has no permission on the ticket", "userID", claims.Metadata.UserID, "ticket's user ID", doc.UserID)
		return serviceutil.NewServiceFailureError("Users have no permissoin on the ticket")
	}
	if doc.OrderId != "" {
		logger.Errorw("could not update the reserved ticket")
		return serviceutil.NewServiceFailureError("could not update the reserved ticket")
	}

	doc.Title = input.Title
	doc.Price = input.Price
	doc.Version += 1

	if err := s.tr.UpdateByID(ctx, id, doc); err != nil {
		logger.Errorw("Could not update the ticket", "error", err)
		return serviceutil.NewServiceFailureError("Could not update the ticket")
	}
	if err := s.ticketUpdatedPub.Publish(ctx, &streaming.TicketUpdatedMessage{
		TicketID:      doc.ID,
		TicketTitle:   doc.Title,
		TicketPrice:   doc.Price,
		TicketVersion: doc.Version,
	}); err != nil {
		logger.Errorw("Could not publish the tickets:updated event", "error", err)
		return serviceutil.NewServiceFailureError("Could not publish the tickets:updated event")
	}
	return nil
}

func (s *Service) DeleteTicket(ctx context.Context, id string) error {
	logger := sugar.FromContext(ctx)
	claims, err := serviceutil.AuthenticateClaims(ctx)
	if err != nil {
		logger.Errorw("The user is unauthorized")
		return err
	}

	doc, err := s.mustFindTicketByID(ctx, id)
	if err != nil {
		return err
	}

	if claims.Metadata.UserID != doc.UserID {
		logger.Errorw("The user has no permission on the ticket", "user ID", claims.Metadata.UserID, "ticket's user ID", doc.UserID)
		return serviceutil.NewServiceFailureError("No permission on the ticket")
	}
	if err := s.tr.DeleteByID(ctx, id); err != nil {
		logger.Errorw("Could not delete the ticket", "error", err, "ticket ID", id)
		return err
	}
	return nil
}

func (s *Service) mustFindTicketByID(ctx context.Context, id string) (*database.Ticket, error) {
	logger := sugar.FromContext(ctx)
	doc, err := s.tr.FindByID(ctx, id)
	if err != nil {
		logger.Errorw("Could not find the ticket", "error", err, "ticket ID", id)
		return nil, serviceutil.NewServiceFailureError("Could not find the ticket")
	}
	if doc == nil {
		logger.Errorw("The ticket not found", "ticket ID", id)
		return nil, serviceutil.NewServiceFailureError("The ticket not found")
	}
	return doc, nil
}
