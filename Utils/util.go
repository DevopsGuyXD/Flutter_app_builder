package util

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joho/godotenv"
)

var tmpl *template.Template

//========================= Error Checking ==========================
func CheckForNil(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
//======================= END Error Checking=========================





//====================== Env file initiation ========================
func InitEnvFile(){
	err := godotenv.Load(".env"); CheckForNil(err)
}
//==================== END Env file initiation ======================





//============================ HTML =================================
func InitHtmlFiles() *template.Template{

	tmpl = template.Must(template.ParseGlob("./Assets/Static/*.html"))

	return tmpl
}
//========================== END HTML ===============================





//============================ AWS =========================
	//----------------- AWS initialization ------------------
	func InitAwsSession() (*session.Session){
		InitEnvFile()

		exec.Command("aws","configure","set","aws_access_key_id",os.Getenv("ACCESS_KEY_ID")).Output()
		exec.Command("aws","configure","set","aws_secret_access_key",os.Getenv("ACCESS_KEY_SECRET")).Output()
		exec.Command("aws","configure","set","region",os.Getenv("AWS_REGION")).Output()

		sess, err := session.NewSessionWithOptions(session.Options{
			Profile: "default",
			Config: aws.Config{
				Region: aws.String(os.Getenv("AWS_REGION")),
			},
		}); CheckForNil(err)
			
		return sess
	}

	//------------------------ SQS ---------------------------
	func SQS() (*sqs.SQS, *sqs.GetQueueUrlOutput){
		InitEnvFile()

		sess := InitAwsSession()
		queue_name := os.Getenv("QUEUE_NAME")

		sqsClient := sqs.New(sess)
		queue_url, err := sqsClient.GetQueueUrl(
			&sqs.GetQueueUrlInput{
				QueueName: &queue_name,
			},
		); CheckForNil(err)

		return sqsClient,
			   queue_url
	}

	//-------------------------- S3 ---------------------------
	func S3(folder_name string, reader *bytes.Reader) *s3manager.UploadOutput{
		InitEnvFile()

		sess := InitAwsSession()
		uploader := s3manager.NewUploader(sess)

		res, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET")),
			Key: aws.String(folder_name),
			Body: reader,
		}); CheckForNil(err)

		return res
	}
//=========================== END AWS ===========================

//======================== Github Login =========================
func GithubLogin(){
	_,err := exec.Command("git","config","--global","user.name",os.Getenv("GIT_USERNAME")).Output(); CheckForNil(err)
	_,err = exec.Command("git","config","--global","user.email",os.Getenv("GIT_EMAIL_ID")).Output(); CheckForNil(err)
	_,err = exec.Command("git","remote","set-url","origin",os.Getenv("GIT_REPO_PRIVATE_ENDPOINT")).Output(); CheckForNil(err)
	_,err = exec.Command("git", "pull").Output(); CheckForNil(err)
}
//====================== END Github Login =======================