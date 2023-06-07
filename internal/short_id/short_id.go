package short_id

import (
	"github.com/teris-io/shortid"
)

func init() {
	sid, _ := shortid.New(1, shortid.DefaultABC, 2342)
	shortid.SetDefault(sid)
}

func GenerateShortId() string {
	id, _ := shortid.Generate()
	return id
}
