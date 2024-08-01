package repository_test

import (
	"os"
	"reflect"
	"testing"

	"boilerplate-api/internal/config"
	"boilerplate-api/tests"
	"go.uber.org/fx"
)

type RepoTests []RepoTest
type RepoTest interface {
	SetupRepoTest(*testing.T)
}

func (repoTest RepoTests) InternalTestSetup() (tests []testing.InternalTest) {
	for _, test := range repoTest {
		name := reflect.TypeOf(test).Name()
		tests = append(tests, testing.InternalTest{
			Name: name,
			F:    test.SetupRepoTest,
		})
	}
	return tests
}

func NewRepoTests() RepoTests {
	return RepoTests{}
}

var RepoTestModules = fx.Options(
	config.TestENVModule,
	config.BaseModule,
	fx.Supply(config.EnvPath("../../.test.env")),
	fx.Provide(NewRepoTests),
	fx.Invoke(bootstrapRepoTest),
)

func bootstrapRepoTest(
	repoTests RepoTests,
	migrations config.Migrations,
) {
	migrations.MigrateUp()

	os.Exit(testing.MainStart(
		tests.MatchStringOnly{},
		repoTests.InternalTestSetup(),
		nil,
		nil,
		nil,
	).Run())
}

func TestMain(m *testing.M) {
	fx.New(RepoTestModules).Run()
}
