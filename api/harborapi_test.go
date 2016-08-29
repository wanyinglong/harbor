//These APIs provide services for manipulating Harbor project.

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"path/filepath"
	"runtime"

	"github.com/vmware/harbor/tests/apitests/apilib"
	"github.com/vmware/harbor/dao"
	"github.com/vmware/harbor/models"

	"github.com/astaxie/beego"
	"github.com/dghubble/sling"

	//for test env prepare
	_ "github.com/vmware/harbor/auth/db"
	_ "github.com/vmware/harbor/auth/ldap"
)

type api struct {
	basePath string
}

func newHarborAPI() *api {
	return &api{
		basePath: "",
	}
}

func newHarborAPIWithBasePath(basePath string) *api {
	return &api{
		basePath: basePath,
	}
}

type usrInfo struct {
	Name   string
	Passwd string
}

func init() {
	dao.InitDB()
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.TestBeegoInit(apppath)

	beego.Router("/api/search/", &SearchAPI{})
	beego.Router("/api/projects/", &ProjectAPI{}, "get:List;post:Post")
	beego.Router("/api/users/:id([0-9]+)/password", &UserAPI{}, "put:ChangePassword")

	_ = updateInitPassword(1, "Harbor12345")

}

//Search for projects and repositories
//Implementation Notes
//The Search endpoint returns information about the projects and repositories
//offered at public status or related to the current logged in user.
//The response includes the project and repository list in a proper display order.
//@param q Search parameter for project and repository name.
//@return []Search
//func (a api) SearchGet (q string) (apilib.Search, error) {
func (a api) SearchGet(q string) (apilib.Search, error) {

	_sling := sling.New().Get(a.basePath)

	// create path and map variables
	path := "/api/search"
	_sling = _sling.Path(path)

	type QueryParams struct {
		Query string `url:"q,omitempty"`
	}

	_sling = _sling.QueryStruct(&QueryParams{Query: q})

	accepts := []string{"application/json", "text/plain"}
	for key := range accepts {
		_sling = _sling.Set("Accept", accepts[key])
		break // only use the first Accept
	}

	req, err := _sling.Request()

	w := httptest.NewRecorder()

	beego.BeeApp.Handlers.ServeHTTP(w, req)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		// handle error
	}

	var successPayload = new(apilib.Search)
	err = json.Unmarshal(body, &successPayload)
	return *successPayload, err
}

//Create a new project.
//Implementation Notes
//This endpoint is for user to create a new project.
//@param project New created project.
//@return void
//func (a api) ProjectsPost (prjUsr usrInfo, project apilib.Project) (int, error) {
func (a api) ProjectsPost(prjUsr usrInfo, project apilib.Project) (int, error) {

	_sling := sling.New().Post(a.basePath)

	// create path and map variables
	path := "/api/projects/"

	_sling = _sling.Path(path)

	// accept header
	accepts := []string{"application/json", "text/plain"}
	for key := range accepts {
		_sling = _sling.Set("Accept", accepts[key])
		break // only use the first Accept
	}

	// body params
	_sling = _sling.BodyJSON(project)
	//fmt.Printf("project post req: %+v\n", _sling)
	req, err := _sling.Request()
	req.SetBasicAuth(prjUsr.Name, prjUsr.Passwd)

	w := httptest.NewRecorder()

	beego.BeeApp.Handlers.ServeHTTP(w, req)

	return w.Code, err
}

//Change password
//Implementation Notes
//Change the password on a user that already exists.
//@param userID user ID
//@param password user old and new password
//@return error
//func (a api) UsersUserIDPasswordPut (user usrInfo, userID int32, password apilib.Password) int {
func (a api) UsersUserIDPasswordPut(user usrInfo, userID int32, password apilib.Password) int {

	_sling := sling.New().Put(a.basePath)

	// create path and map variables
	path := "/api/users/" + fmt.Sprintf("%d", userID) + "/password"
	fmt.Printf("change passwd path: %s\n", path)
	fmt.Printf("password %+v\n", password)
	_sling = _sling.Path(path)

	// accept header
	accepts := []string{"application/json", "text/plain"}
	for key := range accepts {
		_sling = _sling.Set("Accept", accepts[key])
		break // only use the first Accept
	}

	// body params
	_sling = _sling.BodyJSON(password)
	fmt.Printf("project post req: %+v\n", _sling)
	req, err := _sling.Request()
	req.SetBasicAuth(user.Name, user.Passwd)
	fmt.Printf("project post req: %+v\n", req)
	if err != nil {
		// handle error
	}
	w := httptest.NewRecorder()

	beego.BeeApp.Handlers.ServeHTTP(w, req)

	return w.Code

}

