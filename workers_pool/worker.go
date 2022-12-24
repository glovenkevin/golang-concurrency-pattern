package workers_pool

type WorkerLauncherInterface interface {
	LaunchWorker(in chan Request)
}
