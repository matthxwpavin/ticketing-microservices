package ordermongo

import (
	"context"

	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/orders/internal/database"
	"github.com/matthxwpavin/ticketing/orders/internal/database/impl/ordermongo/orderscollection"
	"github.com/matthxwpavin/ticketing/orderstatus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type orderRepository struct {
	*mongo.Collection[database.TicketIdOrder]
}

func (r *orderRepository) FindTicketOrdersByUserId(ctx context.Context, userId string) ([]*database.TicketOrder, error) {
	filterStage := bson.D{{"$match", bson.D{{"user_id", userId}}}}
	return r.findTicketOrders(ctx, filterStage)
}

func (r *orderRepository) FindTicketOrderByOrderID(ctx context.Context, id string) (*database.TicketOrder, error) {
	filterStage := bson.D{{"$match", bson.D{{"_id", id}}}}
	res, err := r.findTicketOrders(ctx, filterStage)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return nil, nil
}

func (r *orderRepository) findTicketOrders(ctx context.Context, stages ...primitive.D) ([]*database.TicketOrder, error) {
	logger := sugar.FromContext(ctx)

	lookupStage := bson.D{{
		"$lookup", bson.D{
			{"from", "tickets"},
			{"localField", "ticket_id"},
			{"foreignField", "_id"},
			{"as", "ticket"},
		}}}
	unwindStage := bson.D{{
		"$unwind", bson.D{
			{"path", "$ticket"},
			{"preserveNullAndEmptyArrays", false},
		}}}

	rest := append([]primitive.D{unwindStage}, stages...)
	cursor, err := r.Aggregate(ctx, lookupStage, rest...)
	if err != nil {
		logger.Errorw("unable to perfom an aggregation", "error", err)
		return nil, err
	}

	var res []*database.TicketOrder
	if err := cursor.All(ctx, &res); err != nil {
		logger.Errorw("could not decode the result", "error", err)
		return nil, err
	}
	return res, nil
}

func (r *orderRepository) FindByTicketIdAndStatuses(
	ctx context.Context,
	ticketId string,
	statuses []string,
) ([]*database.TicketIdOrder, error) {

	filter := bson.D{{"$and", bson.A{
		bson.D{{orderscollection.PropTicket.Name, ticketId}},
		bson.D{{orderscollection.PropStatus.Name, bson.D{{"$in", statuses}}}},
	}}}
	return r.Find(ctx, filter)
}

func (r *orderRepository) IsTicketReserved(ctx context.Context, ticketId string) (bool, error) {
	logger := sugar.FromContext(ctx)
	filter := bson.D{{"$and", bson.A{
		bson.D{{orderscollection.PropTicket.Name, ticketId}},
		bson.D{{orderscollection.PropStatus.Name, bson.D{{"$in", []string{
			orderstatus.Created,
			orderstatus.WaitingPayment,
			orderstatus.Complete,
		}}}}},
	}}}
	existingOrders, err := r.Find(ctx, filter)
	if err != nil {
		logger.Errorw("could not find orders", "error", err)
		return false, err
	}
	return len(existingOrders) > 0, nil
}
