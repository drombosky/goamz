package iam_test

// http://goo.gl/EUIvl
var CreateUserExample = `
<CreateUserResponse>
   <CreateUserResult>
      <User>
         <Path>/division_abc/subdivision_xyz/</Path>
         <UserName>Bob</UserName>
         <UserId>AIDACKCEVSQ6C2EXAMPLE</UserId>
         <Arn>arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/Bob</Arn>
     </User>
   </CreateUserResult>
   <ResponseMetadata>
      <RequestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</RequestId>
   </ResponseMetadata>
</CreateUserResponse>
`

var DuplicateUserExample = `
<ErrorResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
  <Error>
    <Type>Sender</Type>
    <Code>EntityAlreadyExists</Code>
    <Message>User with name Bob already exists.</Message>
  </Error>
  <RequestId>1d5f5000-1316-11e2-a60f-91a8e6fb6d21</RequestId>
</ErrorResponse>
`

var GetUserExample = `
<GetUserResponse>
   <GetUserResult>
      <User>
         <Path>/division_abc/subdivision_xyz/</Path>
         <UserName>Bob</UserName>
         <UserId>AIDACKCEVSQ6C2EXAMPLE</UserId>
         <Arn>arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/Bob</Arn>
      </User>
   </GetUserResult>
   <ResponseMetadata>
      <RequestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</RequestId>
   </ResponseMetadata>
</GetUserResponse>
`

var CreateGroupExample = `
<CreateGroupResponse>
   <CreateGroupResult>
      <Group>
         <Path>/admins/</Path>
         <GroupName>Admins</GroupName>
         <GroupId>AGPACKCEVSQ6C2EXAMPLE</GroupId>
         <Arn>arn:aws:iam::123456789012:group/Admins</Arn>
      </Group>
   </CreateGroupResult>
   <ResponseMetadata>
      <RequestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</RequestId>
   </ResponseMetadata>
</CreateGroupResponse>
`

var ListGroupsExample = `
<ListGroupsResponse>
   <ListGroupsResult>
      <Groups>
         <member>
            <Path>/division_abc/subdivision_xyz/</Path>
            <GroupName>Admins</GroupName>
            <GroupId>AGPACKCEVSQ6C2EXAMPLE</GroupId>
            <Arn>arn:aws:iam::123456789012:group/Admins</Arn>
         </member>
         <member>
            <Path>/division_abc/subdivision_xyz/product_1234/engineering/</Path>
            <GroupName>Test</GroupName>
            <GroupId>AGP2MAB8DPLSRHEXAMPLE</GroupId>
            <Arn>arn:aws:iam::123456789012:group/division_abc/subdivision_xyz/product_1234/engineering/Test</Arn>
         </member>
         <member>
            <Path>/division_abc/subdivision_xyz/product_1234/</Path>
            <GroupName>Managers</GroupName>
            <GroupId>AGPIODR4TAW7CSEXAMPLE</GroupId>
            <Arn>arn:aws:iam::123456789012:group/division_abc/subdivision_xyz/product_1234/Managers</Arn>
         </member>
      </Groups>
      <IsTruncated>false</IsTruncated>
   </ListGroupsResult>
   <ResponseMetadata>
      <RequestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</RequestId>
   </ResponseMetadata>
</ListGroupsResponse>
`

var RequestIdExample = `
<AddUserToGroupResponse>
   <ResponseMetadata>
      <RequestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</RequestId>
   </ResponseMetadata>
</AddUserToGroupResponse>
`

var CreateAccessKeyExample = `
<CreateAccessKeyResponse>
   <CreateAccessKeyResult>
     <AccessKey>
         <UserName>Bob</UserName>
         <AccessKeyId>AKIAIOSFODNN7EXAMPLE</AccessKeyId>
         <Status>Active</Status>
         <SecretAccessKey>wJalrXUtnFEMI/K7MDENG/bPxRfiCYzEXAMPLEKEY</SecretAccessKey>
      </AccessKey>
   </CreateAccessKeyResult>
   <ResponseMetadata>
      <RequestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</RequestId>
   </ResponseMetadata>
</CreateAccessKeyResponse>
`

