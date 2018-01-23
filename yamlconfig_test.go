package yamlconfig

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

func TestParsingYamlFile(t *testing.T) {
	r, err := New([]byte(exampleFileStr))
	if err != nil {
		log.Fatal(err)
	}

	Convey("Get", t, func() {
		Convey("It returns value set for env", func() {
			val, err := r.Get("database")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "development")
		})

		Convey("It returns error for arg that's not defined", func() {
			_, err := r.Get("random")
			So(err, ShouldEqual, ErrKeyNotFound)
		})

		Convey("It falls back on default if key not defined in env", func() {
			val, err := r.Get("adapter")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "postgresql")
		})
	})

	Convey("GetString", t, func() {
		Convey("It returns string value set for env", func() {
			val, err := r.GetString("database")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "development")
			So(val, ShouldHaveSameTypeAs, "")
		})
	})

	Convey("MustGet", t, func() {
		Convey("It returns value set for env", func() {
			So(r.MustGet("database"), ShouldEqual, "development")
		})

		Convey("It panics if key is not found", func() {
			So(func() { r.MustGet("nonexistent") }, ShouldPanic)
		})
	})

	Convey("MustGetString", t, func() {
		Convey("It returns value set for env", func() {
			So(r.MustGetString("database"), ShouldEqual, "development")
		})

		Convey("It panics if key is not found", func() {
			So(func() { r.MustGetString("nonexistent") }, ShouldPanic)
		})
	})

	Convey("SetEnv", t, func() {
		Convey("It returns Development as default if env is not set", func() {
			So(r.GetEnv(), ShouldEqual, Development)
		})

		Convey("It returns Test env after Test is set", func() {
			r.SetEnv(Test)
			So(r.GetEnv(), ShouldEqual, Test)

			val, err := r.Get("database")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "test")

			r.SetEnv(Development) // cleanup
		})
	})
}
