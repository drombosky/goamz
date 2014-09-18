package iam_test

import (
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/iam"
	"github.com/crowdmob/goamz/testutil"
	"gopkg.in/check.v1"
	"strconv"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct {
	iam *iam.IAM
}

var _ = check.Suite(&S{})

var testServer = testutil.NewHTTPServer()

func (s *S) SetUpSuite(c *check.C) {
	testServer.Start()
	auth := aws.Auth{AccessKey: "abc", SecretKey: "123"}
	s.iam = iam.New(auth, aws.Region{IAMEndpoint: testServer.URL})
}

func (s *S) TearDownTest(c *check.C) {
	testServer.Flush()
}

func (s *S) TestCreateUser(c *check.C) {
	testServer.Response(200, nil, CreateUserExample)
	resp, err := s.iam.CreateUser("Bob", "/division_abc/subdivision_xyz/")
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "CreateUser")
	c.Assert(values.Get("UserName"), check.Equals, "Bob")
	c.Assert(values.Get("Path"), check.Equals, "/division_abc/subdivision_xyz/")
	c.Assert(err, check.IsNil)
	c.Assert(resp.RequestId, check.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
	expected := iam.User{
		Path: "/division_abc/subdivision_xyz/",
		Name: "Bob",
		Id:   "AIDACKCEVSQ6C2EXAMPLE",
		Arn:  "arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/Bob",
	}
	c.Assert(resp.User, check.DeepEquals, expected)
}

func (s *S) TestCreateUserConflict(c *check.C) {
	testServer.Response(409, nil, DuplicateUserExample)
	resp, err := s.iam.CreateUser("Bob", "/division_abc/subdivision_xyz/")
	testServer.WaitRequest()
	c.Assert(resp, check.IsNil)
	c.Assert(err, check.NotNil)
	e, ok := err.(*iam.Error)
	c.Assert(ok, check.Equals, true)
	c.Assert(e.Message, check.Equals, "User with name Bob already exists.")
	c.Assert(e.Code, check.Equals, "EntityAlreadyExists")
}

func (s *S) TestGetUser(c *check.C) {
	testServer.Response(200, nil, GetUserExample)
	resp, err := s.iam.GetUser("Bob")
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "GetUser")
	c.Assert(values.Get("UserName"), check.Equals, "Bob")
	c.Assert(err, check.IsNil)
	c.Assert(resp.RequestId, check.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
	expected := iam.User{
		Path: "/division_abc/subdivision_xyz/",
		Name: "Bob",
		Id:   "AIDACKCEVSQ6C2EXAMPLE",
		Arn:  "arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/Bob",
	}
	c.Assert(resp.User, check.DeepEquals, expected)
}

func (s *S) TestDeleteUser(c *check.C) {
	testServer.Response(200, nil, RequestIdExample)
	resp, err := s.iam.DeleteUser("Bob")
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "DeleteUser")
	c.Assert(values.Get("UserName"), check.Equals, "Bob")
	c.Assert(err, check.IsNil)
	c.Assert(resp.RequestId, check.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
}

func (s *S) TestCreateGroup(c *check.C) {
	testServer.Response(200, nil, CreateGroupExample)
	resp, err := s.iam.CreateGroup("Admins", "/admins/")
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "CreateGroup")
	c.Assert(values.Get("GroupName"), check.Equals, "Admins")
	c.Assert(values.Get("Path"), check.Equals, "/admins/")
	c.Assert(err, check.IsNil)
	c.Assert(resp.Group.Path, check.Equals, "/admins/")
	c.Assert(resp.Group.Name, check.Equals, "Admins")
	c.Assert(resp.Group.Id, check.Equals, "AGPACKCEVSQ6C2EXAMPLE")
	c.Assert(resp.RequestId, check.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
}

func (s *S) TestCreateGroupWithoutPath(c *check.C) {
	testServer.Response(200, nil, CreateGroupExample)
	_, err := s.iam.CreateGroup("Managers", "")
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "CreateGroup")
	c.Assert(err, check.IsNil)
	_, ok := map[string][]string(values)["Path"]
	c.Assert(ok, check.Equals, false)
}

