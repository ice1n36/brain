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
	session_id string
	user_id    string
}

type OnePelotonWorkoutDetails struct {
	total_output   string
	total_distance string
	total_calories string
}

type OnePelotonAPIClient interface {
	Login(ctx context.Context) (OnePelotonSession, error)
	GetWorkoutIds(ctx context.Context, ops OnePelotonSession, amount int) ([]string, error)
	GetWorkoutDetails(ctx context.Context, ops OnePelotonSession, workouts []string) ([]OnePelotonWorkoutDetails, error)
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

	ops.session_id = respBody.SessionId
	ops.user_id = respBody.UserId
	return ops, nil
}

func (o onepelotonclient) GetWorkoutIds(ctx context.Context, ops OnePelotonSession, amount int) ([]string, error) {
	workouts := []string{}
	return workouts, nil
}

func (o onepelotonclient) GetWorkoutDetails(ctx context.Context, ops OnePelotonSession, workouts []string) ([]OnePelotonWorkoutDetails, error) {
	workoutDetails := []OnePelotonWorkoutDetails{}
	return workoutDetails, nil

}
