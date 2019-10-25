# Title     : TODO
# Objective : TODO
# Created by: martinomburajr
# Created on: 2019/10/24

library("rjson")
result <- fromJSON(file = "run.json")
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


# Internal Variance of Best of All Time
png('plot-ultimate.png', width=8, height=4, units='in', res=300)
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


############################## SPEC #################################
png('plot-spec.png', width=7, height=4, units='in', res=300)
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

# Generational Averages


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
dev.off()