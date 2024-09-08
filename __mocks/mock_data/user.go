package mock_data

import (
	"time"

	"boilerplate-api/database/dao"
	"github.com/brianvoe/gofakeit/v7"
)

type FakeUser struct {
	dao.User
}

var createdDate = time.Now()

func (u *FakeUser) Fake(f *gofakeit.Faker) (any, error) {
	u.ID = f.Uint32()
	u.FullName = f.Name()
	u.Email = f.Email()
	u.Phone = f.Phone()
	u.Gender = f.Gender()
	u.Password = f.Word()
	u.CreatedAt = createdDate
	u.UpdatedAt = &u.CreatedAt

	createdDate = createdDate.Add(1 * time.Minute)
	return *u, nil
}

var Users = make([]FakeUser, 5)

func init() {
	gofakeit.Slice(&Users)
}
