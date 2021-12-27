package clients

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/config"
	"gopkg.in/square/go-jose.v2/json"
)

const (
	_onePelotonAPIURL                        = "https://api.onepeloton.com"
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
	GetWorkoutDetails(ctx context.Context, ops OnePelotonSession, workouts []string) ([]OnePeletonWorkoutDetails, error)
}

type onepelotonclient struct {
	username string
	password string
}

type onepelotonconfig struct {
	UsernameOrEmail string `yaml:"username_or_email"`
	Password        string `yaml:"password"`
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
	}, nil
}

func (o onepelotonclient) Login(ctx context.Context) (OnePelotonSession, error) {
	creds := map[string]string{
		"username_or_email": o.username,
		"password":          o.password,
	}
	creds_json, err := json.Marshal(creds)
	if err != nil {
		return nil, errors.New("json marshalling error on creds")
	}

	resp, err := http.Post(_onePelotonAPIURL+"/"+_onePelotonLoginAPI, "application/json", bytes.NewBuffer(creds_json))
	if err != nil {
		return nil, errors.New("post error")
	}
	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	session_id := res["session_id"]
	user_id := res["user_id"]

	if session_id == nil || len(session_id) == 0 {
		return nil, errors.New("nil or empty session_id")
	}
	if user_id == nil || len(user_id) == 0 {
		return nil, errors.New("nil or empty user_id")
	}

	return &OnePelotonSession{
		session_id: session_id,
		user_id:    user_id,
	}, nil
}

func (o onepelotonclient) GetWorkoutIds(ctx context.Context, ops OnePelotonSession, amount int) ([]string, error) {

}

func (o onepelotonclient) GetWorkoutDetails(ctx context.Context, ops OnePelotonSession, workouts []string) ([]OnePeletonWorkoutDetails, error) {

}
