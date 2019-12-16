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
    title = sprintf("%s","Generation Average and Best Individual Per Run"),
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


################################################ GENERATION #################
################################################ GENERATION #################
################################################ GENERATION #################
################################################ GENERATION #################
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

######################################################### STRATEGY ############################
######################################################### STRATEGY ############################
######################################################### STRATEGY ############################
######################################################### STRATEGY ############################
######################################################### STRATEGY ############################
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


######################################################### TABLE ###############
######################################################### TABLE ###############
######################################################### TABLE ###############
######################################################### TABLE ###############


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

######################################################### EXECUTION ###############
######################################################### EXECUTION ###############
######################################################### EXECUTION ###############
######################################################### EXECUTION ###############
######################################################### EXECUTION ###############

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
        #     generational_histogram_plot(generationalData, file)
        #     generational_density_plot(generationalData,  file)
        #     generational_density_histogram_plot(generationalData,  file)
        #     generational_average_plot(generationalData,  file)
        #     count <- count + 1
        # }
        # if (grepl("epochal", file)) {
        #     epochalFileNames[epochalcount] <- file
        #
        #     print(file)
        #     filePath <- paste(workDir, file, sep="/")
        #     epochalData = read_csv(filePath)
        #
        #     epochal_plot(epochalData, file)
        #     epochal_pDelta_plot(epochalData, file)
        #     epochal_aDelta_plot(epochalData, file)
        #
        #     epochalcount <- epochalcount + 1
        # }
        if (grepl("best-all", file)) {
            bestAllFileNames[bestAllCount] <- file
            filePath <- paste(workDir, file, sep="/")
            bestAllData = read_csv(filePath)
            best_all_function_plot(bestAllData, file)

            bestAllCount <- bestAllCount + 1
        }
        # if (grepl("best-combined", file)) {
        #     bestCombinedFileNames[bestCombinedCount] <- file
        #     filePath <- paste(workDir, file, sep="/")
        #     bestCombinedData = read_csv(filePath)
        #
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
