# Title     : TODO
# Objective : TODO
# Created by: martinomburajr
# Created on: 2019/10/24

redColor <- "#ff5252"

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


epochalFileNames <- c()
####################################### CODE BEGINS ##########################
######################################## EPOCH
epochal_plot <- function(result, fileName) {
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

epochal_aDelta_plot <- function(result, fileName) {
    data = data.frame(
    value = result$epoch,
    A = result$ADelta
    )

    gg <- ggplot(data, aes(x=value))
    gg <- gg + geom_line(aes(y=A, color = "BugDelta", linetype = 'BugDelta'), size = 1) # setup color name
    gg <- gg + geom_point(aes(y=A), size=0.6)

    gg <- gg + scale_linetype_manual(values=c(BugDelta='solid'), name =
    "Line Type")
    gg <- gg + scale_colour_manual(values=c(BugDelta="red"), name = "Plot Color")

    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=8),
    plot.caption = element_text(size=6))
    gg <- gg + labs(
    color = 'Individuals',
    title = sprintf("%s","Epoch Based Delta Value Variation of Bug"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    caption = sprintf("%s\n%s",
    "Delta: Average difference between spec and individual's value",
    "Bugs: Attempt to maximize delta"),
    x = "Epoch",
    y = "Delta")

    fileName <- paste(fileName, "epochal-delta-A.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
}

epochal_pDelta_plot <- function(result, fileName) {
    data = data.frame(
    value = result$epoch,
    P = result$PDelta
    )

    gg <- ggplot(data, aes(x=value))
    gg <- gg + geom_line(aes(y=P, color = "TestDelta", linetype = 'TestDelta'),  size = 1)
    gg <- gg + geom_point(aes(y=P), size=0.6)

    gg <- gg + scale_linetype_manual(values=c(TestDelta='solid'), name =
    "Line Type")
    gg <- gg + scale_colour_manual(values=c(TestDelta="green"), name = "Plot Color")

    gg <- gg + guides(color = guide_legend(title="Legend"), linetype = guide_legend(title="Legend"))

    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=8),
    plot.caption = element_text(size=6))
    gg <- gg + labs(
    color = 'Individuals',
    title = sprintf("%s","Epoch Based Delta Value Variation of Test"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    caption = sprintf("%s\n%s",
    "Delta: Average difference between spec and individual's value",
    "Tests: Attempt to minimize delta"),
    x = "Epoch",
    y = "Delta")

    fileName <- paste(fileName, "epochal-delta-P.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
}

epochal_density_plot <- function(result, fileName) {
    data = data.frame(
    value = result$A,
    A = result$A,
    P = result$P
    )
    dataA = data.frame(A = result$A)
    dataP = data.frame(A = result$P)

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
    title = sprintf("%s","Epochal Fitness Density Distribution of Bug and Test"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    caption = sprintf("Bug Mean: %.2f | Test Mean: %.2f\nBug SDev: %.2f | Test SDev: %.2f",
    mean(result$A), mean(result$P), sd(result$A), sd(result$P)),
    x = "Fitness",
    y = "Frequency")

    fileName <- paste(fileName, "epochal_density.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
}

epochal_density_histogram_plot <- function(result, fileName) {
    data = data.frame(
    value = result$A,
    A = result$A,
    P = result$P
    )
    dataP = data.frame(A = result$P)
    dataA = data.frame(A = result$A)

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
    title = sprintf("%s","Epochal Fitness Histogram Density Distribution of Bug and Test"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    caption = sprintf("Bug Mean: %.2f | Test Mean: %.2f\nBug SDev: %.2f | Test SDev: %.2f",
    mean(result$A), mean(result$P), sd(result$A), sd(result$P)),
    x = "Fitness",
    y = "Frequency")

    fileName <- paste(fileName, "epochal_density_histogram.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
}

epochal_all_bug_runs_boxplot <- function(result, fileName) {
    # Result is a dataframe containing Runs on X axis and values on Y
    data = data.frame(
        A = result$A,
        P = result$P,
        run = result$run,
        AColor = result$AColor
    )
    data$discreteX = as.character(result$run)

    gg <- ggplot(data, aes(x=discreteX, y=A, fill="BestBug"))
    gg <- gg + geom_boxplot(
        outlier.colour="#A4A4A4",
        outlier.shape=16,
        outlier.size=1,
        notch=FALSE,
        fill="tomato"
        )
    gg <- gg + geom_dotplot(binaxis='y', stackdir='center', dotsize=0.8)
    gg <- gg + stat_summary(fun.y=mean, geom="point", shape=23, size=3, aes(x=discreteX, fill="BestBug"))
    gg <- gg + scale_fill_brewer(palette="YlOrRd") + theme_minimal()
    gg <- gg + scale_colour_manual(values=c(BestBug="red"), name = "Plot Color")

    gg <- gg + guides(
        fill=guide_legend(title="Legend"),
        linetype = guide_legend(title="Legend")
    )
    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=8),
    plot.caption = element_text(size=6))

    gg <- gg + labs(
    color = 'Individuals',
    title = sprintf("%s","Epochal Average of Best Bug in Each Run"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    caption = sprintf("Cumulative Bug Mean: %.2f\nCumulative Bug SDev: %.2f",
    mean(result$A), sd(result$A)),
    x = "Run",
    y = "Fitness")

    fileName <- paste(fileName, ".png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
}

epochal_all_test_runs_boxplot <- function(result, fileName) {
    # Result is a dataframe containing Runs on X axis and values on Y
    data = data.frame(
    A = result$A,
    P = result$P,
    run = result$run,
    AColor = result$AColor
    )
    data$discreteX = as.character(result$run)

    gg <- ggplot(data, aes(x=discreteX, y=P, fill="BestTest"))
    gg <- gg + geom_boxplot(
    outlier.colour="#A4A4A4",
    outlier.shape=16,
    outlier.size=1,
    notch=FALSE,
    fill="green"
    )
    gg <- gg + geom_dotplot(binaxis='y', stackdir='center', dotsize=0.8)
    gg <- gg + stat_summary(fun.y=mean, geom="point", shape=23, size=3, aes(x=discreteX, fill="BestTest"))
    gg <- gg + scale_fill_brewer(palette="YlOrRd") + theme_minimal()
    gg <- gg + scale_colour_manual(values=c(BestBug="green"), name = "Plot Color")

    gg <- gg + guides(
    fill=guide_legend(title="Legend"),
    linetype = guide_legend(title="Legend")
    )
    gg <- gg + theme(
    plot.title = element_text(size=16),
    plot.subtitle = element_text(size=8),
    plot.caption = element_text(size=6))

    gg <- gg + labs(
    color = 'Individuals',
    title = sprintf("%s","Epochal Average of Best Test in Each Run"),
    subtitle = sprintf("%s%d", "Run:", result$run),
    caption = sprintf("Cumulative Test Mean: %.2f\nCumulative Test SDev: %.2f",
    mean(result$P), sd(result$P)),
    x = "Run",
    y = "Fitness")

    fileName <- paste(fileName, ".png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
}

######################################################### EXECUTION ###############
getAllFiles <- function(workDir) {
    files <- list.files(workDir)

    combinedRuns <- data.frame(
        A = double(),
        P = double(),
        run = integer(0),
        AColor = character(),
        PColor = character()
    )

    epochalcount <- 0
    for (file in files) {
        if (grepl("epochal", file)) {
            epochalFileNames[epochalcount] <- file

            print(file)
            filePath <- paste(workDir, file, sep="/")
            epochalData = read_csv(filePath)

            epochal_plot(epochalData, file)
            epochal_pDelta_plot(epochalData, file)
            epochal_aDelta_plot(epochalData, file)
            epochal_density_plot(epochalData, file)
            epochal_density_histogram_plot(epochalData, file)

            meanDataFrame <- data.frame(AColor=rep(as.character(mean(epochalData$A)), nrow(epochalData)))
            meanDataFrameP <- data.frame(PColor=rep(as.character(mean(epochalData$P)), nrow(epochalData)))


            epochalData <- cbind(epochalData, meanDataFrame)
            epochalData <- cbind(epochalData, meanDataFrameP)

            combinedRuns <- rbind(combinedRuns, epochalData)
            epochalcount <- epochalcount + 1
        }
    }

    epochal_all_bug_runs_boxplot(combinedRuns, "epochal_all_bug_runs_boxplot")
    epochal_all_test_runs_boxplot(combinedRuns, "epochal_all_test_runs_boxplot")
}


getAllFiles(workDir)

