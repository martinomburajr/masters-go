package analysis

import (
	"github.com/martinomburajr/masters-go/evolution"
	"github.com/martinomburajr/masters-go/simulation"
)

type CSVBestAll struct {
	FileID                  string                                `csv:"ID"`
	bestIndividualStatistic simulation.RunBestIndividualStatistic `csv:"bestIndividualStatistic"`
	params                  evolution.EvolutionParams             `csv:"evolutionaryParams"`

	//BEST INDIVIDUAL
	SpecEquation string `csv:"specEquation"`
	SpecRange    int    `csv:"range"`
	SpecSeed     int    `csv:"seed"`

	AntagonistID                string  `csv:"AID"`
	ProtagonistID               string  `csv:"PID"`
	Antagonist                  float64 `csv:"AAvg"`
	Protagonist                 float64 `csv:"PAvg"`
	AntagonistBestFitness       float64 `csv:"ABestFit"`
	ProtagonistBestFitness      float64 `csv:"PBestFit"`
	AntagonistStdDev            float64 `csv:"AStdDev"`
	ProtagonistStdDev           float64 `csv:"PStdDev"`
	AntagonistAverageDelta      float64 `csv:"AAvgDelta"`
	ProtagonistAverageDelta     float64 `csv:"PAvgDelta"`
	AntagonistBestDelta         float64 `csv:"ABestDelta"`
	ProtagonistBestDelta        float64 `csv:"PBestDelta"`
	AntagonistEquation          string  `csv:"AEquation"`
	ProtagonistEquation         string  `csv:"PEquation"`
	AntagonistStrategy          string  `csv:"AStrat"`
	ProtagonistStrategy         string  `csv:"PStrat"`
	AntagonistDominantStrategy  string  `csv:"ADomStrat"`
	ProtagonistDominantStrategy string  `csv:"PDomStrat"`
	AntagonistGeneration        int     `csv:"AGen"`
	ProtagonistGeneration       int     `csv:"PGen"`
	AntagonistBirthGen          int     `csv:"ABirthGen"`
	ProtagonistBirthGen         int     `csv:"PBirthGen"`
	AntagonistRun int `csv:"ARun"`
	ProtagonistRun int `csv:"PRun"`
	AntagonistAge               int     `csv:"AAge"`
	ProtagonistAge              int     `csv:"PAge"`
	AntagonistNoOComp		int `csv:"ANoC"`
	ProtagonistNoOComp		int `csv:"PNoC"`

	FinalAntagonist                  float64 `csv:"finAAvg"`
	FinalProtagonist                 float64 `csv:"finPAvg"`
	FinalAntagonistBestFitness       float64 `csv:"finABestFit"`
	FinalProtagonistBestFitness      float64 `csv:"finPBestFit"`
	FinalAntagonistStdDev            float64 `csv:"finAStdDev"`
	FinalProtagonistStdDev           float64 `csv:"finPStdDev"`
	FinalAntagonistAverageDelta      float64 `csv:"finAAvgDelta"`
	FinalProtagonistAverageDelta     float64 `csv:"finPAvgDelta"`
	FinalAntagonistBestDelta         float64 `csv:"finABestDelta"`
	FinalProtagonistBestDelta        float64 `csv:"finPBestDelta"`
	FinalAntagonistEquation          string  `csv:"finAEquation"`
	FinalProtagonistEquation         string  `csv:"finPEquation"`
	FinalAntagonistStrategy          string  `csv:"finAStrat"`
	FinalProtagonistStrategy         string  `csv:"finPStrat"`
	FinalAntagonistDominantStrategy  string  `csv:"finADomStrat"`
	FinalProtagonistDominantStrategy string  `csv:"finPDomStrat"`
	FinalAntagonistBirthGen          int     `csv:"finABirthGen"`
	FinalProtagonistBirthGen         int     `csv:"finPBirthGen"`
	FinalAntagonistAge               int     `csv:"finAAge"`
	FinalProtagonistAge              int     `csv:"finPAge"`
	FinalAntagonistNoOComp		int `csv:"finANoC"`
	FinalProtagonistNoOComp		int `csv:"finPNoC"`

	Run int `csv:"run"`

	// PARAMS
	GenerationCount    int     `csv:"genCount"`
	EachPopulationSize int     `csv:"popCount"`
	TopologyType string `csv:"topology"`
	AntStratCount      int     `csv:"antStratCount"`
	ProStratCount      int     `csv:"proStratCount"`
	AntStrat           string  `csv:"antStrat"`
	ProStrat           string  `csv:"proStrat"`
	RandTreeDepth      int     `csv:"randTreeDepth"`
	AntThreshMult      float64 `csv:"antThreshMult"`
	ProThresMult       float64 `csv:"proThresMult"`
	CrossPercent       float64 `csv:"crossPercent"`
	ProbMutation       float64 `csv:"probMutation"`
	ParentSelect       string  `csv:"parentSelect"`
	TournamentSize     int     `csv:"tournamentSize"`
	SurvivorSelect     string  `csv:"survivorSelect"`
	SurvivorPercent    float64 `csv:"survivorPercent"`
	DivByZero          string  `csv:"d0"`
	DivByZeroPen       float64 `csv:"d0Pen"`
}