func (s *S) TestDeleteGroup(c *check.C) {
	testServer.Response(200, nil, RequestIdExample)
	resp, err := s.iam.DeleteGroup("Admins")
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "DeleteGroup")
	c.Assert(values.Get("GroupName"), check.Equals, "Admins")
	c.Assert(err, check.IsNil)
	c.Assert(resp.RequestId, check.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
}

func (s *S) TestListGroups(c *check.C) {
	testServer.Response(200, nil, ListGroupsExample)
	resp, err := s.iam.Groups("/division_abc/")
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "ListGroups")
	c.Assert(values.Get("PathPrefix"), check.Equals, "/division_abc/")
	c.Assert(err, check.IsNil)
	c.Assert(resp.RequestId, check.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
	expected := []iam.Group{
		{
			Path: "/division_abc/subdivision_xyz/",
			Name: "Admins",
			Id:   "AGPACKCEVSQ6C2EXAMPLE",
			Arn:  "arn:aws:iam::123456789012:group/Admins",
		},
		{
			Path: "/division_abc/subdivision_xyz/product_1234/engineering/",
			Name: "Test",
			Id:   "AGP2MAB8DPLSRHEXAMPLE",
			Arn:  "arn:aws:iam::123456789012:group/division_abc/subdivision_xyz/product_1234/engineering/Test",
		},
		{
			Path: "/division_abc/subdivision_xyz/product_1234/",
			Name: "Managers",
			Id:   "AGPIODR4TAW7CSEXAMPLE",
			Arn:  "arn:aws:iam::123456789012:group/division_abc/subdivision_xyz/product_1234/Managers",
		},
	}
	c.Assert(resp.Groups, check.DeepEquals, expected)
}

func (s *S) TestListGroupsWithoutPathPrefix(c *check.C) {
	testServer.Response(200, nil, ListGroupsExample)
	_, err := s.iam.Groups("")
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "ListGroups")
	c.Assert(err, check.IsNil)
	_, ok := map[string][]string(values)["PathPrefix"]
	c.Assert(ok, check.Equals, false)
}

func (s *S) TestCreateAccessKey(c *check.C) {
	testServer.Response(200, nil, CreateAccessKeyExample)
	resp, err := s.iam.CreateAccessKey("Bob")
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "CreateAccessKey")
	c.Assert(values.Get("UserName"), check.Equals, "Bob")
	c.Assert(err, check.IsNil)
	c.Assert(resp.AccessKey.UserName, check.Equals, "Bob")
	c.Assert(resp.AccessKey.Id, check.Equals, "AKIAIOSFODNN7EXAMPLE")
	c.Assert(resp.AccessKey.Secret, check.Equals, "wJalrXUtnFEMI/K7MDENG/bPxRfiCYzEXAMPLEKEY")
	c.Assert(resp.AccessKey.Status, check.Equals, "Active")
}

func (s *S) TestDeleteAccessKey(c *check.C) {
	testServer.Response(200, nil, RequestIdExample)
	resp, err := s.iam.DeleteAccessKey("ysa8hasdhasdsi", "Bob")
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "DeleteAccessKey")
	c.Assert(values.Get("AccessKeyId"), check.Equals, "ysa8hasdhasdsi")
	c.Assert(values.Get("UserName"), check.Equals, "Bob")
	c.Assert(err, check.IsNil)
	c.Assert(resp.RequestId, check.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
}

func (s *S) TestDeleteAccessKeyBlankUserName(c *check.C) {
	testServer.Response(200, nil, RequestIdExample)
	_, err := s.iam.DeleteAccessKey("ysa8hasdhasdsi", "")
	c.Assert(err, check.IsNil)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "DeleteAccessKey")
	c.Assert(values.Get("AccessKeyId"), check.Equals, "ysa8hasdhasdsi")
	_, ok := map[string][]string(values)["UserName"]
	c.Assert(ok, check.Equals, false)
}

