package impl

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/ttacon/pouch"
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

func TestMapPouch_FindEntities(t *testing.T) {
	Convey("when using a map based pouch", t, func() {
		Convey("given a map of foods", func() {
			foods := make(map[string]interface{})
			foods["food:1"] = &mapFood{
				&Food{
					ID:   1,
					Name: "map based",
				},
			}
			foods["food:2"] = &mapFood{
				&Food{
					ID:   2,
					Name: "kale",
				},
			}
			foods["food:3"] = &mapFood{
				&Food{
					ID:   3,
					Name: "spinach",
				},
			}
			p := MapPouch(foods)
			Convey("so we can find a bunch foods the normal pouch way", func() {
				m0 := &mapFood{
					&Food{
						ID: 1,
					},
				}
				m1 := &mapFood{
					&Food{
						ID: 2,
					},
				}
				m2 := &mapFood{
					&Food{
						ID: 3,
					},
				}
				var m []pouch.Findable = []pouch.Findable{m0, m1, m2}

				err := p.FindAll(m)
				So(err, ShouldBeNil)
				So(m0.Name, ShouldEqual, "map based")
				So(m1.Name, ShouldEqual, "kale")
				So(m2.Name, ShouldEqual, "spinach")
			})

			Convey("when a map pouch has a limit, it only queries for the first n entities", func() {
				m0 := &mapFood{
					&Food{
						ID: 1,
					},
				}
				m1 := &mapFood{
					&Food{
						ID: 2,
					},
				}
				m2 := &mapFood{
					&Food{
						ID: 3,
					},
				}
				var m []pouch.Findable = []pouch.Findable{m0, m1, m2}
				err := p.Limit(2).FindAll(m)
				So(err, ShouldBeNil)
				So(m0.Name, ShouldEqual, "map based")
				So(m1.Name, ShouldEqual, "kale")
				So(m2.Name, ShouldEqual, "")
			})
		})
	})
}
