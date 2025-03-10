package eventstore

import (
	"context"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	cache "github.com/zitadel/zitadel/internal/auth_request/repository"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	iam_view_model "github.com/zitadel/zitadel/internal/iam/repository/view/model"
	"github.com/zitadel/zitadel/internal/id"
	project_view_model "github.com/zitadel/zitadel/internal/project/repository/view/model"
	"github.com/zitadel/zitadel/internal/query"
	user_repo "github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	user_model "github.com/zitadel/zitadel/internal/user/model"
	user_view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

const unknownUserID = "UNKNOWN"

type AuthRequestRepo struct {
	Command      *command.Commands
	Query        *query.Queries
	AuthRequests cache.AuthRequestCache
	View         *view.View
	Eventstore   v1.Eventstore
	UserCodeAlg  crypto.EncryptionAlgorithm

	LabelPolicyProvider       labelPolicyProvider
	UserSessionViewProvider   userSessionViewProvider
	UserViewProvider          userViewProvider
	UserCommandProvider       userCommandProvider
	UserEventProvider         userEventProvider
	OrgViewProvider           orgViewProvider
	LoginPolicyViewProvider   loginPolicyViewProvider
	LockoutPolicyViewProvider lockoutPolicyViewProvider
	PrivacyPolicyProvider     privacyPolicyProvider
	IDPProviderViewProvider   idpProviderViewProvider
	UserGrantProvider         userGrantProvider
	ProjectProvider           projectProvider
	ApplicationProvider       applicationProvider

	IdGenerator id.Generator
}

type labelPolicyProvider interface {
	ActiveLabelPolicyByOrg(context.Context, string) (*query.LabelPolicy, error)
}

type privacyPolicyProvider interface {
	PrivacyPolicyByOrg(context.Context, bool, string) (*query.PrivacyPolicy, error)
}

type userSessionViewProvider interface {
	UserSessionByIDs(string, string, string) (*user_view_model.UserSessionView, error)
	UserSessionsByAgentID(string, string) ([]*user_view_model.UserSessionView, error)
}
type userViewProvider interface {
	UserByID(string, string) (*user_view_model.UserView, error)
}

type loginPolicyViewProvider interface {
	LoginPolicyByID(context.Context, bool, string) (*query.LoginPolicy, error)
}

type lockoutPolicyViewProvider interface {
	LockoutPolicyByOrg(context.Context, bool, string) (*query.LockoutPolicy, error)
}

type idpProviderViewProvider interface {
	IDPProvidersByAggregateIDAndState(string, string, iam_model.IDPConfigState) ([]*iam_view_model.IDPProviderView, error)
}

type userEventProvider interface {
	UserEventsByID(ctx context.Context, id string, sequence uint64) ([]*es_models.Event, error)
}

type userCommandProvider interface {
	BulkAddedUserIDPLinks(ctx context.Context, userID, resourceOwner string, externalIDPs []*domain.UserIDPLink) error
}

type orgViewProvider interface {
	OrgByID(context.Context, bool, string) (*query.Org, error)
	OrgByDomainGlobal(context.Context, string) (*query.Org, error)
}

type userGrantProvider interface {
	ProjectByOIDCClientID(context.Context, string) (*query.Project, error)
	UserGrantsByProjectAndUserID(context.Context, string, string) ([]*query.UserGrant, error)
}

type projectProvider interface {
	ProjectByOIDCClientID(context.Context, string) (*query.Project, error)
	OrgProjectMappingByIDs(orgID, projectID, instanceID string) (*project_view_model.OrgProjectMapping, error)
}

type applicationProvider interface {
	AppByOIDCClientID(context.Context, string) (*query.App, error)
}

func (repo *AuthRequestRepo) Health(ctx context.Context) error {
	return repo.AuthRequests.Health(ctx)
}

func (repo *AuthRequestRepo) CreateAuthRequest(ctx context.Context, request *domain.AuthRequest) (_ *domain.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	reqID, err := repo.IdGenerator.Next()
	if err != nil {
		return nil, err
	}
	request.ID = reqID
	project, err := repo.ProjectProvider.ProjectByOIDCClientID(ctx, request.ApplicationID)
	if err != nil {
		return nil, err
	}
	projectIDQuery, err := query.NewAppProjectIDSearchQuery(project.ID)
	if err != nil {
		return nil, err
	}
	appIDs, err := repo.Query.SearchClientIDs(ctx, &query.AppSearchQueries{Queries: []query.SearchQuery{projectIDQuery}})
	if err != nil {
		return nil, err
	}
	request.Audience = appIDs
	request.AppendAudIfNotExisting(project.ID)
	request.ApplicationResourceOwner = project.ResourceOwner
	request.PrivateLabelingSetting = project.PrivateLabelingSetting
	if err := setOrgID(repo.OrgViewProvider, request); err != nil {
		return nil, err
	}
	if request.LoginHint != "" {
		err = repo.checkLoginName(ctx, request, request.LoginHint)
		logging.LogWithFields("EVENT-aG311", "login name", request.LoginHint, "id", request.ID, "applicationID", request.ApplicationID, "traceID", tracing.TraceIDFromCtx(ctx)).OnError(err).Debug("login hint invalid")
	}
	if request.UserID == "" && request.LoginHint == "" && domain.IsPrompt(request.Prompt, domain.PromptNone) {
		err = repo.tryUsingOnlyUserSession(request)
		logging.LogWithFields("EVENT-SDf3g", "id", request.ID, "applicationID", request.ApplicationID, "traceID", tracing.TraceIDFromCtx(ctx)).OnError(err).Debug("unable to select only user session")
	}

	err = repo.AuthRequests.SaveAuthRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (repo *AuthRequestRepo) AuthRequestByID(ctx context.Context, id, userAgentID string) (_ *domain.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return repo.getAuthRequestNextSteps(ctx, id, userAgentID, false)
}

func (repo *AuthRequestRepo) AuthRequestByIDCheckLoggedIn(ctx context.Context, id, userAgentID string) (_ *domain.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return repo.getAuthRequestNextSteps(ctx, id, userAgentID, true)
}

func (repo *AuthRequestRepo) SaveAuthCode(ctx context.Context, id, code, userAgentID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, id, userAgentID)
	if err != nil {
		return err
	}
	request.Code = code
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) AuthRequestByCode(ctx context.Context, code string) (_ *domain.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.AuthRequests.GetAuthRequestByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	err = repo.fillPolicies(ctx, request)
	if err != nil {
		return nil, err
	}
	steps, err := repo.nextSteps(ctx, request, true)
	if err != nil {
		return nil, err
	}
	request.PossibleSteps = steps
	return request, nil
}

func (repo *AuthRequestRepo) DeleteAuthRequest(ctx context.Context, id string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return repo.AuthRequests.DeleteAuthRequest(ctx, id)
}

func (repo *AuthRequestRepo) CheckLoginName(ctx context.Context, id, loginName, userAgentID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, id, userAgentID)
	if err != nil {
		return err
	}
	err = repo.checkLoginName(ctx, request, loginName)
	if err != nil {
		return err
	}
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) SelectExternalIDP(ctx context.Context, authReqID, idpConfigID, userAgentID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	err = repo.checkSelectedExternalIDP(request, idpConfigID)
	if err != nil {
		return err
	}
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) CheckExternalUserLogin(ctx context.Context, authReqID, userAgentID string, externalUser *domain.ExternalUser, info *domain.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	err = repo.checkExternalUserLogin(ctx, request, externalUser.IDPConfigID, externalUser.ExternalUserID)
	if errors.IsNotFound(err) {
		if err := repo.setLinkingUser(ctx, request, externalUser); err != nil {
			return err
		}
		return err
	}
	if err != nil {
		return err
	}

	err = repo.Command.UserIDPLoginChecked(ctx, request.UserOrgID, request.UserID, request.WithCurrentInfo(info))
	if err != nil {
		return err
	}
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) SetExternalUserLogin(ctx context.Context, authReqID, userAgentID string, externalUser *domain.ExternalUser) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}

	err = repo.setLinkingUser(ctx, request, externalUser)
	if err != nil {
		return err
	}
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) setLinkingUser(ctx context.Context, request *domain.AuthRequest, externalUser *domain.ExternalUser) error {
	request.LinkingUsers = append(request.LinkingUsers, externalUser)
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) SelectUser(ctx context.Context, id, userID, userAgentID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, id, userAgentID)
	if err != nil {
		return err
	}
	user, err := activeUserByID(ctx, repo.UserViewProvider, repo.UserEventProvider, repo.OrgViewProvider, repo.LockoutPolicyViewProvider, userID, false)
	if err != nil {
		return err
	}
	if request.RequestedOrgID != "" && request.RequestedOrgID != user.ResourceOwner {
		return errors.ThrowPreconditionFailed(nil, "EVENT-fJe2a", "Errors.User.NotAllowedOrg")
	}
	username := user.UserName
	if request.RequestedOrgID == "" {
		username = user.PreferredLoginName
	}
	request.SetUserInfo(user.ID, username, user.PreferredLoginName, user.DisplayName, user.AvatarKey, user.ResourceOwner)
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) VerifyPassword(ctx context.Context, authReqID, userID, resourceOwner, password, userAgentID string, info *domain.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequestEnsureUser(ctx, authReqID, userAgentID, userID)
	if err != nil {
		if isIgnoreUserNotFoundError(err, request) {
			return errors.ThrowInvalidArgument(nil, "EVENT-SDe2f", "Errors.User.UsernameOrPassword.Invalid")
		}
		return err
	}
	policy, err := repo.getLockoutPolicy(ctx, resourceOwner)
	if err != nil {
		return err
	}
	err = repo.Command.HumanCheckPassword(ctx, resourceOwner, userID, password, request.WithCurrentInfo(info), lockoutPolicyToDomain(policy))
	if isIgnoreUserInvalidPasswordError(err, request) {
		return errors.ThrowInvalidArgument(nil, "EVENT-Jsf32", "Errors.User.UsernameOrPassword.Invalid")
	}
	return err
}

