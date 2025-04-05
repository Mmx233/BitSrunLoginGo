package webhook

import (
	"container/list"
	"context"
	"github.com/Mmx233/BackoffCli/backoff"
	"github.com/Mmx233/BitSrunLoginGo/internal/config/keys"
	log "github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
	"time"
)

func NewEventQueue(logger log.FieldLogger, webhook Webhook) EventQueue {
	queue := EventQueue{
		lock:     &sync.Mutex{},
		list:     list.New(),
		evChan:   make(chan Event),
		activate: make(chan struct{}),
		running:  &atomic.Bool{},
		webhook:  webhook,
		Logger:   logger,
	}

	go queue._LoopReceive()
	go queue._LoopConsume()

	return queue
}

type EventQueue struct {
	lock *sync.Mutex
	list *list.List // Event

	evChan   chan Event
	activate chan struct{}
	running  *atomic.Bool

	webhook Webhook
	Logger  log.FieldLogger
}

func (q EventQueue) _PushList(ev Event) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.list.PushBack(ev)
	select {
	case q.activate <- struct{}{}:
	default:
	}
}

func (q EventQueue) _PopList() Event {
	q.lock.Lock()
	el := q.list.Front()
	if el == nil {
		q.lock.Unlock()
		<-q.activate
		return q._PopList()
	}
	q.list.Remove(el)
	q.running.Store(true)
	q.lock.Unlock()
	return el.Value.(Event)
}

func (q EventQueue) _LoopReceive() {
	for {
		ev, ok := <-q.evChan
		if !ok {
			return
		}
		q._PushList(ev)
	}
}

func (q EventQueue) _LoopConsume() {
	for {
		ev := q._PopList()

		backoffInstance := backoff.NewInstance(func(ctx context.Context) error {
			return q.webhook.Send(ctx, ev)
		}, backoff.Conf{
			Logger: q.Logger.WithFields(log.Fields{
				keys.LogLoginModule: "backoff",
				"eventID":           ev.GetID(),
			}),
			InitialDuration: time.Second,
			MaxDuration:     time.Second * 30,
			ExponentFactor:  1,
		})
		err := backoffInstance.Run(context.Background())
		if err != nil {
			// should always be nil
			panic(err)
		}
		q.running.Store(false)
	}
}

func (q EventQueue) AddEvent(ev Event) {
	q.evChan <- ev
}

func (q EventQueue) Close(ctx context.Context) error {
	close(q.evChan)
wait:
	select {
	case q.activate <- struct{}{}:
		q.lock.Lock()
		isClosed := q.list.Len() == 0 && !q.running.Load()
		q.lock.Unlock()
		if isClosed {
			return nil
		}
		goto wait
	case <-ctx.Done():
		return ctx.Err()
	}
}