func (s *S) TestAccessKeys(c *check.C) {
	testServer.Response(200, nil, ListAccessKeyExample)
	resp, err := s.iam.AccessKeys("Bob")
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "ListAccessKeys")
	c.Assert(values.Get("UserName"), check.Equals, "Bob")
	c.Assert(err, check.IsNil)
	c.Assert(resp.RequestId, check.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
	c.Assert(resp.AccessKeys, check.HasLen, 2)
	c.Assert(resp.AccessKeys[0].Id, check.Equals, "AKIAIOSFODNN7EXAMPLE")
	c.Assert(resp.AccessKeys[0].UserName, check.Equals, "Bob")
	c.Assert(resp.AccessKeys[0].Status, check.Equals, "Active")
	c.Assert(resp.AccessKeys[1].Id, check.Equals, "AKIAI44QH8DHBEXAMPLE")
	c.Assert(resp.AccessKeys[1].UserName, check.Equals, "Bob")
	c.Assert(resp.AccessKeys[1].Status, check.Equals, "Inactive")
}

func (s *S) TestAccessKeysBlankUserName(c *check.C) {
	testServer.Response(200, nil, ListAccessKeyExample)
	_, err := s.iam.AccessKeys("")
	c.Assert(err, check.IsNil)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "ListAccessKeys")
	_, ok := map[string][]string(values)["UserName"]
	c.Assert(ok, check.Equals, false)
}

func (s *S) TestGetUserPolicy(c *check.C) {
	testServer.Response(200, nil, GetUserPolicyExample)
	resp, err := s.iam.GetUserPolicy("Bob", "AllAccessPolicy")
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "GetUserPolicy")
	c.Assert(values.Get("UserName"), check.Equals, "Bob")
	c.Assert(values.Get("PolicyName"), check.Equals, "AllAccessPolicy")
	c.Assert(err, check.IsNil)
	c.Assert(resp.Policy.UserName, check.Equals, "Bob")
	c.Assert(resp.Policy.Name, check.Equals, "AllAccessPolicy")
	c.Assert(strings.TrimSpace(resp.Policy.Document), check.Equals, `{"Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}`)
	c.Assert(resp.RequestId, check.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
}

func (s *S) TestPutUserPolicy(c *check.C) {
	document := `{
		"Statement": [
		{
			"Action": [
				"s3:*"
			],
			"Effect": "Allow",
			"Resource": [
				"arn:aws:s3:::8shsns19s90ajahadsj/*",
				"arn:aws:s3:::8shsns19s90ajahadsj"
			]
		}]
	}`
	testServer.Response(200, nil, RequestIdExample)
	resp, err := s.iam.PutUserPolicy("Bob", "AllAccessPolicy", document)
	req := testServer.WaitRequest()
	c.Assert(req.Method, check.Equals, "POST")
	c.Assert(req.FormValue("Action"), check.Equals, "PutUserPolicy")
	c.Assert(req.FormValue("PolicyName"), check.Equals, "AllAccessPolicy")
	c.Assert(req.FormValue("UserName"), check.Equals, "Bob")
	c.Assert(req.FormValue("PolicyDocument"), check.Equals, document)
	c.Assert(req.FormValue("Version"), check.Equals, "2010-05-08")
	c.Assert(err, check.IsNil)
	c.Assert(resp.RequestId, check.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
}

func (s *S) TestDeleteUserPolicy(c *check.C) {
	testServer.Response(200, nil, RequestIdExample)
	resp, err := s.iam.DeleteUserPolicy("Bob", "AllAccessPolicy")
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "DeleteUserPolicy")
	c.Assert(values.Get("PolicyName"), check.Equals, "AllAccessPolicy")
	c.Assert(values.Get("UserName"), check.Equals, "Bob")
	c.Assert(err, check.IsNil)
	c.Assert(resp.RequestId, check.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
}

func (s *S) TestGetGroup(c *check.C) {
	testServer.Response(200, nil, GetGroupExample)
	expected := iam.GetGroupResp{
		Group: iam.Group{
			Arn:  "arn:aws:iam::123456789012:group/Admins",
			Id:   "AGPACKCEVSQ6C2EXAMPLE",
			Name: "Admins",
			Path: "/",
		},
		IsTruncated: false,
		Marker:      "",
		Users: []iam.User{
			iam.User{
				Arn:  "arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/Bob",
				Path: "/division_abc/subdivision_xyz/",
				Id:   "AIDACKCEVSQ6C2EXAMPLE",
				Name: "Bob",
			},
			iam.User{
				Arn:  "arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/Susan",
				Path: "/division_abc/subdivision_xyz/",
				Id:   "AIDACKCEVSQ6C2EXAMPLE",
				Name: "Susan",
			},
		},
		RequestId: "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE",
	}
	groupName := "Admins"
	marker := "ignore"
	maxItems := 10
	resp, err := s.iam.GetGroup(groupName, marker, maxItems)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "GetGroup")
	c.Assert(values.Get("GroupName"), check.Equals, groupName)
	c.Assert(values.Get("Marker"), check.Equals, marker)
	c.Assert(values.Get("MaxItems"), check.Equals, strconv.Itoa(maxItems))
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, expected)
}