func isIgnoreUserNotFoundError(err error, request *domain.AuthRequest) bool {
	return request != nil && request.LoginPolicy != nil && request.LoginPolicy.IgnoreUnknownUsernames && errors.IsNotFound(err) && errors.Contains(err, "Errors.User.NotFound")
}

func isIgnoreUserInvalidPasswordError(err error, request *domain.AuthRequest) bool {
	return request != nil && request.LoginPolicy != nil && request.LoginPolicy.IgnoreUnknownUsernames && errors.IsErrorInvalidArgument(err) && errors.Contains(err, "Errors.User.Password.Invalid")
}

func lockoutPolicyToDomain(policy *query.LockoutPolicy) *domain.LockoutPolicy {
	return &domain.LockoutPolicy{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:   policy.ID,
			Sequence:      policy.Sequence,
			ResourceOwner: policy.ResourceOwner,
			CreationDate:  policy.CreationDate,
			ChangeDate:    policy.ChangeDate,
		},
		Default:             policy.IsDefault,
		MaxPasswordAttempts: policy.MaxPasswordAttempts,
		ShowLockOutFailures: policy.ShowFailures,
	}
}

func (repo *AuthRequestRepo) VerifyMFAOTP(ctx context.Context, authRequestID, userID, resourceOwner, code, userAgentID string, info *domain.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequestEnsureUser(ctx, authRequestID, userAgentID, userID)
	if err != nil {
		return err
	}
	return repo.Command.HumanCheckMFAOTP(ctx, userID, code, resourceOwner, request.WithCurrentInfo(info))
}

func (repo *AuthRequestRepo) BeginMFAU2FLogin(ctx context.Context, userID, resourceOwner, authRequestID, userAgentID string) (login *domain.WebAuthNLogin, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	request, err := repo.getAuthRequestEnsureUser(ctx, authRequestID, userAgentID, userID)
	if err != nil {
		return nil, err
	}
	return repo.Command.HumanBeginU2FLogin(ctx, userID, resourceOwner, request, true)
}

func (repo *AuthRequestRepo) VerifyMFAU2F(ctx context.Context, userID, resourceOwner, authRequestID, userAgentID string, credentialData []byte, info *domain.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequestEnsureUser(ctx, authRequestID, userAgentID, userID)
	if err != nil {
		return err
	}
	return repo.Command.HumanFinishU2FLogin(ctx, userID, resourceOwner, credentialData, request, true)
}

