# Title     : TODO
# Objective : TODO
# Created by: martinomburajr
# Created on: 2019/10/24

args = commandArgs(trailingOnly=TRUE)
library(ggplot2)
library(readr)
library(knitr)


workDir <- ""
statsDir <- ""

if (length(args)==0) {
    stop("At least one argument must be supplied (input file).n", call.=FALSE)
} else {
    # 1 - Path to Generational File
    # 2 - Path to Epochal File
    # 3 - Stats File
    print(args)
    workDir = args[1]
    # epochalFilePath = args[4]
    # statsDir = args[5]
    print("SET WORKING DIRECTORY")
    print(workDir)
    statsDir <- workDir
    dir.create(file.path(statsDir, ""), showWarnings = FALSE)
    setwd(file.path(statsDir, ""))
}


generationalFileNames <- c()
generationalFileNames2 <- c()
epochalFileNames <- c()
bestFileNames <- c()
bestCombinedFileNames <- c()
bestAllFileNames <- c()
strategyFileNames <- c()

####################################### CODE BEGINS ##########################

######################################## EPOCH
epochal_plot <- function(result, fileName) {
    p <- ggplot(data = result,
    mapping = aes(
        x=result$epoch,
        y=result$A))

    p + labs(
        title = sprintf("%s","Epoch Based Fitness Variation of Bug and Test"),
        subtitle = sprintf("%s%d", "Run:", result$run),
        caption = sprintf("%s%d", "Run #", result$run),
        x = "Epoch", y = "Fitness"
        # tag="Legend"
        ) +

        geom_line(aes(x=result$epoch,y=result$A, colour="red")) +
        geom_line(aes(x=result$epoch,y=result$P, colour="green")) +
        geom_line(aes(x=result$epoch,y=result$finA, colour="orange"), linetype="dashed") +
        geom_line(aes(x=result$epoch,y=result$finP, colour="blue"), linetype="dashed") +

        scale_color_manual(values = c(
        'red' = 'red',
        'green' = 'green'),
        'orange' = 'orange',
        'blue' = 'blue')) +
        labs(color = 'Y series') +

        # scale_color_discrete(name = "Individuals",
        #     labels = c("Best Bug", "Best Test", "FinGen Bug", "FinGen Test")) +

        theme(
            plot.title = element_text(size=16),
            plot.subtitle = element_text(size=12)
        )

        # theme(
        #     # Legend title and text labels
        #     #:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
        #     # Title font color size and face
        #     # legend.title = element_text(color, size, face),
        #     # Title alignment. Number from 0 (left) to 1 (right)
        #     legend.title.align = NULL,
        #     # Text label font color size and face
        #     # legend.text = element_text(color, size, face),
        #     # Text label alignment. Number from 0 (left) to 1 (right)
        #     legend.text.align = NULL,
        #
        #     # Legend position, margin and background
        #     #:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
        #     # Legend position: right, left, bottom, top, none
        #     legend.position = "bottom",
        #     # Margin around each legend
        #     legend.margin = margin(0.2, 0.2, 0.2, 0.2, "cm"),
        #     # Legend background
        #     # legend.background = element_rect(fill, color, size, linetype),
        #
        #     # Legend direction and justification
        #     #:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
        #     # Layout of items in legends ("horizontal" or "vertical")
        #     legend.direction = "horizontal",
        #     # Positioning legend inside or outside plot
        #     # ("center" or two-element numeric vector)
        #     legend.justification = "center",
        #
        #     # Background underneath legend keys
        #     #:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
        #     # legend.key = element_rect(fill, color),  # Key background
        #     legend.key.size = unit(1.2, "lines"),    # key size (unit)
        #     legend.key.height = NULL,                # key height (unit)
        #     legend.key.width = NULL,                 # key width (unit)
        #
        #     # Spacing between legends.
        #     #:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
        #     legend.spacing = unit(0.4, "cm"),
        #     legend.spacing.x = NULL,                 # Horizontal spacing
        #     legend.spacing.y = NULL,                 # Vertical spacing
        #
        #     # Legend box
        #     #:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
        #     # Arrangement of multiple legends ("horizontal" or "vertical")
        #     legend.box = NULL,
        #     # Margins around the full legend area
        #     legend.box.margin = margin(0, 0, 0, 0, "cm"),
        #     # Background of legend area: element_rect()
        #     legend.box.background = element_blank(),
        #     # The spacing between the plotting area and the legend box
        #     legend.box.spacing = unit(0.4, "cm")
        # )


    fileName <- paste(fileName, "epochal.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
    # dev.off()
}


