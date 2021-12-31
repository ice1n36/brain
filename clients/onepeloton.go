package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"go.uber.org/config"
)

const (
	_onePelotonAPIURL                        = "api.onepeloton.com"
	_onePelotonLoginAPI                      = "auth/login"
	_onePelotonGetWorkoutsAPI                = "api/user/%s/workouts"
	_onePelotonGetWorkoutPerformanceGraphAPI = "api/workout/%s/performance_graph"
)

type OnePelotonSession struct {
	SessionId string
	UserId    string
}

type OnePelotonWorkoutSummary struct {
	WorkoutsPerMonth map[string]int
}

type OnePelotonWorkoutDetails struct {
	TotalOutput   float64
	TotalDistance float64
	TotalCalories float64
}

type OnePelotonWorkout struct {
	Id   string
	Date time.Time
}

type OnePelotonAPIClient interface {
	Login(ctx context.Context) (OnePelotonSession, error)
	GetWorkoutSummary(ctx context.Context, ops OnePelotonSession) (OnePelotonWorkoutSummary, error)
	GetWorkouts(ctx context.Context, ops OnePelotonSession, amount int) ([]OnePelotonWorkout, error)
	GetWorkoutDetails(ctx context.Context, ops OnePelotonSession, workoutId string) (OnePelotonWorkoutDetails, error)
}

type onepelotonclient struct {
	username string
	password string
	client   *http.Client
}

type onepelotonconfig struct {
	UsernameOrEmail string `yaml:"username_or_email"`
	Password        string `yaml:"password"`
}

type loginresponse struct {
	UserId    string `json:"user_id"`
	SessionId string `json:"session_id"`
}

//"data": [
//{
//"created_at": 1640622696,
//"device_type": "home_bike_v1",
//"end_time": 1640623956,
//"fitness_discipline": "cycling",
//"has_pedaling_metrics": true,
//"has_leaderboard_metrics": true,
//"id": "",
//"is_total_work_personal_record": false,
//"metrics_type": "cycling",
//"name": "Cycling Workout",
//"peloton_id": "",
//"platform": "home_bike",
//"start_time": 1640622758,
//"status": "COMPLETE",
//"timezone": "Etc/GMT+8",
//"title": null,
//"total_work": 200806.97,
//"user_id": "",
//"workout_type": "class",
//"total_video_watch_time_seconds": 0,
//"total_video_buffering_seconds": 0,
//"v2_total_video_watch_time_seconds": 1379,
//"v2_total_video_buffering_seconds": 1,
//"total_music_audio_play_seconds": null,
//"total_music_audio_buffer_seconds": null,
//"created": 1640622696,
//"device_time_created_at": 1640593896,
//"strava_id": null,
//"fitbit_id": null,
//"effort_zones": null
//}
//],
type workoutdata struct {
	Id          string `json:"id"`
	MetricsType string `json:"metrics_type"`
	Created     int64  `json:"created"`
	Status      string `json:"status"`
}

type workoutsresponse struct {
	Summary map[string]int `json:"summary"`
	Data    []workoutdata  `json:"data"`
}

type perfsummary struct {
	DisplayName string  `json:"display_name"`
	DisplayUnit string  `json:"display_unit"`
	Value       float64 `json:"value"`
	Slug        string  `json:"slug"`
}

type workoutperformancegraphresponse struct {
	Summaries []perfsummary `json:"summaries"`
}

func NewOnePelotonAPIClient(cfg config.Provider) (OnePelotonAPIClient, error) {
	var opcfg onepelotonconfig
	if err := cfg.Get("onepeloton").Populate(&opcfg); err != nil {
		// fallback to using environment variables
		opcfg.UsernameOrEmail = os.Getenv("ONEPELOTON_USER")
		opcfg.Password = os.Getenv("ONEPELOTON_PW")
	}
	if len(opcfg.UsernameOrEmail) == 0 {
		return nil, errors.New("Cannot retrieve username from config or environment variable ONEPELOTON_USER")
	}
	if len(opcfg.Password) == 0 {
		return nil, errors.New("Cannot retrieve password from config or environment variable ONEPELOTON_PW")
	}

	return &onepelotonclient{
		username: opcfg.UsernameOrEmail,
		password: opcfg.Password,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}, nil
}

func (o onepelotonclient) Login(ctx context.Context) (OnePelotonSession, error) {
	creds := map[string]string{
		"username_or_email": o.username,
		"password":          o.password,
	}
	creds_json, err := json.Marshal(creds)
	var ops OnePelotonSession
	if err != nil {
		return ops, errors.New("json marshalling error on creds")
	}
	loginURI := fmt.Sprintf("https://%s/%s", _onePelotonAPIURL, _onePelotonLoginAPI)

	request, err := http.NewRequest("POST", loginURI, bytes.NewBuffer(creds_json))
	request.Header.Add("Content-Type", "application/json")
	resp, err := o.client.Do(request)
	if err != nil {
		fmt.Println(err)
		return ops, errors.New("error in auth/login")
	}
	respBody := loginresponse{}
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ops, err
	}

	err = json.Unmarshal(respBodyData, &respBody)
	if err != nil {
		return ops, err
	}

	if len(respBody.SessionId) == 0 {
		return ops, errors.New("empty session_id")
	}
	if len(respBody.UserId) == 0 {
		return ops, errors.New("empty user_id")
	}

	ops.SessionId = respBody.SessionId
	ops.UserId = respBody.UserId
	return ops, nil
}

