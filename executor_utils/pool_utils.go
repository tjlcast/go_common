package executor_utils

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Task struct {
	Id string
	Handler  func(v ...interface{})
	Params   []interface{}
	Interupt bool
}

type Pool struct {
	capacity       uint64
	runningWorkers uint64
	status         int64
	chTask         chan *Task
	sync.Mutex

	PanicHandler func(interface{})
}

var ErrInvalidPoolCap = errors.New("invalid pool cap")

const (
	RUNNING = 1
	STOPED  = 0
)

func NewPool(capacity uint64) (*Pool, error) {
	if capacity <= 0 {
		return nil, ErrInvalidPoolCap
	}
	return &Pool{
		capacity: capacity,
		status:   RUNNING,
		// 初始化任务队列, 队列长度为容量
		chTask: make(chan *Task, capacity),
	}, nil
}

func (p *Pool) run() {
	p.incRunning()

	go func() {
		defer func() {
			p.decRunning()
			if r := recover(); r != nil {
				buf := make([]byte, 1<<16)
				runtime.Stack(buf, true)
				if p.PanicHandler != nil {
					p.PanicHandler(r)
				} else {
					sprintf := fmt.Sprintf("Pool-worker panic: error: %s.\n stack: %s.\n", string(buf), r)
					r = WrapError(sprintf, r.(error))
					log.Printf("%s", r)
				}
			}
			p.checkWorker() // worker 退出时检测是否有可运行的 worker
		}()

		for {
			select {
			case task, ok := <-p.chTask:
				if !ok {
					return
				}
				task.Handler(task.Params...)
			}
		}
	}()
}

func (p *Pool) incRunning() { // runningWorkers + 1
	atomic.AddUint64(&p.runningWorkers, 1)
}

func (p *Pool) decRunning() { // runningWorkers - 1
	atomic.AddUint64(&p.runningWorkers, ^uint64(0))
}

func (p *Pool) GetRunningWorkers() uint64 {
	return atomic.LoadUint64(&p.runningWorkers)
}

func (p *Pool) GetCap() uint64 {
	return p.capacity
}

func (p *Pool) setStatus(status int64) bool {
	p.Lock()
	defer p.Unlock()

	if p.status == status {
		return false
	}

	p.status = status

	return true
}

var ErrPoolAlreadyClosed = errors.New("pool already closed")

/**
If there is no valiable thread, `put` will block the current.
 */
func (p *Pool) Put(task *Task) error {
	p.Lock()
	defer p.Unlock()

	if p.status == STOPED { // 如果任务池处于关闭状态, 再 put 任务会返回 ErrPoolAlreadyClosed 错误
		return ErrPoolAlreadyClosed
	}

	// run worker
	if p.GetRunningWorkers() < p.GetCap() {
		// 此时有一个 task (上一次 Put) panic，worker 退出了
		p.run()
	}

	// send task
	if p.status == RUNNING {
		p.chTask <- task // 当前的 task 推送到 chTask，但是没有一个 worker 可以消费到，deadlock!
	}

	return nil
}

func (p *Pool) Close() {
	p.setStatus(STOPED) // 设置 status 为已停止

	for len(p.chTask) > 0 { // 阻塞等待所有任务被 worker 消费
		time.Sleep(1e6) // 防止等待任务清空 cpu 负载突然变大, 这里小睡一下
	}

	close(p.chTask) // 关闭任务队列
}

func (p *Pool) checkWorker() {
	p.Lock()
	defer p.Unlock()

	// 当前没有 worker 且有任务存在，运行一个 worker 消费任务
	// 没有任务无需考虑 (当前 Put 不会阻塞，下次 Put 会启动 worker)
	if p.runningWorkers == 0 && len(p.chTask) > 0 {
		p.run()
	}
}

func WrapError(wrapMsg string, err error) error {
	pc, file, line, ok := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	if !ok {
		return errors.New("WrapError 方法获取堆栈失败")
	}
	if err == nil {
		return nil
	} else {
		errMsg := fmt.Sprintf("%s \n\tat %s:%d (Method %s)\nCause by: %s\n", wrapMsg, file, line, f.Name(), err.Error())
		return errors.New(errMsg)
	}
}
