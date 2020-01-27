package evolution

import (
	"fmt"
	"github.com/martinomburajr/masters-go/evolog"
	"strings"
)

type EvolutionParams struct {
	Name string

	Topology Topology `json:"topology"`
	// StartIndividual - Output Only - This is set by the SpecParam Expression. Do not set it manually
	StartIndividual Program
	// Spec - Output Only - This is set by the SpecParam Expression. Do not set it manually
	Spec      SpecMulti `json:"spec"`
	SpecParam SpecParam `json:"specParam"`
	// MaxGenerations activates the ability for a variable number of generations before the simulation ends.
	// The value must be greater than 9 for the activation to begin, if not,
	// the simulation will default to GenerationsCount number of generations. Once this variable is set,
	// MinimumTopProtagonistMeanBeforeTerminate and ProtagonistMinGenAvgFit will come into effect. If no adequate solution is found,
	// MaxGenerations will terminate. This value should default to about 300.
	MaxGenerations int `json:"maxGenerationsCount",csv:"maxGenerationsCount"`
	// MinimumTopProtagonistMeanBeforeTerminate specifies the percentage of consecutive generations where the
	// ProtMinGenFitnessAvg has been hit by the best protaginist in the generation before the simulation can end.
	MinimumTopProtagonistMeanBeforeTerminate float64 `json:"minimumTopProtagonistMeanBeforeTerminate"`

	// MinimumGenerationMeanBeforeTerminate specifies the percentage of consecutive generations where the
	// ProtMinGenFitnessAvg has been hit by the average of all protagonists in the generation before the simulation can
	// end.
	MinimumGenerationMeanBeforeTerminate float64 `json:"minimumGenerationMeanBeforeTerminate"`

	// ProtagonistMinimumGenFitness specifies the average value of fitness of the best individual after a completed
	// generation. This individual must obtain this fitness value or greater e.g. an average of 0.75
	// for MinimumTopProtagonistMeanBeforeTerminate number of consecutive generations before the simulation can end.
	ProtagonistMinGenAvgFit float64 `json:"protagonistMinGenAvgFit"`
	GenerationsCount        int     `json:"generationCount",csv:"generationCount"`
	// EachPopulationSize represents the size of each protagonist or antagonist population.
	// This value must be even otherwise pairwise operations such as crossover will fail
	EachPopulationSize int  `json:"eachPopulationSize",csv:"eachPopulationSize"`
	EnableParallelism  bool `json:"enableParallelism",csv:"enableParallelism"`

	Strategies Strategies `json:"strategies",csv:"strategies"`

	FitnessStrategy FitnessStrategy `json:"fitnessStrategy",csv:"fitnessStrategy"`
	Reproduction    Reproduction    `json:"reproduction",csv:"reproduction"`
	Selection       Selection       `json:"selection",csv:"selection"`

	// FitnessCalculatorType allows user to select the fitness calculator.
	// The more complex the function 1 is better but slower. 0 for simple polynomials with single digit constants e.
	// g. x*x*x or x*x+4
	FitnessCalculatorType int `json:"fitnessCalculatorType"`
	// ShouldRunInteractiveTerminal ensures the interactive terminal is run at the end of the evolution that allows
	// users to query all individuals in all generations.
	ShouldRunInteractiveTerminal bool             `json:"shouldRunInteractiveTerminal"`
	StatisticsOutput             StatisticsOutput `json:"statisticsOutput"`
	// InternalCount - Output Only (Helps with file name assignments)
	InternalCount int

	EnableLogging bool `json:"-"`
	RunStats      bool `json:"-"`

	// FinalGeneration records if the simulation ended early by fulfilling the maxGen requirements.
	FinalGeneration       int    `json:"finalGeneration",csv:"finalGeneration"`
	FinalGenerationReason string `json:"finalGenerationReason",csv:"finalGenerationReason"`

	//Channels
	LoggingChan chan evolog.Logger `json:"-"`
	ErrorChan   chan error         `json:"-"`
	DoneChan    chan bool          `json:"-"`
	ParamFile   string             `json:"-"`

	//FolderPercentages help track progress when a certain percentage is reached
	FolderPercentages []float64
}

