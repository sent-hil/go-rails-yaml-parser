package railsyamlparser

import (
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var exampleFileStr = `
defaults: &defaults
 adapter: postgresql
 encoding: unicode
 pool: 5
 port: 5432
 host: localhost

development:
 <<: *defaults
 database: development
 username: development
 password: development

test:
 <<: *defaults
 database: test
 username: test
 password: test
`

func TestParsingRailsYamlFile(t *testing.T) {
	r, err := New([]byte(exampleFileStr))
	if err != nil {
		log.Fatal(err)
	}

	Convey("Get", t, func() {
		Convey("It returns value set for env", func() {
			val, err := r.Get("Database")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "development")
		})

		Convey("It returns error for arg that's not defined", func() {
			_, err := r.Get("random")
			So(err, ShouldEqual, ErrKeyNotFound)
		})
	})

	Convey("SetEnv", t, func() {
		Convey("It returns Development as default if env is not set", func() {
			So(r.GetEnv(), ShouldEqual, Development)
		})

		Convey("It sets given env", func() {
			r.SetEnv(Production)
			So(r.GetEnv(), ShouldEqual, Production)
		})
	})
}