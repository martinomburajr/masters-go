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

strategyFileNames <- c()
strategyCombFileNames <- c()

######################################################### STRATEGY ############################

strategy_run_histogram_plot <- function(result, fileName) {
    data = data.frame(
    A = result$A,
    P = result$P
    )
    dataA = data.frame(A = result$A)
    dataP = data.frame(A = result$P)
    alpha <- 0.2

    gg <- ggplot(data, aes(A))
    gg <- gg + geom_bar(data=dataA, stat="count", aes(color = "Bug"), alpha = alpha, fill="red", size=0.8)
    gg <- gg + geom_bar(data=dataP, stat="count", aes(color = "Test"), alpha = alpha, fill="green", size=0.8)
    gg <- gg + scale_colour_manual(values=c(Bug="red", Test="green"), name = "Plot Color")
    #
    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=8),
    plot.caption = element_text(size=6))
    gg <- gg + labs(
    color = 'Individuals',
    title = sprintf("%s","Histogram of Bug and Test Strategy Selection"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    x = "Strategy",
    y = "Frequency")

    fileName <- paste(fileName, "strat_bar.png", sep="-")
    ggsave(fileName, width=16, height=6, units='in', dpi="retina")
}

strategy_best_histogram_plot <- function(result, fileName) {
    data = data.frame(
    A = result$A,
    P = result$P
    )
    dataA = data.frame(A = result$A)
    dataP = data.frame(A = result$P)
    alpha <- 0.2

    gg <- ggplot(data, aes(A))
    gg <- gg + geom_bar(data=dataA, stat="count", aes(color = "Bug"), alpha = alpha, fill="red", size=0.8)
    gg <- gg + geom_bar(data=dataP, stat="count", aes(color = "Test"), alpha = alpha, fill="green", size=0.8)
    gg <- gg + scale_colour_manual(values=c(Bug="red", Test="green"), name = "Plot Color")
    #
    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=8),
    plot.caption = element_text(size=6))
    gg <- gg + labs(
    color = 'Individuals',
    title = sprintf("%s","Histogram of Best Bug and Test Strategy in All Runs"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    x = "Strategy",
    y = "Frequency")

    fileName <- paste(fileName, "strat_bar.png", sep="-")
    ggsave(fileName, width=16, height=6, units='in', dpi="retina")
}


strategy_cum_run_histogram_plot <- function(result, fileName) {
    data = data.frame(
    A = result$A,
    P = result$P
    )
    dataA = data.frame(A = result$A)
    dataP = data.frame(A = result$P)
    alpha <- 0.2

    gg <- ggplot(data, aes(A))
    gg <- gg + geom_bar(data=dataA, stat="count", aes(color = "Bug"), alpha = alpha, fill="red", size=0.8)
    gg <- gg + geom_bar(data=dataP, stat="count", aes(color = "Test"), alpha = alpha, fill="green", size=0.8)
    gg <- gg + scale_colour_manual(values=c(Bug="red", Test="green"), name = "Plot Color")
    #
    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=8),
    plot.caption = element_text(size=6))
    gg <- gg + labs(
    color = 'Individuals',
    title = sprintf("%s","Cumulative Histogram of Best Bug and Test Strategy Selection (All Runs)"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    caption = sprintf("%s", "Best is defined as the best in all generations (not necessarily the final individual)"),
    x = "Strategy",
    y = "Frequency")

    fileName <- paste(fileName, ".png", sep="-")
    ggsave(fileName, width=20, height=6, units='in', dpi="retina")
}

strategy_cum_finalIndividual_run_histogram_plot <- function(result, fileName) {
    data = data.frame(
    A = result$AFinal,
    P = result$PFinal
    )
    dataA = data.frame(A = result$AFinal)
    dataP = data.frame(A = result$PFinal)
    alpha <- 0.2

    gg <- ggplot(data, aes(A))
    gg <- gg + geom_bar(data=dataA, stat="count", aes(color = "Bug"), alpha = alpha, fill="red", size=0.8)
    gg <- gg + geom_bar(data=dataP, stat="count", aes(color = "Test"), alpha = alpha, fill="green", size=0.8)
    gg <- gg + scale_colour_manual(values=c(Bug="red", Test="green"), name = "Plot Color")
    #
    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=8),
    plot.caption = element_text(size=6))
    gg <- gg + labs(
    color = 'Individuals',
    title = sprintf("%s","Histogram of Final Bug and Test Strategy Selection (All Runs)"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    caption = sprintf("%s", "Final is defined as the individuals in the final generation for a given run."),
    x = "Strategy",
    y = "Frequency")

    fileName <- paste(fileName, ".png", sep="-")
    ggsave(fileName, width=20, height=6, units='in', dpi="retina")
}

getAllFiles <- function(workDir) {
    files <- list.files(workDir)
    strategyCount <- 1

    combinedStrategies <- data.frame(
    A = character(),
    P = character(),
    AFinal = character(),
    PFinal = character(),
    AGen = integer(0),
    PGen =  integer(0),
    count =  integer(0),
    run =  integer(0)
    )
    for (file in files) {
        if (grepl("strategy-", file)) {
            strategyFileNames[strategyCount] <- file
            filePath <- paste(workDir, file, sep="/")
            strategyData = read_csv(filePath)

            combinedStrategies <- rbind(combinedStrategies, strategyData)
            strategy_run_histogram_plot(strategyData, file)
            strategyCount <- strategyCount + 1
        }
        if (grepl("strategybest", file)) {
            strategyFileNames[strategyCount] <- file
            filePath <- paste(workDir, file, sep="/")
            strategyData = read_csv(filePath)

            strategy_best_histogram_plot(strategyData, file)
            strategyCount <- strategyCount + 1
        }
    }
    strategy_cum_run_histogram_plot(combinedStrategies, "combined-strategies")
    strategy_cum_finalIndividual_run_histogram_plot(combinedStrategies, "combined-strategies-finalindividual")
}


getAllFiles(workDir)