func (s *S) TestGetGroupPolicy(c *check.C) {
	testServer.Response(200, nil, GetGroupPolicyExample)
	expected := iam.GetGroupPolicyResp{
		GroupName:      "Admins",
		PolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}`,
		PolicyName:     "AdminRoot",
		RequestId:      "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE",
	}
	groupName := "Admins"
	policyName := "AdminRoot"
	resp, err := s.iam.GetGroupPolicy(groupName, policyName)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "GetGroupPolicy")
	c.Assert(values.Get("GroupName"), check.Equals, groupName)
	c.Assert(values.Get("PolicyName"), check.Equals, policyName)
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, expected)
}

func (s *S) TestGetInstanceProfile(c *check.C) {
	testServer.Response(200, nil, GetInstanceProfileExample)
	expected := iam.GetInstanceProfileResp{
		Profile: iam.InstanceProfile{
			Arn:        "arn:aws:iam::123456789012:instance-profile/application_abc/component_xyz/Webserver",
			CreateDate: "2012-05-09T16:11:10Z",
			Id:         "AIPAD5ARO2C5EXAMPLE3G",
			Name:       "Webserver",
			Path:       "/application_abc/component_xyz/",
			Roles: []iam.Role{
				iam.Role{
					Arn: "arn:aws:iam::123456789012:role/application_abc/component_xyz/S3Access",
					AssumeRolePolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}`,
					CreateDate:               "2012-05-09T15:45:35Z",
					Path:                     "/application_abc/component_xyz/",
					Id:                       "AROACVYKSVTSZFEXAMPLE",
					Name:                     "S3Access",
				},
			},
		},
		RequestId: "37289fda-99f2-11e1-a4c3-27EXAMPLE804",
	}
	instanceProfileName := "Webserver"
	resp, err := s.iam.GetInstanceProfile(instanceProfileName)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "GetInstanceProfile")
	c.Assert(values.Get("InstanceProfileName"), check.Equals, instanceProfileName)
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, expected)
}

