package services

import (
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
	"net/http"
)

type sLogger interface {
	Info(args ...interface{})
}

type StripeConfig struct {
	stripeSecretKey string
	stripeProductID string
	logger          sLogger
}

type CustomerSubscription struct {
	StripePriceID    string
	StripeCustomerID string
}

type StripeService struct {
	*client.API
	stripeProductID string
}

// StripeErrorResponse struct
type StripeErrorResponse struct {
	Message   string `json:"message"`
	ErrorType int    `json:"error_type"`
}

func NewStripeService(
	stripeConfig StripeConfig,
) StripeService {
	_client := &client.API{}
	_client.Init(stripeConfig.stripeSecretKey, nil)
	stripeConfig.logger.Info("✅ Stripe client created.")
	return StripeService{
		API:             _client,
		stripeProductID: stripeConfig.stripeProductID,
	}
}

func (service StripeService) CreateCustomer(name, email string) (*stripe.Customer, error) {
	stripeCustomer, err := service.Customers.New(&stripe.CustomerParams{
		Name:  &name,
		Email: &email,
	})

	if err != nil {
		return nil, err
	}
	return stripeCustomer, err
}

func (service StripeService) CreateSubscription(
	customer CustomerSubscription, // FIXME
	backdateStartDate *int64,
) (subs *stripe.Subscription, err error) {

	itemsParams := []*stripe.SubscriptionItemsParams{
		{
			Price:    &customer.StripePriceID,
			Quantity: stripe.Int64(1),
		},
	}

	paymentSettings := stripe.SubscriptionPaymentSettingsParams{
		SaveDefaultPaymentMethod: stripe.String("on_subscription"),
	}
	subscriptionParams := stripe.SubscriptionParams{
		Customer:        &customer.StripeCustomerID,
		Items:           itemsParams,
		PaymentSettings: &paymentSettings,
		PaymentBehavior: stripe.String("default_incomplete"),
		BillingCycleAnchorConfig: &stripe.SubscriptionBillingCycleAnchorConfigParams{
			DayOfMonth: stripe.Int64(1),
		},
		BackdateStartDate: backdateStartDate,
	}
	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	subscriptionParams.AddExpand("pending_setup_intent")
	subs, err = service.Subscriptions.New(&subscriptionParams)
	if err != nil {
		return nil, err
	}
	return subs, err
}

func (service StripeService) UpdateSubscription(
	stripeSubscriptionID string,
	stripeParams *stripe.SubscriptionParams,
) *StripeErrorResponse {
	_, err := service.Subscriptions.Update(stripeSubscriptionID, stripeParams)
	if err != nil {
		return &StripeErrorResponse{
			Message:   "Errors while updating subscription",
			ErrorType: http.StatusInternalServerError,
		}
	}
	return nil
}

func (service StripeService) CancelSubscription(
	stripeSubscriptionID string,
	stripeParams *stripe.SubscriptionCancelParams,
) *StripeErrorResponse {
	_, err := service.Subscriptions.Cancel(stripeSubscriptionID, stripeParams)
	if err != nil {
		return &StripeErrorResponse{
			Message:   "Errors while canceling subscription",
			ErrorType: http.StatusInternalServerError,
		}
	}
	return nil
}

func (service StripeService) CreatePrices(
	title string, price int64,
) (*stripe.Price, *StripeErrorResponse) {
	priceParams := stripe.PriceParams{
		Product:    stripe.String(service.stripeProductID),
		Currency:   stripe.String(string(stripe.CurrencyJPY)),
		Nickname:   stripe.String(title),
		UnitAmount: stripe.Int64(price),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String(string(stripe.PriceRecurringIntervalMonth)),
		},
	}
	prices, err := service.Prices.New(&priceParams)
	if err != nil {
		return nil, &StripeErrorResponse{
			ErrorType: http.StatusInternalServerError,
			Message:   "Error while creating price",
		}
	}
	return prices, nil

}

func (service StripeService) UpdatePrices(
	stripePriceID string,
	priceParams *stripe.PriceParams,
) (prices *stripe.Price, errResponse *StripeErrorResponse) {
	prices, err := service.Prices.Update(
		stripePriceID,
		priceParams,
	)
	if err != nil {
		return nil, &StripeErrorResponse{
			ErrorType: http.StatusInternalServerError,
			Message:   "Error while updating price",
		}
	}
	return prices, nil

}

func (service StripeService) CreatePaymentIntent(
	paymentParams *stripe.PaymentIntentParams,
) (payment *stripe.PaymentIntent, errResponse *StripeErrorResponse) {
	paymentMethod := stripe.PaymentIntentAutomaticPaymentMethodsParams{
		Enabled: stripe.Bool(true),
	}
	paymentParams.AutomaticPaymentMethods = &paymentMethod
	payment, err := service.PaymentIntents.New(paymentParams)
	if err != nil {
		return nil, &StripeErrorResponse{
			ErrorType: http.StatusInternalServerError,
			Message:   "Error while creating payment intent",
		}
	}
	return payment, nil

}

func (service StripeService) VoidInvoice(invoiceID string) *StripeErrorResponse {
	params := &stripe.InvoiceVoidInvoiceParams{}
	if _, err := service.Invoices.VoidInvoice(invoiceID, params); err != nil {
		return &StripeErrorResponse{
			Message:   "Error while voiding invoice",
			ErrorType: http.StatusInternalServerError,
		}
	}
	return nil
}
