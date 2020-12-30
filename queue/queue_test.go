package queue

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
)

// Data ...
type Data struct {
	opts     options
	UUID     string
	CreateAt time.Time
	ExpireAt *time.Time
	Value    interface{}
}

// NewData ...
func NewData(v interface{}) *Data {
	uuid, _ := uuid.NewUUID()
	return &Data{
		opts:     newOptions(Name("Default")),
		UUID:     uuid.String(),
		CreateAt: time.Now(),
		Value:    v,
	}
}

var que Server

type Unit struct {
	Name   string
	Price  float64
	Amount int
	UID    int
	ID     int
}

func TestMain(t *testing.M) {
	que = NewQueue(Name("BUY"))

	t.Run()
}

func TestPush(t *testing.T) {
	rand.Seed(time.Now().Unix())
	node := NewData(&Unit{
		Name:   "qwe",
		Amount: rand.Intn(1000),
		Price:  1.0,
		UID:    0,
		ID:     0,
	})
	que.PushBack(node)

	t.Run("TestPrint", TestQueuePrint)
}

func TestPushFront(t *testing.T) {
	rand.Seed(time.Now().Unix())
	data := NewData(&Unit{
		Name:   "asd",
		Amount: int(rand.Intn(1000)),
		Price:  2.0,
		UID:    1,
		ID:     1,
	})
	que.PushFront(data)
	t.Run("TestPrint", TestQueuePrint)
}

func TestQueuePrint(t *testing.T) {
	for node := que.Front(); node != nil; node = node.Next() {
		data := node.Value
		t.Log(data)
	}
}

func BenchmarkQueue(t *testing.B) {
	t.Run("PushFront", BenchmarkPushFront)
	t.Run("PushBack", BenchmarkPushBack)
	t.RunParallel(func(p *testing.PB) {
		for p.Next() {
			price, _ := strconv.ParseFloat(strconv.Itoa(rand.Intn(10000000000)), 64)
			que.PushFront(NewData(&Unit{
				Name:   "xlj",
				Amount: int(rand.Intn(10000000000)),
				Price:  price,
				UID:    int(rand.Intn(10000000000)),
				ID:     int(rand.Intn(10000000000)),
			}))
		}
	})
}

func BenchmarkPushFront(t *testing.B) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < t.N; i++ {
		price, _ := strconv.ParseFloat(strconv.Itoa(i), 64)
		que.PushFront(NewData(&Unit{
			Name:   "xlj",
			Amount: int(rand.Intn(10000000000)),
			Price:  price,
			UID:    int(i),
			ID:     int(i),
		}))
	}
	t.Log(t.N, que.Len(), "success")
}

func BenchmarkPushBack(t *testing.B) {
	for i := 0; i < t.N; i++ {
		price, _ := strconv.ParseFloat(strconv.Itoa(i), 64)
		que.PushBack(NewData(&Unit{
			Name:   "xlj",
			Amount: int(rand.Intn(10000000000)),
			Price:  price,
			UID:    int(i),
			ID:     int(i),
		}))
	}
	t.Log(t.N, que.Len(), "success")
}
