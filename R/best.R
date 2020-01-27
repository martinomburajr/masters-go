# Title     : TODO
# Objective : TODO
# Created by: martinomburajr
# Created on: 2019/10/24

args = commandArgs(trailingOnly=TRUE)
library(ggplot2)
library(readr)
library(knitr)
library(dplyr)

workDir <- ""
statsDir <- ""

if (length(args)==0) {
    stop("At least one argument must be supplied (input file).n", call.=FALSE)
} else {
    workDir = args[1]
    statsDir <- workDir
    dir.create(file.path(statsDir, ""), showWarnings = FALSE)
    setwd(file.path(statsDir, ""))
}

bestAllFileNames <- c()


################################################### BEST ##################
################################################### BEST ##################
################################################### BEST ##################
################################################### BEST ##################
################################################### BEST ##################

best_p_spec_plot <- function(result, fileName) {
    data = data.frame(
    value = result$epoch,
    A = result$A,
    P = result$P,
    finA = result$finA,
    finP = result$finP
    )

    gg <- ggplot(data, aes(x=value))
    gg <- gg + geom_line(aes(y=A, color = "BestBug", linetype = 'BestBug'), size = 1) # setup color name
    gg <- gg + geom_line(aes(y=P, color = "BestTest", linetype = 'BestTest'),  size = 1)
    gg <- gg + geom_line(aes(y=finA, color = "FinalBug", linetype = 'FinalBug'),  size = 1.2)
    gg <- gg + geom_line(aes(y=finP, color = "FinalTest", linetype = 'FinalTest'),  size = 1.2)
    gg <- gg + geom_point(aes(y=A), size=0.6)
    gg <- gg + geom_point(aes(y=P), size=0.6)
    gg <- gg + geom_point(aes(y=finP), size=0.6)
    gg <- gg + geom_point(aes(y=finA), size=0.6)
    gg <- gg + scale_linetype_manual(values=c(BestBug='solid', BestTest='solid', FinalBug="dotted", FinalTest="dotted"), name =
    "Line Type")
    gg <- gg + scale_colour_manual(values=c(BestBug="red", BestTest="green", FinalBug="red", FinalTest="green"), name = "Plot Color")

    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=8),
    plot.caption = element_text(size=6))
    gg <- gg + labs(
    color = 'Individuals',
    title = sprintf("%s","Epoch Based Fitness Variation of Bug and Test"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    caption = sprintf("%s\n%s\n%s",
    "Best: The fittest Bug and Test in the Run",
    "Final: The last generations best bug and test",
    "*More Fitness Is Better"),
    x = "Epoch",
    y = "Fitness")

    fileName <- paste(fileName, "epochal.png", sep="-")
    ggsave(fileName, width=8, height=5, units='in', dpi="retina")
}

## Pass in best-all.csv
best_all_function_plot <- function(result, fileName) {
    gg <- ggplot(data = data.frame(x = 0), mapping = aes(x = result$range))

    specEquation <- function(x){eval(parse(text=result$specEquation))}
    ultAntagonistEquation <- function(x){eval(parse(text=result$AEquation))}
    ultProtagonistEquation <- function(x){eval(parse(text=result$PEquation))}

    #spec
    gg <- gg + stat_function(
    stat = "function",
    fun = specEquation,
    mapping = aes(color="Spec", linetype="Spec"),
    size=1.3
    )
    gg <- gg + stat_function(
    stat = "function",
    fun = ultAntagonistEquation,
    mapping = aes(color = "BestBugEquation", linetype="BestBugEquation")
    )
    gg <- gg + stat_function(
    stat = "function",
    fun = ultProtagonistEquation,
    mapping = aes(color = "BestTestEquation", linetype="BestTestEquation")
    )

    gg <- gg + scale_x_continuous(limits=c(result$seed, result$seed + result$range))
    gg <- gg + scale_linetype_manual(
    name = "Line Type",
    values=c(Spec='dotted', BestBugEquation='solid', BestTestEquation="solid")
    )
    gg <- gg + scale_color_manual(
    name = "Functions",
    values = c(Spec="black", BestBugEquation="red", BestTestEquation="green")
    )

    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=12),
    plot.caption = element_text(size=10))
    gg <- gg + labs (
    color = 'Individuals',
    title = sprintf("%s","Resulting Best Equation For Bug and Test against Spec"),
    subtitle = sprintf("Spec: %s", result$specEquation),
    caption = sprintf(
    "BestBug: %s\nBestTest: %s\nRange: [%d, %d]\n%s",
    toString(result$AEquation),
    toString(result$PEquation),
    result$seed,
    result$seed + result$range,
    "*Closer mapping on to spec is better"),
    x = "X",
    y = "Y"
    )

    fileName <- paste(fileName, "best-all.png", sep="-")
    ggsave(fileName, width=10, height=6, units='in', dpi="retina")
}