var ListAccessKeyExample = `
<ListAccessKeysResponse>
   <ListAccessKeysResult>
      <UserName>Bob</UserName>
      <AccessKeyMetadata>
         <member>
            <UserName>Bob</UserName>
            <AccessKeyId>AKIAIOSFODNN7EXAMPLE</AccessKeyId>
            <Status>Active</Status>
         </member>
         <member>
            <UserName>Bob</UserName>
            <AccessKeyId>AKIAI44QH8DHBEXAMPLE</AccessKeyId>
            <Status>Inactive</Status>
         </member>
      </AccessKeyMetadata>
      <IsTruncated>false</IsTruncated>
   </ListAccessKeysResult>
   <ResponseMetadata>
      <RequestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</RequestId>
   </ResponseMetadata>
</ListAccessKeysResponse>
`

var GetUserPolicyExample = `
<GetUserPolicyResponse>
   <GetUserPolicyResult>
      <UserName>Bob</UserName>
      <PolicyName>AllAccessPolicy</PolicyName>
      <PolicyDocument>
      {"Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}
      </PolicyDocument>
   </GetUserPolicyResult>
   <ResponseMetadata>
      <RequestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</RequestId>
   </ResponseMetadata>
</GetUserPolicyResponse>
`

var GetGroupExample = `
<GetGroupResponse>
 <GetGroupResult>
    <Group>
       <Path>/</Path>
       <GroupName>Admins</GroupName>
       <GroupId>AGPACKCEVSQ6C2EXAMPLE</GroupId>
       <Arn>arn:aws:iam::123456789012:group/Admins</Arn>
    </Group>
    <Users>
       <member>
          <Path>/division_abc/subdivision_xyz/</Path>
          <UserName>Bob</UserName>
          <UserId>AIDACKCEVSQ6C2EXAMPLE</UserId>
          <Arn>arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/Bob</Arn>
       </member>
       <member>
          <Path>/division_abc/subdivision_xyz/</Path>
          <UserName>Susan</UserName>
          <UserId>AIDACKCEVSQ6C2EXAMPLE</UserId>
          <Arn>arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/Susan</Arn>
       </member>
    </Users>
    <IsTruncated>false</IsTruncated>
 </GetGroupResult>
 <ResponseMetadata>
    <RequestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</RequestId>
 </ResponseMetadata>
</GetGroupResponse>
`

var GetGroupPolicyExample = `
<GetGroupPolicyResponse>
 <GetGroupPolicyResult>
    <GroupName>Admins</GroupName>
    <PolicyName>AdminRoot</PolicyName>
    <PolicyDocument>{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}</PolicyDocument>
 </GetGroupPolicyResult>
 <ResponseMetadata>
    <RequestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</RequestId>
 </ResponseMetadata>
</GetGroupPolicyResponse>
`

var GetInstanceProfileExample = `
<GetInstanceProfileResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
<GetInstanceProfileResult>
  <InstanceProfile>
    <InstanceProfileId>AIPAD5ARO2C5EXAMPLE3G</InstanceProfileId>
    <Roles>
      <member>
        <Path>/application_abc/component_xyz/</Path>
        <Arn>arn:aws:iam::123456789012:role/application_abc/component_xyz/S3Access</Arn>
        <RoleName>S3Access</RoleName>
        <AssumeRolePolicyDocument>{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}</AssumeRolePolicyDocument>
        <CreateDate>2012-05-09T15:45:35Z</CreateDate>
        <RoleId>AROACVYKSVTSZFEXAMPLE</RoleId>
      </member>
    </Roles>
    <InstanceProfileName>Webserver</InstanceProfileName>
    <Path>/application_abc/component_xyz/</Path>
    <Arn>arn:aws:iam::123456789012:instance-profile/application_abc/component_xyz/Webserver</Arn>
    <CreateDate>2012-05-09T16:11:10Z</CreateDate>
  </InstanceProfile>
</GetInstanceProfileResult>
<ResponseMetadata>
  <RequestId>37289fda-99f2-11e1-a4c3-27EXAMPLE804</RequestId>
</ResponseMetadata>
</GetInstanceProfileResponse>
`

var GetRoleExample = `
<GetRoleResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
<GetRoleResult>
  <Role>
    <Path>/application_abc/component_xyz/</Path>
    <Arn>arn:aws:iam::123456789012:role/application_abc/component_xyz/S3Access</Arn>
    <RoleName>S3Access</RoleName>
    <AssumeRolePolicyDocument>{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}</AssumeRolePolicyDocument>
    <CreateDate>2012-05-08T23:34:01Z</CreateDate>
    <RoleId>AROADBQP57FF2AEXAMPLE</RoleId>
  </Role>
</GetRoleResult>
<ResponseMetadata>
  <RequestId>df37e965-9967-11e1-a4c3-270EXAMPLE04</RequestId>
</ResponseMetadata>
</GetRoleResponse>
`

