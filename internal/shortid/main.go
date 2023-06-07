package shortid

import (
	"github.com/teris-io/shortid"
)

func init() {
	sid, _ := shortid.New(1, shortid.DefaultABC, 2342)
	shortid.SetDefault(sid)
}

// GenerateShortID generates Random Short ID
func GenerateShortID() string {
	id, _ := shortid.Generate()
	return id
}
