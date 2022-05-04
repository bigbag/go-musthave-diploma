package order

import (
	"runtime"
	"time"

	"github.com/mborders/artifex"
	"github.com/sirupsen/logrus"

	"github.com/bigbag/go-musthave-diploma/internal/accrual"
)

var (
	workersLimit = runtime.NumCPU()
	tasksLimit   = 100
)

type Task struct {
	OrderID string
}

func NewTask(orderID string) *Task {
	return &Task{OrderID: orderID}
}

type Worker struct {
	l  logrus.FieldLogger
	d  *artifex.Dispatcher
	or *Repository
	ar *accrual.Repository
}

func NewWorker(
	l logrus.FieldLogger,
	or *Repository,
	ar *accrual.Repository,
) *Worker {
	d := artifex.NewDispatcher(workersLimit, tasksLimit)
	d.Start()
	return &Worker{l: l, d: d, or: or, ar: ar}
}

func (w *Worker) Init() error {
	orderIDs, err := w.or.GetAllForChecking()
	for _, orderID := range orderIDs {
		w.Add(NewTask(orderID))
	}
	return err
}

func (w *Worker) Add(t *Task) error {
	var err error
	w.d.Dispatch(func() {
		err = w.process(t)
	})
	return err
}

func (w *Worker) process(t *Task) error {
	info, err := w.ar.Get(t.OrderID)
	if err != nil {
		return err
	}
	w.l.Info("worker: accrual info ", info)

	if info.IsFinal {
		w.l.Info("worker: save final state of order ", t.OrderID)

		err = w.or.Update(
			info.OrderID,
			info.Status,
			info.Amount,
			info.IsFinal,
		)
		if err != nil {
			return err
		}
	}

	w.d.DispatchIn(func() {
		err = w.process(t)
	}, time.Second*info.Timeout)

	return err

}

func (w *Worker) Close() {
	w.d.Stop()
}
