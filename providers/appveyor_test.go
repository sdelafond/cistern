package providers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/nbedos/cistern/utils"
)

func TestParseAppVeyorURL(t *testing.T) {
	u := "https://ci.appveyor.com/project/nbedos/cistern/builds/29070120"
	owner, repo, id, err := parseAppVeyorURL(u)
	if err != nil {
		t.Fatal(err)
	}
	if owner != "nbedos" || repo != "cistern" || id != 29070120 {
		t.Fail()
	}
}

func TestAppVeyorJob_ToCacheJob(t *testing.T) {
	j := appVeyorJob{
		ID:           "id",
		Name:         "name",
		AllowFailure: true,
		Status:       "success",
		CreatedAt:    "2019-11-23T12:24:26.9181871+00:00",
		StartedAt:    "2019-11-23T12:24:31.8145735+00:00",
		FinishedAt:   "2019-11-23T12:24:34.5646724+00:00",
	}

	expectedJob := Step{
		ID:        "id",
		Type:      StepJob,
		State:     "passed",
		Name:      "name",
		CreatedAt: time.Date(2019, 11, 23, 12, 24, 26, 918187100, time.UTC),
		StartedAt: utils.NullTime{
			Valid: true,
			Time:  time.Date(2019, 11, 23, 12, 24, 31, 814573500, time.UTC),
		},
		FinishedAt: utils.NullTime{
			Valid: true,
			Time:  time.Date(2019, 11, 23, 12, 24, 34, 564672400, time.UTC),
		},
		Duration: utils.NullDuration{
			Valid:    true,
			Duration: 2750098900 * time.Nanosecond,
		},
		AllowFailure: true,
		WebURL: utils.NullString{
			String: "buildURL/job/id",
			Valid:  true,
		},
	}

	job, err := j.toCacheStep(expectedJob.ID, "buildURL")
	if err != nil {
		t.Fatal(err)
	}

	if diff := expectedJob.Diff(job); len(diff) > 0 {
		t.Fatal(diff)
	}
}

func TestAppVeyorBuild_ToCacheBuild(t *testing.T) {
	b := appVeyorBuild{
		ID:          42,
		Jobs:        nil,
		Number:      42,
		Version:     "1.0.42",
		Message:     "message",
		Branch:      "feature/appveyor",
		IsTag:       false,
		Sha:         "fd4c4ae5a4005e38c66566e2480087072620e9de",
		Author:      "nbedos",
		CommittedAt: "2019-11-23T12:23:09+00:00",
		Status:      "failed",
		CreatedAt:   "2019-11-23T12:24:25.5900258+00:00",
		StartedAt:   "2019-11-23T12:24:31.8145735+00:00",
		FinishedAt:  "2019-11-23T12:24:34.5646724+00:00",
		UpdatedAt:   "2019-11-23T12:24:34.5646724+00:00",
	}

	expectedBuild := Pipeline{
		Number: "42",
		GitReference: GitReference{
			SHA:   "fd4c4ae5a4005e38c66566e2480087072620e9de",
			Ref:   "feature/appveyor",
			IsTag: false,
		},
		Step: Step{
			ID:    "42",
			State: "failed",
			CreatedAt: time.Date(2019, 11, 23, 12, 24, 25, 590025800, time.UTC),
			StartedAt: utils.NullTime{
				Valid: true,
				Time:  time.Date(2019, 11, 23, 12, 24, 31, 814573500, time.UTC),
			},
			FinishedAt: utils.NullTime{
				Valid: true,
				Time:  time.Date(2019, 11, 23, 12, 24, 34, 564672400, time.UTC),
			},
			UpdatedAt: time.Date(2019, 11, 23, 12, 24, 34, 564672400, time.UTC),
			Duration: utils.NullDuration{
				Valid:    true,
				Duration: 2750098900 * time.Nanosecond,
			},
			WebURL: utils.NullString{
				String: "https://ci.appveyor.com/project/owner/repo/builds/42",
				Valid:  true,
			},
		},
	}

	t.Run("git ref is branch", func(t *testing.T) {
		build, err := b.toCachePipeline("owner", "repo")
		if err != nil {
			t.Fatal(err)
		}

		if diff := expectedBuild.Diff(build); len(diff) > 0 {
			t.Fatal(diff)
		}
	})

	t.Run("git ref is tag", func(t *testing.T) {
		b := b
		b.IsTag = true
		b.Tag = "0.1.0"
		build, err := b.toCachePipeline("owner", "repo")
		if err != nil {
			t.Fatal(err)
		}

		expectedBuild := expectedBuild
		expectedBuild.IsTag = true
		expectedBuild.Ref = "0.1.0"
		if diff := expectedBuild.Diff(build); len(diff) > 0 {
			t.Fatal(diff)
		}
	})

}

func TestAppVeyorClient_BuildFromURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := ""
		switch {
		case r.Method == "GET" && r.URL.Path == "/api/projects/nbedos/cistern/history":
			filename = "appveyor/appveyor_history_29070120.json"
		case r.Method == "GET" && r.URL.Path == "/api/projects/nbedos/cistern/build/1.0.22":
			filename = "appveyor/appveyor_build_1_0_22.json"
		default:
			w.WriteHeader(404)
			return
		}

		bs, err := ioutil.ReadFile(fmt.Sprintf("test_data/%s", filename))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := fmt.Fprint(w, string(bs)); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	tsu, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	tsu.Path += "/api"
	tsu.RawPath += "/api"

	client := AppVeyorClient{
		url:         *tsu,
		client:      &http.Client{Timeout: 10 * time.Second},
		rateLimiter: time.Tick(time.Millisecond),
		token:       "token",
		provider: Provider{
			ID:   "id",
			Name: "name",
		},
	}

	buildURL := "https://ci.appveyor.com/project/nbedos/cistern/builds/29070120"
	build, err := client.BuildFromURL(context.Background(), buildURL)
	if err != nil {
		t.Fatal(err)
	}
	if build.ID != "29070120" {
		t.Fatal()
	}
}

func TestAppVeyorClient_Log(t *testing.T) {
	expectedLog := "log\n"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/api/buildjobs/jobId/log":
			if _, err := fmt.Fprint(w, expectedLog); err != nil {
				t.Fatal(err)
			}
		default:
			w.WriteHeader(404)
			return
		}
	}))
	defer ts.Close()

	tsu, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	tsu.Path += "/api"
	tsu.RawPath += "/api"

	client := AppVeyorClient{
		url:         *tsu,
		client:      &http.Client{Timeout: 10 * time.Second},
		rateLimiter: time.Tick(time.Millisecond),
		token:       "token",
		provider: Provider{
			ID:   "id",
			Name: "name",
		},
	}

	job := Step{
		ID:   "jobId",
		Type: StepJob,
	}
	log, err := client.Log(context.Background(), job)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(log, expectedLog); len(diff) > 0 {
		t.Fatal(diff)
	}
}
