# Title     : TODO
# Objective : TODO
# Created by: martinomburajr
# Created on: 2019/10/24

args = commandArgs(trailingOnly=TRUE)
# library("rjson")
library("expss")
library(jsonlite)
library("dplyr")
filePath <- ""



if (length(args)==0) {
    stop("At least one argument must be supplied (input file).n", call.=FALSE)
} else {
    filePath = args[2]
    print(filePath)
    statsDir = args[3]
    dir.create(file.path(statsDir, ""), showWarnings = FALSE)
    setwd(file.path(statsDir, ""))
}

toVector <- function(result) {
    g <- result$averages$antagonistCoordinates$independentCoordinate
}

toCSV <- function() {
    raw <- readLines(filePath)
    # get rid of the "/* 0 */" lines
    json <- grep("^/\\* [0-9]* \\*/", raw, value = TRUE, invert = TRUE)
    n <- length(json)
    json[-n] <- gsub("^}$", "},", json[-n])
    json <- c("[", json, "]")
    table <- fromJSON(json)
    flatten(table)[1:3, c(1, 6, 12)]

    tab_list <- lapply(1:nrow(table), function(i) data.frame(table[i, 1], table[i, 1],
    stringsAsFactors = FALSE))
    library(dplyr)
    flat_table <- bind_rows(tab_list)
    write.csv(flat_table, file = "mydata.csv")
}

# result <- fromJSON(file = filePath)
toCSV()

average_plot <- function(result) {
    png('averages.png', width=8, height=4, units='in', res=300)
    print("AVERAGES")
    # AVERAGES + TOP Individual Per Generation + Bottom Individual Per Generation
    plot(
    result$averages$antagonistCoordinates$independentCoordinate,
    result$averages$antagonistCoordinates$dependentCoordinate,
    xlim=c(0,50),
    ylim=c(-1,1),
    main=result$averages$title,
    ylab="Fitness",
    xlab="Generation")

    legend("topleft",
    c("tests", "bugs"),
    fill=c("green", "red"))

    # Generational Averages
    lines(
    result$averages$antagonistCoordinates$independentCoordinate,
    result$averages$antagonistCoordinates$dependentCoordinate,
    col="red")

    lines(
    result$averages$protagonistCoordinates$independentCoordinate,
    result$averages$protagonistCoordinates$dependentCoordinate,
    col="green")
    dev.off()
}
#top_individual returns averages as well as performance of the top individuals
top_individual <- function(result) {
    png('plot.png', width=8, height=4, units='in', res=300)

    # AVERAGES + TOP Individual Per Generation + Bottom Individual Per Generation
    plot(
    result$averages$antagonistCoordinates$independentCoordinate,
    result$averages$antagonistCoordinates$dependentCoordinate,
    xlim=c(0,50),
    ylim=c(-1,1),
    main=result$averages$title,
    ylab="Fitness",
    xlab="Generation")

    legend("topleft",
    c("tests", "bugs"),
    fill=c("green", "red"))

    # Generational Averages
    lines(
    result$averages$antagonistCoordinates$independentCoordinate,
    result$averages$antagonistCoordinates$dependentCoordinate,
    col="red")

    lines(
    result$averages$protagonistCoordinates$independentCoordinate,
    result$averages$protagonistCoordinates$dependentCoordinate,
    col="green")

    # Top Individual
    lines(
    result$topPerGeneration$antagonistCoordinates$independentCoordinate,
    result$topPerGeneration$antagonistCoordinates$dependentCoordinate,
    col="red", pch=22, lty=5)

    lines(
    result$topPerGeneration$protagonistCoordinates$independentCoordinate,
    result$topPerGeneration$protagonistCoordinates$dependentCoordinate,
    col="green", pch=22, lty=5)

    # Bottom Individual
    lines(
    result$bottomPerGeneration$antagonistCoordinates$independentCoordinate,
    result$averages$antagonistCoordinates$dependentCoordinate,
    col="red", pch=23, lty=3)

    lines(
    result$bottomPerGeneration$protagonistCoordinates$independentCoordinate,
    result$bottomPerGeneration$protagonistCoordinates$dependentCoordinate,
    col="green", pch=23, lty=3)
    dev.off()
}