var GetRolePolicyExample = `
<GetRolePolicyResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
<GetRolePolicyResult>
  <PolicyName>S3AccessPolicy</PolicyName>
  <RoleName>S3Access</RoleName>
  <PolicyDocument>{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:*"],"Resource":["*"]}]}</PolicyDocument>
</GetRolePolicyResult>
<ResponseMetadata>
  <RequestId>7e7cd8bc-99ef-11e1-a4c3-27EXAMPLE804</RequestId>
</ResponseMetadata>
</GetRolePolicyResponse>
`

var AccountAliasesExample = `
<ListAccountAliasesResponse>
<ListAccountAliasesResult>
  <IsTruncated>false</IsTruncated>
  <AccountAliases>
    <member>foocorporation</member>
    <member>barcorporation</member>
  </AccountAliases>
</ListAccountAliasesResult>
<ResponseMetadata>
  <RequestId>c5a076e9-f1b0-11df-8fbe-45274EXAMPLE</RequestId>
</ResponseMetadata>
</ListAccountAliasesResponse>
`

var GroupPoliciesExample = `
<ListGroupPoliciesResponse>
 <ListGroupPoliciesResult>
    <PolicyNames>
       <member>AdminRoot</member>
       <member>KeyPolicy</member>
    </PolicyNames>
    <IsTruncated>false</IsTruncated>
 </ListGroupPoliciesResult>
 <ResponseMetadata>
    <RequestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</RequestId>
 </ResponseMetadata>
</ListGroupPoliciesResponse>
`

var GroupsForUserExample = `
<ListGroupsForUserResponse>
 <ListGroupsForUserResult>
    <Groups>
       <member>
          <Path>/</Path>
          <GroupName>Admins</GroupName>
          <GroupId>AGPACKCEVSQ6C2EXAMPLE</GroupId>
          <Arn>arn:aws:iam::123456789012:group/Admins</Arn>
       </member>
    </Groups>
    <IsTruncated>false</IsTruncated>
 </ListGroupsForUserResult>
 <ResponseMetadata>
    <RequestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</RequestId>
 </ResponseMetadata>
</ListGroupsForUserResponse>
`

var InstanceProfilesExample = `
<ListInstanceProfilesResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
<ListInstanceProfilesResult>
  <IsTruncated>false</IsTruncated>
  <InstanceProfiles>
    <member>
      <InstanceProfileId>AIPACIFN4OZXG7EXAMPLE</InstanceProfileId>
      <Roles>
        <member>
          <Path>/application_abc/component_xyz/</Path>
          <Arn>arn:aws:iam::123456789012:role/application_abc/component_xyz/S3Access</Arn>
          <RoleName>S3Access</RoleName>
          <AssumeRolePolicyDocument>{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}</AssumeRolePolicyDocument>
          <CreateDate>2012-05-09T15:45:35Z</CreateDate>
          <RoleId>AROACVSVTSZYK3EXAMPLE</RoleId>
        </member>
      </Roles>
      <InstanceProfileName>Database</InstanceProfileName>
      <Path>/application_abc/component_xyz/</Path>
      <Arn>arn:aws:iam::123456789012:instance-profile/application_abc/component_xyz/Database</Arn>
      <CreateDate>2012-05-09T16:27:03Z</CreateDate>
    </member>
    <member>
      <InstanceProfileId>AIPACZLSXM2EYYEXAMPLE</InstanceProfileId>
      <Roles/>
      <InstanceProfileName>Webserver</InstanceProfileName>
      <Path>/application_abc/component_xyz/</Path>
      <Arn>arn:aws:iam::123456789012:instance-profile/application_abc/component_xyz/Webserver</Arn>
      <CreateDate>2012-05-09T16:27:11Z</CreateDate>
    </member>
  </InstanceProfiles>
</ListInstanceProfilesResult>
<ResponseMetadata>
  <RequestId>fd74fa8d-99f3-11e1-a4c3-27EXAMPLE804</RequestId>
</ResponseMetadata>
</ListInstanceProfilesResponse>
`

