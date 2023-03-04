package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	util "github.com/DevopsGuyXD/Bizapp/Utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//========================= Variables =========================
var wg sync.WaitGroup
var sqsClient, sqsQueueUrl = util.SQS()
//======================= END Variables =======================





//=================================== Form Structs =====================================
//---------------------------------- Complete form -------------------------------------
type FullForm struct{
	FormData *FormData `json:"formData,omitempty"`
	UploadLogoImage *ImageFile `json:"uploadlogoimage,omitempty"`
}

//----------------------------------- Text Inputs --------------------------------------
type FormData struct {
	Id                    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AppName 			  string             `json:"appname,omitempty"`
	FPTag                 string	         `json:"fptag,omitempty"`
	FPRootAliasURI        string	         `json:"fprootaliasuri,omitempty"`
	FPLogoURI             string	         `json:"fplogouri,omitempty"`
	FPLogoBackgroundColor string	         `json:"fplogobackgroundcolor,omitempty"`
	FPLogoForegroundColor string	         `json:"fplogoforegroundcolor,omitempty"`	
	Version               string	         `json:"version,omitempty"`
}

//----------------------------------- Image Inputs --------------------------------------
type ImageFile struct{
	LogoImage []byte `json:"uploadLogoImage,omitempty"`
}
//=============================== End Form Structs ======================================





//=================================== Form-Methods ======================================
//------------------------------- Generating .env file ----------------------------------
func (fd *FullForm) GenerateEnvFile(messagedetails FormData){
	wg.Done()
	
	folder_name := fmt.Sprintf("%v/V%v/.env" ,messagedetails.FPTag, messagedetails.Version)

	envFile := map[string]string{
		"FPTAG": strings.ToUpper(messagedetails.FPTag),
		"APPNAME": strings.Title(strings.ToLower(messagedetails.AppName)),
		"FPROOTALIASURI": messagedetails.FPRootAliasURI,
		"FPLOGOURI": messagedetails.FPLogoURI,
		"FPLOGOBACKGROUNDCOLOR": messagedetails.FPLogoBackgroundColor,
		"FPLOGOFOREGROUNDCOLOR": messagedetails.FPLogoForegroundColor,
		"VERSION": messagedetails.Version,
		"IMAGE_NAME": "icon.png",
	}

	envFileInBytes, err := json.Marshal(envFile); util.CheckForNil(err)
	replacer1 := strings.NewReplacer(`{`,"",`}`,"",`,`,"\n",`:`,"=",`"`,"")
	replacer2 := strings.NewReplacer(`s=`,"s:")
	envFileModified1 := replacer1.Replace(string(envFileInBytes))
	envFileModified2 := replacer2.Replace(string(envFileModified1))
	

	fmt.Printf("✔️  %v: %v\n",messagedetails.FPTag,"Preparing .envfile...")

	reader := bytes.NewReader([]byte(envFileModified2))

	util.S3(folder_name, reader)
}

//------------------------------- Generating pubspec.yaml file -----------------------------
func (fd *FullForm) GeneratePubSpec(messagedetails FormData){
	wg.Done()

	message_version := messagedetails.Version
	folder_name := fmt.Sprintf("%v/V%v/pubspec.yaml" ,messagedetails.FPTag, message_version)

	pubspecfile := util.YamlGenerator(message_version)
	fmt.Printf("✔️  %v: %v\n",messagedetails.FPTag,"Preparing pubspecfile...")
	
	reader := bytes.NewReader(pubspecfile)

	util.S3(folder_name, reader)
}

//------------------------------- Generating name and version file --------------------------
func (fd *FullForm) NameAndVersion(messagedetails FormData){
	wg.Done()

	name_version := messagedetails.FPTag+"\nV"+messagedetails.Version
	folder_name := fmt.Sprintf("%v/V%v/name_version.txt" ,messagedetails.FPTag, messagedetails.Version)
	
	replacer := strings.NewReplacer(`"`,"")
	name_version = replacer.Replace(name_version)
	fmt.Printf("✔️  %v: %v\n",messagedetails.FPTag,"Generating metadata...")
	
	reader := bytes.NewReader([]byte(name_version))

	util.S3(folder_name, reader)
}

//--------------------------------- Getting image for further use ----------------------------
func (fd *FullForm) Uploadimage(messagedetails FullForm){
	wg.Done()

	folder_name := fmt.Sprintf("%v/V%v/icon.png" ,messagedetails.FormData.FPTag, messagedetails.FormData.Version)

	fmt.Printf("✔️  %v: %v\n",messagedetails.FormData.FPTag,"Uploading image...")
	
	reader := bytes.NewReader(messagedetails.UploadLogoImage.LogoImage)

	util.S3(folder_name, reader)
}
//========================================= End Form-Methods ==================================





