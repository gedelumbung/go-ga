package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/analytics/v3"
)

type GaAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GaProperty struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GaProfile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	ActiveUsers [][]string `json:"active_users"`
	NewUsers    [][]string `json:"new_users"`
	Meta        Meta       `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Meta struct {
	Profiles   interface{} `json:"profiles"`
	Accounts   interface{} `json:"accounts"`
	Properties interface{} `json:"properties"`
}

func main() {
	http.HandleFunc("/active-users", getCurrentActiveUsers)
	http.ListenAndServe(":9090", nil)
}

func getCurrentActiveUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		startDate, eSd := r.URL.Query()["start_date"]

		if !eSd || len(startDate[0]) < 1 {
			errSdResponse, _ := json.Marshal(ErrorResponse{
				Message: "Start Date is required",
			})

			w.WriteHeader(200)
			w.Write(errSdResponse)
			return
		}
		startDateString := strings.Join(startDate, "")

		endDate, eEd := r.URL.Query()["end_date"]

		if !eEd || len(endDate[0]) < 1 {
			errEdResponse, _ := json.Marshal(ErrorResponse{
				Message: "End Date is required",
			})

			w.WriteHeader(200)
			w.Write(errEdResponse)
			return
		}
		endDateString := strings.Join(startDate, "")

		key, _ := ioutil.ReadFile("credential.json")

		jwtConf, err := google.JWTConfigFromJSON(
			key,
			analytics.AnalyticsReadonlyScope,
		)
		p(err)

		httpClient := jwtConf.Client(oauth2.NoContext)
		svc, err := analytics.New(httpClient)
		p(err)

		accountResponse, err := svc.Management.Accounts.List().Do()
		p(err)

		var accountID string
		var gaAccounts []GaAccount

		for i, acc := range accountResponse.Items {

			if i == 0 {
				accountID = acc.Id
			}

			gaAccounts = append(gaAccounts, GaAccount{
				ID:   acc.Id,
				Name: acc.Name,
			})
		}

		webProps, err := svc.Management.Webproperties.List(accountID).Do()
		p(err)

		var wpID string
		var gaProperties []GaProperty

		for i, wp := range webProps.Items {

			if i == 0 {
				wpID = wp.Id
			}

			gaProperties = append(gaProperties, GaProperty{
				ID:   wp.Id,
				Name: wp.Name,
			})
		}

		profiles, err := svc.Management.Profiles.List(accountID, wpID).Do()
		p(err)

		var viewID string
		var gaProfiles []GaProfile

		for i, p := range profiles.Items {

			if i == 0 {
				viewID = "ga:" + p.Id
			}

			gaProfiles = append(gaProfiles, GaProfile{
				ID:   p.Id,
				Name: p.Name,
			})
		}

		au, err := svc.Data.Realtime.Get(viewID, "rt:activeUsers").Do()
		p(err)

		nu, err := svc.Data.Ga.Get(viewID, startDateString, endDateString, "ga:newUsers").Do()
		p(err)

		response, _ := json.Marshal(Response{
			ActiveUsers: au.Rows,
			NewUsers:    nu.Rows,
			Meta: Meta{
				Profiles:   gaProfiles,
				Accounts:   gaAccounts,
				Properties: gaProperties,
			},
		})

		w.WriteHeader(200)
		w.Write(response)
		return
	}
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}
