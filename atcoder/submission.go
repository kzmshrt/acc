package atcoder

type Judge string

const (
	JudgeNA  Judge = "Unknown"
	JudgeWJ  Judge = "WJ"  // Waiting Judge
	JudgeAC  Judge = "AC"  // Accepted
	JudgeWA  Judge = "WA"  // Wrong Answer
	JudgeTLE Judge = "TLE" // Time Limit Exceed
	JudgeMLE Judge = "MLE" // Memory Limit Exceed
	JudgeRE  Judge = "RE"  // Runtime Error
	JudgeCE  Judge = "CE"  // Compile Error
)

type Submission struct {
	Judge       Judge
	TimeScore   int
	MemoryScore int
	CodeSize    int
	DetailURL   string
}