func (repo *AuthRequestRepo) BeginPasswordlessSetup(ctx context.Context, userID, resourceOwner string, authenticatorPlatform domain.AuthenticatorAttachment) (login *domain.WebAuthNToken, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return repo.Command.HumanAddPasswordlessSetup(ctx, userID, resourceOwner, true, authenticatorPlatform)
}

func (repo *AuthRequestRepo) VerifyPasswordlessSetup(ctx context.Context, userID, resourceOwner, userAgentID, tokenName string, credentialData []byte) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	_, err = repo.Command.HumanHumanPasswordlessSetup(ctx, userID, resourceOwner, tokenName, userAgentID, credentialData)
	return err
}

func (repo *AuthRequestRepo) BeginPasswordlessInitCodeSetup(ctx context.Context, userID, resourceOwner, codeID, verificationCode string, preferredPlatformType domain.AuthenticatorAttachment) (login *domain.WebAuthNToken, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	passwordlessInitCode, err := repo.Query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypePasswordlessInitCode, repo.UserCodeAlg)
	if err != nil {
		return nil, err
	}
	return repo.Command.HumanAddPasswordlessSetupInitCode(ctx, userID, resourceOwner, codeID, verificationCode, preferredPlatformType, passwordlessInitCode)
}

func (repo *AuthRequestRepo) VerifyPasswordlessInitCodeSetup(ctx context.Context, userID, resourceOwner, userAgentID, tokenName, codeID, verificationCode string, credentialData []byte) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	passwordlessInitCode, err := repo.Query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypePasswordlessInitCode, repo.UserCodeAlg)
	if err != nil {
		return err
	}
	_, err = repo.Command.HumanPasswordlessSetupInitCode(ctx, userID, resourceOwner, tokenName, userAgentID, codeID, verificationCode, credentialData, passwordlessInitCode)
	return err
}

func (repo *AuthRequestRepo) BeginPasswordlessLogin(ctx context.Context, userID, resourceOwner, authRequestID, userAgentID string) (login *domain.WebAuthNLogin, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequestEnsureUser(ctx, authRequestID, userAgentID, userID)
	if err != nil {
		return nil, err
	}
	return repo.Command.HumanBeginPasswordlessLogin(ctx, userID, resourceOwner, request, true)
}

func (repo *AuthRequestRepo) VerifyPasswordless(ctx context.Context, userID, resourceOwner, authRequestID, userAgentID string, credentialData []byte, info *domain.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequestEnsureUser(ctx, authRequestID, userAgentID, userID)
	if err != nil {
		return err
	}
	return repo.Command.HumanFinishPasswordlessLogin(ctx, userID, resourceOwner, credentialData, request, true)
}

func (repo *AuthRequestRepo) LinkExternalUsers(ctx context.Context, authReqID, userAgentID string, info *domain.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	err = linkExternalIDPs(ctx, repo.UserCommandProvider, request)
	if err != nil {
		return err
	}
	err = repo.Command.UserIDPLoginChecked(ctx, request.UserOrgID, request.UserID, request.WithCurrentInfo(info))
	if err != nil {
		return err
	}
	request.LinkingUsers = nil
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) ResetLinkingUsers(ctx context.Context, authReqID, userAgentID string) error {
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	request.LinkingUsers = nil
	request.SelectedIDPConfigID = ""
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) AutoRegisterExternalUser(ctx context.Context, registerUser *domain.Human, externalIDP *domain.UserIDPLink, orgMemberRoles []string, authReqID, userAgentID, resourceOwner string, metadatas []*domain.Metadata, info *domain.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	initCodeGenerator, err := repo.Query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeInitCode, repo.UserCodeAlg)
	if err != nil {
		return err
	}
	phoneCodeGenerator, err := repo.Query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyPhoneCode, repo.UserCodeAlg)
	if err != nil {
		return err
	}
	human, err := repo.Command.RegisterHuman(ctx, resourceOwner, registerUser, externalIDP, orgMemberRoles, initCodeGenerator, phoneCodeGenerator)
	if err != nil {
		return err
	}
	request.UserID = human.AggregateID
	request.UserOrgID = human.ResourceOwner
	request.SelectedIDPConfigID = externalIDP.IDPConfigID
	request.LinkingUsers = nil
	err = repo.Command.UserIDPLoginChecked(ctx, request.UserOrgID, request.UserID, request.WithCurrentInfo(info))
	if err != nil {
		return err
	}
	if len(metadatas) > 0 {
		_, err = repo.Command.BulkSetUserMetadata(ctx, request.UserID, request.UserOrgID, metadatas...)
		if err != nil {
			return err
		}
	}
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) getAuthRequestNextSteps(ctx context.Context, id, userAgentID string, checkLoggedIn bool) (*domain.AuthRequest, error) {
	request, err := repo.getAuthRequest(ctx, id, userAgentID)
	if err != nil {
		return nil, err
	}
	steps, err := repo.nextSteps(ctx, request, checkLoggedIn)
	if err != nil {
		return request, err
	}
	request.PossibleSteps = steps
	return request, nil
}

func (repo *AuthRequestRepo) getAuthRequestEnsureUser(ctx context.Context, authRequestID, userAgentID, userID string) (*domain.AuthRequest, error) {
	request, err := repo.getAuthRequest(ctx, authRequestID, userAgentID)
	if err != nil {
		return nil, err
	}
	if request.UserID != userID {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-GBH32", "Errors.User.NotMatchingUserID")
	}
	_, err = activeUserByID(ctx, repo.UserViewProvider, repo.UserEventProvider, repo.OrgViewProvider, repo.LockoutPolicyViewProvider, request.UserID, false)
	if err != nil {
		return request, err
	}
	return request, nil
}

