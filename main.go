package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// rate of mutation
var MutationRate = 0.005

// size of the population
var PopSize = 500

type Organism struct {
	DNA     []byte
	Fitness float64
}

// create a random population
func createOrganism(target []byte) (org Organism) {
	ba := []byte
	for i := 0; i < len(target); i++ {
		ba[i] = byte(rand.Intn(95) + 32)
	}
	org = Organism{ DNA: ba, Fitness: 0 }
	org.calcFitness(target)
	return
}

func (org *Organism) calcFitness(target []byte) {
	score := 0
	for i := 0; i < len(org.DNA); i++ {
		if org.DNA[i] == target[i] {
			score++
		}
	}
	org.Fitness = float64(score) / float64(len(org.DNA))
}

func createPopulation(target []byte) (population []Organism) {
	population = make([]Organism, PopSize)
	for i := 0; i < PopSize; i++ {
		population[i] = createOrganism(target)
	}
	return
}

// breeding pool
func createPool(population []Organism, target []byte, maxFiteness float64) (poll []Organism) {
	poll = make([]Organism, 0)
	// create a poll of next generation
	for i := 0; i < len(population); i++ {
		population[i].calcFitness(target)
		num := int((population[i].Fitness / maxFiteness) * 100)
		for j := 0; j < num; j++ {
			poll = append(poll, population[i])
		}
	}
	return
}

// bring two parents together and mate so that we can make the natural selection
func naturalSelection(poll []Organism, population []Organism, target []byte) []Organism {
	next := make([]Organism, len(population))
	for i := 0; i < len(population); i++ {
		rand1, rand2 := rand.Intn(len(poll)), rand.Intn(len(poll))
		a := poll[rand1]
		b := poll[rand2]
		child := crossover(a, b)
		child.mutate()
		child.calcFitness(target)
		next[i] = child
	}
	return next
}

// it executes the crossover of 2 parents inheriting the DNA
func crossover(p1 Organism, p2 Organism) Organism {
	child := Organism{
		DNA:     make([]byte, len(p1.DNA)),
		Fitness: 0,
	}
	mid := rand.Intn(len(p1.DNA))
	for i := 0; i < len(p1.DNA); i++ {
		if i > mid {
			child.DNA[i] = p1.DNA[i]
		} else {
			child.DNA[i] = p2.DNA[i]
		}
	}
	return child
}

// mutate the generation so for example if the letter t is not found in the initial population at
// all, we will never be able to come up with the quote no matter how many generation we go throgh
func (org *Organism) mutate() {
	for i := 0; i < len(org.DNA); i++ {
		if rand.Float64() < MutationRate {
			org.DNA[i] = byte(rand.Intn(95) + 32)
		}
	}
}

func getBest(population []Organism) Organism {
	best := 0.0
	index := 0
	for i := 0; i < len(population); i++ {
		if population[i].Fitness > best {
			index = i
			best = population[i].Fitness
		}
	}
	return population[index]
}

func main() {
	start := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())
	target := []byte(os.Args[1])
	population := createPopulation(target)

	found := false
	generation := 0
	for !found {
		generation++
		bestOrganism := getBest(population)
		fmt.Printf("\r generation: %d | %s | fitness: %2f", generation, string(bestOrganism.DNA), bestOrganism.Fitness)
		if bytes.Compare(bestOrganism.DNA, target) == 0 {
			found = true
		} else {
			maxFiteness := bestOrganism.Fitness
			poll := createPool(population, target, maxFiteness)
			population = naturalSelection(poll, population, target)
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("\nTime taken: %s\n", elapsed)
}