func (s *S) TestGetRole(c *check.C) {
	testServer.Response(200, nil, GetRoleExample)
	expected := iam.GetRoleResp{
		Role: iam.Role{
			Arn: "arn:aws:iam::123456789012:role/application_abc/component_xyz/S3Access",
			AssumeRolePolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}`,
			CreateDate:               "2012-05-08T23:34:01Z",
			Path:                     "/application_abc/component_xyz/",
			Id:                       "AROADBQP57FF2AEXAMPLE",
			Name:                     "S3Access",
		},
		RequestId: "df37e965-9967-11e1-a4c3-270EXAMPLE04",
	}
	roleName := "S3Access"
	resp, err := s.iam.GetRole(roleName)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "GetRole")
	c.Assert(values.Get("RoleName"), check.Equals, roleName)
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, expected)
}

func (s *S) TestGetRolePolocy(c *check.C) {
	testServer.Response(200, nil, GetRolePolicyExample)
	expected := iam.GetRolePolicyResp{
		PolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:*"],"Resource":["*"]}]}`,
		PolicyName:     "S3AccessPolicy",
		RoleName:       "S3Access",
		RequestId:      "7e7cd8bc-99ef-11e1-a4c3-27EXAMPLE804",
	}
	roleName := "S3Access"
	policyName := "S3AccessPolicy"
	resp, err := s.iam.GetRolePolicy(roleName, policyName)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "GetRolePolicy")
	c.Assert(values.Get("RoleName"), check.Equals, roleName)
	c.Assert(values.Get("PolicyName"), check.Equals, policyName)
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, expected)
}

func (s *S) TestAccountAliases(c *check.C) {
	testServer.Response(200, nil, AccountAliasesExample)
	expected := iam.AccountAliasesResp{
		Aliases:     []string{"foocorporation", "barcorporation"},
		IsTruncated: false,
		Marker:      "",
		RequestId:   "c5a076e9-f1b0-11df-8fbe-45274EXAMPLE",
	}
	marker := "ignore"
	maxItems := 10
	resp, err := s.iam.AccountAliases(marker, maxItems)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "ListAccountAliases")
	c.Assert(values.Get("Marker"), check.Equals, marker)
	c.Assert(values.Get("MaxItems"), check.Equals, strconv.Itoa(maxItems))
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, expected)
}

func (s *S) TestGroupPolicies(c *check.C) {
	testServer.Response(200, nil, GroupPoliciesExample)
	expected := iam.GroupPoliciesResp{
		Names:       []string{"AdminRoot", "KeyPolicy"},
		IsTruncated: false,
		Marker:      "",
		RequestId:   "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE",
	}
	groupName := "Admins"
	marker := "ignore"
	maxItems := 10
	resp, err := s.iam.GroupPolicies(groupName, marker, maxItems)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "ListGroupPolicies")
	c.Assert(values.Get("GroupName"), check.Equals, groupName)
	c.Assert(values.Get("Marker"), check.Equals, marker)
	c.Assert(values.Get("MaxItems"), check.Equals, strconv.Itoa(maxItems))
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, expected)
}

func (s *S) TestGroupsForUser(c *check.C) {
	testServer.Response(200, nil, GroupsForUserExample)
	expected := iam.GroupsForUserResp{
		Groups: []iam.Group{
			iam.Group{
				Arn:  "arn:aws:iam::123456789012:group/Admins",
				Id:   "AGPACKCEVSQ6C2EXAMPLE",
				Name: "Admins",
				Path: "/",
			},
		},
		IsTruncated: false,
		Marker:      "",
		RequestId:   "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE",
	}
	userName := "Bob"
	marker := "ignore"
	maxItems := 10
	resp, err := s.iam.GroupsForUser(userName, marker, maxItems)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "ListGroupsForUser")
	c.Assert(values.Get("UserName"), check.Equals, userName)
	c.Assert(values.Get("Marker"), check.Equals, marker)
	c.Assert(values.Get("MaxItems"), check.Equals, strconv.Itoa(maxItems))
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, expected)
}