func (repo *AuthRequestRepo) getAuthRequest(ctx context.Context, id, userAgentID string) (*domain.AuthRequest, error) {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if request.AgentID != userAgentID {
		return nil, errors.ThrowPermissionDenied(nil, "EVENT-adk13", "Errors.AuthRequest.UserAgentNotCorresponding")
	}
	err = repo.fillPolicies(ctx, request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (repo *AuthRequestRepo) getLoginPolicyAndIDPProviders(ctx context.Context, orgID string) (*query.LoginPolicy, []*domain.IDPProvider, error) {
	policy, err := repo.LoginPolicyViewProvider.LoginPolicyByID(ctx, false, orgID)
	if err != nil {
		return nil, nil, err
	}
	if !policy.AllowExternalIDPs {
		return policy, nil, nil
	}
	idpProviders, err := getLoginPolicyIDPProviders(repo.IDPProviderViewProvider, authz.GetInstance(ctx).InstanceID(), orgID, policy.IsDefault)
	if err != nil {
		return nil, nil, err
	}

	providers := iam_model.IdpProviderViewsToDomain(idpProviders)
	return policy, providers, nil
}

func (repo *AuthRequestRepo) fillPolicies(ctx context.Context, request *domain.AuthRequest) error {
	orgID := request.RequestedOrgID
	if orgID == "" {
		orgID = request.UserOrgID
	}
	if orgID == "" {
		orgID = authz.GetInstance(ctx).InstanceID()
	}

	loginPolicy, idpProviders, err := repo.getLoginPolicyAndIDPProviders(ctx, orgID)
	if err != nil {
		return err
	}
	request.LoginPolicy = queryLoginPolicyToDomain(loginPolicy)
	if idpProviders != nil {
		request.AllowedExternalIDPs = idpProviders
	}
	lockoutPolicy, err := repo.getLockoutPolicy(ctx, orgID)
	if err != nil {
		return err
	}
	request.LockoutPolicy = lockoutPolicyToDomain(lockoutPolicy)
	privacyPolicy, err := repo.GetPrivacyPolicy(ctx, orgID)
	if err != nil {
		return err
	}
	request.PrivacyPolicy = privacyPolicy
	privateLabelingOrgID := authz.GetInstance(ctx).InstanceID()
	if request.PrivateLabelingSetting != domain.PrivateLabelingSettingUnspecified {
		privateLabelingOrgID = request.ApplicationResourceOwner
	}
	if request.PrivateLabelingSetting == domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy || request.PrivateLabelingSetting == domain.PrivateLabelingSettingUnspecified {
		if request.UserOrgID != "" {
			privateLabelingOrgID = request.UserOrgID
		}
	}
	if request.RequestedOrgID != "" {
		privateLabelingOrgID = request.RequestedOrgID
	}
	labelPolicy, err := repo.getLabelPolicy(ctx, privateLabelingOrgID)
	if err != nil {
		return err
	}
	request.LabelPolicy = labelPolicy
	defaultLoginTranslations, err := repo.getLoginTexts(ctx, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return err
	}
	request.DefaultTranslations = defaultLoginTranslations
	orgLoginTranslations, err := repo.getLoginTexts(ctx, orgID)
	if err != nil {
		return err
	}
	request.OrgTranslations = orgLoginTranslations
	return nil
}

func (repo *AuthRequestRepo) tryUsingOnlyUserSession(request *domain.AuthRequest) error {
	userSessions, err := userSessionsByUserAgentID(repo.UserSessionViewProvider, request.AgentID, request.InstanceID)
	if err != nil {
		return err
	}
	if len(userSessions) == 1 {
		user := userSessions[0]
		username := user.UserName
		if request.RequestedOrgID == "" {
			username = user.LoginName
		}
		request.SetUserInfo(user.UserID, username, user.LoginName, user.DisplayName, user.AvatarKey, user.ResourceOwner)
	}
	return nil
}

func (repo *AuthRequestRepo) checkLoginName(ctx context.Context, request *domain.AuthRequest, loginName string) (err error) {
	user := new(user_view_model.UserView)
	preferredLoginName := loginName
	if request.RequestedOrgID != "" {
		if request.RequestedOrgID != "" {
			preferredLoginName += "@" + request.RequestedPrimaryDomain
		}
		user, err = repo.View.UserByLoginNameAndResourceOwner(preferredLoginName, request.RequestedOrgID, request.InstanceID)
	} else {
		user, err = repo.View.UserByLoginName(loginName, request.InstanceID)
		if err == nil {
			err = repo.checkLoginPolicyWithResourceOwner(ctx, request, user)
			if err != nil {
				return err
			}
		}
	}
	if request.LoginPolicy.IgnoreUnknownUsernames {
		if errors.IsNotFound(err) || (user != nil && user.State == int32(domain.UserStateInactive)) {
			if request.LabelPolicy.HideLoginNameSuffix {
				preferredLoginName = loginName
			}
			request.SetUserInfo(unknownUserID, preferredLoginName, preferredLoginName, preferredLoginName, "", request.RequestedOrgID)
			return nil
		}
	}
	if err != nil {
		return err
	}
	if user.State == int32(domain.UserStateInactive) {
		return errors.ThrowPreconditionFailed(nil, "AUTH-2n8fs", "Errors.User.Inactive")
	}
	request.SetUserInfo(user.ID, loginName, user.PreferredLoginName, "", "", user.ResourceOwner)
	return nil
}

func (repo AuthRequestRepo) checkLoginPolicyWithResourceOwner(ctx context.Context, request *domain.AuthRequest, user *user_view_model.UserView) error {
	loginPolicy, idpProviders, err := repo.getLoginPolicyAndIDPProviders(ctx, user.ResourceOwner)
	if err != nil {
		return err
	}
	if len(request.LinkingUsers) != 0 && !loginPolicy.AllowExternalIDPs {
		return errors.ThrowInvalidArgument(nil, "LOGIN-s9sio", "Errors.User.NotAllowedToLink")
	}
	if len(request.LinkingUsers) != 0 {
		exists := linkingIDPConfigExistingInAllowedIDPs(request.LinkingUsers, idpProviders)
		if !exists {
			return errors.ThrowInvalidArgument(nil, "LOGIN-Dj89o", "Errors.User.NotAllowedToLink")
		}
	}
	request.LoginPolicy = queryLoginPolicyToDomain(loginPolicy)
	request.AllowedExternalIDPs = idpProviders
	return nil
}

func queryLoginPolicyToDomain(policy *query.LoginPolicy) *domain.LoginPolicy {
	return &domain.LoginPolicy{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:   policy.OrgID,
			Sequence:      policy.Sequence,
			ResourceOwner: policy.OrgID,
			CreationDate:  policy.CreationDate,
			ChangeDate:    policy.ChangeDate,
		},
		Default:                    policy.IsDefault,
		AllowUsernamePassword:      policy.AllowUsernamePassword,
		AllowRegister:              policy.AllowRegister,
		AllowExternalIDP:           policy.AllowExternalIDPs,
		ForceMFA:                   policy.ForceMFA,
		SecondFactors:              policy.SecondFactors,
		MultiFactors:               policy.MultiFactors,
		PasswordlessType:           policy.PasswordlessType,
		HidePasswordReset:          policy.HidePasswordReset,
		IgnoreUnknownUsernames:     policy.IgnoreUnknownUsernames,
		PasswordCheckLifetime:      policy.PasswordCheckLifetime,
		ExternalLoginCheckLifetime: policy.ExternalLoginCheckLifetime,
		MFAInitSkipLifetime:        policy.MFAInitSkipLifetime,
		SecondFactorCheckLifetime:  policy.SecondFactorCheckLifetime,
		MultiFactorCheckLifetime:   policy.MultiFactorCheckLifetime,
	}
}

func (repo *AuthRequestRepo) checkSelectedExternalIDP(request *domain.AuthRequest, idpConfigID string) error {
	for _, externalIDP := range request.AllowedExternalIDPs {
		if externalIDP.IDPConfigID == idpConfigID {
			request.SelectedIDPConfigID = idpConfigID
			return nil
		}
	}
	return errors.ThrowNotFound(nil, "LOGIN-Nsm8r", "Errors.User.ExternalIDP.NotAllowed")
}

func (repo *AuthRequestRepo) checkExternalUserLogin(ctx context.Context, request *domain.AuthRequest, idpConfigID, externalUserID string) (err error) {
	externalIDP := new(user_view_model.ExternalIDPView)
	if request.RequestedOrgID != "" {
		externalIDP, err = repo.View.ExternalIDPByExternalUserIDAndIDPConfigIDAndResourceOwner(externalUserID, idpConfigID, request.RequestedOrgID, request.InstanceID)
	} else {
		externalIDP, err = repo.View.ExternalIDPByExternalUserIDAndIDPConfigID(externalUserID, idpConfigID, request.InstanceID)
	}
	if err != nil {
		return err
	}
	user, err := activeUserByID(ctx, repo.UserViewProvider, repo.UserEventProvider, repo.OrgViewProvider, repo.LockoutPolicyViewProvider, externalIDP.UserID, false)
	if err != nil {
		return err
	}
	username := user.UserName
	if request.RequestedOrgID == "" {
		username = user.PreferredLoginName
	}
	request.SetUserInfo(user.ID, username, user.PreferredLoginName, user.DisplayName, user.AvatarKey, user.ResourceOwner)
	return nil
}

func (repo *AuthRequestRepo) nextSteps(ctx context.Context, request *domain.AuthRequest, checkLoggedIn bool) ([]domain.NextStep, error) {
	if request == nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-ds27a", "Errors.Internal")
	}
	steps := make([]domain.NextStep, 0)
	if !checkLoggedIn && domain.IsPrompt(request.Prompt, domain.PromptNone) {
		return append(steps, &domain.RedirectToCallbackStep{}), nil
	}
	if request.UserID == "" {
		if request.LinkingUsers != nil && len(request.LinkingUsers) > 0 {
			steps = append(steps, new(domain.ExternalNotFoundOptionStep))
			return steps, nil
		}
		steps = append(steps, new(domain.LoginStep))
		if domain.IsPrompt(request.Prompt, domain.PromptCreate) {
			return append(steps, &domain.RegistrationStep{}), nil
		}
		if len(request.Prompt) == 0 || domain.IsPrompt(request.Prompt, domain.PromptSelectAccount) {
			users, err := repo.usersForUserSelection(request)
			if err != nil {
				return nil, err
			}
			if domain.IsPrompt(request.Prompt, domain.PromptSelectAccount) {
				steps = append(steps, &domain.SelectUserStep{Users: users})
			}
			if request.SelectedIDPConfigID != "" {
				steps = append(steps, &domain.RedirectToExternalIDPStep{})
			}
			if len(request.Prompt) == 0 && len(users) > 0 {
				steps = append(steps, &domain.SelectUserStep{Users: users})
			}
		}
		return steps, nil
	}
	user, err := activeUserByID(ctx, repo.UserViewProvider, repo.UserEventProvider, repo.OrgViewProvider, repo.LockoutPolicyViewProvider, request.UserID, request.LoginPolicy.IgnoreUnknownUsernames)
	if err != nil {
		return nil, err
	}
	if user.PreferredLoginName != "" {
		request.LoginName = user.PreferredLoginName
	}
	userSession, err := userSessionByIDs(ctx, repo.UserSessionViewProvider, repo.UserEventProvider, request.AgentID, user)
	if err != nil {
		return nil, err
	}

	isInternalLogin := request.SelectedIDPConfigID == "" && userSession.SelectedIDPConfigID == ""
	if !isInternalLogin && len(request.LinkingUsers) == 0 && !checkVerificationTimeMaxAge(userSession.ExternalLoginVerification, request.LoginPolicy.ExternalLoginCheckLifetime, request) {
		selectedIDPConfigID := request.SelectedIDPConfigID
		if selectedIDPConfigID == "" {
			selectedIDPConfigID = userSession.SelectedIDPConfigID
		}
		return append(steps, &domain.ExternalLoginStep{SelectedIDPConfigID: selectedIDPConfigID}), nil
	}
	if isInternalLogin || (!isInternalLogin && len(request.LinkingUsers) > 0) {
		step := repo.firstFactorChecked(request, user, userSession)
		if step != nil {
			return append(steps, step), nil
		}
	}

	step, ok, err := repo.mfaChecked(userSession, request, user)
	if err != nil {
		return nil, err
	}
	if !ok {
		return append(steps, step), nil
	}

	if user.PasswordChangeRequired {
		steps = append(steps, &domain.ChangePasswordStep{})
	}
	if !user.IsEmailVerified {
		steps = append(steps, &domain.VerifyEMailStep{})
	}
	if user.UsernameChangeRequired {
		steps = append(steps, &domain.ChangeUsernameStep{})
	}

	if user.PasswordChangeRequired || !user.IsEmailVerified || user.UsernameChangeRequired {
		return steps, nil
	}

	if request.LinkingUsers != nil && len(request.LinkingUsers) != 0 {
		return append(steps, &domain.LinkUsersStep{}), nil
	}
	//PLANNED: consent step

	missing, err := projectRequired(ctx, request, repo.ProjectProvider)
	if err != nil {
		return nil, err
	}
	if missing {
		return append(steps, &domain.ProjectRequiredStep{}), nil
	}

	missing, err = userGrantRequired(ctx, request, user, repo.UserGrantProvider)
	if err != nil {
		return nil, err
	}
	if missing {
		return append(steps, &domain.GrantRequiredStep{}), nil
	}

	ok, err = repo.hasSucceededPage(ctx, request, repo.ApplicationProvider)
	if err != nil {
		return nil, err
	}
	if ok {
		steps = append(steps, &domain.LoginSucceededStep{})
	}
	return append(steps, &domain.RedirectToCallbackStep{}), nil
}

