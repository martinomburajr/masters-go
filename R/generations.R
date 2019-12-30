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

generationalFileNames <- c()

################################################ GENERATION #################

# Plots out the average between the average of all antagonists in a given geernation, and the average of all
# protagonists in the same generation.
generational_average_plot <- function(result, fileName) {
    data = data.frame(
    value = result$gen,
    A = result$AGenFitAvg,
    P = result$PGenFitAvg,
    bestA = result$AGenBestFitAvg,
    bestP = result$PGenBestFitAvg
    )

    gg <- ggplot(data, aes(x=value))
    gg <- gg + geom_line(aes(y=A, color = "AvgBug", linetype = 'AvgBug'), size = 1) # setup color name
    gg <- gg + geom_line(aes(y=P, color = "AvgTest", linetype = 'AvgTest'),  size = 1)
    gg <- gg + geom_line(aes(y=bestA, color = "BestAvgBug", linetype = 'BestAvgBug'),  size = 1.2)
    gg <- gg + geom_line(aes(y=bestP, color = "BestAvgTest", linetype = 'BestAvgTest'),  size = 1.2)
    gg <- gg + geom_point(aes(y=A), size=0.6)
    gg <- gg + geom_point(aes(y=P), size=0.6)
    gg <- gg + geom_point(aes(y=bestA), size=0.6)
    gg <- gg + geom_point(aes(y=bestP), size=0.6)
    gg <- gg + scale_linetype_manual(values=c(AvgBug='solid', AvgTest='solid', BestAvgBug="dotted",
    BestAvgTest="dotted"), name =
    "Line Type")
    gg <- gg + scale_colour_manual(values=c(AvgBug="red", AvgTest="green", BestAvgBug="red", BestAvgTest="green"), name = "Plot Color")

    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=8),
    plot.caption = element_text(size=6))
    gg <- gg + labs(
    color = 'Individuals',
    title = sprintf("%s","Generation Based Fitness Variation of Bug and Test"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    caption = sprintf("%s\n%s\n%s",
    "Avg: Avg of all grouped individuals",
    "BestAvg: Best individuals fitness avg",
    "*More Fitness Is Better"),
    x = "Generation",
    y = "Fitness")

    fileName <- paste(fileName, "generational.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
    # dev.off()
}

generational_histogram_plot <- function(result, fileName) {
    data = data.frame(
    value = result$AGenFitAvg,
    A = result$AGenFitAvg,
    P = result$PGenFitAvg
    )
    dataP = data.frame(A = result$PGenFitAvg)
    dataA = data.frame(A = result$AGenFitAvg)

    gg <- ggplot(data, aes(A))
    gg <- gg + geom_histogram(data=dataA, binwidth=0.002, aes(color = "Bug", linetype = 'Bug'), alpha = 0.2)
    gg <- gg + geom_histogram(data=dataP, binwidth=0.002, aes(color = "Test", linetype = 'Test'), alpha = 0.2)
    gg <- gg + scale_colour_manual(values=c(Bug="red", Test="green"), name = "Plot Color")
    #
    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=8),
    plot.caption = element_text(size=6))
    gg <- gg + labs(
    color = 'Individuals',
    title = sprintf("%s","Histogram of Bug and Test Fitness"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    x = "Fitness",
    y = "Frequency")

    fileName <- paste(fileName, "gen_histogram.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
}

generational_density_plot <- function(result, fileName) {
    data = data.frame(
    value = result$AGenFitAvg,
    A = result$AGenFitAvg,
    P = result$PGenFitAvg
    )
    dataP = data.frame(A = result$PGenFitAvg)
    dataA = data.frame(A = result$AGenFitAvg)

    gg <- ggplot(data, aes(A))
    gg <- gg + geom_density(data=dataA, kernel = "gaussian", aes(color = "Bug", linetype = 'Bug'), alpha = 0.2)
    gg <- gg + geom_density(data=dataP, kernel = "gaussian", aes(color = "Test", linetype = 'Test'), alpha = 0.2)

    gg <- gg + geom_vline(aes(xintercept=mean(A), color = "Bug"), linetype = 'dotted', size=0.7)
    gg <- gg + geom_vline(aes(xintercept=mean(P), color = "Test"), linetype = 'dotted', size=0.7)

    gg <- gg + scale_linetype_manual(values=c(Bug='solid', Test='solid'), name = "Line Type")
    gg <- gg + scale_colour_manual(values=c(Bug="red", Test="green"), name = "Plot Color")
    #
    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=8),
    plot.caption = element_text(size=6))
    gg <- gg + labs(
    color = 'Individuals',
    title = sprintf("%s","Density Distribution of Bug and Test Fitness"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    x = "Fitness",
    y = "Frequency")

    fileName <- paste(fileName, "gen_density.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
}

generational_density_histogram_plot <- function(result, fileName) {
    data = data.frame(
    value = result$AGenFitAvg,
    A = result$AGenFitAvg,
    P = result$PGenFitAvg
    )
    dataP = data.frame(A = result$PGenFitAvg)
    dataA = data.frame(A = result$AGenFitAvg)

    gg <- ggplot(data, aes(A))
    gg <- gg + geom_density(data=dataA, kernel = "gaussian", aes(color = "Bug", linetype = 'Bug'), alpha = 0.2)
    gg <- gg + geom_density(data=dataP, kernel = "gaussian", aes(color = "Test", linetype = 'Test'), alpha = 0.2)
    gg <- gg + geom_histogram(data=dataA, binwidth=0.018, aes(color = "Bug", linetype = 'Bug'), alpha = 0.2)
    gg <- gg + geom_histogram(data=dataP, binwidth=0.018, aes(color = "Test", linetype = 'Test'), alpha = 0.2)
    gg <- gg + scale_colour_manual(values=c(Bug="red", Test="green"), name = "Plot Color")
    #
    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=8),
    plot.caption = element_text(size=6))
    gg <- gg + labs(
    color = 'Individuals',
    title = sprintf("%s","Histogram Density Distribution of Bug and Test Fitness"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    x = "Fitness",
    y = "Frequency")

    fileName <- paste(fileName, "gen_density_histogram.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
}

getAllFiles <- function(workDir) {
    files <- list.files(workDir)
    count <- 1
    for (file in files) {
        if (grepl("generational", file)) {
            generationalFileNames[count] <- file

            filePath <- paste(workDir, file, sep="/")
            print(filePath)
            generationalData = read_csv(filePath)

            generational_histogram_plot(generationalData, file)
            generational_density_plot(generationalData,  file)
            generational_density_histogram_plot(generationalData,  file)
            generational_average_plot(generationalData,  file)
            count <- count + 1
        }
    }

    print(length(generationalFileNames))
}


getAllFiles(workDir)

