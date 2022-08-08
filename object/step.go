package object

import (
	"fmt"
	"log"
	"time"

	"github.com/open-ct/openitem/util"
	"xorm.io/builder"
	"xorm.io/core"
)

type Step struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"created_time"`

	ProjectId   string             `json:"project_id"`
	Index       int                `json:"index"`
	Description string             `json:"description"`
	Requirement string             `json:"requirement"`
	Status      int                `json:"status"`
	Deadline    int64              `json:"deadline"`
	Timetable   []ProjectTimePoint `xorm:"mediumtext" json:"timetable"`
	Creator     string             `json:"creator"`
	Attachments []string           `xorm:"mediumtext" json:"attachments"` // uuid of files

	CreateAt  time.Time `xorm:"created" json:"create_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

type AddStepAttachment struct {
	StepId   string   `json:"step_id"`
	FilesIds []string `json:"files_ids"`
	Uploader string   `json:"uploader"`
}

type SetStepTimePointRequest struct {
	StepId     string `json:"step_id"`
	PointIndex int    `json:"point_index"`
	// index < 0 || index >= len  -> create a new time point
	Info TimePoint `json:"info"`
}

type TimePoint struct {
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Notice    string    `json:"notice"`
	Comment   string    `json:"comment"`
}

type DeleteStepTimePointRequest struct {
	StepId     string `json:"step_id"`
	PointIndex int    `json:"point_index"`
}

type StepDataStatistic struct {
	Total     float64 `json:"total"`
	PassRate  float64 `json:"pass_rate"`
	Pass      float64 `json:"pass"`
	Returned  float64 `json:"returned"`
	ToUpload  float64 `json:"to_upload"`
	ToAudit   float64 `json:"to_audit"`
	ToCorrect float64 `json:"to_correct"`
}

func getStep(owner string, name string) *Step {
	step := Step{Owner: owner, Name: name}
	existed, err := adapter.engine.Get(&step)
	if err != nil {
		panic(err)
	}

	if existed {
		return &step
	} else {
		return nil
	}
}

func GetStep(id string) *Step {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getStep(owner, name)
}

func UpdateStep(id string, step *Step) bool {
	owner, name := util.GetOwnerAndNameFromId(id)
	if getStep(owner, name) == nil {
		return false
	}

	_, err := adapter.engine.ID(core.PK{owner, name}).AllCols().Update(step)
	if err != nil {
		panic(err)
	}

	// return affected != 0
	return true
}

func AddStep(step *Step) error {
	_, err := adapter.engine.Insert(step)
	if err != nil {
		panic(err)
	}

	return nil
}

func CreateOneStep(req *Step) (string, error) {
	newStep := Step{
		Owner:       req.Owner,
		Name:        req.Name,
		CreatedTime: util.GetCurrentTime(),

		ProjectId:   req.ProjectId,
		Index:       req.Index,
		Description: req.Description,
		Requirement: req.Requirement,
		Deadline:    req.Deadline,
		Status:      0,
		Creator:     req.Creator,
	}

	var newTimeTable []ProjectTimePoint
	for _, timePoint := range req.Timetable {
		point := ProjectTimePoint{
			Title:     timePoint.Title,
			StartTime: timePoint.StartTime,
			EndTime:   timePoint.EndTime,
			Notice:    timePoint.Notice,
			Comment:   timePoint.Comment,
		}
		newTimeTable = append(newTimeTable, point)
	}
	newStep.Timetable = newTimeTable

	err := AddStep(&newStep)
	if err != nil {
		log.Printf("[File upload (create new step  failed)] %s\n", err.Error())
		return "", err
	}

	insertedStepId := fmt.Sprintf("%s/%s", newStep.Owner, newStep.Name)

	log.Printf("[Insert] %s\n", insertedStepId)
	return insertedStepId, nil
}

func GetStepInfo(sid string) (*Step, error) {
	var step Step

	owner, name := util.GetOwnerAndNameFromId(sid)

	_, err := adapter.engine.ID(core.PK{owner, name}).Get(&step)
	if err != nil {
		log.Printf("err: %s\n", err.Error())
		return nil, err
	}
	return &step, nil
}

func GetAllStepsInProject(pid string) (*[]Step, error) {
	var steps []Step

	err := adapter.engine.Where(builder.Eq{"project_id": pid}).Find(&steps)
	if err != nil {
		log.Printf("err: %s\n", err.Error())
		return nil, err
	}
	return &steps, nil
}

func UploadStepAttachments(req *AddStepAttachment) error {
	var updateStep Step

	updateStep.Attachments = req.FilesIds

	owner, name := util.GetOwnerAndNameFromId(req.StepId)
	_, err := adapter.engine.ID(core.PK{owner, name}).Cols("attachments").Update(&updateStep)
	if err != nil {
		log.Printf("add attachments for step: %s\n", err.Error())
		return err
	}
	return nil
}

