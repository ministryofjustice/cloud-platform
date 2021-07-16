# 24. Use SSO for CP Team to access AWS Console

Date: 2021-03-31

## Status

🤔 Proposed

## Proposal

We will enable AWS Single Sign-On (SSO) log-in to Cloud Platform AWS accounts:

* GitHub will be the Identity Provider, meaning users enter their GitHub username/password/2FA during login
* Auth0 will be the broker
* AWS SSO will be configured to allow the `webops` GitHub team access to two administrative IAM Roles: for read-only access and read+write access

This login method will be used by the Cloud Platform team to access Cloud Platform's AWS account, with both console and command-line access.

This sign-in has a time-out, depending on the role:

* read-only role - 12 hours
* read+write role - 1 hour

This method of authentication replaces the team's existing:

* signing into the AWS Console using current methods:
  * SAML login - using GitHub, Auth0 and AWS SAML Federation
  * Console password and MFA (for people who SAML didn't work), for their IAM User
* access key for their IAM User

We'll transition CP team to use SSO, switch off SAML auth and delete their IAM Users.

## Context

Administrative access to the Cloud Platform's AWS account is used by the Cloud Platform team's engineers. This is so they can:

* access the AWS Console, to understand the AWS infrastructure and occasionally to do administrative actions (such as emergency restoring a backup EBS)
* authenticate on the command-line, to run 'terraform plan' while developing or reviewing terraform changes, for example

Existing administrative access has been provided to these people like this:

* signing into the AWS Console using current methods:
  * SAML login - using GitHub (as Identity Provider), Auth0 (broker) and AWS SAML Federation
  * Console password and MFA (for people who SAML didn't work), for their IAM User
* access key for their IAM User

With existing SAML login, Auth0 acts as an OIDC/SAML broker, with Rules that ensures the user is a member of the `ministryofjustice` GitHub organization, and sets the AWS Role that can be assumed, corresponding to their GitHub teams. A couple of team members have not been able to get this method to work, and have been given username/password/MFA access instead.

* Ownership of the is not with the CP team any more - we need PR approval of the AWS organization root account to change the GitHub group(s) and their permissions. However admin permissions don't change much. And we can describe the [change in terraform](https://github.com/ministryofjustice/aws-root-account/blob/89c75d19663ce7d65588c3ee91cde724caba68c5/terraform/sso-admin-account-assignments.tf#L77-L86) so it really is just a PR approval we'd need.

## Factors in the decision

Pros:

* SSO makes it easy to login to different roles and accounts without having to store different credentials. This allows us to have a read-only role, which the team uses by default, which is safer than using full permissions all the time - a "least privilege" security benefit. And it makes it easier to administer other AWS accounts that use SSO, like DSD. In comparison, the SAML solution doesn't appear to have a way to select different roles, and selecting between AWS accounts means using different login URLs, which is less convenient.
* SAML solution is flakey. A couple of people haven't been able to get the SAML method working, falling back to user a username/password. This is not ideal as the team have to remember to remove them in their leaver process, rather than being able to rely on the standard IT-managed leaver process that removes them from the GitHub organization.
* With SSO we don't have to create and remove IAM Users, like we did with SAML.
* Temporary AWS credentials are more secure than long-lived ones, and whilst this is already the case for console access, SSO provides this for command-line access too

Cons:

* The SSO login is only temporary, so login needs to be easy and get the team used to doing it. TODO What is existing SAML login timeout? TODO Can AWS Vault help with this?

## Security

TODO: Is this 2FA? Can/should we need that? GitHub login requires 2FA.
TODO: readonly / timeout
TODO: revoking creds of engineers that leave
TODO: credentials in plain text on hard disk
TODO: update runbooks with use of creds? e.g. export AWS_PROFILE=moj-cp
TODO: do we have auditability of AWS actions, since everyone is sharing a role?

## How it will work

The proposed login works like this:

1. Users go to https://moj.awsapps.com/start#/ which redirects you to login.
2. Auth0 log-in selection page. User selects "Continue with GitHub".
3. GitHub login page. If already logged into GitHub, it's skipped.
4. AWS page to select account, role and either console or command-line access.
    For CP users their roles could be: AdministratorAccess or Readonly
5. Select either:

    **AWS Console** The user is logged in with a federated role, such as: `AWSReservedSSO_AdministratorAccess_bffb036637b86425/davidread@digital.justice.gov.uk`.
    After the creds time-out, the console prompts you to "relogin", and a single click you're back in.

    OR

    **CLI** It shows commands to paste into a command-line, to set a temporary access key.
    After the time-out, they need to redo this from step 1.

Users may prefer the added convenience of using AWS Vault, which works like this:

1. User sets up a profile in their `.aws/config`:

        [profile cp]
        sso_start_url=https://moj.awsapps.com/start
        sso_region=eu-west-2
        sso_account_id=754256621582
        sso_role_name=AdministratorAccess

2. For **AWS Console**:

        aws-vault login cp

    OR for **CLI**:

        aws-vault exec cp -- aws s3 ls

    You have to click a link in the browser to login.

3. They are prompted for their storage password to save/retrieve a token (this is the case for MacOS keychain, and may differ for other secret storage mechanisms)

4. The browser will check you're logged into GitHub and invite you to click through an AWS login page. Both store cookies, so this step rarely needs doing after the first time.

Access for both Console and CLI expires after the time-out. The user re-establishes the connection from the command again (step 2). TODO: check this is correct