type Topology struct {
	Type     string `json:"type"`
	KRandomK int    `json:"kRandomK"`
}

type Generations struct {
}

type StatisticsOutput struct {
	OutputPath string `json:"outputPath"`
	Name       string `json:"name"`
	OutputDir  string `json:"outputDir"`
}

type AvailableVariablesAndOperators struct {
	Constants []string `json:"constants"`
	Variables []string `json:"variables"`
	Operators []string `json:"operators"`
}

type AvailableSymbolicExpressions struct {
	//Constants []SymbolicExpression
	NonTerminals []SymbolicExpression
	Terminals    []SymbolicExpression
}

type Strategies struct {
	//AvailableStrategies            []Strategy `json:"availableStrategies"`
	AntagonistAvailableStrategies  []Strategy `json:"antagonistAvailableStrategies"`
	ProtagonistAvailableStrategies []Strategy `json:"protagonistAvailableStrategies"`

	AntagonistStrategyCount  int `json:"antagonistStrategyCount"`
	ProtagonistStrategyCount int `json:"protagonistStrategyCount"`

	DepthOfRandomNewTrees int `json:"depthOfRandomNewTrees"`
}

type FitnessStrategy struct {
	Type string `json:"type"`
	// AntagonistThresholdMultiplier is the multiplier applied to the antagonist delta when calculating fitness.
	// A large value means that antagonists have to attain a greater delta from the spec in order to gain adequate
	// fitness, conversely a smaller value gives the antagonists more slack to not manipulate the program excessively.
	// For good results set it to a value greater than that of the protagonist delta.
	// This value is only used when using DualThresholdedRatioFitness.
	AntagonistThresholdMultiplier float64 `json:"antagonistThresholdMultiplier"`

	// ProtagonistThresholdMultiplier is the multiplier applied to the protagonist delta when calculating fitness.
	// A large value means that protagonist can be less precise and gain adequate fitness,
	// conversely a smaller value gives the protagonist little room for mistake between its delta and that of the spec.
	// this value is used in both DualThresholdedRatioFitness and ThresholdedRatioFitness as a fitness value for
	// both antagonist and protagonists thresholds.
	ProtagonistThresholdMultiplier float64 `json:"protagonistThresholdMultiplier"`
}

type SpecParam struct {
	// SpecRange defines a range of variables on either side of the X axis. A range of 4 will include -2, -1,
	// 0 and 1.
	Range int `json:"range"`
	//Expression is the actual expression being tested.
	// It is the initial function that is converted to the startIndividual
	Expression string `json:"expression"`
	//OUTPUT
	ExpressionParsed               string `json:"expressionParsed"`
	Seed                           int    `json:"seed"`
	AvailableVariablesAndOperators AvailableVariablesAndOperators
	// AvailableSymbolicExpressions - Output Only
	AvailableSymbolicExpressions AvailableSymbolicExpressions
	DivideByZeroStrategy         string  `json:"divideByZeroStrategy",csv:"divideByZeroStrategy"`
	DivideByZeroPenalty          float64 `json:"divideByZeroPenalty",csv:"divideByZeroPenalty"`
}

type Reproduction struct {
	CrossoverStrategy string `json:"crossoverStrategy",csv:"crossoverStrategy"`
	// CrossoverPercentrage pertains to the amount of genetic material crossed-over. FOR SPX
	// This is a percentage represented as a float64. A value of 1 means all material is swapped.
	// A value of 0 means no material is swapped (which in effect are the same thing).
	// Avoid 0 or 1 use values in between
	CrossoverPercentage   float64 `json:"crossoverPercentage",csv:"crossoverPercentage"`
	ProbabilityOfMutation float64 `json:"probabilityOfMutation",csv:"probabilityOfMutation"`
	KPointCrossover       int     `json:"kPointCrossover",csv:"kPointCrossover"`
}
type Selection struct {
	Parent   ParentSelection   `json:"parentSelection",csv:"parentSelection"`
	Survivor SurvivorSelection `json:"survivorSelection",csv:"survivorSelection"`
}

