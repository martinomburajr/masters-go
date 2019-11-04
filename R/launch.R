# Title     : TODO
# Objective : TODO
# Created by: martinomburajr
# Created on: 2019/10/24

args = commandArgs(trailingOnly=TRUE)
# library("rjson")
library(expss)
library(jsonlite)
library(dplyr)
library(ggplot2)
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


####################################### CODE BEGINS ##########################

dataset <- read.csv(filePath)

# Plots out the average between the average of all antagonists in a given geernation, and the average of all
# protagonists in the same generation.
average_plot <- function(result) {
    png('averages.png', width=8, height=4, units='in', res=300)
    p <- ggplot(result,
                aes(x=result$generation,
                    y=result$averageAntagonist,
                    color=result$averageAntagonist))
    
    p + geom_point() + 
        geom_smooth() + 
        geom_point(
            aes(y=result$averageProtagonist,
                       color=result$averageProtagonist)) +
        geom_smooth() +

        # topAntagonistReference Plot
        geom_line(
            aes(y=result$topAntagonist, color="red")) +

        # topProtagonistReference Plot
        geom_line(
            aes(y=result$topProtagonist, color="green"))

    ggsave('averages.png', width=8, height=4, units='in', dpi="retina")
}

run_stats <- function(result) {
    average_plot(result)
}

run_stats(dataset)

print("done")
dev.off()