## Pass in best-all.csv
best_bug_spec_function_plot <- function(result, fileName) {
    gg <- ggplot(data = data.frame(x = 0), mapping = aes(x = result$range))

    specEquation <- function(x){eval(parse(text=result$specEquation))}
    ultAntagonistEquation <- function(x){eval(parse(text=result$AEquation))}

    #spec
    gg <- gg + stat_function(
    stat = "function",
    fun = specEquation,
    mapping = aes(color="Spec", linetype="Spec"),
    size=1.3
    )
    gg <- gg + stat_function(
    stat = "function",
    fun = ultAntagonistEquation,
    mapping = aes(color = "BestBugEquation", linetype="BestBugEquation")
    )

    gg <- gg + scale_x_continuous(limits=c(result$seed, result$seed + result$range))
    gg <- gg + scale_linetype_manual(
    name = "Line Type",
    values=c(Spec='dotted', BestBugEquation='solid')
    )
    gg <- gg + scale_color_manual(
    name = "Functions",
    values = c(Spec="black", BestBugEquation="red")
    )

    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=12),
    plot.caption = element_text(size=10))
    gg <- gg + labs (
    color = 'Individuals',
    title = sprintf("%s","Resulting Best Equation For Bug against Spec"),
    subtitle = sprintf("Spec: %s", result$specEquation),
    caption = sprintf(
    "BestBug: %s\nRange: [%d, %d]\n%s",
    toString(result$AEquation),
    result$seed,
    result$seed + result$range,
    "*Closer mapping on to spec is better"),
    x = "X",
    y = "Y"
    )

    fileName <- paste(fileName, "best-bug-spec-all.png", sep="-")
    ggsave(fileName, width=10, height=6, units='in', dpi="retina")
}

## Pass in best-all.csv
best_test_spec_function_plot <- function(result, fileName) {
    gg <- ggplot(data = data.frame(x = 0), mapping = aes(x = result$range))

    specEquation <- function(x){eval(parse(text=result$specEquation))}
    ultProtagonistEquation <- function(x){eval(parse(text=result$PEquation))}

    #spec
    gg <- gg + stat_function(
    stat = "function",
    fun = specEquation,
    mapping = aes(color="Spec", linetype="Spec"),
    size=1.3
    )
    gg <- gg + stat_function(
    stat = "function",
    fun = ultProtagonistEquation,
    mapping = aes(color = "BestTestEquation", linetype="BestTestEquation")
    )

    gg <- gg + scale_x_continuous(limits=c(result$seed, result$seed + result$range))
    gg <- gg + scale_linetype_manual(
    name = "Line Type",
    values=c(Spec='dotted', BestTestEquation="solid")
    )
    gg <- gg + scale_color_manual(
    name = "Functions",
    values = c(Spec="black", BestTestEquation="green")
    )

    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=12),
    plot.caption = element_text(size=10))
    gg <- gg + labs (
    color = 'Individuals',
    title = sprintf("%s","Resulting Best Equation for Test aginst Spec"),
    subtitle = sprintf("Spec: %s", result$specEquation),
    caption = sprintf(
    "BestTest: %s\nRange: [%d, %d]\n%s",
    toString(result$PEquation),
    result$seed,
    result$seed + result$range,
    "*Closer mapping on to spec is better"),
    x = "X",
    y = "Y"
    )

    fileName <- paste(fileName, "best-test-spec-all.png", sep="-")
    ggsave(fileName, width=10, height=6, units='in', dpi="retina")
}