////Delete a repository or a tag in a repository.
////Delete a repository or a tag in a repository.
////This endpoint let user delete repositories and tags with repo name and tag.\n
////@param repoName The name of repository which will be deleted.
////@param tag Tag of a repository.
////@return void
////func (a HarborAPI) RepositoriesDelete(prjUsr UsrInfo, repoName string, tag string) (int, error) {
//func (a HarborAPI) RepositoriesDelete(prjUsr UsrInfo, repoName string, tag string) (int, error) {
//	_sling := sling.New().Delete(a.basePath)

//	// create path and map variables
//	path := "/api/repositories"

//	_sling = _sling.Path(path)

//	type QueryParams struct {
//		RepoName string `url:"repo_name,omitempty"`
//		Tag      string `url:"tag,omitempty"`
//	}

//	_sling = _sling.QueryStruct(&QueryParams{RepoName: repoName, Tag: tag})
//	// accept header
//	accepts := []string{"application/json", "text/plain"}
//	for key := range accepts {
//		_sling = _sling.Set("Accept", accepts[key])
//		break // only use the first Accept
//	}

//	req, err := _sling.Request()
//	req.SetBasicAuth(prjUsr.Name, prjUsr.Passwd)
//	//fmt.Printf("request %+v", req)

//	client := &http.Client{}
//	httpResponse, err := client.Do(req)
//	defer httpResponse.Body.Close()

//	if err != nil {
//		// handle error
//	}
//	return httpResponse.StatusCode, err
//}

//Return projects created by Harbor
//func (a HarborApi) ProjectsGet (projectName string, isPublic int32) ([]Project, error) {
//    }

//Check if the project name user provided already exists.
//func (a HarborApi) ProjectsHead (projectName string) (error) {
//}

//Get access logs accompany with a relevant project.
//func (a HarborApi) ProjectsProjectIdLogsFilterPost (projectId int32, accessLog AccessLog) ([]AccessLog, error) {
//}

//Return a project&#39;s relevant role members.
//func (a HarborApi) ProjectsProjectIdMembersGet (projectId int32) ([]Role, error) {
//}

//Add project role member accompany with relevant project and user.
//func (a HarborApi) ProjectsProjectIdMembersPost (projectId int32, roles RoleParam) (error) {
//}

//Delete project role members accompany with relevant project and user.
//func (a HarborApi) ProjectsProjectIdMembersUserIdDelete (projectId int32, userId int32) (error) {
//}

//Return role members accompany with relevant project and user.
//func (a HarborApi) ProjectsProjectIdMembersUserIdGet (projectId int32, userId int32) ([]Role, error) {
//}

//Update project role members accompany with relevant project and user.
//func (a HarborApi) ProjectsProjectIdMembersUserIdPut (projectId int32, userId int32, roles RoleParam) (error) {
//}

//Update properties for a selected project.
//func (a HarborApi) ProjectsProjectIdPut (projectId int32, project Project) (error) {
//}

//Get repositories accompany with relevant project and repo name.
//func (a HarborApi) RepositoriesGet (projectId int32, q string) ([]Repository, error) {
//}

//Get manifests of a relevant repository.
//func (a HarborApi) RepositoriesManifestGet (repoName string, tag string) (error) {
//}

//Get tags of a relevant repository.
//func (a HarborApi) RepositoriesTagsGet (repoName string) (error) {
//}

//Get registered users of Harbor.
//func (a HarborApi) UsersGet (userName string) ([]User, error) {
//}

//Creates a new user account.
//func (a HarborApi) UsersPost (user User) (error) {
//}

//Mark a registered user as be removed.
//func (a HarborApi) UsersUserIdDelete (userId int32) (error) {
//}

//Update a registered user to change to be an administrator of Harbor.
//func (a HarborApi) UsersUserIdPut (userId int32) (error) {
//}

func updateInitPassword(userID int, password string) error {
	queryUser := models.User{UserID: userID}
	user, err := dao.GetUser(queryUser)
	if err != nil {
		return fmt.Errorf("Failed to get user, userID: %d %v", userID, err)
	}
	if user == nil {
		return fmt.Errorf("User id: %d does not exist.", userID)
	}
	if user.Salt == "" {
		salt, err := dao.GenerateRandomString()
		if err != nil {
			return fmt.Errorf("Failed to generate salt for encrypting password, %v", err)
		}

		user.Salt = salt
		user.Password = password
		err = dao.ChangeUserPassword(*user)
		if err != nil {
			return fmt.Errorf("Failed to update user encrypted password, userID: %d, err: %v", userID, err)
		}

	} else {
	}
	return nil
}