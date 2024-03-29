package database

import (
	"context"
	"log"

	"github.com/stretchr/testify/assert"
)

func (suite *TestSuite) TestUserService_CreateUser() {
	ctx := context.Background()
	tx, err := suite.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer tx.Rollback()

	userID, err := createUser(ctx, tx, "test", "test")
	if err != nil {
		log.Println(err)
		return
	}

	assert.Equal(suite.T(), userID, 1, "they should be equal because it's the first user")
}

func (suite *TestSuite) Test_getUser() {
	ctx := context.Background()
	tx, err := suite.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer tx.Rollback()

	user, err := getUser(ctx, tx, 1)
	if err != nil {
		log.Println(err)
		return
	}

	assert.Nil(suite.T(), user, "should be nil because no user exists")
}
