package fetcher

import (
	"context"
	"time"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/indexing-engine/worker"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
)

// Manager represents a fetcher manager
type Manager struct {
	cfg    *config.Config
	pool   *worker.Pool
	store  store.Jobs
	client figmentclient.Client

	backoffs map[model.JobID]*worker.Backoff
}

// NewManager creates a fetcher manager
func NewManager(cfg *config.Config, pool *worker.Pool, store store.Jobs, client figmentclient.Client) (*Manager, error) {
	manager := Manager{
		cfg:    cfg,
		pool:   pool,
		store:  store,
		client: client,

		backoffs: make(map[model.JobID]*worker.Backoff),
	}

	return &manager, nil
}

// Run starts the fetcher manager
func (m *Manager) Run() error {
	defer m.pool.Stop()

	m.pool.Run(m.handleResponse)

	for {
		jobs, err := m.getJobs()
		if err != nil {
			return err
		}

		for _, job := range jobs {
			if err := m.scheduleJob(&job); err != nil {
				return err
			}
		}

		m.pool.Wait()
	}
}

func (m *Manager) getJobs() ([]model.Job, error) {
	jobs, err := m.store.FindAllUnfinished()
	if err != nil {
		return nil, err
	}

	if len(jobs) == 0 {
		return m.createJobs()
	}

	return jobs, nil
}

func (m *Manager) createJobs() ([]model.Job, error) {
	var jobs []model.Job

	hr, err := m.getHeightRange()
	if err != nil {
		return nil, err
	}

	for h := hr.StartHeight(); h <= hr.EndHeight(); h++ {
		height := h
		jobs = append(jobs, model.Job{Height: &height})
	}

	err = m.store.Create(jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (m *Manager) getHeightRange() (*pipeline.HeightRange, error) {
	chainStatus, err := m.client.GetChainStatus(context.Background())
	if err != nil {
		return nil, err
	}
	latestHeight := chainStatus.LastBlockHeight

	lastHeight, err := m.store.LastFinishedHeight()
	if err != nil {
		lastHeight = -1
	}

	hr := pipeline.HeightRange{
		LatestHeight:  latestHeight,
		LastHeight:    lastHeight,
		InitialHeight: m.cfg.FirstBlockHeight,
		BatchSize:     m.cfg.DefaultBatchSize,
	}

	err = hr.Validate(false /* checkLength */)
	if err != nil {
		return nil, err
	}

	return &hr, nil
}

func (m *Manager) scheduleJob(job *model.Job) error {
	if m.isJobDelayed(job) {
		return nil
	}

	m.pool.Process(*job.Height)

	now := time.Now()

	job.StartedAt = &now
	job.RunCount++

	err := m.store.Update(job, "started_at", "run_count")
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) isJobDelayed(job *model.Job) bool {
	if job.StartedAt == nil {
		return false
	}

	return time.Since(*job.StartedAt) < m.jobBackoff(job).Delay()
}

func (m *Manager) jobBackoff(job *model.Job) *worker.Backoff {
	if m.backoffs[job.ID] == nil {
		m.backoffs[job.ID] = &worker.Backoff{}
	}

	return m.backoffs[job.ID]
}

func (m *Manager) handleResponse(res worker.Response) {
	job, err := m.store.FindByHeight(res.Height)
	if err != nil {
		panic(err)
	}

	if res.Success {
		m.handleSuccess(job)
	} else {
		m.handleFailure(job, res)
	}
}

func (m *Manager) handleSuccess(job *model.Job) {
	now := time.Now()

	job.FinishedAt = &now

	err := m.store.Update(job, "finished_at")
	if err != nil {
		panic(err)
	}
}

func (m *Manager) handleFailure(job *model.Job, res worker.Response) {
	job.LastError = &res.Error

	err := m.store.Update(job, "last_error")
	if err != nil {
		panic(err)
	}

	m.jobBackoff(job).Attempt()
}
