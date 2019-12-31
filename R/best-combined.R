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

bestCombinedFileNames <- c()

################################################### BEST_COMBINED ##################
################################################### BEST_COMBINED ##################
################################################### BEST_COMBINED ##################
################################################### BEST_COMBINED ##################


best_combined_avg_plot <- function(result, fileName) {
    data = data.frame(
    value = result$ARun,
    A = result$AAvg,
    P = result$PAvg,
    finA = result$ABestFit,
    finP = result$PBestFit
    )

    gg <- ggplot(data, aes(x=value))
    gg <- gg + geom_line(aes(y=A, color = "AvgBug", linetype = 'AvgBug'), size = 1) # setup color name
    gg <- gg + geom_line(aes(y=P, color = "AvgTest", linetype = 'AvgTest'),  size = 1)
    gg <- gg + geom_line(aes(y=finA, color = "BestBug", linetype = 'BestBug'),  size = 1.2)
    gg <- gg + geom_line(aes(y=finP, color = "BestTest", linetype = 'BestTest'),  size = 1.2)
    gg <- gg + geom_point(aes(y=A), size=0.6)
    gg <- gg + geom_point(aes(y=P), size=0.6)
    gg <- gg + geom_point(aes(y=finP), size=0.6)
    gg <- gg + geom_point(aes(y=finA), size=0.6)
    gg <- gg + scale_linetype_manual(values=c(AvgBug='solid', AvgTest='solid', BestBug="dotted", BestTest="dotted"), name =
    "Line Type")
    gg <- gg + scale_colour_manual(values=c(AvgBug="red", AvgTest="green", BestBug="red", BestTest="green"), name = "Plot Color")

    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=8),
    plot.caption = element_text(size=6))
    gg <- gg + labs(
    color = 'Individuals',
    title = sprintf("%s","Average and Best Individual Per Run"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    caption = sprintf("%s\n%s\n%s",
    "Avg: The run average for every generation for each kind of individual",
    "Best: The fittest Bug and Test in the Run",
    "*More Fitness Is Better"),
    x = "Run",
    y = "Fitness")

    fileName <- paste(fileName, "best-combined_avg.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
}

getAllFiles <- function(workDir) {
    files <- list.files(workDir)
    bestCombinedCount <- 1

    for (file in files) {
        if (grepl("best-combined", file)) {
            bestCombinedFileNames[bestCombinedCount] <- file
            filePath <- paste(workDir, file, sep="/")
            bestCombinedData = read_csv(filePath)

            best_combined_avg_plot(bestCombinedData, file)

            bestCombinedCount <- bestCombinedCount + 1
        }
    }
}


getAllFiles(workDir)

