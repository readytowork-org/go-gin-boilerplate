package services

import (
	"boilerplate-api/internal/api_errors"
	"boilerplate-api/internal/config"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
)

type CustomerSubscription struct {
	StripePriceID    string
	StripeCustomerID string
}

type StripeService struct {
	*client.API
	env config.Env
}

func NewStripeService(
	env config.Env,
	logger config.Logger,
) StripeService {
	_client := &client.API{}
	_client.Init(env.StripeSecretKey, nil)
	logger.Info("✅ Stripe client created.")
	return StripeService{
		API: _client,
		env: env,
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
) error {
	_, err := service.Subscriptions.Update(stripeSubscriptionID, stripeParams)
	if err != nil {
		return api_errors.InternalError.New("Errors while updating subscription")
	}
	return nil
}

func (service StripeService) CancelSubscription(
	stripeSubscriptionID string,
	stripeParams *stripe.SubscriptionCancelParams,
) error {
	_, err := service.Subscriptions.Cancel(stripeSubscriptionID, stripeParams)
	if err != nil {
		return api_errors.InternalError.New("Errors while updating subscription")
	}
	return nil
}

func (service StripeService) CreatePrices(
	title string, price int64,
) (prices *stripe.Price, err error) {
	priceParams := stripe.PriceParams{
		Product:    stripe.String(service.env.StripeProductID),
		Currency:   stripe.String(string(stripe.CurrencyJPY)),
		Nickname:   stripe.String(title),
		UnitAmount: stripe.Int64(price),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String(string(stripe.PriceRecurringIntervalMonth)),
		},
	}
	prices, err = service.Prices.New(&priceParams)
	if err != nil {
		return nil, err
	}
	return prices, err

}

func (service StripeService) UpdatePrices(
	stripePriceID string,
	priceParams *stripe.PriceParams,
) (prices *stripe.Price, err error) {
	prices, err = service.Prices.Update(
		stripePriceID,
		priceParams,
	)
	if err != nil {
		return nil, err
	}
	return prices, err

}

func (service StripeService) CreatePaymentIntent(
	paymentParams *stripe.PaymentIntentParams,
) (payment *stripe.PaymentIntent, err error) {
	paymentMethod := stripe.PaymentIntentAutomaticPaymentMethodsParams{
		Enabled: stripe.Bool(true),
	}
	paymentParams.AutomaticPaymentMethods = &paymentMethod
	payment, err = service.PaymentIntents.New(paymentParams)
	if err != nil {
		return nil, err
	}
	return payment, err

}

func (service StripeService) VoidInvoice(invoiceID string) error {
	params := &stripe.InvoiceVoidInvoiceParams{}
	_, err := service.Invoices.VoidInvoice(invoiceID, params)
	return err
}
