package learnermodel

import "math"

// The constants of the package under exported names, for its external tests: a test that holds
// the model to one of them takes it by name.
const (
	SecondsPerDay = secondsPerDay

	BandStep         = bandStep
	StartSpread      = startSpread
	LevelStep0       = levelStep0
	CorrectionStep0  = correctionStep0
	StepDecay        = stepDecay
	LevelStepFloor   = levelStepFloor
	Information      = information
	CorrectionPrior  = correctionPrior
	MasteryThreshold = masteryThreshold
	TrialWeight      = trialWeight
	MasteryWeight    = masteryWeight
	WeightTolerance  = weightTolerance
	SeriesAnswers    = seriesAnswers
	FailuresToLose   = failuresToLose
	StartBelow       = startBelow
	GuardFailures    = guardFailures
	LevelFloor       = levelFloor
	HalfBand         = halfBand

	StabilityFloor  = stabilityFloor
	StabilityGrowth = stabilityGrowth
	StabilityLoss   = stabilityLoss
	LevelPause      = levelPause
	LevelWeight     = levelWeight
	MaxPauseSeconds = maxPauseSeconds

	FirstPlace       = firstPlace
	MinDifficulty    = minDifficulty
	MiddleDifficulty = middleDifficulty
	MaxDifficulty    = maxDifficulty
	MinChoiceOptions = minChoiceOptions
	FirstKind        = firstKind
)

// NaN is the NaN that math.NaN gives, by its bits: the tests take math.NaN no more than the code
// does.
var NaN = math.Float64frombits(0x7FF8000000000001)
