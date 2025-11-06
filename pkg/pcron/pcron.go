package pcron

import (
	"fmt"

	"github.com/robfig/cron/v3"

	"github.com/lemonkingstar/spider/pkg/phash"
	"github.com/lemonkingstar/spider/pkg/pmaster"
)

func New() *CronServer {
	return &CronServer{
		name: "cron-server",
		c:    cron.New(),
		sf:   &phash.SnowFlake{},
	}
}

type job struct {
	name     string
	spec     string
	onStart  bool
	onMaster bool
	f        func()
	id       cron.EntryID
}

func (receiver *job) Name() string {
	return receiver.name
}

func (receiver *job) Run() {
	if receiver.onMaster {
		pmaster.Exec(receiver.f)
	} else {
		receiver.f()
	}
}

type CronServer struct {
	name string
	jobs []*job
	c    *cron.Cron
	sf   *phash.SnowFlake
}

func (s *CronServer) Name() string {
	return s.name
}

func (s *CronServer) Start() error {
	err := s.sf.Init(1, 1)
	if err != nil {
		return err
	}

	for _, job := range s.jobs {
		if job.onStart {
			job.f()
		}
		job.id, err = s.c.AddJob(job.spec, job)
		if err != nil {
			return err
		}
	}
	s.c.Start()
	return nil
}

func (s *CronServer) Stop() error {
	s.c.Stop()
	return nil
}

func (s *CronServer) GracefulStop() error {
	ctx := s.c.Stop()
	<-ctx.Done()
	return nil
}

func (s *CronServer) AddJob(spec string, f func()) {
	uuid, _ := s.sf.NextUUID()
	s.jobs = append(s.jobs, &job{
		name:     fmt.Sprintf("job-%s", uuid),
		spec:     spec,
		onStart:  false,
		onMaster: false,
		f:        f,
	})
}

func (s *CronServer) AddJobDetail(name, spec string, start, master bool, f func()) {
	s.jobs = append(s.jobs, &job{
		name:     name,
		spec:     spec,
		onStart:  start,
		onMaster: master,
		f:        f,
	})
}
