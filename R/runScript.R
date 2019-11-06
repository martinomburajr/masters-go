# Title     : TODO
# Objective : TODO
# Created by: martinomburajr
# Created on: 2019/10/24

args = commandArgs(trailingOnly=TRUE)
# library("rjson")
# library(expss)
# library(dplyr)
library(ggplot2)
library(readr)
library(knitr)
# library(kableExtra)
# devtools::install_github("haozhu233/kableExtra")

generationalFilePath <- ""
epochalFilePath <- ""

if (length(args)==0) {
    stop("At least one argument must be supplied (input file).n", call.=FALSE)
} else {
    # 1 - Path to Generational File
    # 2 - Path to Epochal File
    # 3 - Stats File
    generationalFilePath = args[3]
    epochalFilePath = args[4]
    statsDir = args[5]
    print(statsDir)
    dir.create(file.path(statsDir, ""), showWarnings = FALSE)
    setwd(file.path(statsDir, ""))
}

####################################### CODE BEGINS ##########################
datasetGenerational <- read_csv(generationalFilePath)
datasetEpochal <- read.csv(epochalFilePath)

# Plots out the average between the average of all antagonists in a given geernation, and the average of all
# protagonists in the same generation.
generational_average_plot <- function(result) {
    # result <- data.frame()
    p <- ggplot(data = result,
                mapping = aes(
                    x=result$gen,
                    y=result$avgA))

    p + labs(title = sprintf("%s %d", "Averages for ", result$run),
        x = "Generation", y = "Fitness") +

    # geom_point(
    #     aes(y=result$avgP,
    #            color=result$avgP)) +
        # # topAntagonistReference Plot
    geom_line(
        aes(y=result$avgA, colour="red")) +

        # topProtagonistReference Plot
    geom_line(
        aes(y=result$avgP, colour="green")) +

    geom_line(colour="red",
        aes(x=result$gen,y=result$topA)) +

    # topProtagonistReference Plot
    geom_line(colour="green",
        aes(x=result$gen,y=result$topP))

    ggsave('averages_generational.png', width=8, height=4, units='in', dpi="retina")
    # dev.off()
}


generational_density_plot <- function(result) {
    p <- ggplot(data=result, mapping=aes(x=result$avgA, y=result$gen))
    p + geom_density(kernel="gaussian", mapping=aes(x=result$avgA, y=result$gen))
    ggsave('generational_density.png', width=8, height=4, units='in', dpi="retina")
}

generational_histogram_plot <- function(result) {
    # plotAvgA <- ggplot(data=result, mapping=aes(x=result$avgA))
    # plotAvgA + geom_histogram(binwidth=0.1, mapping=aes(x=result$avgA))
    # ggsave('generational_histogram-avgA.png', plot=plotAvgA,  width=8, height=4, units='in', dpi="retina")


    plotAvgP <- ggplot(data=result, mapping=aes(x=result$avgP))
    plotAvgP + geom_histogram(binwidth=0.1, mapping=aes(x=result$avgP), fill="green", colour="black")
                # geom_histogram(binwidth=0.1, mapping=aes(x=result$avgA), fill="red")
    ggsave('generational_histogram-avgP.png',  width=8, height=4, units='in', dpi="retina")
}

plot_table <- function(result) {
    #Avg
    varAvgA <- var(result$avgA)
    varAvgP <- var(result$avgP)
    sdAvgA <- sd(result$avgA)
    sdAvgP <- sd(result$avgP)
    avgAvgA <- mean(result$avgA)
    avgAvgP <- mean(result$avgP)
    #Cor
    corAvgAP <- cor(result$avgA, result$avgP)
    #Cov
    covAvgAP <- cov(result$avgA, result$avgP)

    # Top
    varTopA <- var(result$topA)
    varTopP <- var(result$topP)
    sdTopP <- sd(result$topP)
    sdTopA <- sd(result$topA)
    avgTopA <- mean(result$topA)
    avgTopP <- mean(result$topP)
    #Cor
    corTopAP <- cor(result$topA, result$topP)
    covTopAP <- cov(result$topA, result$topP)

    #Delta
    varDeltaA <- var(result$topADelta)
    varDeltaP <- var(result$topPDelta)
    sdDeltaA <- sd(result$topADelta)
    sdDeltaP <- sd(result$topPDelta)
    avgDeltaA <- mean(result$topADelta)
    avgDeltaP <- mean(result$topPDelta)
    #Cor
    corDeltaAP <- cor(result$topADelta, result$topPDelta)
    covDeltaAP <- cov(result$topADelta, result$topPDelta)


    # data(result)
    # huxResult <- as_hux(result)
    # ht <- hux(
    #     AntagonistAvg     = result$topA,
    #     ProtagonistAvg       = result$topP,
    #     add_colnames = TRUE
    # )
    # print_screen(ht)

    Stats <- c("Average", "Standard Deviation", "Variance", "Correlation", "Covariance")
    Antagonists <-  c(avgAvgA, sdAvgA, varAvgA, corAvgAP, covAvgAP)
    Protagonists <- c(avgAvgP, sdAvgP, varAvgP, corAvgAP, covAvgAP)
    TopAntagonist <- c(avgTopA, sdTopA, varTopA, corTopAP, covTopAP)
    TopProtagonist <- c(avgTopP, sdTopP, varTopP, corTopAP, covTopAP)
    DeltaAntagonist <- c(avgDeltaA, sdDeltaA, varDeltaA, corDeltaAP, covDeltaAP)
    DeltaProtagonist <- c(avgDeltaP, sdDeltaP, varDeltaP, corDeltaAP, covDeltaAP)
    summaryS <- data.frame(
        Antagonists,
        Protagonists,
        TopAntagonist,
        TopProtagonist,
        DeltaAntagonist,
        DeltaProtagonist
    )
    headings <- c("Antagonist", "Protagonist", "TopAntagonist", "TopProtagonist", "DeltaAntagonist", "DeltaProtagonist")
    names(summaryS) <- headings
    str(summaryS)

    summaryS + kable(x=summaryS) + kable_styling(bootstrap_options = c("striped", "hover"))
    kable(summaryS)

    # print_screen(huxResult)
    # print_screen(summaryS)
    # quick_pdf(summaryS, file="summary.pdf")
    # print_rtf(summaryS)
    # print_md(summaryS, file = "summary.md")
    # print_html(summaryS, file = "summary.html")


    # t <- as.data.frame(x=result$gen, row.names=result$avgA)
    # print(t)
    # p <- ggplot(,
    # mapping = aes(
    # x=result$gen,
    # y=result$avgA))
    #
    # p + labs(title = sprintf("%s %d", "Averages for ", result$run),
    # x = "Generation", y = "Fitness") + geom_bar(stat=result$avgP)
    #
    # ggsave('data.png', width=8, height=4, units='in', dpi="retina")
}

run_stats <- function(datasetGenerational) {
    # generational_average_plot(datasetGenerational)
    generational_histogram_plot(datasetGenerational)
    # generational_density_plot(datasetGenerational)
    # plot_table(datasetGenerational)
}

run_stats(datasetGenerational)


