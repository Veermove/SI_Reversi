package main

func main() {
	RunNormaln()

	// bWeights, wWeights := RunLearn()
	// StartGame(BestPlayer(6, &wWeights), BestPlayer(6, &bWeights), false)
}

func RunNormaln() {
	StartGame(BestPlayer(4, nil), BestPlayer(4, nil), false)
}
