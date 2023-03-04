package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	config "github.com/DevopsGuyXD/Bizapp/Configs"
	util "github.com/DevopsGuyXD/Bizapp/Utils"
)

// ====================== Home Handler ======================
func HomeHandler(w http.ResponseWriter, r *http.Request){

	tmpl := util.InitHtmlFiles()
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

// ====================== Form handler ======================
func CreateApp(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")

	var completeForm config.FullForm

	r.ParseMultipartForm(10 << 20)

	//--------------------- FORM: Text Inputs ---------------------
	formDataRaw := map[string]string{
		"APPNAME": r.FormValue("APPNAME"),
		"FPTAG": r.FormValue("FPTAG"),
		"FPROOTALIASURI": r.FormValue("FPROOTALIASURI"),
		"FPLOGOURI": r.FormValue("FPLOGOURI"),
		"FPLOGOBACKGROUNDCOLOR": r.FormValue("FPLOGOBACKGROUNDCOLOR"),
		"FPLOGOFOREGROUNDCOLOR": r.FormValue("FPLOGOFOREGROUNDCOLOR"),
		"VERSION": r.FormValue("VERSION"),
		"IMAGE_NAME": r.FormValue("IMAGE_NAME"),
	}

	//--------------------- FORM: Image input ----------------------
	GetFormFile, _ , err := r.FormFile("UPLOADLOGOIMAGE"); util.CheckForNil(err)
	defer GetFormFile.Close()
	image_file, err := ioutil.ReadAll(GetFormFile); util.CheckForNil(err)

	formFile := map[string][]byte{
		"UPLOADLOGOIMAGE": image_file,
	}

	//--------- FORM: Converting input values to bytes -------------
	formDataToByte, err:= json.Marshal(formDataRaw); util.CheckForNil(err)
	formFileToByte, err:= json.Marshal(formFile); util.CheckForNil(err)

	//--------- FORM: Converting input values to json --------------
	json.Unmarshal(formDataToByte, &completeForm.FormData)
	json.Unmarshal(formFileToByte, &completeForm.UploadLogoImage)

	//--------- FORM: Converting complete "form-data" to bytes -----
	CompleteFormToBytes, err :=  json.Marshal(&completeForm); util.CheckForNil(err)

	//--------- FORM: Sending "form-data" to Queue -----------------
	queue_res := config.SendToQueue(CompleteFormToBytes)

	json.NewEncoder(w).Encode(string(queue_res))
}

// ====================== Health Check ======================
func HealthCheck(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")

	currentTime := time.Now()
	timeFormatted := currentTime.Format("15:04 UTC Monday")

	fmt.Fprintf(w,"Checking at %q...\n\n%v", timeFormatted,"Healthy ✔️") 
}