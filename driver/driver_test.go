package driver

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	. "github.com/pingcap/check"
)

// Use docker mysql to test, mysql is 3306
var testHost = os.Getenv("MYSQL_HOST")
var testPort = os.Getenv("MYSQL_PORT")
var testUser = os.Getenv("MYSQL_USER")
var testPassword = os.Getenv("MYSQL_PASSWORD")
var testDB = os.Getenv("MYSQL_DATABASE")

func TestDriver(t *testing.T) {
	TestingT(t)
}

type testDriverSuite struct {
	db *sqlx.DB
}

var _ = Suite(&testDriverSuite{})

func (s *testDriverSuite) SetUpSuite(c *C) {
	dsn := fmt.Sprintf("%s:%s@%s:%s?%s", testUser, testPassword, testHost, testPort, testDB)

	var err error
	s.db, err = sqlx.Open("mysql", dsn)
	c.Assert(err, IsNil)
}

func (s *testDriverSuite) TearDownSuite(c *C) {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *testDriverSuite) TestConn(c *C) {
	var n int
	err := s.db.Get(&n, "SELECT 1")
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 1)

	_, err = s.db.Exec("USE test")
	c.Assert(err, IsNil)
}

func (s *testDriverSuite) TestStmt(c *C) {
	stmt, err := s.db.Preparex("SELECT ? + ?")
	c.Assert(err, IsNil)

	var n int
	err = stmt.Get(&n, 1, 1)
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 2)

	err = stmt.Close()
	c.Assert(err, IsNil)
}

func (s *testDriverSuite) TestTransaction(c *C) {
	tx, err := s.db.Beginx()
	c.Assert(err, IsNil)

	var n int
	err = tx.Get(&n, "SELECT 1")
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 1)

	err = tx.Commit()
	c.Assert(err, IsNil)
}