func (repo *AuthRequestRepo) usersForUserSelection(request *domain.AuthRequest) ([]domain.UserSelection, error) {
	userSessions, err := userSessionsByUserAgentID(repo.UserSessionViewProvider, request.AgentID, request.InstanceID)
	if err != nil {
		return nil, err
	}
	users := make([]domain.UserSelection, 0)
	for _, session := range userSessions {
		if request.RequestedOrgID == "" || request.RequestedOrgID == session.ResourceOwner {
			users = append(users, domain.UserSelection{
				UserID:            session.UserID,
				DisplayName:       session.DisplayName,
				UserName:          session.UserName,
				LoginName:         session.LoginName,
				ResourceOwner:     session.ResourceOwner,
				AvatarKey:         session.AvatarKey,
				UserSessionState:  session.State,
				SelectionPossible: request.RequestedOrgID == "" || request.RequestedOrgID == session.ResourceOwner,
			})
		}
	}
	return users, nil
}

func (repo *AuthRequestRepo) firstFactorChecked(request *domain.AuthRequest, user *user_model.UserView, userSession *user_model.UserSessionView) domain.NextStep {
	if user.InitRequired {
		return &domain.InitUserStep{PasswordSet: user.PasswordSet}
	}

	var step domain.NextStep
	if request.LoginPolicy.PasswordlessType != domain.PasswordlessTypeNotAllowed && user.IsPasswordlessReady() {
		if checkVerificationTimeMaxAge(userSession.PasswordlessVerification, request.LoginPolicy.MultiFactorCheckLifetime, request) {
			request.AuthTime = userSession.PasswordlessVerification
			return nil
		}
		step = &domain.PasswordlessStep{
			PasswordSet: user.PasswordSet,
		}
	}

	if user.PasswordlessInitRequired {
		return &domain.PasswordlessRegistrationPromptStep{}
	}

	if user.PasswordInitRequired {
		return &domain.InitPasswordStep{}
	}

	if checkVerificationTimeMaxAge(userSession.PasswordVerification, request.LoginPolicy.PasswordCheckLifetime, request) {
		request.PasswordVerified = true
		request.AuthTime = userSession.PasswordVerification
		return nil
	}
	if step != nil {
		return step
	}
	return &domain.PasswordStep{}
}

