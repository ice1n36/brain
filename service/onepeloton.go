package service

import (
	"context"
	"fmt"
	"os"

	"github.com/ice1n36/brain/clients"
	"go.uber.org/config"
)

type OnePelotonService interface {
	GetTotalDistanceTravelledInYear() (map[int]float64, error)
}

type onepelotonservice struct {
	pelotonClient clients.OnePelotonAPIClient
}

func NewOnePelotonService(cfg config.Provider) (OnePelotonService, error) {
	pelotonClient, err := clients.NewOnePelotonAPIClient(cfg)
	if err != nil {
		return nil, err
	}
	return &onepelotonservice{
		pelotonClient: pelotonClient,
	}, nil
}

func (o onepelotonservice) GetTotalDistanceTravelledInYear() (map[int]float64, error) {
	distanceTravelledByYear := map[int]float64{}
	session, err := o.pelotonClient.Login(context.TODO())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//fmt.Println(session)

	summary, err := o.pelotonClient.GetWorkoutSummary(context.TODO(), session)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//fmt.Println(summary)

	workoutsToCount := 0

	for _, numWorkouts := range summary.WorkoutsPerMonth {
		workoutsToCount = workoutsToCount + numWorkouts
	}

	workouts, err := o.pelotonClient.GetWorkouts(context.TODO(), session, workoutsToCount)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//fmt.Println(workoutIds)

	for _, workout := range workouts {
		year := workout.Date.Year()
		if year == 0 {
			fmt.Println("Unable to get year of workout so skipping:", workout.Id)
			continue
		}
		workoutDetail, err := o.pelotonClient.GetWorkoutDetails(context.TODO(), session, workout.Id)
		if err != nil {
			fmt.Println("Unable to get workout details for id: ", workout.Id)
			fmt.Println(err)
			continue
		}

		if distance, ok := distanceTravelledByYear[year]; ok {
			distanceTravelledByYear[year] = distance + workoutDetail.TotalDistance
		} else {
			distanceTravelledByYear[year] = workoutDetail.TotalDistance
		}
	}

	return distanceTravelledByYear, nil
}
