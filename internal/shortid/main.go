package shortid

import (
	"math/rand"

	"github.com/teris-io/shortid"
)

func init() {
	sid, _ := shortid.New(1, shortid.DefaultABC, rand.Uint64())
	shortid.SetDefault(sid)
}

// GenerateShortID generates Random Short ID
func GenerateShortID() string {
	id, _ := shortid.Generate()
	return id
}
