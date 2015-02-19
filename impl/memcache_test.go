package impl

import (
	"fmt"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/ttacon/pouch"
)

func cleanMemcache(loc string) error {
	c := memcache.New(loc)
	err := c.DeleteAll()
	return err
}

func TestMemcachePouch(t *testing.T) {
	Convey("given a memcache pouch", t, func() {
		mC := memcache.New(*memcacheLoc)
		m := NewMemcachePouch(mC)

		Convey("we can store things", func() {
			f := &Food{
				ID:   1,
				Name: "spinach",
			}
			err := m.Create(&memcachedFood{f})
			So(err, ShouldBeNil)
		})

		Convey("we can also store more than one thing at once", func() {
			f0 := &Food{
				ID:   2,
				Name: "spinach",
			}
			f1 := &Food{
				ID:   3,
				Name: "kale",
			}
			err := m.CreateAll([]pouch.Createable{
				&memcachedFood{f0},
				&memcachedFood{f1},
			})
			So(err, ShouldBeNil)
		})

		Convey("we can also retrieve things", func() {
			var found *Food = &Food{
				ID: 3,
			}
			err := m.Find(&memcachedFood{found})
			So(err, ShouldBeNil)
			So(found.Name, ShouldEqual, "kale")

			Convey("we can even retrieve more than one at once", func() {
				f0 := &Food{
					ID: 1,
				}
				f1 := &Food{
					ID: 2,
				}
				err := m.FindAll([]pouch.Findable{
					&memcachedFood{f0},
					&memcachedFood{f1},
				})
				So(err, ShouldBeNil)
				So(f0.Name, ShouldEqual, "spinach")
				So(f1.Name, ShouldEqual, "spinach")
			})
		})

	})
}

type memcachedFood struct {
	*Food
}

func (m *memcachedFood) MemcacheKey() string {
	return fmt.Sprintf("memcache:%d", m.ID)
}
