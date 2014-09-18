// The iam package provides types and functions for interaction with the AWS
// Identity and Access Management (IAM) service.
package iam

import (
	"encoding/xml"
	"github.com/crowdmob/goamz/aws"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// The IAM type encapsulates operations operations with the IAM endpoint.
type IAM struct {
	aws.Auth
	aws.Region
}

// New creates a new IAM instance.
func New(auth aws.Auth, region aws.Region) *IAM {
	return &IAM{auth, region}
}

func (iam *IAM) query(params map[string]string, resp interface{}) error {
	params["Version"] = "2010-05-08"
	params["Timestamp"] = time.Now().In(time.UTC).Format(time.RFC3339)
	endpoint, err := url.Parse(iam.IAMEndpoint)
	if err != nil {
		return err
	}
	sign(iam.Auth, "GET", "/", params, endpoint.Host)
	endpoint.RawQuery = multimap(params).Encode()
	r, err := http.Get(endpoint.String())
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode > 200 {
		return buildError(r)
	}
	return xml.NewDecoder(r.Body).Decode(resp)
}

func (iam *IAM) postQuery(params map[string]string, resp interface{}) error {
	endpoint, err := url.Parse(iam.IAMEndpoint)
	if err != nil {
		return err
	}
	params["Version"] = "2010-05-08"
	params["Timestamp"] = time.Now().In(time.UTC).Format(time.RFC3339)
	sign(iam.Auth, "POST", "/", params, endpoint.Host)
	encoded := multimap(params).Encode()
	body := strings.NewReader(encoded)
	req, err := http.NewRequest("POST", endpoint.String(), body)
	if err != nil {
		return err
	}
	req.Header.Set("Host", endpoint.Host)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(encoded)))
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode > 200 {
		return buildError(r)
	}
	return xml.NewDecoder(r.Body).Decode(resp)
}

func buildError(r *http.Response) error {
	var (
		err    Error
		errors xmlErrors
	)
	xml.NewDecoder(r.Body).Decode(&errors)
	if len(errors.Errors) > 0 {
		err = errors.Errors[0]
	}
	err.StatusCode = r.StatusCode
	if err.Message == "" {
		err.Message = r.Status
	}
	return &err
}

func multimap(p map[string]string) url.Values {
	q := make(url.Values, len(p))
	for k, v := range p {
		q[k] = []string{v}
	}
	return q
}

// Response to a CreateUser request.
//
// See http://goo.gl/JS9Gz for more details.
type CreateUserResp struct {
	RequestId string `xml:"ResponseMetadata>RequestId"`
	User      User   `xml:"CreateUserResult>User"`
}

// User encapsulates a user managed by IAM.
//
// See http://goo.gl/BwIQ3 for more details.
type User struct {
	Arn  string
	Path string
	Id   string `xml:"UserId"`
	Name string `xml:"UserName"`
}