func (s *S) TestInstanceProfiles(c *check.C) {
	testServer.Response(200, nil, InstanceProfilesExample)
	expected := iam.InstanceProfilesResp{
		Profiles: []iam.InstanceProfile{
			iam.InstanceProfile{
				Arn:        "arn:aws:iam::123456789012:instance-profile/application_abc/component_xyz/Database",
				CreateDate: "2012-05-09T16:27:03Z",
				Id:         "AIPACIFN4OZXG7EXAMPLE",
				Name:       "Database",
				Path:       "/application_abc/component_xyz/",
				Roles: []iam.Role{
					iam.Role{
						Arn: "arn:aws:iam::123456789012:role/application_abc/component_xyz/S3Access",
						AssumeRolePolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}`,
						CreateDate:               "2012-05-09T15:45:35Z",
						Path:                     "/application_abc/component_xyz/",
						Id:                       "AROACVSVTSZYK3EXAMPLE",
						Name:                     "S3Access",
					},
				},
			},
			iam.InstanceProfile{
				Arn:        "arn:aws:iam::123456789012:instance-profile/application_abc/component_xyz/Webserver",
				CreateDate: "2012-05-09T16:27:11Z",
				Id:         "AIPACZLSXM2EYYEXAMPLE",
				Name:       "Webserver",
				Path:       "/application_abc/component_xyz/",
				Roles:      nil,
			},
		},
		IsTruncated: false,
		Marker:      "",
		RequestId:   "fd74fa8d-99f3-11e1-a4c3-27EXAMPLE804",
	}
	pathPrefix := "/application_abc/"
	marker := "ignore"
	maxItems := 10
	resp, err := s.iam.InstanceProfiles(pathPrefix, marker, maxItems)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "ListInstanceProfiles")
	c.Assert(values.Get("PathPrefix"), check.Equals, pathPrefix)
	c.Assert(values.Get("Marker"), check.Equals, marker)
	c.Assert(values.Get("MaxItems"), check.Equals, strconv.Itoa(maxItems))
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, expected)
}

func (s *S) TestInstanceProfilesForRole(c *check.C) {
	testServer.Response(200, nil, InstanceProfilesForRoleExample)
	expected := iam.InstanceProfilesForRoleResp{
		Profiles: []iam.InstanceProfile{
			iam.InstanceProfile{
				Arn:        "arn:aws:iam::123456789012:instance-profile/application_abc/component_xyz/Webserver",
				CreateDate: "2012-05-09T16:27:11Z",
				Id:         "AIPACZLS2EYYXMEXAMPLE",
				Name:       "Webserver",
				Path:       "/application_abc/component_xyz/",
				Roles: []iam.Role{
					iam.Role{
						Arn: "arn:aws:iam::123456789012:role/application_abc/component_xyz/S3Access",
						AssumeRolePolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}`,
						CreateDate:               "2012-05-09T15:45:35Z",
						Path:                     "/application_abc/component_xyz/",
						Id:                       "AROACVSVTSZYK3EXAMPLE",
						Name:                     "S3Access",
					},
				},
			},
		},
		IsTruncated: false,
		Marker:      "",
		RequestId:   "6a8c3992-99f4-11e1-a4c3-27EXAMPLE804",
	}
	roleName := "S3Access"
	marker := "ignore"
	maxItems := 10
	resp, err := s.iam.InstanceProfilesForRole(roleName, marker, maxItems)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "ListInstanceProfilesForRole")
	c.Assert(values.Get("RoleName"), check.Equals, roleName)
	c.Assert(values.Get("Marker"), check.Equals, marker)
	c.Assert(values.Get("MaxItems"), check.Equals, strconv.Itoa(maxItems))
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, expected)
}

func (s *S) TestRolePolicies(c *check.C) {
	testServer.Response(200, nil, RolePoliciesExample)
	expected := iam.RolePoliciesResp{
		Names: []string{
			"CloudwatchPutMetricPolicy",
			"S3AccessPolicy",
		},
		IsTruncated: false,
		Marker:      "",
		RequestId:   "8c7e1816-99f0-11e1-a4c3-27EXAMPLE804",
	}
	roleName := "S3Access"
	marker := "ignore"
	maxItems := 10
	resp, err := s.iam.RolePolicies(roleName, marker, maxItems)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "ListRolePolicies")
	c.Assert(values.Get("RoleName"), check.Equals, roleName)
	c.Assert(values.Get("Marker"), check.Equals, marker)
	c.Assert(values.Get("MaxItems"), check.Equals, strconv.Itoa(maxItems))
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, expected)
}

