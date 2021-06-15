package launcher

import "log"

type Launcher struct {
	fnchan chan func()
}

func New(n int) *Launcher {
	return &Launcher{
		fnchan: make(chan func(), n),
	}
}

func (l *Launcher) worker(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			l.fnchan <- fn
		}
	}()

	fn()
}

func (l *Launcher) WithFuncs(fns ...func()) {
	for _, fn := range fns {
		l.fnchan <- fn
	}
}

func (l *Launcher) Launch() {
	for {
		select {
		case fn := <-l.fnchan:
			go l.worker(fn)
		}
	}
}
