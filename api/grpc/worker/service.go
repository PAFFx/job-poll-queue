package worker

import (
	"context"

	"github.com/PAFFx/job-poll-queue/proto/worker"
	"github.com/PAFFx/job-poll-queue/queue"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service implements the WorkerService gRPC service
type Service struct {
	worker.UnimplementedWorkerServiceServer
	jobQueue *queue.Queue
}

// NewService creates a new worker service with the given job queue
func NewService(jobQueue *queue.Queue) *Service {
	return &Service{
		jobQueue: jobQueue,
	}
}

// RequestJob handles job requests from workers
func (s *Service) RequestJob(ctx context.Context, req *worker.JobRequest) (*worker.Job, error) {
	// Pop a job from the queue
	msg, err := s.jobQueue.Pop()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve job: %v", err)
	}

	// No jobs available
	if msg == nil {
		return nil, status.Errorf(codes.NotFound, "no jobs available")
	}

	// Convert queue.Message to worker.Job
	job := &worker.Job{
		Id:      msg.ID,
		Payload: msg.Payload,
		Headers: msg.Headers,
	}

	return job, nil
}

// CompleteJob handles job completion reports from workers
func (s *Service) CompleteJob(ctx context.Context, result *worker.JobResult) (*worker.CompleteResponse, error) {
	// Get the job status manager
	statusMgr := s.jobQueue.GetStatusManager()

	// Submit the result to the status manager
	err := statusMgr.SubmitResult(result.JobId, result.Payload, nil)
	if err != nil {
		return &worker.CompleteResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "failed to save job result: %v", err)
	}

	return &worker.CompleteResponse{
		Success: true,
	}, nil
}
