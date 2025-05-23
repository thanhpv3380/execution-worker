package services

type RunnerService interface {
	Run(code string) error
}