func UpdateStepInfo(req *Step) error {
	var oldStep Step

	_, err := adapter.engine.ID(core.PK{req.Owner, req.Name}).Get(&oldStep)
	if err != nil {
		log.Printf("update step information: %s\n", err.Error())
		return err
	}

	_, err = adapter.engine.ID(core.PK{oldStep.Owner, oldStep.Name}).Update(req)
	if err != nil {
		log.Printf("update step information: %s\n", err.Error())
		return err
	}
	return nil
}

func SetStepStatus(req *Step) error {
	_, err := adapter.engine.ID(core.PK{req.Owner, req.Name}).Cols("status").Update(req)
	if err != nil {
		log.Printf("change step status: %s\n", err.Error())
		return err
	}
	return nil
}

func SetStepTimePoint(req *SetStepTimePointRequest) (*[]ProjectTimePoint, error) {
	// get step data:
	var step Step

	owner, name := util.GetOwnerAndNameFromId(req.StepId)
	_, err := adapter.engine.ID(core.PK{owner, name}).Get(&step)
	if err != nil {
		log.Printf("address the step error: %s\n", err.Error())
		return nil, err
	}
	newTimeTable := step.Timetable
	newTimePoint := ProjectTimePoint{
		Title:     req.Info.Title,
		StartTime: req.Info.StartTime,
		EndTime:   req.Info.EndTime,
		Notice:    req.Info.Notice,
		Comment:   req.Info.Comment,
	}
	if req.PointIndex < 0 || req.PointIndex >= len(step.Timetable) {
		newTimeTable = append(newTimeTable, newTimePoint)

		_, err = adapter.engine.ID(core.PK{owner, name}).Cols("timetable").Update(&Step{Timetable: newTimeTable})
		if err != nil {
			log.Printf("append step time point error: %s\n" + err.Error())
			return nil, err
		}
		return &newTimeTable, nil
	}
	newTimeTable[req.PointIndex] = newTimePoint

	_, err = adapter.engine.ID(core.PK{owner, name}).Cols("timetable").Update(&Step{Timetable: newTimeTable})
	if err != nil {
		log.Printf("append step time point error: %s\n" + err.Error())
		return nil, err
	}
	return &newTimeTable, nil
}

func DeleteStepTimePoint(req *DeleteStepTimePointRequest) error {
	var step Step

	owner, name := util.GetOwnerAndNameFromId(req.StepId)
	_, err := adapter.engine.ID(core.PK{owner, name}).Get(&step)
	if err != nil {
		log.Printf("address the step error: %s\n", err.Error())
		return err
	}
	if req.PointIndex >= len(step.Timetable) || req.PointIndex < 0 {
		return err
	}
	newTimeTable := step.Timetable
	// delete array element
	newTimeTable = append(newTimeTable[:req.PointIndex], newTimeTable[req.PointIndex+1:]...)
	_, err = adapter.engine.ID(core.PK{owner, name}).Cols("timetable").Update(&Step{Timetable: newTimeTable})
	if err != nil {
		log.Printf("delete time point error: %s\n", err.Error())
		return err
	}
	return nil
}

func GetStepDataStatistic(stepId string) (*StepDataStatistic, error) {
	var dataStatistic StepDataStatistic
	var submits []Submit

	err := adapter.engine.Where(builder.Eq{"step_id": stepId}).Find(&submits)
	if err != nil {
		log.Printf("find submits info err: " + err.Error())
		return nil, err
	}
	dataStatistic.Total = float64(len(submits))
	for _, submit := range submits {
		contentCount := len(submit.Contents)
		if contentCount == 0 {
			dataStatistic.ToUpload += 1
			continue
		}
		lastContentId := submit.Contents[contentCount-1].Uuid
		var lastAudit Audit
		_, err := adapter.engine.Where(builder.Eq{"submit_content": lastContentId}).Get(&lastAudit)
		if err != nil {
			dataStatistic.ToAudit += 1
			continue
		}
		if lastAudit.Result == "not pass" {
			dataStatistic.Returned += 1
		}
		if lastAudit.Result == "pass" {
			dataStatistic.Pass += 1
		}
		if lastAudit.Result == "need correct" {
			dataStatistic.ToCorrect += 1
		}
	}
	if dataStatistic.Total == 0 {
		dataStatistic.PassRate = 0
	} else {
		dataStatistic.PassRate = dataStatistic.Pass / dataStatistic.Total
	}
	return &dataStatistic, nil
}

func DeleteStep(stepId string) error {
	owner, name := util.GetOwnerAndNameFromId(stepId)

	_, err := adapter.engine.ID(core.PK{owner, name}).Delete(&Step{})
	if err != nil {
		log.Printf("delete step error: %s\n" + err.Error())
		return err
	}
	return nil
}
