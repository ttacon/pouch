package impl

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDynamicPouch(t *testing.T) {
	Convey("given a dynamic pouch", t, func() {
		d := NewDynamicPouch(nil)
		Convey("with no defined functions", func() {
			Convey("it should error out when trying to do anything", func() {
				err := d.Find(nil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "no Find function has been defined")
			})
		})
	})
}