func (o onepelotonclient) GetWorkoutSummary(ctx context.Context, ops OnePelotonSession) (OnePelotonWorkoutSummary, error) {
	summary := OnePelotonWorkoutSummary{}
	getWorkoutsURI := fmt.Sprintf(_onePelotonGetWorkoutsAPI, ops.UserId)
	uri := fmt.Sprintf("https://%s/%s?limit=1", _onePelotonAPIURL, getWorkoutsURI)

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println(err)
		return summary, errors.New("error creating request")
	}
	sessionCookie := &http.Cookie{
		Name:   "peloton_session_id",
		Value:  ops.SessionId,
		MaxAge: 300,
	}
	request.AddCookie(sessionCookie)
	request.Header.Add("Content-Type", "application/json")
	resp, err := o.client.Do(request)
	if err != nil {
		fmt.Println(err)
		return summary, errors.New("error in http call to get workout summary")
	}
	respBody := workoutsresponse{}
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return summary, err
	}

	err = json.Unmarshal(respBodyData, &respBody)
	if err != nil {
		return summary, err
	}

	summary.WorkoutsPerMonth = respBody.Summary

	return summary, nil
}

func (o onepelotonclient) GetWorkouts(ctx context.Context, ops OnePelotonSession, amount int) ([]OnePelotonWorkout, error) {
	workouts := []OnePelotonWorkout{}
	pages := amount / 100
	getWorkoutsURI := fmt.Sprintf(_onePelotonGetWorkoutsAPI, ops.UserId)
	for i := 0; i <= pages; i++ {
		uri := fmt.Sprintf("https://%s/%s?limit=100&page=%d&sort_by=-created", _onePelotonAPIURL, getWorkoutsURI, i)

		request, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			fmt.Println(err)
			return workouts, errors.New("error creating request")
		}
		sessionCookie := &http.Cookie{
			Name:   "peloton_session_id",
			Value:  ops.SessionId,
			MaxAge: 300,
		}
		request.AddCookie(sessionCookie)
		request.Header.Add("Content-Type", "application/json")
		resp, err := o.client.Do(request)
		if err != nil {
			fmt.Println(err)
			return workouts, errors.New("error in http call to get workout summary")
		}
		respBody := workoutsresponse{}
		respBodyData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return workouts, err
		}

		err = json.Unmarshal(respBodyData, &respBody)
		if err != nil {
			return workouts, err
		}

		for _, data := range respBody.Data {
			workout := OnePelotonWorkout{
				Id:   data.Id,
				Date: time.Unix(data.Created, 0),
			}
			workouts = append(workouts, workout)
		}
	}
	return workouts, nil
}

func (o onepelotonclient) GetWorkoutDetails(ctx context.Context, ops OnePelotonSession, workoutId string) (OnePelotonWorkoutDetails, error) {
	workoutDetails := OnePelotonWorkoutDetails{}
	getWorkoutPerformanceURI := fmt.Sprintf(_onePelotonGetWorkoutPerformanceGraphAPI, workoutId)
	uri := fmt.Sprintf("https://%s/%s?every_n=1000", _onePelotonAPIURL, getWorkoutPerformanceURI)
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println(err)
		return workoutDetails, errors.New("error creating request")
	}
	sessionCookie := &http.Cookie{
		Name:   "peloton_session_id",
		Value:  ops.SessionId,
		MaxAge: 300,
	}
	request.AddCookie(sessionCookie)
	request.Header.Add("Content-Type", "application/json")
	resp, err := o.client.Do(request)
	if err != nil {
		fmt.Println(err)
		return workoutDetails, errors.New("error in http call to get workout details")
	}
	respBody := workoutperformancegraphresponse{}
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return workoutDetails, err
	}

	err = json.Unmarshal(respBodyData, &respBody)
	if err != nil {
		return workoutDetails, err
	}
	for _, perfsummary := range respBody.Summaries {
		if perfsummary.Slug == "total_output" {
			workoutDetails.TotalOutput = perfsummary.Value
		} else if perfsummary.Slug == "distance" {
			workoutDetails.TotalDistance = perfsummary.Value
		} else if perfsummary.Slug == "calories" {
			workoutDetails.TotalCalories = perfsummary.Value
		} else {
			fmt.Println("WARNING - unknown slug: ", perfsummary.Slug)
		}
	}
	return workoutDetails, nil

}
