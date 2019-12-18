package route

import (
	ldap "../../ldap-authentication"
	"../config"
	"github.com/common-go/auth"
	"github.com/common-go/mongo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type EvaRoutes struct {
	Router *echo.Echo
}

func NewEvaluationRoutes(e *echo.Echo, mongoConfig mongo.MongoConfig, ldapConfig ldap.LDAPConfig, tokenConfig auth.TokenConfig) (*EvaRoutes, error) {
	//httpMethods := []string{"GET", "POST"}
	//applicationContext, _ := config.NewApplicationContext(msLandingUrl)
	applicationContext, err := config.NewApplicationContext(mongoConfig, ldapConfig, tokenConfig)
	if err != nil {
		return nil, err
	}

	//middle for all routes
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	authenticationController := applicationContext.AuthenticationController
	e.POST("/authentication/authenticate", authenticationController.Authenticate())
	batchController := applicationContext.BatchController
	batchPath := "/evaluation/batch"
	e.GET(batchPath, batchController.GetAll())
	e.POST(batchPath, batchController.Insert())
	e.GET(batchPath+"/:id", batchController.GetById())
	e.POST(batchPath+"/search", batchController.Search())
	e.PUT(batchPath+"/:id", batchController.Update())

	//group routes with /user prefix
	//g := e.Group("/user")
	//use middleware for group
	//g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	//	if username == "admin" && password == "123" {
	//		return true, nil
	//	}
	//	return false, nil
	//}))
	//g.Match(httpMethods, "/getFeed", applicationContext.LandingController.GetFeedContent)
	//e.POST("/user/getFeed", applicationContext.LandingController.GetFeedContent)
	//scheme
	schemeController := applicationContext.SchemeController
	schemePath :=  "/evaluation/scheme"
	e.GET(schemePath, schemeController.GetAll())
	e.POST(schemePath, schemeController.Insert())
	e.GET(schemePath+"/:id", schemeController.GetById())
	e.POST(schemePath+"/search", schemeController.Search())
	e.PUT(schemePath+"/:id", schemeController.Update())

	//candidate
	candidateController := applicationContext.CandidateController
	candidatePath := "/evaluation/candidate"
	candidateImportPath := "/evaluation/candidate/import"
	e.POST(candidateImportPath, candidateController.ImportArrayObject())

	e.GET(candidatePath, candidateController.GetAll())
	e.POST(candidatePath, candidateController.Insert())
	e.GET(candidatePath+"/:id", candidateController.GetById())
	e.POST(candidatePath+"/search", candidateController.Search())
	e.PUT(candidatePath+"/:id", candidateController.Update())

	//candidateEvaluation
	candidateEvaluationController := applicationContext.CandidateEvaluationController
	candidateEvaluationPath := "/evaluation/candidateEvaluation"
	e.GET(candidateEvaluationPath, candidateEvaluationController.GetAll())
	e.POST(candidateEvaluationPath, candidateEvaluationController.Insert())
	e.GET(candidateEvaluationPath+"/:id", candidateEvaluationController.GetById())
	e.POST(candidateEvaluationPath+"/search", candidateEvaluationController.Search())
	e.PUT(candidateEvaluationPath+"/:id", candidateEvaluationController.Update())

	//signOut
	signOutController := applicationContext.SignOutController
	signOutPath := "/authentication/signout/:userName"
	e.GET(signOutPath, signOutController.SignOut())


	return &EvaRoutes{e}, nil
}