#### EPOCH-DELTA

epochal_aDelta_plot <- function(result, fileName) {
    p <- ggplot(data = result,
    mapping = aes(x=result$epoch, y=result$aDelta))

    p + labs(title = sprintf("%s %d", "Epoch for ", result$run),
    x = "Epoch", y = "Antagonist Delta") +

        geom_line(
        aes(x=result$epoch,y=result$ADelta),
        colour="red"
        ) +

        geom_line(
        aes(x=result$epoch,y=result$finADelta),
        linetype="dashed",
        colour="red"
        )

    fileName <- paste(fileName, "epochal-delta-A.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
}

epochal_pDelta_plot <- function(result, fileName) {
    p <- ggplot(data = result,
    mapping = aes(x=result$epoch, y=result$PDelta))

    p + labs(title = sprintf("%s %d", "Epoch for ", result$run),
    x = "Epoch", y = "Protagonist Delta") +

    geom_line(
        aes(x=result$epoch,y=result$PDelta),
        colour="green"
    ) +

    geom_line(
        aes(x=result$epoch,y=result$finPDelta),
        linetype="dashed",
        colour="green"
    )

    fileName <- paste(fileName, "epochal-delta-P.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
}

################################################### BEST
best_combined_avg_plot <- function(result, fileName) {
    p <- ggplot(data=result, mapping = aes(
        x=result$ARun,
        y=result$AAvg)
    )

    p + labs(title = sprintf("%s", "Best Individual per Run "),
    x = "Run", y = "Fitness") +

    geom_line(aes(x=result$ARun,y=result$AAvg),
    color="red"
    ) +

    # topProtagonistReference Plot
    geom_line(aes(x=result$ARun,y=result$PAvg),
    color="green"
    ) +

    geom_line(aes(x=result$ARun,y=result$ABestFit),
    color="red",
    linetype="dashed"
    ) +

    # topProtagonistReference Plot
    geom_line(aes(x=result$ARun,y=result$PBestFit),
    color="green",
    linetype="dashed"
    )

    fileName <- paste(fileName, "best-combined_avg.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
}

    ## Pass in best-all.csv
best_all_function_plot <- function(result, fileName) {
    p <- ggplot(data = data.frame(x = result$seed), mapping = aes(x = x))

    specEquation <- function(x){eval(parse(text=result$spec))}
    ultAntagonistEquation <- function(x){eval(parse(text=result$AEquation))}
    ultProtagonistEquation <- function(x){eval(parse(text=result$PEquation))}

    #spec
    p + layer(geom = "path",        # Default. Can be omitted.
            stat = "function",
            fun = specEquation,          # Give function
            mapping = aes(color = "specEquation") # Give a meaningful name to color
            ) +

    layer(geom = "path",        # Default. Can be omitted.
        stat = "function",
        fun = ultAntagonistEquation,          # Give function
        mapping = aes(color = "ultAntagonistEquation") # Give a meaningful name to color
    ) +

    layer(geom = "path",        # Default. Can be omitted.
        stat = "function",
        fun = ultProtagonistEquation,          # Give function
        mapping = aes(color = "ultProtagonistEquation") # Give a meaningful name to color
    ) +

    scale_x_continuous(limits=c(result$seed,result$seed + result$range)) +
    scale_color_manual(name = "Functions",
        values = c("black", "red", "green"),
        labels = c(specEquation, ultAntagonistEquation, ultProtagonistEquation))

    fileName <- paste(fileName, "best-all.png", sep="-")
    ggsave(fileName, width=8, height=4, units='in', dpi="retina")
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

    geom_line(aes(x=result$gen,y=result$AGenFitAvg),
        color="red"
    ) +

        # topProtagonistReference Plot
    geom_line(aes(x=result$gen,y=result$PGenFitAvg),
        color="green"
    ) +

    geom_line(aes(x=result$gen,y=result$AGenBestFitAvg),
        color="red",
        linetype="dashed"
    ) +

    # topProtagonistReference Plot
    geom_line(aes(x=result$gen,y=result$PGenBestFitAvg),
        color="green",
        linetype="dashed"
    )

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

}

################################################## STRATEGY

strategy_run_histogram_plot <- function(result, filename) {
    p <- ggplot(data=result, mapping=aes(x=result$A))
    p + geom_histogram(
        alpha=0.7,
        stat="count",
        position="identity",
        # mapping=aes(y=result$P, color="yellow")
    ) +
    geom_histogram(
        alpha=0.7,
        stat="count",
        position="identity",
        # mapping=aes(y=result$P, color="yellow")
    ) +
    geom_density(alpha=0.4) +
    # geom_vline(
    #     aes(xintercept=7.5),
    #     color="black",
    #     linetype="dashed",
    #     size=1) +
    # labs(x=feature, y = "Density")

    #     fill="red", colour="black", alpha=0.2) +
    #     # geom_histogram(stat="count", mapping=aes(result$A), fill="red", colour="black", alpha=0.2) +
    # geom_density() +
    labs(title="Frequency of Strategy in Best Individuals", x="Strategy", y="Frequency")

    outputPath <- paste(filename, "histogram-bestP.png", sep="-")
    ggsave(outputPath,  width=12, height=6, units='in', dpi="retina")
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
    bestAllCount <- 1
    bestCombinedCount <- 1
    strategyCount <- 1
    for (file in files) {
        # if (grepl("generational", file)) {
        #     generationalFileNames[count] <- file
        #
        #     filePath <- paste(workDir, file, sep="/")
        #     print(filePath)
        #     generationalData = read_csv(filePath)
        #
        #     generational_average_plot(generationalData,  file)
        #     count <- count + 1
        # }
        if (grepl("epochal", file)) {
            epochalFileNames[epochalcount] <- file

            print(file)
            filePath <- paste(workDir, file, sep="/")
            epochalData = read_csv(filePath)
            epochal_plot(epochalData, file)
            # epochal_pDelta_plot(epochalData, file)
            # epochal_aDelta_plot(epochalData, file)

            epochalcount <- epochalcount + 1
        }
        # if (grepl("best-all", file)) {
        #     bestAllFileNames[bestAllCount] <- file
        #     filePath <- paste(workDir, file, sep="/")
        #     bestAllData = read_csv(filePath)
        #     best_all_function_plot(bestAllData, file)
        #
        #     bestAllCount <- bestAllCount + 1
        # }
        # if (grepl("best-combined", file)) {
        #     bestCombinedFileNames[bestCombinedCount] <- file
        #     filePath <- paste(workDir, file, sep="/")
        #     bestCombinedData = read_csv(filePath)
        #     best_combined_avg_plot(bestCombinedData, file)
        #
        #     bestCombinedCount <- bestCombinedCount + 1
        # }
        # if (grepl("strategy", file)) {
        #     strategyFileNames[strategyCount] <- file
        #     filePath <- paste(workDir, file, sep="/")
        #     strategyData = read_csv(filePath)
        #
        #     strategy_run_histogram_plot(strategyData, file)
        #     strategyCount <- strategyCount + 1
        # }

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


