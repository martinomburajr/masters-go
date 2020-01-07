package evolution

import (
	"fmt"
	"github.com/martinomburajr/masters-go/utils"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
)

type EvolutionParams struct {
	Name string
	// StartIndividual - Output Only - This is set by the SpecParam Expression. Do not set it manually
	StartIndividual Program
	// Spec - Output Only - This is set by the SpecParam Expression. Do not set it manually
	Spec             SpecMulti `json:"spec"`
	SpecParam        SpecParam `json:"specParam"`
	GenerationsCount int       `json:"generationCount"`
	// EachPopulationSize represents the size of each protagonist or antagonist population.
	// This value must be even otherwise pairwise operations such as crossover will fail
	EachPopulationSize int  `json:"eachPopulationSize"`
	EnableParallelism  bool `json:"enableParallelism"`

	Strategies Strategies `json:"strategies"`

	FitnessStrategy FitnessStrategy `json:"fitnessStrategy"`
	Reproduction    Reproduction    `json:"reproduction"`
	Selection       Selection       `json:"selection"`

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

	//Channels
	LoggingChan chan string `json:"-"`
	ErrorChan   chan error  `json:"-"`
	DoneChan    chan bool   `json:"-"`
	ParamFile   string `json:"-"`
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
	Expression                     string `json:"expression"`
	//OUTPUT
	ExpressionParsed string `json:"expressionParsed"`
	Seed                           int    `json:"seed"`
	AvailableVariablesAndOperators AvailableVariablesAndOperators
	// AvailableSymbolicExpressions - Output Only
	AvailableSymbolicExpressions AvailableSymbolicExpressions
	DivideByZeroStrategy         string  `json:"divideByZeroStrategy"`
	DivideByZeroPenalty          float64 `json:"divideByZeroPenalty"`
}

type Reproduction struct {
	// CrossoverPercentage pertains to the amount of genetic material crossed-over.
	// This is a percentage represented as a float64. A value of 1 means all material is swapped.
	// A value of 0 means no material is swapped (which in effect are the same thing).
	// Avoid 0 or 1 use values in between
	CrossoverPercentage   float64 `json:"crossoverPercentage"`
	ProbabilityOfMutation float64 `json:"probabilityOfMutation"`
}
type Selection struct {
	Parent   ParentSelection   `json:"parentSelection"`
	Survivor SurvivorSelection `json:"survivorSelection"`
}

type ParentSelection struct {
	Type           string `json:"type"`
	TournamentSize int    `json:"tournamentSize"`
}

type SurvivorSelection struct {
	Type string `json:"type"`
	// SurvivorPercentage represents how many individulas in the parent vs child population should continue.
	// 1 means all parents move on. 0 means only children move on. Any number in betwee is a percentage value.
	// It cannot be greater than 1 or less than 0.
	SurvivorPercentage float64 `json:"survivorPercentage"`
}

type EvolutionEngine struct {
	Generations []*Generation   `json:"generations"`
	Parameters  EvolutionParams `json:"parameters"`
}

