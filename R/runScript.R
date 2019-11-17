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

workDir <- ""
statsDir <- ""

if (length(args)==0) {
    stop("At least one argument must be supplied (input file).n", call.=FALSE)
} else {
    # 1 - Path to Generational File
    # 2 - Path to Epochal File
    # 3 - Stats File
    workDir = args[3]
    # epochalFilePath = args[4]
    # statsDir = args[5]
    print(workDir)
    statsDir <- workDir
    dir.create(file.path(statsDir, ""), showWarnings = FALSE)
    setwd(file.path(statsDir, ""))
}


generationalFileNames <- c()
generationalFileNames2 <- c()
epochalFileNames <- c()
bestFileNames <- c()




####################################### CODE BEGINS ##########################

######################################## EPOCH
epochal_plot <- function(result, fileName) {
    p <- ggplot(data = result,
    mapping = aes(
    x=result$epoch,
    y=result$A))

    p + labs(title = sprintf("%s %d", "Epoch for ", result$run),
    x = "Epoch", y = "Fitness") +

        geom_line(
        aes(y=result$A, colour="red")) +

        geom_line(
        aes(y=result$P, colour="green"))


    # finalAntagonist Plot
        geom_line(colour="red",
        aes(x=result$epoch,y=result$finA)) +

    # finalProtagonist Plot
        geom_line(colour="green",
        aes(x=result$epoch,y=result$finP))


    fileName <- paste(fileName, "epochal.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
    # dev.off()
}


######################################## EPOCH-DELTA

epochal_plot_aDelta <- function(result, fileName) {
    p <- ggplot(data = result,
    mapping = aes(
    x=result$epoch,
    y=result$ADelta))

    p + labs(title = sprintf("%s %d", "Epoch for ", result$run),
    x = "Epoch", y = "Delta") +

    geom_line(
    aes(y=result$ADelta, colour="red")) +

    geom_line(
    aes(y=result$finADelta, colour="blue"))
    aes(x=result$epoch,y=result$finA))

    fileName <- paste(fileName, "epochal-delta-A.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
    # dev.off()
}


epochal_plot_pDelta <- function(result, fileName) {
    p <- ggplot(data = result,
    mapping = aes(
    x=result$epoch,
    y=result$PDelta))

    p + labs(title = sprintf("%s %d", "Epoch for ", result$run),
    x = "Epoch", y = "Delta") +

    geom_line(
    aes(y=result$pDelta, colour="red")) +

    geom_line(
    aes(y=result$finPDelta, colour="blue"))


    aes(x=result$epoch,y=result$finP))

    fileName <- paste(fileName, "epochal-delta-P.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
    # dev.off()
}

# Plots out the average between the average of all antagonists in a given geernation, and the average of all
# protagonists in the same generation.
generational_average_plot <- function(result, fileName) {
    # result <- data.frame()
    p <- ggplot(data = result,
                mapping = aes(
                    x=result$gen,
                    y=result$AGenFitAvg))

    p + labs(title = sprintf("%s %d", "Averages for ", result$run),
        x = "Generation", y = "Fitness") +

    geom_line(
        aes(y=result$AGenFitAvg, colour="red")) +

        # topProtagonistReference Plot
    geom_line(
        aes(y=result$PGenFitAvg, colour="green")) +

    geom_line(colour="red",
        aes(x=result$gen,y=result$AGenBestFitAvg)) +

    # topProtagonistReference Plot
    geom_line(colour="green",
        aes(x=result$gen,y=result$PGenBestFitAvg))


    fileName <- paste(fileName, "generational.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
    # dev.off()
}


generational_density_plot <- function(result, filename) {
    p <- ggplot(data=result, mapping=aes(x=result$AGenFitAvg, y=result$gen))
    p + geom_density(kernel="gaussian", mapping=aes(x=result$AGenFitAvg, y=result$gen))

    outputPath <- paste(filename, "generational_density.png.png", sep="-")
    ggsave('outputPath', width=8, height=4, units='in', dpi="retina")
}

generational_histogram_plot <- function(result, filename) {
    # plotAvgA <- ggplot(data=result, mapping=aes(x=result$avgA))
    # plotAvgA + geom_histogram(binwidth=0.1, mapping=aes(x=result$avgA))
    # ggsave('generational_histogram-avgA.png', plot=plotAvgA,  width=8, height=4, units='in', dpi="retina")


    plotAvgP <- ggplot(data=result, mapping=aes(x=result$PGenFitAvg))
    plotAvgP + geom_histogram(binwidth=0.1, mapping=aes(x=result$PGenFitAvg), fill="green", colour="black")
                # geom_histogram(binwidth=0.1, mapping=aes(x=result$avgA), fill="red")
    outputPath <- paste(filename, "generational_histogram-avgP.png", sep="-")
    ggsave(outputPath,  width=8, height=4, units='in', dpi="retina")
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

runGenerational <- function(generationalFiles) {
    print("Running Generational Files")
    print(length(generationalFiles))
    for (generationalFile in generationalFiles) {
        filePath <- paste(workDir, generationalFile)
        print(filePath)
        generationalData = read_csv(filePath)

        # functions
        generational_average_plot(generationalData,  generationalFile)
    }
}

# runGenerational(generationalFileNames2)

getAllFiles <- function(workDir) {
    files <- list.files(workDir)
    count <- 1
    epochalcount <- 1
    bestcount <- 1
    for (file in files) {
        if (grepl("generational", file)) {
            generationalFileNames[count] <- file

            filePath <- paste(workDir, file, sep="/")
            print(filePath)
            generationalData = read_csv(filePath)

            generational_average_plot(generationalData,  file)
            count <- count + 1
        }
        if (grepl("epochal", file)) {
            epochalFileNames[epochalcount] <- file

            filePath <- paste(workDir, file, sep="/")
            epochalData = read_csv(filePath)
            epochal_plot(epochalData, file)
            epochal_plot_pDelta(epochalData, file)
            epochal_plot_aDelta(epochalData, file)

            epochalcount <- epochalcount + 1
        }
        if (grepl("best", file)) {
            bestFileNames[bestcount] <- file
            bestcount <- bestcount + 1
        }
    }
    generationalFileNames2 <- generationalFileNames
    print(length(generationalFileNames))
}

getAllFiles(workDir)

# run_stats <- function(datasetGenerational) {
#     generational_average_plot(datasetGenerational)
#     generational_histogram_plot(datasetGenerational)
#     # generational_density_plot(datasetGenerational)
#     # plot_table(datasetGenerational)
# }

# run_stats(datasetGenerational)