// CreateUser creates a new user in IAM.
//
// See http://goo.gl/JS9Gz for more details.
func (iam *IAM) CreateUser(name, path string) (*CreateUserResp, error) {
	params := map[string]string{
		"Action":   "CreateUser",
		"Path":     path,
		"UserName": name,
	}
	resp := new(CreateUserResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a GetUser request.
//
// See http://goo.gl/ZnzRN for more details.
type GetUserResp struct {
	RequestId string `xml:"ResponseMetadata>RequestId"`
	User      User   `xml:"GetUserResult>User"`
}

// GetUser gets a user from IAM.
//
// See http://goo.gl/ZnzRN for more details.
func (iam *IAM) GetUser(name string) (*GetUserResp, error) {
	params := map[string]string{
		"Action": "GetUser",
	}
	if name != "" {
		params["UserName"] = name
	}
	resp := new(GetUserResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a GetGroup request.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_GetGroupResult.html for more details.
type GetGroupResp struct {
	Group       Group  `xml:"GetGroupResult>Group"`
	IsTruncated bool   `xml:"GetGroupResult>IsTruncated"`
	Marker      string `xml:"GetGroupResult>Marker"`
	Users       []User `xml:"GetGroupResult>Users>member"`
	RequestId   string `xml:"ResponseMetadata>RequestId"`
}

// GetGroup gets a group from IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_GetGroup.html for more details.
func (iam *IAM) GetGroup(groupName, marker string, maxItems int) (*GetGroupResp, error) {
	params := map[string]string{
		"Action":    "GetGroup",
		"GroupName": groupName,
	}
	if marker != "" {
		params["Marker"] = marker
	}
	if maxItems != 0 {
		params["MaxItems"] = strconv.Itoa(maxItems)
	}
	resp := new(GetGroupResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a GetGroupPolicy request.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_GetGroupPolicyResult.html for more details.
type GetGroupPolicyResp struct {
	GroupName      string `xml:"GetGroupPolicyResult>GroupName"`
	PolicyDocument string `xml:"GetGroupPolicyResult>PolicyDocument"`
	PolicyName     string `xml:"GetGroupPolicyResult>PolicyName"`
	RequestId      string `xml:"ResponseMetadata>RequestId"`
}

// GetGroupPolicy gets a policy for a group from IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_GetGroupPolicy.html for more details.
func (iam *IAM) GetGroupPolicy(groupName, policyName string) (*GetGroupPolicyResp, error) {
	params := map[string]string{
		"Action":     "GetGroupPolicy",
		"GroupName":  groupName,
		"PolicyName": policyName,
	}
	resp := new(GetGroupPolicyResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Role encapsulates a role managed by IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_Role.html for more details.
type Role struct {
	Arn                      string
	AssumeRolePolicyDocument string
	CreateDate               string
	Path                     string
	Id                       string `xml:"RoleId"`
	Name                     string `xml:"RoleName"`
}

// InstanceProfile encapsulates an instance profile managed by IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_InstanceProfile.html for more details.
type InstanceProfile struct {
	Arn        string
	CreateDate string
	Id         string `xml:"InstanceProfileId"`
	Name       string `xml:"InstanceProfileName"`
	Path       string
	Roles      []Role `xml:"Roles>member"`
}

// Response to a GetInstanceProfile request.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_GetInstanceProfileResult.html for more details.
type GetInstanceProfileResp struct {
	Profile   InstanceProfile `xml:"GetInstanceProfileResult>InstanceProfile"`
	RequestId string          `xml:"ResponseMetadata>RequestId"`
}

// GetInstanceProfile gets a profile for an instance from IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_GetInstanceProfile.html for more details.
func (iam *IAM) GetInstanceProfile(instanceProfileName string) (*GetInstanceProfileResp, error) {
	params := map[string]string{
		"Action":              "GetInstanceProfile",
		"InstanceProfileName": instanceProfileName,
	}
	resp := new(GetInstanceProfileResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a GetRole request.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_GetRoleResult.html for more details.
type GetRoleResp struct {
	Role      Role   `xml:"GetRoleResult>Role"`
	RequestId string `xml:"ResponseMetadata>RequestId"`
}

// GetRole gets a role from IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_GetRole.html for more details.
func (iam *IAM) GetRole(roleName string) (*GetRoleResp, error) {
	params := map[string]string{
		"Action":   "GetRole",
		"RoleName": roleName,
	}
	resp := new(GetRoleResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a GetRolePolicy request.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_GetRolePolicyResult.html for more details.
type GetRolePolicyResp struct {
	PolicyDocument string `xml:"GetRolePolicyResult>PolicyDocument"`
	PolicyName     string `xml:"GetRolePolicyResult>PolicyName"`
	RoleName       string `xml:"GetRolePolicyResult>RoleName"`
	RequestId      string `xml:"ResponseMetadata>RequestId"`
}

// GetRolePolicy gets a policy got a group from IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_GetRolePolicy.html for more details.
func (iam *IAM) GetRolePolicy(roleName, policyName string) (*GetRolePolicyResp, error) {
	params := map[string]string{
		"Action":     "GetRolePolicy",
		"RoleName":   roleName,
		"PolicyName": policyName,
	}
	resp := new(GetRolePolicyResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a AccountAliases request.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListAccountAliasesResult.html for more details.
type AccountAliasesResp struct {
	Aliases     []string `xml:"ListAccountAliasesResult>AccountAliases>member"`
	IsTruncated bool     `xml:"ListAccountAliasesResult>IsTruncated"`
	Marker      string   `xml:"ListAccountAliasesResult>Marker"`
	RequestId   string   `xml:"ResponseMetadata>RequestId"`
}

// AccountAliases list the associated aliases with the account from IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListAccountAliases.html for more details.
func (iam *IAM) AccountAliases(marker string, maxItems int) (*AccountAliasesResp, error) {
	params := map[string]string{
		"Action": "ListAccountAliases",
	}
	if marker != "" {
		params["Marker"] = marker
	}
	if maxItems != 0 {
		params["MaxItems"] = strconv.Itoa(maxItems)
	}
	resp := new(AccountAliasesResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a GroupPolicies request.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListGroupPoliciesResult.html for more details.
type GroupPoliciesResp struct {
	Names       []string `xml:"ListGroupPoliciesResult>PolicyNames>member"`
	IsTruncated bool     `xml:"ListGroupPoliciesResult>IsTruncated"`
	Marker      string   `xml:"ListGroupPoliciesResult>Marker"`
	RequestId   string   `xml:"ResponseMetadata>RequestId"`
}

// GroupPolicies gets the policies associated with a group from IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListGroupPolicies.html for more details.
func (iam *IAM) GroupPolicies(groupName, marker string, maxItems int) (*GroupPoliciesResp, error) {
	params := map[string]string{
		"Action":    "ListGroupPolicies",
		"GroupName": groupName,
	}
	if marker != "" {
		params["Marker"] = marker
	}
	if maxItems != 0 {
		params["MaxItems"] = strconv.Itoa(maxItems)
	}
	resp := new(GroupPoliciesResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a GroupsForUser request.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListGroupsForUserResult.html for more details.
type GroupsForUserResp struct {
	Groups      []Group `xml:"ListGroupsForUserResult>Groups>member"`
	IsTruncated bool    `xml:"ListGroupsForUserResult>IsTruncated"`
	Marker      string  `xml:"ListGroupsForUserResult>Marker"`
	RequestId   string  `xml:"ResponseMetadata>RequestId"`
}

// GroupsForUser gets the groups a user belongs from IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListGroupsForUser.html for more details.
func (iam *IAM) GroupsForUser(userName, marker string, maxItems int) (*GroupsForUserResp, error) {
	params := map[string]string{
		"Action":   "ListGroupsForUser",
		"UserName": userName,
	}
	if marker != "" {
		params["Marker"] = marker
	}
	if maxItems != 0 {
		params["MaxItems"] = strconv.Itoa(maxItems)
	}
	resp := new(GroupsForUserResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a InstanceProfiles request.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListInstanceProfilesResult.html for more details.
type InstanceProfilesResp struct {
	Profiles    []InstanceProfile `xml:"ListInstanceProfilesResult>InstanceProfiles>member"`
	IsTruncated bool              `xml:"ListInstanceProfilesResult>IsTruncated"`
	Marker      string            `xml:"ListInstanceProfilesResult>Marker"`
	RequestId   string            `xml:"ResponseMetadata>RequestId"`
}

// InstanceProfiles gets the instance profiles from IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListInstanceProfiles.html for more details.
func (iam *IAM) InstanceProfiles(pathPrefix, marker string, maxItems int) (*InstanceProfilesResp, error) {
	params := map[string]string{
		"Action": "ListInstanceProfiles",
	}
	if pathPrefix != "" {
		params["PathPrefix"] = pathPrefix
	}
	if marker != "" {
		params["Marker"] = marker
	}
	if maxItems != 0 {
		params["MaxItems"] = strconv.Itoa(maxItems)
	}
	resp := new(InstanceProfilesResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a InstanceProfilesForRole request.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListInstanceProfilesForRoleResult.html for more details.
type InstanceProfilesForRoleResp struct {
	Profiles    []InstanceProfile `xml:"ListInstanceProfilesForRoleResult>InstanceProfiles>member"`
	IsTruncated bool              `xml:"ListInstanceProfilesForRoleResult>IsTruncated"`
	Marker      string            `xml:"ListInstanceProfilesForRoleResult>Marker"`
	RequestId   string            `xml:"ResponseMetadata>RequestId"`
}

// InstanceProfilesForRole gets the instance profiles for a role from IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListInstanceProfilesForRole.html for more details.
func (iam *IAM) InstanceProfilesForRole(roleName, marker string, maxItems int) (*InstanceProfilesForRoleResp, error) {
	params := map[string]string{
		"Action":   "ListInstanceProfilesForRole",
		"RoleName": roleName,
	}
	if marker != "" {
		params["Marker"] = marker
	}
	if maxItems != 0 {
		params["MaxItems"] = strconv.Itoa(maxItems)
	}
	resp := new(InstanceProfilesForRoleResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a RolePolicies request.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListRolePoliciesResult.html for more details.
type RolePoliciesResp struct {
	Names       []string `xml:"ListRolePoliciesResult>PolicyNames>member"`
	IsTruncated bool     `xml:"ListGroupPoliciesResult>IsTruncated"`
	Marker      string   `xml:"ListGroupPoliciesResult>Marker"`
	RequestId   string   `xml:"ResponseMetadata>RequestId"`
}

// RolePolicies gets the policies associated with a role from IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListRolePolicies.html for more details.
func (iam *IAM) RolePolicies(roleName, marker string, maxItems int) (*RolePoliciesResp, error) {
	params := map[string]string{
		"Action":   "ListRolePolicies",
		"RoleName": roleName,
	}
	if marker != "" {
		params["Marker"] = marker
	}
	if maxItems != 0 {
		params["MaxItems"] = strconv.Itoa(maxItems)
	}
	resp := new(RolePoliciesResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a Roles request.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListRolesResult.html for more details.
type RolesResp struct {
	Roles       []Role `xml:"ListRolesResult>Roles>member"`
	IsTruncated bool   `xml:"ListRolesResult>IsTruncated"`
	Marker      string `xml:"ListRolesResult>Marker"`
	RequestId   string `xml:"ResponseMetadata>RequestId"`
}

// Roles gets the roles from IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListRoles.html for more details.
func (iam *IAM) Roles(pathPrefix, marker string, maxItems int) (*RolesResp, error) {
	params := map[string]string{
		"Action": "ListRoles",
	}
	if pathPrefix != "" {
		params["PathPrefix"] = pathPrefix
	}
	if marker != "" {
		params["Marker"] = marker
	}
	if maxItems != 0 {
		params["MaxItems"] = strconv.Itoa(maxItems)
	}
	resp := new(RolesResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteUser deletes a user from IAM.
//
// See http://goo.gl/jBuCG for more details.
func (iam *IAM) DeleteUser(name string) (*SimpleResp, error) {
	params := map[string]string{
		"Action":   "DeleteUser",
		"UserName": name,
	}
	resp := new(SimpleResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a CreateGroup request.
//
// See http://goo.gl/n7NNQ for more details.
type CreateGroupResp struct {
	Group     Group  `xml:"CreateGroupResult>Group"`
	RequestId string `xml:"ResponseMetadata>RequestId"`
}

// Group encapsulates a group managed by IAM.
//
// See http://goo.gl/ae7Vs for more details.
type Group struct {
	Arn  string
	Id   string `xml:"GroupId"`
	Name string `xml:"GroupName"`
	Path string
}

// CreateGroup creates a new group in IAM.
//
// The path parameter can be used to identify which division or part of the
// organization the user belongs to.
//
// If path is unset ("") it defaults to "/".
//
// See http://goo.gl/n7NNQ for more details.
func (iam *IAM) CreateGroup(name string, path string) (*CreateGroupResp, error) {
	params := map[string]string{
		"Action":    "CreateGroup",
		"GroupName": name,
	}
	if path != "" {
		params["Path"] = path
	}
	resp := new(CreateGroupResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a ListGroups request.
//
// See http://goo.gl/W2TRj for more details.
type GroupsResp struct {
	Groups    []Group `xml:"ListGroupsResult>Groups>member"`
	RequestId string  `xml:"ResponseMetadata>RequestId"`
}

// Groups list the groups that have the specified path prefix.
//
// The parameter pathPrefix is optional. If pathPrefix is "", all groups are
// returned.
//
// See http://goo.gl/W2TRj for more details.
func (iam *IAM) Groups(pathPrefix string) (*GroupsResp, error) {
	params := map[string]string{
		"Action": "ListGroups",
	}
	if pathPrefix != "" {
		params["PathPrefix"] = pathPrefix
	}
	resp := new(GroupsResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a UserPolicies request.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListUserPoliciesResult.html for more details.
type UserPoliciesResp struct {
	Names       []string `xml:"ListUserPoliciesResult>PolicyNames>member"`
	IsTruncated bool     `xml:"ListUserPoliciesResult>IsTruncated"`
	Marker      string   `xml:"ListUserPoliciesResult>Marker"`
	RequestId   string   `xml:"ResponseMetadata>RequestId"`
}

// UserPolicies gets the policies associated with a user from IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListUserPolicies.html for more details.
func (iam *IAM) UserPolicies(userName, marker string, maxItems int) (*UserPoliciesResp, error) {
	params := map[string]string{
		"Action":   "ListUserPolicies",
		"UserName": userName,
	}
	if marker != "" {
		params["Marker"] = marker
	}
	if maxItems != 0 {
		params["MaxItems"] = strconv.Itoa(maxItems)
	}
	resp := new(UserPoliciesResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a Users request.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListUsersResult.html for more details.
type UsersResp struct {
	Users       []User `xml:"ListUsersResult>Users>member"`
	IsTruncated bool   `xml:"ListUsersResult>IsTruncated"`
	Marker      string `xml:"ListUsersResult>Marker"`
	RequestId   string `xml:"ResponseMetadata>RequestId"`
}

// Users gets the users from IAM.
//
// See http://docs.aws.amazon.com/IAM/latest/APIReference/API_ListUsers.html for more details.
func (iam *IAM) Users(pathPrefix, marker string, maxItems int) (*UsersResp, error) {
	params := map[string]string{
		"Action": "ListUsers",
	}
	if pathPrefix != "" {
		params["PathPrefix"] = pathPrefix
	}
	if marker != "" {
		params["Marker"] = marker
	}
	if maxItems != 0 {
		params["MaxItems"] = strconv.Itoa(maxItems)
	}
	resp := new(UsersResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteGroup deletes a group from IAM.
//
// See http://goo.gl/d5i2i for more details.
func (iam *IAM) DeleteGroup(name string) (*SimpleResp, error) {
	params := map[string]string{
		"Action":    "DeleteGroup",
		"GroupName": name,
	}
	resp := new(SimpleResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a CreateAccessKey request.
//
// See http://goo.gl/L46Py for more details.
type CreateAccessKeyResp struct {
	RequestId string    `xml:"ResponseMetadata>RequestId"`
	AccessKey AccessKey `xml:"CreateAccessKeyResult>AccessKey"`
}

// AccessKey encapsulates an access key generated for a user.
//
// See http://goo.gl/LHgZR for more details.
type AccessKey struct {
	UserName string
	Id       string `xml:"AccessKeyId"`
	Secret   string `xml:"SecretAccessKey,omitempty"`
	Status   string
}

// CreateAccessKey creates a new access key in IAM.
//
// See http://goo.gl/L46Py for more details.
func (iam *IAM) CreateAccessKey(userName string) (*CreateAccessKeyResp, error) {
	params := map[string]string{
		"Action":   "CreateAccessKey",
		"UserName": userName,
	}
	resp := new(CreateAccessKeyResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a AccessKeys request.
//
// See http://goo.gl/Vjozx for more details.
type AccessKeysResp struct {
	RequestId  string      `xml:"ResponseMetadata>RequestId"`
	AccessKeys []AccessKey `xml:"ListAccessKeysResult>AccessKeyMetadata>member"`
}

// AccessKeys lists all acccess keys associated with a user.
//
// The userName parameter is optional. If set to "", the userName is determined
// implicitly based on the AWS Access Key ID used to sign the request.
//
// See http://goo.gl/Vjozx for more details.
func (iam *IAM) AccessKeys(userName string) (*AccessKeysResp, error) {
	params := map[string]string{
		"Action": "ListAccessKeys",
	}
	if userName != "" {
		params["UserName"] = userName
	}
	resp := new(AccessKeysResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteAccessKey deletes an access key from IAM.
//
// The userName parameter is optional. If set to "", the userName is determined
// implicitly based on the AWS Access Key ID used to sign the request.
//
// See http://goo.gl/hPGhw for more details.
func (iam *IAM) DeleteAccessKey(id, userName string) (*SimpleResp, error) {
	params := map[string]string{
		"Action":      "DeleteAccessKey",
		"AccessKeyId": id,
	}
	if userName != "" {
		params["UserName"] = userName
	}
	resp := new(SimpleResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Response to a GetUserPolicy request.
//
// See http://goo.gl/BH04O for more details.
type GetUserPolicyResp struct {
	Policy    UserPolicy `xml:"GetUserPolicyResult"`
	RequestId string     `xml:"ResponseMetadata>RequestId"`
}

// UserPolicy encapsulates an IAM group policy.
//
// See http://goo.gl/C7hgS for more details.
type UserPolicy struct {
	Name     string `xml:"PolicyName"`
	UserName string `xml:"UserName"`
	Document string `xml:"PolicyDocument"`
}

// GetUserPolicy gets a user policy in IAM.
//
// See http://goo.gl/BH04O for more details.
func (iam *IAM) GetUserPolicy(userName, policyName string) (*GetUserPolicyResp, error) {
	params := map[string]string{
		"Action":     "GetUserPolicy",
		"UserName":   userName,
		"PolicyName": policyName,
	}
	resp := new(GetUserPolicyResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
	return nil, nil
}

// PutUserPolicy creates a user policy in IAM.
//
// See http://goo.gl/ldCO8 for more details.
func (iam *IAM) PutUserPolicy(userName, policyName, policyDocument string) (*SimpleResp, error) {
	params := map[string]string{
		"Action":         "PutUserPolicy",
		"UserName":       userName,
		"PolicyName":     policyName,
		"PolicyDocument": policyDocument,
	}
	resp := new(SimpleResp)
	if err := iam.postQuery(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteUserPolicy deletes a user policy from IAM.
//
// See http://goo.gl/7Jncn for more details.
func (iam *IAM) DeleteUserPolicy(userName, policyName string) (*SimpleResp, error) {
	params := map[string]string{
		"Action":     "DeleteUserPolicy",
		"PolicyName": policyName,
		"UserName":   userName,
	}
	resp := new(SimpleResp)
	if err := iam.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type SimpleResp struct {
	RequestId string `xml:"ResponseMetadata>RequestId"`
}

type xmlErrors struct {
	Errors []Error `xml:"Error"`
}

// Error encapsulates an IAM error.
type Error struct {
	// HTTP status code of the error.
	StatusCode int

	// AWS code of the error.
	Code string

	// Message explaining the error.
	Message string
}

func (e *Error) Error() string {
	var prefix string
	if e.Code != "" {
		prefix = e.Code + ": "
	}
	if prefix == "" && e.StatusCode > 0 {
		prefix = strconv.Itoa(e.StatusCode) + ": "
	}
	return prefix + e.Message
}
