package queue

type Queue interface {
	Enqueue(jobID string) error
	Dequeue() string
}

type JobQueue struct {
	q Queue
}

func NewJobQueue(addr string) (*JobQueue, error) {
	client, err := NewRedisClient(addr)
	if err != nil {
		return nil, err
	}
	return &JobQueue{q: client.NewQueue("jobs")}, nil
}

func (j *JobQueue) Enqueue(jobID string) error {
	return j.q.Enqueue(jobID)
}

func (j *JobQueue) Dequeue() string {
	return j.q.Dequeue()
}