type ParentSelection struct {
	Type           string `json:"type",csv:"type"`
	TournamentSize int    `json:"tournamentSize",csv:"tournamentSize"`
}

type SurvivorSelection struct {
	Type string `json:"type",csv:"type"`
	// SurvivorPercentage represents how many individulas in the parent vs child population should continue.
	// 1 means all parents move on. 0 means only children move on. Any number in betwee is a percentage value.
	// It cannot be greater than 1 or less than 0.
	SurvivorPercentage float64 `json:"survivorPercentage",csv:"survivorPercentage"`
}

func (e EvolutionParams) ToString() string {
	builder := strings.Builder{}
	//Input Program
	expressionStr := strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(e.SpecParam.Expression, "*", ""),
				"+", "+"),
			"-", "-"),
		"/", "DIV")
	builder.WriteString(fmt.Sprintf("%sR%dS%d", expressionStr, e.SpecParam.Range, e.SpecParam.Seed))
	builder.WriteString("-")
	// GenCount
	builder.WriteString(fmt.Sprintf("G%d", e.GenerationsCount))
	builder.WriteString("-")
	// Population Size
	builder.WriteString(fmt.Sprintf("P%d", e.EachPopulationSize))
	builder.WriteString("-")
	// Fitness
	fitness := strings.ReplaceAll(e.FitnessStrategy.Type, "Fitness", "")
	builder.WriteString(strings.ReplaceAll(fmt.Sprintf("F%sa%.2fp%.2f",
		fitness[:len(fitness)/2], e.FitnessStrategy.AntagonistThresholdMultiplier,
		e.FitnessStrategy.ProtagonistThresholdMultiplier), ".", ""))
	builder.WriteString("-")
	//Parent
	builder.WriteString(fmt.Sprintf("P%sTornSz%d", e.Selection.Parent.Type[0:2], e.Selection.Parent.TournamentSize))
	builder.WriteString("-")
	builder.WriteString(fmt.Sprintf("Tree%d", e.Strategies.DepthOfRandomNewTrees))
	builder.WriteString("-")
	//Survivor
	builder.WriteString(strings.ReplaceAll(fmt.Sprintf("S%sPr%.2f", e.Selection.Survivor.Type[0:2],
		e.Selection.Survivor.SurvivorPercentage), ".", ""))
	builder.WriteString("-")
	// ReproductionPercentage
	builder.WriteString(strings.ReplaceAll(fmt.Sprintf("Cro%.2fMut%.2f", e.Reproduction.CrossoverPercentage,
		e.Reproduction.ProbabilityOfMutation), ".", ""))
	builder.WriteString("-")
	// StrategyCount
	builder.WriteString(fmt.Sprintf("PSc%dASc%d", e.Strategies.ProtagonistStrategyCount,
		e.Strategies.AntagonistStrategyCount))
	//antStrat := TruncShort(e.Strategies.AntagonistAvailableStrategies)
	//proStrat := TruncShort(e.Strategies.ProtagonistAvailableStrategies)
	//builder.WriteString(fmt.Sprintf("AAvaiSt%sAvaiSt%s", antStrat, proStrat))
	builder.WriteString("-")

	// Spec
	divide0Penalty := fmt.Sprintf("D0P%.2fD0S%s", e.SpecParam.DivideByZeroPenalty,
		e.SpecParam.DivideByZeroStrategy)
	divide0Penalty = strings.ReplaceAll(divide0Penalty, ".", "")
	builder.WriteString(divide0Penalty)
	//builder.WriteString("-")
	//builder.WriteString(fmt.Sprintf("%s", RandString(4)))

	return builder.String()
}
