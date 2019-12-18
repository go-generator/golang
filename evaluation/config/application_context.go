package config

import (
	as "../../ldap-authentication"
	"../controller"
	"../service/impl"
	"context"
	"fmt"
	"github.com/common-go/auth"
	"github.com/common-go/jwt"
	"github.com/common-go/mongo"
	"reflect"
)

type JwtTokenVerifier struct {
}

func (t *JwtTokenVerifier) VerifyToken(tokenString string, secret string) (interface{}, int64, int64, error) {
	payload, c, err := jwt.VerifyToken(tokenString, secret)
	return payload, c.IssuedAt, c.ExpiresAt, err
}

type JwtTokenGenerator struct {
}

func (t *JwtTokenGenerator) GenerateToken(payload interface{}, secret string, expiresIn uint64) (string, error) {
	return jwt.GenerateToken(payload, secret, expiresIn)
}

type ApplicationContext struct {
	BatchController 							*controller.BatchController
	SchemeController 							*controller.SchemeController
	CandidateController 						*controller.CandidateController
	AuthenticationController 					*controller.AuthenticationController
	CandidateEvaluationController 				*controller.CandidateEvaluationController
	SignOutController 							*controller.SignOutController
}

func NewApplicationContext(mongoConfig mongo.MongoConfig, ldapConfig as.LDAPConfig, tokenConfig auth.TokenConfig) (*ApplicationContext, error) {
	ctx := context.Background()
	mongoDb, er1 := mongo.SetupMongo(ctx, mongoConfig)
	if er1 != nil {
		return nil, er1
	}
	//
	//db, er2 := sql.CreatePool(dbConfig)
	//if er2 != nil {
	//	return nil, er2
	//}

	mongoQueryBuilder := &mongo.DefaultQueryBuilder{}
	mongoSortBuilder := &mongo.DefaultSortBuilder{}
	mongoSearchResultBuilder := &mongo.DefaultSearchResultBuilder{
		Database:     mongoDb,
		QueryBuilder: mongoQueryBuilder,
		SortBuilder:  mongoSortBuilder,
	}
	// resultInfoBuilder := &builder.DefaultResultInfoBuilder{}
	activityLogService := impl.NewActivityLogServiceImpl(mongoDb)

	//authentication
	//userModuleService := NewUserServiceImpl(db)
	userInfoService := impl.NewUserInfoServiceImpl()
	//menuService := NewPrivilegeServiceImpl(db)

	tokenGenerator := &JwtTokenGenerator{}
	ldapAuthenticationService := &as.LDAPAuthenticationService{
		LDAPConfig:       ldapConfig,
		UserInfoService:  userInfoService,
		//PrivilegeService: menuService,
		TokenGenerator:   tokenGenerator,
		TokenConfig:      tokenConfig,
	}
	authenticationController := controller.NewAuthenticationController(ldapAuthenticationService, activityLogService)

	batchService := impl.NewBatchServiceImpl(mongoDb, mongoSearchResultBuilder)
	batchController := controller.NewBatchController(batchService, activityLogService)

	//scheme
	schemeService := impl.NewSchemeServiceImpl(mongoDb, mongoSearchResultBuilder)
	schemeController := controller.NewSchemeController(schemeService, activityLogService)

	//candidate
	candidateService := impl.NewCandidateServiceImpl(mongoDb, mongoSearchResultBuilder)
	candidateController := controller.NewCandidateController(candidateService, activityLogService)

	//CandidateEvaluation
	candidateEvaluationService := impl.NewCandidateEvaluationServiceImpl(mongoDb, mongoSearchResultBuilder, candidateService, schemeService)
	candidateEvaluationController := controller.NewCandidateEvaluationController(candidateEvaluationService, activityLogService)

	//SignOut
	// redisUrl := "redis://@localhost:6379"
	// redisService, _ := redis_client.NewRedisService(redisUrl)
	fmt.Println(reflect.TypeOf(tokenConfig.Expires))
	// blacklistTokenService := security.NewDefaultTokenBlacklistTokenService("", tokenConfig.Expires, redisService)
	signOutService := impl.NewSignOutServiceImpl()
	signOutController := controller.NewSignOutController(signOutService, activityLogService)

	return &ApplicationContext{batchController, schemeController, candidateController, authenticationController,candidateEvaluationController,signOutController}, nil
}