type CSVCombinedGenerations struct {
	FileID                  string                                `csv:"ID"`
	params                  evolution.EvolutionParams             `csv:"evolutionaryParams"`

	//BEST INDIVIDUAL
	Generation int `csv:"gen"`
	SpecEquation string `csv:"specEquation"`
	SpecRange    int    `csv:"range"`
	SpecSeed     int    `csv:"seed"`

	TopAEquation string `csv:"topAEquation"`
	TopPEquation string `csv:"topPEquation"`
	Correlation float64 `csv:"correlation"`
	Covariance float64 `csv:"covariance"`


	Antagonist                  float64 `csv:"AMean"`
	Protagonist                 float64 `csv:"PMean"`
	TopAntagonistMean float64 `csv:"topAMean"`
	TopProtagonistMean float64 `csv:"topPMean"`
	TopAntagonistBestFitness       float64 `csv:"topABest"`
	TopProtagonistBestFitness      float64 `csv:"topPBest"`
	TopAntagonistStdDev            float64 `csv:"AStd"`
	TopProtagonistStdDev           float64 `csv:"PStd"`
	TopAntagonistVar float64 `csv:"AVar"`
	TopProtagonistVar float64 `csv:"PVar"`
	TopAntagonistSkew float64 `csv:"ASkew"`
	TopProtagonistSkew float64 `csv:"PSkew"`
	TopAntagonistKurtosis float64 `csv:"AExKur"`
	TopProtagonistKurtosis float64 `csv:"PExKur"`


	TopAntagonistAverageDelta      float64 `csv:"topAMeanDelta"`
	TopProtagonistAverageDelta     float64 `csv:"topPMeanDelta"`
	TopAntagonistBestDelta         float64 `csv:"topABestDelta"`
	TopProtagonistBestDelta        float64 `csv:"topPBestDelta"`
	TopAntagonistStrategy          string  `csv:"topAStrat"`
	TopProtagonistStrategy         string  `csv:"topPStrat"`
	TopAntagonistDominantStrategy  string  `csv:"topADomStrat"`
	TopProtagonistDominantStrategy string  `csv:"topPDomStrat"`
	TopAntagonistGeneration        int     `csv:"topAGen"`
	TopProtagonistGeneration       int     `csv:"topPGen"`
	TopAntagonistBirthGen          int     `csv:"topABirthGen"`
	TopProtagonistBirthGen         int     `csv:"topPBirthGen"`
	TopAntagonistAge               int     `csv:"topAAge"`
	TopProtagonistAge              int     `csv:"topPAge"`

	Run int `csv:"run"`

	// PARAMS
	GenerationCount    int     `csv:"genCount"`
	EachPopulationSize int     `csv:"popCount"`
	TopologyType string `csv:"topology"`
	AntStratCount      int     `csv:"antStratCount"`
	ProStratCount      int     `csv:"proStratCount"`
	AntStrat           string  `csv:"antStrat"`
	ProStrat           string  `csv:"proStrat"`
	RandTreeDepth      int     `csv:"randTreeDepth"`
	AntThreshMult      float64 `csv:"antThreshMult"`
	ProThresMult       float64 `csv:"proThresMult"`
	CrossPercent       float64 `csv:"crossPercent"`
	ProbMutation       float64 `csv:"probMutation"`
	ParentSelect       string  `csv:"parentSelect"`
	TournamentSize     int     `csv:"tournamentSize"`
	SurvivorSelect     string  `csv:"survivorSelect"`
	SurvivorPercent    float64 `csv:"survivorPercent"`
	DivByZero          string  `csv:"d0"`
	DivByZeroPen       float64 `csv:"d0Pen"`
}