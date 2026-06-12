package apps

import (
	"GU/refs"
	"GU/utils"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

var Ns = NewGoCronScheduler(time.FixedZone("Asia/Tokyo", 9*60*60))

// Scheduler init は環境変数を読み込む
/*func init() {
	envLoader := &utils.DotenvLoader{}
	envLoader.LoadEnv("channel.env")

	channelId = map[string]string{
		"a": os.Getenv("TEAM_A"),
		"b": os.Getenv("TEAM_B"),
		"c": os.Getenv("TEAM_C"),
		"d": os.Getenv("TEAM_D"),
		"e": os.Getenv("TEAM_E"),
	}
}
*/
// Scheduler スケジューラを管理するインターフェース
type Scheduler interface {
	RegisterJob(cronExpr string, jobFunc interface{}, params ...interface{})
	Start()
	Jobs() []gocron.Job
}

// GoCronScheduler gocronを使用したスケジューラの実装
type GoCronScheduler struct {
	scheduler gocron.Scheduler
	err       error
}

// NewGoCronScheduler GoCronSchedulerのインスタンスを作成
func NewGoCronScheduler(location *time.Location) *GoCronScheduler {
	s, err := gocron.NewScheduler(gocron.WithLocation(location))
	return &GoCronScheduler{
		scheduler: s,
		err:       err,
	}
}

// RegisterJob ジョブをスケジューラに登録
func (s *GoCronScheduler) RegisterJob(cronTime time.Time, jobFunc interface{}, YURUBO refs.JobData, id string) (string, int) {
	job, err := s.scheduler.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(cronTime)),
		gocron.NewTask(jobFunc, YURUBO),
		gocron.WithTags(func() string {
			if id != "" {
				return id
			}
			return getID()
		}()),
	)
	if err != nil {
		utils.Log(err, "", "RegisterJob")
		return "", 1
	}
	return job.Tags()[0], 0
}

// Start スケジューラを開始
func (s *GoCronScheduler) Start() {
	s.scheduler.Start()
}

// Jobs 登録されているジョブを返す
func (s *GoCronScheduler) Jobs() []gocron.Job {
	return s.scheduler.Jobs()
}

// RemoveJob ジョブをスケジューラから削除
func (s *GoCronScheduler) RemoveJob(jobID string) (uint8, string, []refs.JobData) {
	jobList := s.searchJobFromInstance(jobID)
	var res = ""
	if len(jobList) == 0 {
		utils.Log(nil, fmt.Sprintf("Not Found jobID In The Instance : %s", jobID), "RemoveJob")
		return 16, "", []refs.JobData{}
	}
	var status uint8 = 0
	for _, job := range jobList {
		err := s.scheduler.RemoveJob(job.ID())
		if err != nil {
			utils.Log(err, "", "RemoveJob")
			run, err := job.NextRun()
			if err != nil {
				utils.Log(err, "", "RemoveJob")
				continue
			}
			res += job.Name() + "\n" + run.Format(time.DateTime) + "\n"
			continue
		}
		status++
	}
	return status, res, utils.JSONFM.SearchJobFromJSON(jobID)
}

func (s *GoCronScheduler) InitializeSchedule() {
	for _, jobData := range utils.JobDataSlice {
		date, _ := time.Parse(time.DateTime, jobData.Cron)
		date = date.Add(time.Minute * time.Duration(jobData.Gap))
		s.RegisterJob(date, utils.SendYURUBOItem, refs.JobData{
			Id:     jobData.Id,
			Title:  jobData.Title,
			Rank:   jobData.Rank,
			Number: jobData.Number,
			Cron:   jobData.Cron,
			Role:   jobData.Role,
			Gap:    jobData.Gap,
			Party:  jobData.Party,
		}, "")
	}
	s.Start()
}

func (s *GoCronScheduler) searchJobFromInstance(jobID string) []gocron.Job {
	var jobList []gocron.Job
	for _, j := range Ns.Jobs() {
		if j.Tags()[0] == jobID {
			jobList = append(jobList, j)
		}
	}
	return jobList
}

func getID() string {
	utils.IDChannel <- ""
	return <-utils.IDChannel
}