func (s *S) TestRoles(c *check.C) {
	testServer.Response(200, nil, RolesExample)
	expected := iam.RolesResp{
		Roles: []iam.Role{
			iam.Role{
				Arn: "arn:aws:iam::123456789012:role/application_abc/component_xyz/S3Access",
				AssumeRolePolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}`,
				CreateDate:               "2012-05-09T15:45:35Z",
				Path:                     "/application_abc/component_xyz/",
				Id:                       "AROACVSVTSZYEXAMPLEYK",
				Name:                     "S3Access",
			},
			iam.Role{
				Arn: "arn:aws:iam::123456789012:role/application_abc/component_xyz/SDBAccess",
				AssumeRolePolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}`,
				CreateDate:               "2012-05-09T15:45:45Z",
				Path:                     "/application_abc/component_xyz/",
				Id:                       "AROAC2ICXG32EXAMPLEWK",
				Name:                     "SDBAccess",
			},
		},
		IsTruncated: false,
		Marker:      "",
		RequestId:   "20f7279f-99ee-11e1-a4c3-27EXAMPLE804",
	}
	pathPrefix := "/application_abc/"
	marker := "ignore"
	maxItems := 10
	resp, err := s.iam.Roles(pathPrefix, marker, maxItems)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "ListRoles")
	c.Assert(values.Get("PathPrefix"), check.Equals, pathPrefix)
	c.Assert(values.Get("Marker"), check.Equals, marker)
	c.Assert(values.Get("MaxItems"), check.Equals, strconv.Itoa(maxItems))
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, expected)
}

func (s *S) TestUserPolicies(c *check.C) {
	testServer.Response(200, nil, UserPoliciesExample)
	expected := iam.UserPoliciesResp{
		Names:       []string{"AllAccessPolicy", "KeyPolicy"},
		IsTruncated: false,
		Marker:      "",
		RequestId:   "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE",
	}
	userName := "Bob"
	marker := "ignore"
	maxItems := 10
	resp, err := s.iam.UserPolicies(userName, marker, maxItems)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "ListUserPolicies")
	c.Assert(values.Get("UserName"), check.Equals, userName)
	c.Assert(values.Get("Marker"), check.Equals, marker)
	c.Assert(values.Get("MaxItems"), check.Equals, strconv.Itoa(maxItems))
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, expected)
}

func (s *S) TestUsers(c *check.C) {
	testServer.Response(200, nil, UsersExample)
	expected := iam.UsersResp{
		Users: []iam.User{
			iam.User{
				Arn:  "arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/engineering/Andrew",
				Path: "/division_abc/subdivision_xyz/engineering/",
				Id:   "AID2MAB8DPLSRHEXAMPLE",
				Name: "Andrew",
			},
			iam.User{
				Arn:  "arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/engineering/Jackie",
				Path: "/division_abc/subdivision_xyz/engineering/",
				Id:   "AIDIODR4TAW7CSEXAMPLE",
				Name: "Jackie",
			},
		},
		IsTruncated: false,
		Marker:      "",
		RequestId:   "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE",
	}
	pathPrefix := "/division_abc/subdivision_xyz/product_1234/engineering/"
	marker := "ignore"
	maxItems := 10
	resp, err := s.iam.Users(pathPrefix, marker, maxItems)
	values := testServer.WaitRequest().URL.Query()
	c.Assert(values.Get("Action"), check.Equals, "ListUsers")
	c.Assert(values.Get("PathPrefix"), check.Equals, pathPrefix)
	c.Assert(values.Get("Marker"), check.Equals, marker)
	c.Assert(values.Get("MaxItems"), check.Equals, strconv.Itoa(maxItems))
	c.Assert(err, check.IsNil)
	c.Assert(*resp, check.DeepEquals, expected)
}
