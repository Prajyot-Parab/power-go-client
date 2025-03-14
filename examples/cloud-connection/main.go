package main

import (
	"fmt"
	"log"
	"time"

	v "github.com/IBM-Cloud/power-go-client/clients/instance"
	ps "github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

const (
	JOBCOMPLETED = "completed"
	JOBFAILED    = "failed"
)

func main() {

	//session Inputs
	token := " < IAM TOKEN > "
	region := " < REGION > "
	accountID := " < ACCOUNT ID > "

	// Image inputs
	name := " < NAME OF THE CONNECTION > "
	piID := " < POWER INSTANCE ID > "
	var speed int64 = 5000

	session, err := ps.New(token, region, true, 50000000, accountID, region)
	if err != nil {
		log.Fatal(err)
	}
	ccClient := v.NewIBMPICloudConnectionClient(session, piID)
	if err != nil {
		log.Fatal(err)
	}
	jobClient := v.NewIBMPIJobClient(session, piID)
	if err != nil {
		log.Fatal(err)
	}

	body := &models.CloudConnectionCreate{
		Name:  &name,
		Speed: &speed,
	}
	createRespOk, createRespAccepted, err := ccClient.Create(body, piID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("***************[1]****************** %+v\n", createRespOk)
	log.Printf("***************[1]****************** %+v\n\n", createRespAccepted)

	var ccId, jobId string
	if createRespOk != nil {
		ccId = *createRespOk.CloudConnectionID
	} else {
		ccId = *createRespAccepted.CloudConnectionID
		jobId = *createRespAccepted.JobRef.ID
		waitForJobState(jobClient, jobId, piID, 2000)
		if err != nil {
			log.Fatal(err)
		}
	}

	getResp, err := ccClient.Get(ccId, piID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("***************[2]****************** %+v \n\n", *getResp)

	getAllResp, err := ccClient.GetAll(piID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("***************[3]****************** %+v \n\n", *getAllResp)

	delok, delAccepted, err := ccClient.Delete(piID, ccId)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("***************[4]****************** %+v\n", delok)
	log.Printf("***************[4]****************** %+v\n\n", delAccepted)

	if delAccepted != nil {
		jobId = *delAccepted.ID
		waitForJobState(jobClient, jobId, piID, 2000)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func waitForJobState(jobClient *v.IBMPIJobClient, jobId, cloudinstanceid string, interval time.Duration) error {
	var status string

	for status != JOBCOMPLETED && status != JOBFAILED {
		job, err := jobClient.Get(jobId, cloudinstanceid)
		if err != nil {
			return err
		}
		if job == nil || job.Status == nil {
			return fmt.Errorf("cannot find job status for job id %s with cloud instance %s", jobId, cloudinstanceid)
		}
		time.Sleep(interval)
		status = *job.Status.State
	}
	return nil
}