internal_variance <- function(result) {
    # Internal Variance of Best of All Time
    png('internal_variance.png', width=8, height=4, units='in', res=300)
    plot(
    result$ultimateIndividuals$antagonistCoordinates$independentCoordinate,
    result$ultimateIndividuals$antagonistCoordinates$dependentCoordinate,
    xlim=c(0,9),
    ylim=c(-1,1),
    main=result$ultimateIndividuals$title,
    ylab="Fitness",
    xlab="Epoch")

    legend("topleft",
    c("tests", "bugs"),
    fill=c("green", "red"))

    # Generational Averages
    lines(
    result$ultimateIndividuals$antagonistCoordinates$independentCoordinate,
    result$ultimateIndividuals$antagonistCoordinates$dependentCoordinate,
    col="red")

    lines(
    result$ultimateIndividuals$protagonistCoordinates$independentCoordinate,
    result$ultimateIndividuals$protagonistCoordinates$dependentCoordinate,
    col="green")
    print("done")
    dev.off()
}

spec_vs_solutions <- function(result) {
    png('spec_vs_solutions.png', width=7, height=4, units='in', res=300)
    seed <- result$equations$spec$seed
    range <- result$equations$spec$range

    specExpression <- result$equations$spec$expression
    ultAntagonistExpression <- result$equations$ultimateAntagonist$expression
    ultProtagonistExpression <- result$equations$ultimateProtagonist$expression

    print(specExpression)

    specEquation <- function(x){eval(parse(text=specExpression))}
    ultAntagonistEquation <- function(x){eval(parse(text=ultAntagonistExpression))}
    ultProtagonistEquation <- function(x){eval(parse(text=ultProtagonistExpression))}

    plot(
    specEquation,
    from=seed,
    to=(seed+range),
    ylab="Y",
    xlab="X",
    )

    plot(
    ultAntagonistEquation,
    from=seed,
    to=(seed+range),
    ylab="Y",
    xlab="X",
    col="red",
    add=TRUE
    )

    plot(
    ultProtagonistEquation,
    from=seed,
    to=(seed+range),
    ylab="Y",
    xlab="X",
    col="green",
    add=TRUE
    )

    par(xpd=NA)
    legend("right",
    c("spec", "tests", "bugs"),
    fill=c("black", "green", "red"), title="Top Competitors vs Spec")
    dev.off()
    ############################## SPEC #################################
}

parse_individuals <- function(result) {
    average_plot(result)
    top_individual(result)
    internal_variance(result)
    spec_vs_solutions(result)
}

parse_table <- function(result) {
#     Create Table for all Averages for Generations (Antagonist, Protagonist)
#     averagesTable <- data.frame(
#         generation: c(0:49)
#         antagonist: result$averages$antagonistCoordinates$independentCoordinate
#         protagonist: result$averages$protagonistCoordinates$independentCoordinate
#     )

    print(result)
    data(result$averages)
    averagesTable = apply_labels(result$averages$protagonistCoordinates,
    independentCoordinates = "Protagonist Independent Coordinates",
    dependentCoordinates = "Protagonist Dependent Coordinates")

    # png('table-averages.png', width=7, height=4, units='in', res=300)
    # print(averagesTable)

    # cro(result$averages$antagonistCoordinates, result$averages$protagonistCoordinates)
}
# Generational Averages

# parse_individuals(result)
parse_table(result)

# lines(
#     result$ultimateIndividuals$protagonistCoordinates$independentCoordinate,
#     result$ultimateIndividuals$protagonistCoordinates$dependentCoordinate,
#     col="green")
# # Spec
# lines(
#     result$ultimateIndividuals$protagonistCoordinates$independentCoordinate,
#     result$ultimateIndividuals$protagonistCoordinates$dependentCoordinate,
#     col="green")
print("done")
# dev.off()

