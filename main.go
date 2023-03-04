package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	config "github.com/DevopsGuyXD/Bizapp/Configs"
	router "github.com/DevopsGuyXD/Bizapp/Routers"
	util "github.com/DevopsGuyXD/Bizapp/Utils"
)

var wg sync.WaitGroup

// ====================== Main ======================
func main() {
	util.InitEnvFile()

	server := router.RouterCollection()
	fmt.Printf("%v \n\n","Server listening on port 8000")

	err := http.ListenAndServe(os.Getenv("PORT"), server); util.CheckForNil(err)
}

// =========== Process messages from queue ===========
func ProcessRequest(){
	for{
		time.Sleep(60 * time.Second)
		config.PollFromQueue()
	}
}

// ================= Get GitHub Repo =================
func GetGitRepo(){
	util.InitEnvFile()

	if _, err := os.Stat("boost_msme_app_builder"); os.IsNotExist(err){
		_, err := exec.Command("git", "clone", "-b", "build_requirements", os.Getenv("GIT_REPO_PRIVATE_ENDPOINT")).Output(); util.CheckForNil(err)
	}
}

// ====================== Init =======================
func init(){
	
	wg.Add(3)

	go main()
	go ProcessRequest()
	go GetGitRepo()

	wg.Wait()
}