package fromfile

// run as "pouch -pkg=github.com/ttacon/pouch/pouch/examples/fromfile"

type CoolStruct struct {
	Name string
	ID   int
	Yolo float64 `db:"yolo"`
}
