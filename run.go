package main

import (
	"fmt"
	"math/rand"
	"sync"
)


func RunLearn() ([5]float64, [5]float64) {
	wWeights := [5]float64 {1, 1, 1, 1, 1}
	bWeights := [5]float64 {1, 1, 1, 1, 1}

	weightsShift := float64(2)
	depth := 4
	populationSize := 10
	generations := 15

	wWins := 0
	bWins := 0



	for i := 0; i < generations; i++ {
		fmt.Println("White gen", i, ":", wWeights)
		fmt.Println("Black gen", i, ":", bWeights)

		var wg sync.WaitGroup

		resultingWeightsBlack := make([][5]float64, populationSize)
		resultingWeightsWhite := make([][5]float64, populationSize)

		for ind := 0; ind < populationSize; ind++ {
			wg.Add(1)

			go func(ex int) {
				defer wg.Done()

				bCurrentWeights := bWeights
				wCurrentWeights := wWeights

				// offset weights randomly a bit
				for w := 0; w < len(bWeights); w++ {
					bCurrentWeights[w] *= (rand.Float64() + 0.5) * float64(weightsShift)
					wCurrentWeights[w] *= (rand.Float64() + 0.5) * float64(weightsShift)
				}

				result := StartGame(BestPlayer(depth, &wCurrentWeights), BestPlayer(depth, &bCurrentWeights), true)

				if result == BLACK {
					resultingWeightsBlack[ex] = bCurrentWeights
					bWins++
				} else {
					resultingWeightsWhite[ex] = wCurrentWeights
					wWins++
				}
			} (ind)
		}

		wg.Wait()

		var totalsB float64
		var totalsW float64

		// calculate new weights:
		// AVG( CurrentWeight, [Weights of winners] )
		for w := 0; w < len(bWeights); w++ {
			var bCounter int
			var wCounter int
			var newBlackWeight float64
			var newWhiteWeight float64

			for ind := 0; ind < populationSize; ind ++ {
				if len(resultingWeightsBlack[ind]) != 0 && resultingWeightsBlack[ind][w] != 0 {
					newBlackWeight += resultingWeightsBlack[ind][w]
					bCounter++
				}
				if len(resultingWeightsWhite[ind]) != 0 && resultingWeightsWhite[ind][w] != 0 {
					newWhiteWeight += resultingWeightsWhite[ind][w]
					wCounter++
				}
			}

			bWeights[w] = (newBlackWeight + bWeights[w]) / (float64(bCounter) + 1)
			wWeights[w] = (newWhiteWeight + wWeights[w]) / (float64(wCounter) + 1)

			totalsB += bWeights[w]
			totalsW += wWeights[w]
		}

		// normalize weights
		// so that SUM( [ weights ] ) == len( [ weights ] )
		for w := 0; w < len(bWeights); w++ {
			bWeights[w] = (bWeights[w] / totalsB) * float64(len(bWeights))
			wWeights[w] = (wWeights[w] / totalsW) * float64(len(wWeights))
		}
	}

	fmt.Println("Black wins", bWins)
	fmt.Println("White wins", wWins)

	return bWeights, wWeights
}
