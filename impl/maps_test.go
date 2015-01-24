package impl

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMapPouch_Find(t *testing.T) {
	Convey("when using a map based pouch", t, func() {
		Convey("given a map of foods", func() {
			foods := make(map[string]interface{})
			foods["food:1"] = &mapFood{
				&Food{
					ID:   1,
					Name: "map based",
				},
			}
			p := MapPouch(foods)
			Convey("so we can find a food the normal pouch way", func() {
				m := &mapFood{
					&Food{
						ID: 1,
					},
				}
				err := p.Find(m)
				So(err, ShouldBeNil)
				So(m.Name, ShouldEqual, "map based")
			})
		})
	})
}