//======================================== QUEUE ===========================================
//------------------------------------ Send to Queue ---------------------------------------
func SendToQueue(message []byte) string{

	message_to_string := string(message)
	_, err := sqsClient.SendMessage(&sqs.SendMessageInput{QueueUrl: sqsQueueUrl.QueueUrl, MessageBody: aws.String(message_to_string)}); util.CheckForNil(err)

	return "Request made successfully"
}

//------------------------------------ Poll from Queue -------------------------------------
func PollFromQueue(){

	var messagedetails FullForm

	attr, err := sqsClient.GetQueueAttributes(&sqs.GetQueueAttributesInput{QueueUrl: sqsQueueUrl.QueueUrl, AttributeNames: aws.StringSlice([]string{"All"})}); util.CheckForNil(err)
	message_in_queue_count, err := strconv.Atoi(*attr.Attributes["ApproximateNumberOfMessages"]); util.CheckForNil(err)

	if message_in_queue_count > 0{

		max_message_count, err := strconv.ParseInt(os.Getenv("MAX_MESSAGE_COUNT"),10, 64); util.CheckForNil(err)
		message, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{QueueUrl: sqsQueueUrl.QueueUrl, MaxNumberOfMessages: aws.Int64(max_message_count)}); util.CheckForNil(err)

		err = json.Unmarshal([]byte(*message.Messages[0].Body), &messagedetails); util.CheckForNil(err)

		fmt.Printf("=======================================\n")
		fmt.Printf("VERSION: V%v\n\n",messagedetails.FormData.Version)
		
		wg.Add(4)

		go messagedetails.GenerateEnvFile(*messagedetails.FormData)
		go messagedetails.GeneratePubSpec(*messagedetails.FormData)
		go messagedetails.NameAndVersion(*messagedetails.FormData)
		go messagedetails.Uploadimage(messagedetails)

		wg.Wait()

		defer TriggerGithubAction(messagedetails)
		defer DeleteFromQueue(message)
	}else{
		log.Println("No requests to process at the moment")
	}
}

//------------------------------------- Delete from Queue -------------------------------------
func DeleteFromQueue(message *sqs.ReceiveMessageOutput){
	_, err := sqsClient.DeleteMessage(&sqs.DeleteMessageInput{QueueUrl: sqsQueueUrl.QueueUrl, ReceiptHandle: message.Messages[0].ReceiptHandle}); util.CheckForNil(err)
	fmt.Printf("✔️  %v\n","Request processed successfully")
}
//========================================= END QUEUE ==========================================





//==================================== Trigger app builder =====================================
func TriggerGithubAction(messagedetails FullForm){

	util.InitEnvFile()

	message := fmt.Sprintf("%v has been created", messagedetails.FormData.FPTag)
	fileName := "./boost_msme_app_builder/"+messagedetails.FormData.FPTag+"/"+"V"+messagedetails.FormData.Version
	commit_Id := messagedetails.FormData.FPTag+":V"+messagedetails.FormData.Version

	if _, err := os.Stat("./boost_msme_app_builder/"+messagedetails.FormData.FPTag); os.IsNotExist(err) {
		err := os.Mkdir("./boost_msme_app_builder/"+messagedetails.FormData.FPTag, os.ModePerm); util.CheckForNil(err)
	}

	file, err := os.Create(fileName); util.CheckForNil(err)
	defer file.Close()

	_, err = io.WriteString(file, message); util.CheckForNil(err)

	os.Chdir("./boost_msme_app_builder")

	_,err = exec.Command("git","config","--global","user.name",os.Getenv("GIT_USERNAME")).Output(); util.CheckForNil(err)
	_,err = exec.Command("git","config","--global","user.email",os.Getenv("GIT_EMAIL_ID")).Output(); util.CheckForNil(err)
	_,err = exec.Command("git","remote","set-url","origin",os.Getenv("GIT_REPO_PRIVATE_ENDPOINT")).Output(); util.CheckForNil(err)
	_,err = exec.Command("git", "pull").Output(); util.CheckForNil(err)
	_,err = exec.Command("git", "add", ".").Output(); util.CheckForNil(err)
	_,err = exec.Command("git", "commit", "-m", commit_Id).Output(); util.CheckForNil(err)
	_,err = exec.Command("git", "push").Output(); util.CheckForNil(err)

	fmt.Printf("🔥 %v\n","Flutter build triggered")
	fmt.Printf("=======================================\n\n")
	os.Chdir("..")
}
//============================== END Trigger app builder =======================================