best_all_tests_plot <- function(result, fileName) {
    gg <- ggplot(data = data.frame(x = 0), mapping = aes(x = result$range))

    specEquation <- function(x){eval(parse(text=result$specEquation))}
    ultProtagonistEquation <- function(x){eval(parse(text=result$PEquation))}

    #spec
    gg <- gg + stat_function(
    stat = "function",
    fun = specEquation,
    mapping = aes(color="Spec", linetype="Spec"),
    size=1.3
    )

    for (row in nrow(result)) {
        gg <- gg + stat_function(
        stat = "function",
        fun = function(x){eval(parse(text=result$PEquation))},
        mapping = aes(color = "BestTestEquation", linetype="BestTestEquation")
        )
    }

    gg <- gg + scale_x_continuous(limits=c(result$seed, result$seed + result$range))
    gg <- gg + scale_linetype_manual(
    name = "Line Type",
    values=c(Spec='dotted', BestTestEquation="solid")
    )
    gg <- gg + scale_color_manual(
    name = "Functions",
    values = c(Spec="black", BestTestEquation="green")
    )

    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=12),
    plot.caption = element_text(size=10))
    gg <- gg + labs (
    color = 'Individuals',
    title = sprintf("%s","Resulting Best Equation for Test aginst Spec"),
    subtitle = sprintf("Spec: %s", result$specEquation),
    caption = sprintf(
    "BestTest: %s\nRange: [%d, %d]\n%s",
    toString(result$PEquation),
    result$seed,
    result$seed + result$range,
    "*Closer mapping on to spec is better"),
    x = "X",
    y = "Y"
    )

    fileName <- paste(fileName, "best-test-spec-all.png", sep="-")
    ggsave(fileName, width=10, height=6, units='in', dpi="retina")
}

# Plots out the average between the average of all antagonists in a given geernation, and the average of all
# protagonists in the same generation.
covcor_average_plot <- function(result, fileName) {
    data = data.frame(
    value = result$gen,
    A = result$meanCorrInRun,
    P = result$meanCovInRun,
    )

    gg <- ggplot(data, aes(x=value))
    gg <- gg + geom_line(aes(y=A, color = "Corr", linetype = 'Corr'), size = 1) # setup color name
    gg <- gg + geom_line(aes(y=P, color = "Cov", linetype = 'Cov'),  size = 1)
    gg <- gg + geom_point(aes(y=A), size=0.6)
    gg <- gg + geom_point(aes(y=P), size=0.6)
    gg <- gg + scale_linetype_manual(values=c(Corr='solid', Cov='solid'), name =
    "Line Type")
    gg <- gg + scale_colour_manual(values=c(Corr="blue", Cov="green",), name = "Plot Color")

    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=8),
    plot.caption = element_text(size=6))
    gg <- gg + labs(
    color = 'Corr vs. Cov',
    title = sprintf("%s","Mean Correlation and Mean Covariance"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    x = "Run",
    y = "Index")

    fileName <- paste(fileName, "covcor.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
    # dev.off()
}

getAllFiles <- function(workDir) {
    files <- list.files(workDir)
    bestAllCount <- 1

    combinedBest <- data.frame(
        specEquation= character(),
        AEquation = character(),
        PEquation = character(),
        run = integer(0)
    )
    for (file in files) {
        if (grepl("best-all", file)) {
            bestAllFileNames[bestAllCount] <- file
            filePath <- paste(workDir, file, sep="/")
            bestAllData = read_csv(filePath)
            best_all_function_plot(bestAllData, file)
            best_bug_spec_function_plot(bestAllData, file)
            best_test_spec_function_plot(bestAllData, file)
            covcor_average_plot(bestAllData, file)

            combinedBest <- merge(combinedBest, bestAllData)
            print(combinedBest)

            bestAllCount <- bestAllCount + 1
        }
    }
    print(combinedBest)
}


getAllFiles(workDir)

