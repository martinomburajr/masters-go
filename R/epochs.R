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
    y = "Fitness")

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
    y = "Fitness")

    fileName <- paste(fileName, "epochal-delta-P.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
}

######################################################### EXECUTION ###############
getAllFiles <- function(workDir) {
    files <- list.files(workDir)
    epochalcount <- 1
    for (file in files) {
        if (grepl("epochal", file)) {
            epochalFileNames[epochalcount] <- file

            print(file)
            filePath <- paste(workDir, file, sep="/")
            epochalData = read_csv(filePath)

            epochal_plot(epochalData, file)
            epochal_pDelta_plot(epochalData, file)
            epochal_aDelta_plot(epochalData, file)

            epochalcount <- epochalcount + 1
        }
    }
}


getAllFiles(workDir)

