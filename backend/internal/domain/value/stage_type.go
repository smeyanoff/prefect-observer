package value

type StageType string

const (
	ParallelStage   StageType = "PARALLEL"
	SequentialStage StageType = "SEQUENTIAL"
	ObserverStage   StageType = "OBSERVER"
)