func (repo *AuthRequestRepo) mfaChecked(userSession *user_model.UserSessionView, request *domain.AuthRequest, user *user_model.UserView) (domain.NextStep, bool, error) {
	mfaLevel := request.MFALevel()
	allowedProviders, required := user.MFATypesAllowed(mfaLevel, request.LoginPolicy)
	promptRequired := (user.MFAMaxSetUp < mfaLevel) || (len(allowedProviders) == 0 && required)
	if promptRequired || !repo.mfaSkippedOrSetUp(user, request) {
		types := user.MFATypesSetupPossible(mfaLevel, request.LoginPolicy)
		if promptRequired && len(types) == 0 {
			return nil, false, errors.ThrowPreconditionFailed(nil, "LOGIN-5Hm8s", "Errors.Login.LoginPolicy.MFA.ForceAndNotConfigured")
		}
		if len(types) == 0 {
			return nil, true, nil
		}
		return &domain.MFAPromptStep{
			Required:     promptRequired,
			MFAProviders: types,
		}, false, nil
	}
	switch mfaLevel {
	default:
		fallthrough
	case domain.MFALevelNotSetUp:
		if len(allowedProviders) == 0 {
			return nil, true, nil
		}
		fallthrough
	case domain.MFALevelSecondFactor:
		if checkVerificationTimeMaxAge(userSession.SecondFactorVerification, request.LoginPolicy.SecondFactorCheckLifetime, request) {
			request.MFAsVerified = append(request.MFAsVerified, userSession.SecondFactorVerificationType)
			request.AuthTime = userSession.SecondFactorVerification
			return nil, true, nil
		}
		fallthrough
	case domain.MFALevelMultiFactor:
		if checkVerificationTimeMaxAge(userSession.MultiFactorVerification, request.LoginPolicy.MultiFactorCheckLifetime, request) {
			request.MFAsVerified = append(request.MFAsVerified, userSession.MultiFactorVerificationType)
			request.AuthTime = userSession.MultiFactorVerification
			return nil, true, nil
		}
	}
	return &domain.MFAVerificationStep{
		MFAProviders: allowedProviders,
	}, false, nil
}

func (repo *AuthRequestRepo) mfaSkippedOrSetUp(user *user_model.UserView, request *domain.AuthRequest) bool {
	if user.MFAMaxSetUp > domain.MFALevelNotSetUp {
		return true
	}
	return checkVerificationTime(user.MFAInitSkipped, request.LoginPolicy.MFAInitSkipLifetime)
}

func (repo *AuthRequestRepo) GetPrivacyPolicy(ctx context.Context, orgID string) (*domain.PrivacyPolicy, error) {
	policy, err := repo.PrivacyPolicyProvider.PrivacyPolicyByOrg(ctx, false, orgID)
	if errors.IsNotFound(err) {
		return new(domain.PrivacyPolicy), nil
	}
	if err != nil {
		return nil, err
	}
	return privacyPolicyToDomain(policy), err
}

