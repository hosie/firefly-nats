package nats

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/hyperledger/firefly-common/pkg/config"
	"github.com/hyperledger/firefly-common/pkg/ffresty"
	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/log"
	"github.com/hyperledger/firefly/pkg/core"
	"github.com/hyperledger/firefly/pkg/events"
	"github.com/nats-io/nats.go"
)

type Nats struct {
	ctx          context.Context
	capabilities *events.Capabilities
	callbacks    map[string]events.Callbacks
	client       *resty.Client
	connID       string
}

type Factory struct {
}

func (f *Factory) Type() string {
	return "nats"
}

func (f *Factory) NewInstance() events.Plugin {
	//create a new instance and return its address
	return &Nats{}
}

func (f *Factory) InitConfig(config config.Section) {
	ffresty.InitConfig(config)
}

//for backward compatibility
func (f *Nats) InitConfig(config config.Section) {
	factory := Factory{}
	factory.InitConfig(config)
}

type natsRequest struct {
	r         *resty.Request
	url       string
	method    string
	body      fftypes.JSONObject
	forceJSON bool
	replyTx   string
}

type natsResponse struct {
	Status  int                `json:"status"`
	Headers fftypes.JSONObject `json:"headers"`
	Body    *fftypes.JSONAny   `json:"body"`
}

// Implementation of events.Plugin interface
func (n *Nats) Name() string { return "nats" }

func (n *Nats) Init(ctx context.Context, config config.Section) (err error) {
	log.L(ctx).Info("Nats:Init")

	*n = Nats{
		ctx:          ctx,
		capabilities: &events.Capabilities{},
		callbacks:    make(map[string]events.Callbacks),
		client:       ffresty.New(ctx, config),
		connID:       fftypes.ShortID(),
	}
	return nil
}

func (n *Nats) SetHandler(namespace string, handler events.Callbacks) error {
	log.L(n.ctx).Infof("Nats:SetHandler: namespace='%s'", namespace)

	n.callbacks[namespace] = handler
	// We have a single logical connection, that matches all subscriptions
	return handler.RegisterConnection(n.connID, func(sr core.SubscriptionRef) bool { return true })
}

func (n *Nats) Capabilities() *events.Capabilities {
	log.L(n.ctx).Info("Nats:Capabilities")

	return n.capabilities
}

func (n *Nats) ValidateOptions(options *core.SubscriptionOptions) error {
	log.L(n.ctx).Info("Nats:ValidateOptions")

	return nil
}

func (n *Nats) DeliveryRequest(connID string, sub *core.Subscription, event *core.EventDelivery, data core.DataArray) error {
	log.L(n.ctx).Infof("Nats:DeliveryRequest: ConnectionId='%s'", connID)

	// Connect to a server
	nc, err := nats.Connect("nats://nats-server:4222")

	if err != nil {
		log.L(n.ctx).Error("Received a error")
		log.L(n.ctx).Error(err)

	}
	log.L(n.ctx).Info("Nats:DeliveryRequest: connected")

	defer nc.Close()

	// Simple Publisher
	eventJSON, jsonErr := json.Marshal(event)

	if jsonErr != nil {
		log.L(n.ctx).Error("Nats:DeliveryRequest Received a json error")
		log.L(n.ctx).Error(jsonErr)
		return jsonErr
	}
	log.L(n.ctx).Infof("Nats:DeliveryRequest publishing: event='%s'", string(eventJSON))

	err = nc.Publish(event.Topic, eventJSON)
	if err != nil {
		log.L(n.ctx).Error("Nats:DeliveryRequest Received a publishing error")
		log.L(n.ctx).Error(err)

	}
	log.L(n.ctx).Info("Nats:DeliveryRequest published")

	if cb, ok := n.callbacks[sub.Namespace]; ok {
		cb.DeliveryResponse(connID, &core.EventDeliveryResponse{
			ID:           event.ID,
			Rejected:     false,
			Subscription: event.Subscription,
			Reply:        nil,
		})
	}
	return nil
}