var InstanceProfilesForRoleExample = `
<ListInstanceProfilesForRoleResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
<ListInstanceProfilesForRoleResult>
  <IsTruncated>false</IsTruncated>
  <InstanceProfiles>
    <member>
      <InstanceProfileId>AIPACZLS2EYYXMEXAMPLE</InstanceProfileId>
      <Roles>
        <member>
          <Path>/application_abc/component_xyz/</Path>
          <Arn>arn:aws:iam::123456789012:role/application_abc/component_xyz/S3Access</Arn>
          <RoleName>S3Access</RoleName>
          <AssumeRolePolicyDocument>{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}</AssumeRolePolicyDocument>
          <CreateDate>2012-05-09T15:45:35Z</CreateDate>
          <RoleId>AROACVSVTSZYK3EXAMPLE</RoleId>
        </member>
      </Roles>
      <InstanceProfileName>Webserver</InstanceProfileName>
      <Path>/application_abc/component_xyz/</Path>
      <Arn>arn:aws:iam::123456789012:instance-profile/application_abc/component_xyz/Webserver</Arn>
      <CreateDate>2012-05-09T16:27:11Z</CreateDate>
    </member>
  </InstanceProfiles>
</ListInstanceProfilesForRoleResult>
<ResponseMetadata>
  <RequestId>6a8c3992-99f4-11e1-a4c3-27EXAMPLE804</RequestId>
</ResponseMetadata>
</ListInstanceProfilesForRoleResponse>
`

var RolePoliciesExample = `
<ListRolePoliciesResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
<ListRolePoliciesResult>
  <PolicyNames>
    <member>CloudwatchPutMetricPolicy</member>
    <member>S3AccessPolicy</member>
  </PolicyNames>
  <IsTruncated>false</IsTruncated>
</ListRolePoliciesResult>
<ResponseMetadata>
  <RequestId>8c7e1816-99f0-11e1-a4c3-27EXAMPLE804</RequestId>
</ResponseMetadata>
</ListRolePoliciesResponse>
`

var RolesExample = `
<ListRolesResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
<ListRolesResult>
  <IsTruncated>false</IsTruncated>
  <Roles>
    <member>
      <Path>/application_abc/component_xyz/</Path>
      <Arn>arn:aws:iam::123456789012:role/application_abc/component_xyz/S3Access</Arn>
      <RoleName>S3Access</RoleName>
      <AssumeRolePolicyDocument>{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}</AssumeRolePolicyDocument>
      <CreateDate>2012-05-09T15:45:35Z</CreateDate>
      <RoleId>AROACVSVTSZYEXAMPLEYK</RoleId>
    </member>
    <member>
      <Path>/application_abc/component_xyz/</Path>
      <Arn>arn:aws:iam::123456789012:role/application_abc/component_xyz/SDBAccess</Arn>
      <RoleName>SDBAccess</RoleName>
      <AssumeRolePolicyDocument>{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}</AssumeRolePolicyDocument>
      <CreateDate>2012-05-09T15:45:45Z</CreateDate>
      <RoleId>AROAC2ICXG32EXAMPLEWK</RoleId>
    </member>
  </Roles>
</ListRolesResult>
<ResponseMetadata>
  <RequestId>20f7279f-99ee-11e1-a4c3-27EXAMPLE804</RequestId>
</ResponseMetadata>
</ListRolesResponse>
`

var UserPoliciesExample = `
<ListUserPoliciesResponse>
 <ListUserPoliciesResult>
    <PolicyNames>
       <member>AllAccessPolicy</member>
       <member>KeyPolicy</member>
    </PolicyNames>
    <IsTruncated>false</IsTruncated>
 </ListUserPoliciesResult>
 <ResponseMetadata>
    <RequestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</RequestId>
 </ResponseMetadata>
</ListUserPoliciesResponse>
`

var UsersExample = `
<ListUsersResponse>
 <ListUsersResult>
    <Users>
       <member>
          <Path>/division_abc/subdivision_xyz/engineering/</Path>
          <UserName>Andrew</UserName>
          <UserId>AID2MAB8DPLSRHEXAMPLE</UserId>
          <Arn>arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/engineering/Andrew</Arn>
       </member>
       <member>
          <Path>/division_abc/subdivision_xyz/engineering/</Path>
          <UserName>Jackie</UserName>
          <UserId>AIDIODR4TAW7CSEXAMPLE</UserId>
          <Arn>arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/engineering/Jackie</Arn>
       </member>
    </Users>
    <IsTruncated>false</IsTruncated>
 </ListUsersResult>
 <ResponseMetadata>
    <RequestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</RequestId>
 </ResponseMetadata>
</ListUsersResponse>
`