func privacyPolicyToDomain(p *query.PrivacyPolicy) *domain.PrivacyPolicy {
	return &domain.PrivacyPolicy{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:   p.ID,
			Sequence:      p.Sequence,
			ResourceOwner: p.ResourceOwner,
			CreationDate:  p.CreationDate,
			ChangeDate:    p.ChangeDate,
		},
		State:       p.State,
		Default:     p.IsDefault,
		TOSLink:     p.TOSLink,
		PrivacyLink: p.PrivacyLink,
		HelpLink:    p.HelpLink,
	}
}

func (repo *AuthRequestRepo) getLockoutPolicy(ctx context.Context, orgID string) (*query.LockoutPolicy, error) {
	policy, err := repo.LockoutPolicyViewProvider.LockoutPolicyByOrg(ctx, false, orgID)
	if err != nil {
		return nil, err
	}
	return policy, err
}

func (repo *AuthRequestRepo) getLabelPolicy(ctx context.Context, orgID string) (*domain.LabelPolicy, error) {
	policy, err := repo.LabelPolicyProvider.ActiveLabelPolicyByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return labelPolicyToDomain(policy), nil
}

func labelPolicyToDomain(p *query.LabelPolicy) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:   p.ID,
			Sequence:      p.Sequence,
			ResourceOwner: p.ResourceOwner,
			CreationDate:  p.CreationDate,
			ChangeDate:    p.ChangeDate,
		},
		State:               p.State,
		Default:             p.IsDefault,
		PrimaryColor:        p.Light.PrimaryColor,
		BackgroundColor:     p.Light.BackgroundColor,
		WarnColor:           p.Light.WarnColor,
		FontColor:           p.Light.FontColor,
		LogoURL:             p.Light.LogoURL,
		IconURL:             p.Light.IconURL,
		PrimaryColorDark:    p.Dark.PrimaryColor,
		BackgroundColorDark: p.Dark.BackgroundColor,
		WarnColorDark:       p.Dark.WarnColor,
		FontColorDark:       p.Dark.FontColor,
		LogoDarkURL:         p.Dark.LogoURL,
		IconDarkURL:         p.Dark.IconURL,
		Font:                p.FontURL,
		HideLoginNameSuffix: p.HideLoginNameSuffix,
		ErrorMsgPopup:       p.ShouldErrorPopup,
		DisableWatermark:    p.WatermarkDisabled,
	}
}

func (repo *AuthRequestRepo) getLoginTexts(ctx context.Context, aggregateID string) ([]*domain.CustomText, error) {
	loginTexts, err := repo.Query.CustomTextListByTemplate(ctx, aggregateID, domain.LoginCustomText)
	if err != nil {
		return nil, err
	}
	return query.CustomTextsToDomain(loginTexts), err
}

func (repo *AuthRequestRepo) hasSucceededPage(ctx context.Context, request *domain.AuthRequest, provider applicationProvider) (bool, error) {
	app, err := provider.AppByOIDCClientID(ctx, request.ApplicationID)
	if err != nil {
		return false, err
	}
	return app.OIDCConfig.AppType == domain.OIDCApplicationTypeNative, nil
}

func setOrgID(orgViewProvider orgViewProvider, request *domain.AuthRequest) error {
	primaryDomain := request.GetScopeOrgPrimaryDomain()
	if primaryDomain == "" {
		return nil
	}

	org, err := orgViewProvider.OrgByDomainGlobal(context.TODO(), primaryDomain)
	if err != nil {
		return err
	}
	request.RequestedOrgID = org.ID
	request.RequestedOrgName = org.Name
	request.RequestedPrimaryDomain = primaryDomain
	return nil
}

func getLoginPolicyIDPProviders(provider idpProviderViewProvider, iamID, orgID string, defaultPolicy bool) ([]*iam_model.IDPProviderView, error) {
	if defaultPolicy {
		idpProviders, err := provider.IDPProvidersByAggregateIDAndState(iamID, iamID, iam_model.IDPConfigStateActive)
		if err != nil {
			return nil, err
		}
		return iam_view_model.IDPProviderViewsToModel(idpProviders), nil
	}
	idpProviders, err := provider.IDPProvidersByAggregateIDAndState(orgID, iamID, iam_model.IDPConfigStateActive)
	if err != nil {
		return nil, err
	}
	return iam_view_model.IDPProviderViewsToModel(idpProviders), nil
}

func checkVerificationTimeMaxAge(verificationTime time.Time, lifetime time.Duration, request *domain.AuthRequest) bool {
	if !checkVerificationTime(verificationTime, lifetime) {
		return false
	}
	if request.MaxAuthAge == nil {
		return true
	}
	return verificationTime.After(request.CreationDate.Add(-*request.MaxAuthAge))
}

func checkVerificationTime(verificationTime time.Time, lifetime time.Duration) bool {
	return verificationTime.Add(lifetime).After(time.Now().UTC())
}

func userSessionsByUserAgentID(provider userSessionViewProvider, agentID, instanceID string) ([]*user_model.UserSessionView, error) {
	session, err := provider.UserSessionsByAgentID(agentID, instanceID)
	if err != nil {
		return nil, err
	}
	return user_view_model.UserSessionsToModel(session), nil
}

func userSessionByIDs(ctx context.Context, provider userSessionViewProvider, eventProvider userEventProvider, agentID string, user *user_model.UserView) (*user_model.UserSessionView, error) {
	session, err := provider.UserSessionByIDs(agentID, user.ID, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		if !errors.IsNotFound(err) {
			return nil, err
		}
		session = &user_view_model.UserSessionView{UserAgentID: agentID, UserID: user.ID}
	}
	events, err := eventProvider.UserEventsByID(ctx, user.ID, session.Sequence)
	if err != nil {
		logging.Log("EVENT-Hse6s").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("error retrieving new events")
		return user_view_model.UserSessionToModel(session), nil
	}
	sessionCopy := *session
	for _, event := range events {
		switch eventstore.EventType(event.Type) {
		case user_repo.UserV1PasswordCheckSucceededType,
			user_repo.UserV1PasswordCheckFailedType,
			user_repo.UserV1MFAOTPCheckSucceededType,
			user_repo.UserV1MFAOTPCheckFailedType,
			user_repo.UserV1SignedOutType,
			user_repo.UserLockedType,
			user_repo.UserDeactivatedType,
			user_repo.HumanPasswordCheckSucceededType,
			user_repo.HumanPasswordCheckFailedType,
			user_repo.UserIDPLoginCheckSucceededType,
			user_repo.HumanMFAOTPCheckSucceededType,
			user_repo.HumanMFAOTPCheckFailedType,
			user_repo.HumanSignedOutType,
			user_repo.HumanPasswordlessTokenCheckSucceededType,
			user_repo.HumanPasswordlessTokenCheckFailedType,
			user_repo.HumanU2FTokenCheckSucceededType,
			user_repo.HumanU2FTokenCheckFailedType:
			eventData, err := user_view_model.UserSessionFromEvent(event)
			if err != nil {
				logging.Log("EVENT-sdgT3").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("error getting event data")
				return user_view_model.UserSessionToModel(session), nil
			}
			if eventData.UserAgentID != agentID {
				continue
			}
		case user_repo.UserRemovedType:
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dG2fe", "Errors.User.NotActive")
		}
		err := sessionCopy.AppendEvent(event)
		logging.Log("EVENT-qbhj3").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Warn("error appending event")
	}
	return user_view_model.UserSessionToModel(&sessionCopy), nil
}

