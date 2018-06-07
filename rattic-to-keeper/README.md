Migration from Rattic to Keeper
===============================

This subproject will document the migration from the in-house Rattic installation to hosted Keeper

Rationale
---------

Described in [ADR #9](https://github.com/ministryofjustice/cloud-platform/tree/master/architecture-decision-record)

Rattic
------

 - Installed on a local OS instance, isolated via VPN
 - Has a simpleAPI, a curl -H 'Authorization:ApiKey user.name:aaabbbccc' -H 'Accept: application/json' https://rattic/api/v1/cred/108/ gets you a secret
 - This might be a quick way to export all the data without going via SQL, see https://rattic/api/v1/?format=json
 - API and web calls are shown the same in the audit log, so we can't easily see who accessed programatically - might break automation we are not aware of
 - Some developer design decisions should make our task easier: "We didn't include encryption in the application at all. We also tried to make the ACL system as simple as possible, passwords belong to a single group, and users can be in any number of groups."

Keeper
------

 - Managed service, hosted in AWS, publicly accesible at https://keepersecurity.com
 - Google auth integration documented at https://support.google.com/a/answer/7374103?hl=en, requires org 'super admin' privs; Keeper SSO architectural view documented in https://keepersecurity.com/assets/pdf/KeeperSSOConnect_v11.pdf
 - API documented mostly as part of the [Commander CLI client](https://github.com/Keeper-Security/Commander)
 - API has limited number of calls: create shared folder, import secrets, delete everything (?!); important missing call is "create teams" https://github.com/Keeper-Security/Commander/issues/44
 - Client doesn't have a clear or intuitive enough way of dealing with diffs, secret creation invoked multiple times with the same json results in multiple similar secrets
 - There is no clear documentation on why a Folder is different than a Shared Folder and how they are dealt with by the API
 - The "zero-knowledge" design of the platform makes some calls (esp search) quite slow

