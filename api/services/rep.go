package services

const (
	MAXREP = 25
)

func CalculateCurrentAnswerEligibilityRep(currentRep int) int {
	return MAXREP - currentRep
}