func (e *EvolutionEngine) Start() (*EvolutionResult, error) {
	err := e.validate()
	if err != nil {
		return nil, err
	}

	// Set First Generation - TODO Parallelize Individual Creation
	genID := GenerateGenerationID(0)
	gen0 := Generation{
		count:        0,
		GenerationID: genID,
		Protagonists: nil,
		Antagonists:  nil,
		engine:       e,
	}
	e.Generations[0] = &gen0

	antagonists, err := e.Generations[0].GenerateRandomIndividual(IndividualAntagonist,
		e.Parameters.StartIndividual)
	if err != nil {
		return nil, err
	}

	protagonists, err := e.Generations[0].GenerateRandomIndividual(IndividualProtagonist,
		Program{})
	if err != nil {
		return nil, err
	}

	gen0.Protagonists = protagonists
	gen0.Antagonists = antagonists

	// cycle through generationCount
	e.Generations[0] = &gen0
	for i := 0; i < e.Parameters.GenerationsCount-1; i++ {
		started := time.Now()
		protagonistsCleanse, err := CleansePopulation(e.Generations[i].Protagonists, *e.Parameters.StartIndividual.T)
		if err != nil {
			return nil, err
		}
		antagonistsCleanse, err := CleansePopulation(e.Generations[i].Antagonists, *e.Parameters.StartIndividual.T)
		if err != nil {
			return nil, err
		}

		e.Generations[i].Protagonists = protagonistsCleanse
		e.Generations[i].Antagonists = antagonistsCleanse

		// GENERATIONS BEGIN HERE
		nextGeneration, err := e.Generations[i].Start(i)
		if err != nil {
			return nil, err
		}
		e.Generations[i+1] = nextGeneration

		elapsed := utils.TimeTrack(started)
		numGoroutine := runtime.NumGoroutine()
		msg := fmt.Sprintf("\nFile: %s\t | Spec: %s\t | Run: %d | Gen: (%d/%d) | TSz: %d | numG#: %d | Elapsed: %s",
			e.Parameters.ParamFile,
			e.Parameters.SpecParam.ExpressionParsed,
			e.Parameters.InternalCount,
			i+1,
			e.Parameters.GenerationsCount,
			e.Parameters.Strategies.DepthOfRandomNewTrees,
			numGoroutine,
			elapsed.String())
		e.Parameters.LoggingChan <- msg

		if float64(i) == math.Floor(float64(e.Parameters.GenerationsCount) * 0.25) {
			go WriteToDataFolder(e.Parameters.StatisticsOutput.OutputPath,
				"25.txt",
				time.Now().Format(time.RFC3339),
				e.Parameters.LoggingChan,
				e.Parameters.ErrorChan)
		}
		if float64(i) == math.Floor(float64(e.Parameters.GenerationsCount) * 0.5) {
			go WriteToDataFolder(e.Parameters.StatisticsOutput.OutputPath,
				"50.txt",
				time.Now().Format(time.RFC3339),
				e.Parameters.LoggingChan,
				e.Parameters.ErrorChan)
		}
		if float64(i) == math.Floor(float64(e.Parameters.GenerationsCount) * 0.75) {
			go WriteToDataFolder(e.Parameters.StatisticsOutput.OutputPath,
				"75.txt",
				time.Now().Format(time.RFC3339),
				e.Parameters.LoggingChan,
				e.Parameters.ErrorChan)
		}
	}

	evolutionResult := &EvolutionResult{}
	err = evolutionResult.Analyze(e.Generations, true,
		e.Parameters)
	if err != nil {
		return nil, err
	}

	return evolutionResult, nil
}

func WriteToDataFolder(dataFolderPath string, fileName string, fileValue string, logChan chan string,
	errChan chan error) {
	filepath := fmt.Sprintf("%s/%s", dataFolderPath, fileName)
	os.Mkdir(dataFolderPath, 0775)

	file, err := os.Create(filepath)
	if err != nil {
		errChan <- err
	}

	n, err := fmt.Fprintf(file, "%s", fileValue)
	if err != nil {
		errChan <- err
	}else {
		logChan <- fmt.Sprintf("25 PERCENT: => Wrote %d bytes to file %s", n, filepath)
	}
}

// Todo Implement EvolutionProcess validate
func (e *EvolutionEngine) validate() error {
	if e.Parameters.GenerationsCount < 1 {
		return fmt.Errorf("set number of generationCount by calling e.GenerationsCount(x)")
	}
	if e.Parameters.EachPopulationSize%2 != 0 {
		return fmt.Errorf("set number of EachPopulationSize to an Even number")
	}
	//if e.Parameters.SetEqualStrategyLength == true && e.Parameters.EqualStrategiesLength < 1 {
	//	return fmt.Errorf("cannot SetEqualStrategyLength to true and EqualStrategiesLength less than 1")
	//}
	if e.Parameters.StartIndividual.T == nil {
		return fmt.Errorf("start individual cannot have a nil Tree")
	}
	if e.Parameters.Spec == nil {
		return fmt.Errorf("spec cannot be nil")
	}
	if len(e.Parameters.Spec) < 1 {
		return fmt.Errorf("spec cannot be empty")
	}
	if e.Parameters.Selection.Survivor.SurvivorPercentage > 1 || e.Parameters.Selection.Survivor.
		SurvivorPercentage < 0 {
		return fmt.Errorf("SurvivorPercentage cannot be less than 0 or greater than 1. It is a percent value")
	}
	if e.Parameters.Selection.Parent.TournamentSize >= e.Parameters.EachPopulationSize {
		return fmt.Errorf("Tournament Size should not be greater than the population size.")
	}
	//err := e.StartIndividual.Validate()
	//if err != nil {
	//	return err
	//}

	if len(e.Parameters.Spec) < 3 {
		return fmt.Errorf("a small spec will hamper evolutionary accuracy")
	}
	return nil
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
	divide0Penalty = strings.ReplaceAll(divide0Penalty,".","")
	builder.WriteString(divide0Penalty)
	//builder.WriteString("-")
	//builder.WriteString(fmt.Sprintf("%s", RandString(4)))

	return builder.String()
}

func TruncShort(s []Strategy) string {
	sb := strings.Builder{}

	for _, str := range s {
		sb.WriteByte(str[0])
	}

	return sb.String()
}