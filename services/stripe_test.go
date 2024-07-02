package services

import (
	"bytes"
	"testing"

	"boilerplate-api/internal/config"

	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/form"
)

type stripeBackendMock struct {
	mock.Mock
}

func (s *stripeBackendMock) Call(
	method,
	path,
	key string,
	params stripe.ParamsContainer,
	v stripe.LastResponseSetter,
) error {
	args := s.Called(method, path, key, params, v)
	return args.Error(0)
}

func (s *stripeBackendMock) CallStreaming(
	method,
	path,
	key string,
	params stripe.ParamsContainer,
	v stripe.StreamingLastResponseSetter,
) error {
	args := s.Called(method, path, key, params, v)
	return args.Error(0)
}

// MockEnv is a mock implementation of config.Env
type MockEnv struct {
	mock.Mock
}

func (m *MockEnv) GetEnvironment() string {
	args := m.Called()
	return args.String(0)
}

func (s *stripeBackendMock) CallRaw(
	method,
	path,
	key string,
	body *form.Values,
	params *stripe.Params,
	v stripe.LastResponseSetter,
) error {
	args := s.Called(method, path, key, params, v)
	return args.Error(0)
}

func (s *stripeBackendMock) CallMultipart(
	method,
	path,
	key,
	boundary string,
	body *bytes.Buffer,
	params *stripe.Params,
	v stripe.LastResponseSetter,
) error {
	args := s.Called(method, path, key, boundary, body, params, v)
	return args.Error(0)
}

func (s *stripeBackendMock) SetMaxNetworkRetries(maxNetworkRetries int64) {
	s.Called(maxNetworkRetries)
}

func TestCreateCustomer(t *testing.T) {
	stripeBackendMock := new(stripeBackendMock)
	//stripeTestBackends := &stripe.Backends{
	//	API:     stripeBackendMock,
	//	Connect: stripeBackendMock,
	//	Uploads: stripeBackendMock,
	//}

	//stripeClient := client.New("sk_test", stripeTestBackends)

	env := config.NewEnv(
		".env",
	)
	stripeService := NewStripeService(
		env,
		config.GetLogger(),
	)
	stripeBackendMock.On("Call",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		stripeCustomer := args.Get(4).(*stripe.Customer)
		*stripeCustomer = stripe.Customer{
			Email: stripeCustomer.Email,
			Name:  stripeCustomer.Name,
		}
		// reponse := args.Get(5).(stripe.UserStripeInfo.A)
		// *&reponse :=
	}).Return(nil).Once()

	t.Run("Test if user is created is stripe", func(t *testing.T) {
		name := *stripe.String("test")
		email := *stripe.String("test@gmail.com")
		customer, err := stripeService.CreateCustomer(name, email)
		if err != nil {
			t.Error(err)
			return
		}

		if customer.Email != email && customer.Name != name {
			t.Error("Customer details doesn't match")
		}

	})
}

func TestCreateSubscription(t *testing.T) {
	stripeBackendMock := new(stripeBackendMock)
	//stripeTestBackends := &stripe.Backends{
	//	API:     stripeBackendMock,
	//	Connect: stripeBackendMock,
	//	Uploads: stripeBackendMock,
	//}

	//stripeClient := client.New("sk_test", stripeTestBackends)

	env := config.NewEnv(
		".env",
	)
	stripeService := NewStripeService(
		env,
		config.GetLogger(),
	)
	stripeBackendMock.On("Call",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		subs := args.Get(4).(*stripe.Subscription)
		*subs = stripe.Subscription{
			Customer: subs.Customer,
			Items:    subs.Items,
		}
	}).Return(nil).Once()

	stripeService.CreateSubscription(CustomerSubscription{
		StripeCustomerID: "test@gmail.com",
		StripePriceID:    "test@gmail.com",
	}, nil)
}

func TestCreatePrices(t *testing.T) {
	stripeBackendMock := new(stripeBackendMock)
	//stripeTestBackends := &stripe.Backends{
	//	API:     stripeBackendMock,
	//	Connect: stripeBackendMock,
	//	Uploads: stripeBackendMock,
	//}

	//stripeClient := client.New("sk_test", stripeTestBackends)
	env := config.NewEnv(
		".env",
	)
	stripeService := NewStripeService(
		env,
		config.GetLogger(),
	)

	stripeBackendMock.On("Call",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		price := args.Get(4).(*stripe.Price)
		*price = stripe.Price{
			UnitAmount: price.UnitAmount,
			Currency:   price.Currency,
			Recurring:  price.Recurring,
		}
	}).Return(nil).Once()
	stripeService.CreatePrices("Test Price", 1000)
}

func TestPaymentIntent(t *testing.T) {
	stripeBackendMock := new(stripeBackendMock)
	//stripeTestBackends := &stripe.Backends{
	//	API:     stripeBackendMock,
	//	Connect: stripeBackendMock,
	//	Uploads: stripeBackendMock,
	//}

	//stripeClient := client.New("sk_test", stripeTestBackends)

	env := config.NewEnv(
		".env",
	)
	stripeService := NewStripeService(
		env,
		config.GetLogger(),
	)

	stripeBackendMock.On("Call",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		payment := args.Get(4).(*stripe.PaymentIntent)
		*payment = stripe.PaymentIntent{
			Amount:                  payment.Amount,
			Currency:                payment.Currency,
			AutomaticPaymentMethods: payment.AutomaticPaymentMethods,
		}
	}).Return(nil).Once()

	paymentMethod := stripe.PaymentIntentAutomaticPaymentMethodsParams{
		Enabled: stripe.Bool(true),
	}
	stripeService.CreatePaymentIntent(&stripe.PaymentIntentParams{
		Amount:                  stripe.Int64(1000),
		Currency:                stripe.String("USD"),
		AutomaticPaymentMethods: &paymentMethod,
	})
}
