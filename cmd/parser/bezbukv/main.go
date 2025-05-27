package main

import (
	"github.com/Markard/wordka/config"
	"github.com/Markard/wordka/internal/repo"
	"github.com/Markard/wordka/pkg/postgres"
	"github.com/Markard/wordka/pkg/slogext"
	"github.com/PuerkitoBio/goquery"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Job struct {
	Id  int
	Url string
}

type Result struct {
	JobId int
	Words []string
}

type FiveLetterNounParser struct {
	pages   int
	workers int
	baseUrl string
	results chan<- *Result
}

func ParserWithPagination(pages int, workers int, baseUrl string, results chan<- *Result) *FiveLetterNounParser {
	return &FiveLetterNounParser{pages: pages, workers: workers, baseUrl: baseUrl, results: results}
}

func (p *FiveLetterNounParser) Parse() {
	jobs := make(chan *Job, p.pages)

	for w := 1; w <= p.workers; w++ {
		worker := NewWorker(w, jobs, p.results)
		go worker.Start()
	}

	for j := 1; j <= p.pages; j++ {
		job := &Job{
			Id:  j,
			Url: p.baseUrl + strconv.Itoa(j),
		}
		jobs <- job
	}
	close(jobs)
}

type Worker struct {
	Id      int
	Jobs    <-chan *Job
	Results chan<- *Result
}

func NewWorker(id int, jobs <-chan *Job, results chan<- *Result) *Worker {
	return &Worker{
		Id:      id,
		Jobs:    jobs,
		Results: results,
	}
}

func (w *Worker) Start() {
	for job := range w.Jobs {
		w.Results <- w.DoWork(job)
	}
}

func (w *Worker) DoWork(job *Job) *Result {
	logger := slog.Default()
	logger.Info("Fetching words", "url", job.Url)
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Get(job.Url)
	if err != nil {
		slogext.Fatal(logger, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != 200 {
		slogext.Fatal(logger, err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		slogext.Fatal(logger, err)
	}

	result := &Result{
		JobId: job.Id,
		Words: make([]string, 0),
	}
	doc.Find("div.view").Each(func(i int, s *goquery.Selection) {
		re := regexp.MustCompile(`[а-яА-ЯёЁ-]{5}`)
		lines := strings.Split(s.Text(), "\n")
		for _, line := range lines {
			matches := re.FindStringSubmatch(line)
			if len(matches) == 1 {
				word := strings.ToLower(matches[0])
				result.Words = append(result.Words, word)
			}
		}
	})

	return result
}

func main() {
	setup := config.MustLoad()
	logger := slogext.SetupLogger(setup.Env.AppEnv)
	pages := 36
	results := make(chan *Result, pages)
	parser := ParserWithPagination(
		pages,
		5,
		"https://bezbukv.ru/mask/%2A%2A%2A%2A%2A/noun?page=",
		results,
	)
	parser.Parse()

	db := postgres.New(setup.Env.PgDSN, logger)
	defer func() {
		err := db.Close()
		if err != nil {
			slogext.Error(logger, err)
		}
	}()
	gameRepo := repo.NewGameRepository(db)

	for r := 1; r <= pages; r++ {
		result := <-results
		err := gameRepo.SaveWords(result.Words)
		if err != nil {
			slogext.Error(logger, err)
		}
		logger.Info("Successfully fetched", "words", len(result.Words))
	}
	close(results)
}
