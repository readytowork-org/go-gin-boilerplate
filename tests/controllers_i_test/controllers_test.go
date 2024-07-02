package controllers_i_test

import (
	"os"
	"reflect"
	"testing"

	"boilerplate-api/internal/config"
	"boilerplate-api/internal/router"
	"boilerplate-api/services"
	"boilerplate-api/tests"
	"go.uber.org/fx"
)

type ControllerTest interface {
	SetupControllerTest(*testing.T)
}

type ControllerTests []ControllerTest

func (controllerTests ControllerTests) InternalTestSetup() (tests []testing.InternalTest) {
	for _, test := range controllerTests {
		name := reflect.TypeOf(test).Name()
		tests = append(tests, testing.InternalTest{
			Name: name,
			F:    test.SetupControllerTest,
		})
	}
	return tests
}

func NewControllerTests() ControllerTests {
	return ControllerTests{}
}

var ControllerIntegrationTestModules = fx.Options(
	config.TestENVModule,
	config.BaseModule,
	services.Module,
	fx.Supply(config.EnvPath("../../.test.env")),
	fx.Provide(router.NewRouter),
	fx.Provide(NewControllerTests),
	fx.Invoke(bootstrapRepoTest),
)

func bootstrapRepoTest(
	repoTests ControllerTests,
) {
	os.Exit(testing.MainStart(
		tests.MatchStringOnly{},
		repoTests.InternalTestSetup(),
		nil,
		nil,
		nil,
	).Run())
}

func TestMain(m *testing.M) {
	fx.New(ControllerIntegrationTestModules).Run()
}