func activeUserByID(ctx context.Context, userViewProvider userViewProvider, userEventProvider userEventProvider, queries orgViewProvider, lockoutPolicyProvider lockoutPolicyViewProvider, userID string, ignoreUnknownUsernames bool) (user *user_model.UserView, err error) {
	// PLANNED: Check LockoutPolicy
	user, err = userByID(ctx, userViewProvider, userEventProvider, userID)
	if err != nil {
		if ignoreUnknownUsernames && errors.IsNotFound(err) {
			return &user_model.UserView{
				ID:        userID,
				HumanView: &user_model.HumanView{},
			}, nil
		}
		return nil, err
	}

	if user.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Lm69x", "Errors.User.NotHuman")
	}
	if user.State == user_model.UserStateLocked || user.State == user_model.UserStateSuspend {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-FJ262", "Errors.User.Locked")
	}
	if !(user.State == user_model.UserStateActive || user.State == user_model.UserStateInitial) {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-FJ262", "Errors.User.NotActive")
	}
	org, err := queries.OrgByID(ctx, false, user.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if org.State != domain.OrgStateActive {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Zws3s", "Errors.User.NotActive")
	}
	return user, nil
}

func userByID(ctx context.Context, viewProvider userViewProvider, eventProvider userEventProvider, userID string) (*user_model.UserView, error) {
	user, viewErr := viewProvider.UserByID(userID, authz.GetInstance(ctx).InstanceID())
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	} else if user == nil {
		user = new(user_view_model.UserView)
	}
	events, err := eventProvider.UserEventsByID(ctx, userID, user.Sequence)
	if err != nil {
		logging.Log("EVENT-dfg42").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("error retrieving new events")
		return user_view_model.UserToModel(user), nil
	}
	if len(events) == 0 {
		if viewErr != nil {
			return nil, viewErr
		}
		return user_view_model.UserToModel(user), viewErr
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return user_view_model.UserToModel(user), nil
		}
	}
	if userCopy.State == int32(user_model.UserStateDeleted) {
		return nil, errors.ThrowNotFound(nil, "EVENT-3F9so", "Errors.User.NotFound")
	}
	return user_view_model.UserToModel(&userCopy), nil
}

func linkExternalIDPs(ctx context.Context, userCommandProvider userCommandProvider, request *domain.AuthRequest) error {
	externalIDPs := make([]*domain.UserIDPLink, len(request.LinkingUsers))
	for i, linkingUser := range request.LinkingUsers {
		externalIDP := &domain.UserIDPLink{
			ObjectRoot:     es_models.ObjectRoot{AggregateID: request.UserID},
			IDPConfigID:    linkingUser.IDPConfigID,
			ExternalUserID: linkingUser.ExternalUserID,
			DisplayName:    linkingUser.DisplayName,
		}
		externalIDPs[i] = externalIDP
	}
	data := authz.CtxData{
		UserID: "LOGIN",
		OrgID:  request.UserOrgID,
	}
	return userCommandProvider.BulkAddedUserIDPLinks(authz.SetCtxData(ctx, data), request.UserID, request.UserOrgID, externalIDPs)
}

func linkingIDPConfigExistingInAllowedIDPs(linkingUsers []*domain.ExternalUser, idpProviders []*domain.IDPProvider) bool {
	for _, linkingUser := range linkingUsers {
		exists := false
		for _, idp := range idpProviders {
			if idp.IDPConfigID == linkingUser.IDPConfigID {
				exists = true
				continue
			}
		}
		if !exists {
			return false
		}
	}
	return true
}

func userGrantRequired(ctx context.Context, request *domain.AuthRequest, user *user_model.UserView, userGrantProvider userGrantProvider) (_ bool, err error) {
	var project *query.Project
	switch request.Request.Type() {
	case domain.AuthRequestTypeOIDC:
		project, err = userGrantProvider.ProjectByOIDCClientID(ctx, request.ApplicationID)
		if err != nil {
			return false, err
		}
	default:
		return false, errors.ThrowPreconditionFailed(nil, "EVENT-dfrw2", "Errors.AuthRequest.RequestTypeNotSupported")
	}
	if !project.ProjectRoleCheck {
		return false, nil
	}
	grants, err := userGrantProvider.UserGrantsByProjectAndUserID(ctx, project.ID, user.ID)
	if err != nil {
		return false, err
	}
	return len(grants) == 0, nil
}

func projectRequired(ctx context.Context, request *domain.AuthRequest, projectProvider projectProvider) (_ bool, err error) {
	var project *query.Project
	switch request.Request.Type() {
	case domain.AuthRequestTypeOIDC:
		project, err = projectProvider.ProjectByOIDCClientID(ctx, request.ApplicationID)
		if err != nil {
			return false, err
		}
	default:
		return false, errors.ThrowPreconditionFailed(nil, "EVENT-dfrw2", "Errors.AuthRequest.RequestTypeNotSupported")
	}
	if !project.HasProjectCheck {
		return false, nil
	}
	_, err = projectProvider.OrgProjectMappingByIDs(request.UserOrgID, project.ID, request.InstanceID)
	if errors.IsNotFound(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, nil
}
