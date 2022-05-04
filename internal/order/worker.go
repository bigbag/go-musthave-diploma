package order

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/bigbag/go-musthave-diploma/internal/accrual"
)

type Task struct {
	OrderID  string
	RepeatAt int64
}

func NewTask(orderID string) *Task {
	return &Task{
		OrderID:  orderID,
		RepeatAt: time.Now().Unix(),
	}
}

type Queue struct {
	arr  []*Task
	mu   sync.Mutex
	cond *sync.Cond
	stop bool
}

func (q *Queue) close() {
	q.cond.L.Lock()
	q.stop = true
	q.cond.Broadcast()
	q.cond.L.Unlock()
}

func (q *Queue) PopWait() (*Task, bool) {
	q.cond.L.Lock()

	for len(q.arr) == 0 && !q.stop {
		q.cond.Wait()
	}

	if q.stop {
		q.cond.L.Unlock()
		return nil, false
	}

	t := q.arr[0]
	q.arr = q.arr[1:]

	q.cond.L.Unlock()

	return t, true
}

type TaskPool struct {
	l          logrus.FieldLogger
	or         *Repository
	ar         *accrual.Repository
	workerPool []*TaskWorker
	wg         *sync.WaitGroup
	queue      *Queue
	total      chan int
}

func NewTaskPool(
	ctx context.Context,
	l logrus.FieldLogger,
	or *Repository,
	ar *accrual.Repository,
) *TaskPool {

	p := &TaskPool{l: l, or: or, ar: ar}
	p.workerPool = make([]*TaskWorker, 0, runtime.NumCPU())

	orderIDs, _ := or.GetAllForChecking()
	tasks := make([]*Task, 0, 100)
	for _, orderID := range orderIDs {
		tasks = append(tasks, NewTask(orderID))
	}
	p.queue = p.newQueue(tasks)

	for i := 0; i < runtime.NumCPU(); i++ {
		p.workerPool = append(p.workerPool, p.newWorker(i))
	}

	ctx, cancel := context.WithCancel(ctx)
	g, _ := errgroup.WithContext(ctx)
	p.wg = &sync.WaitGroup{}

	for _, w := range p.workerPool {
		p.wg.Add(1)
		worker := w
		f := func() error {
			return worker.loop(ctx)
		}
		g.Go(f)
	}

	go func() {
		if err := g.Wait(); err != nil {
			p.l.Info("worker: pool error ", err)
		}
	}()
	go func() {
		p.wg.Wait()
		close(p.total)
		cancel()
	}()

	p.total = make(chan int)
	go func() {
		total := 0
		for c := range p.total {
			total = total + c
		}
	}()

	return p
}

func (p *TaskPool) newQueue(arr []*Task) *Queue {
	q := Queue{arr: arr, stop: false}
	q.cond = sync.NewCond(&q.mu)
	return &q
}

func (p *TaskPool) newWorker(id int) *TaskWorker {
	p.l.Info("worker: init ", id)
	return &TaskWorker{id, p}
}

func (p *TaskPool) Push(t *Task) error {
	if p.queue.stop {
		return errors.New("worker: queue was stopped")
	}

	p.queue.cond.L.Lock()
	defer p.queue.cond.L.Unlock()

	p.queue.arr = append(p.queue.arr, t)
	p.queue.cond.Signal()
	return nil
}

func (p *TaskPool) Close() {
	p.queue.close()
}

type TaskWorker struct {
	id   int
	pool *TaskPool
}

func (w *TaskWorker) process(t *Task) error {
	delta := t.RepeatAt - time.Now().Unix()
	if delta > 0 {
		time.Sleep(time.Duration(delta) * time.Second)
		w.pool.Push(t)
		return nil
	}

	info, err := w.pool.ar.Get(t.OrderID)
	if err != nil {
		return err
	}
	w.pool.l.Info("worker: accrual response ", info)

	if info.IsFinal {
		w.pool.l.Info("worker: save final state of order ", t.OrderID)

		err = w.pool.or.Update(
			info.OrderID,
			info.Status,
			info.Amount,
			info.IsFinal,
		)
		if err != nil {
			return err
		}
	}

	w.pool.Push(&Task{
		OrderID:  t.OrderID,
		RepeatAt: time.Now().Unix() + info.Timeout,
	})

	return nil
}

func (w *TaskWorker) loop(ctx context.Context) error {
	defer func() {
		w.pool.wg.Done()
		w.pool.queue.close()

		<-ctx.Done()
	}()

	for {
		t, ok := w.pool.queue.PopWait()
		if !ok {
			return nil
		}

		w.pool.l.Info("worker: new task ", w.id)
		err := w.process(t)
		if err != nil {
			w.pool.l.Info("worker: run to out from loop ")
			return err
		}

		w.pool.total <- 1
	}
}
