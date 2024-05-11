# Spyglass: Magnifying the status of almost everything

## Concepts

### Target
Target of a monitoring operation.

Fields:
- Id (in the form /Parent1/Parent2/TargetAlias)
- Description: Description of the target
- Url: URL of the target 
- Status (0 - 100, 0 displays as red, 1-99 yellow, 100 green)
- StatusDescription (String)
- Critical (boolean/yes, no): Indicates that the interface should show RED status of the parent if it is not in 100.
- Monitor (Gets the status of the target, when using a pull strategy, optional)
  - MonitorId (Id of the monitor)
  - Params (Monitor parameters)
    - ParamName: Value
    - ParamName: "${secret:secretName}"

### Monitor
Monitors a target. Returns a value from 0 to 100.

Fields:
- Id: Monitor Id
- Type: Docker (docker run container), K8S (k8s run container), LambdaFunction (for AWS Lambdas), AzureFunction (for Azure Functions)
- Schedule: Frequency described in Cron format. (Default: */30 * * * * - every 30 minutes)
- Definition: Monitor definition, depending on the type.
  - Docker, K8S:
    - Image: Image name
    - Entrypoint: Entry point
    - Env: List of environment variables
      - Name: param_name
        Value: "${param}"
  - LambdaFunction:
    - LambdaArn: ARN of lambda function
    - Event: Full JSON of the event info to pass to the function.
  - AzureFunction: 
    - AzureFunc: URL of azure function
    - Body: Full JSON of the body to pass to the function.

### Permissions

By default, all resources created has:
- Owners  = [user] (User who created the resource)
  - Note: Only the owners can change Owners, Readers and Writers fields, as well as delete the root target or monitor, role or user.
- Readers = [] (empty means all users are readers)
  - Note: Readers have read only access to the existing data.
- Writers = [user] (User who created the resource)
  - For targets, users can create, modify and delete targets under the specified target (not the root target)
  - For monitors, users can modify target configuration.
  - For users, users can modify the password of other users.
  - For roles, users can assign/revoke the roles to other users.
This can be changed at creation time, specifying Readers, Writers and Owners list on the resource definition.

### User

A single user on the system. Once created, a user group named after the user is created as well. A third party Identity
provider will be configured accordingly to provide user validation and registration using the UI.

Fields:
- UserEmail: Serves as email and user account ID.
- FullName: Optional full name of the user.

Special user:
- admin: Full administrator user. Cannot be deleted. Password set in service backend configuration.

### User groups

A user group is created by a user.

Fields:
- Name: Name of the user group.
- Users: List of users (by UserEmail) belonging to the group.

Special user groups (cannot be deleted):
- users (all users are members)
- admins (only admin users are members, configured on the backend)
