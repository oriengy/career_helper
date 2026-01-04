package idgen

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/samber/lo"
)

var node *snowflake.Node

func init() {
	nodeID := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(1024)
	node = lo.Must(snowflake.NewNode(nodeID))
}

type ID struct {
	snowflake.ID
}

func (id ID) Int() int {
	return int(id.Int64())
}

func (id ID) String() string {
	return id.Base36()
}

func NewID() ID {
	return ID{ID: node.Generate()}
}

func Int() int {
	return int(node.Generate().Int64())
}

func Uint() uint {
	return uint(node.Generate().Int64())
}

func Base36() string {
	return node.Generate().Base36()
}

func FromInt(idint int) ID {
	return ID{ID: snowflake.ID(idint)}
}

func FromBase36(idstring string) (ID, error) {
	id, err := snowflake.ParseBase36(idstring)
	if err != nil {
		return ID{}, err
	}
	return ID{ID: id}, nil
}

func FromTime(t time.Time) int64 {
	return (t.UnixMilli() - snowflake.Epoch) << 22
}

func ToTime(id int64) time.Time {
	return time.UnixMilli(id>>22 + snowflake.Epoch)
}
