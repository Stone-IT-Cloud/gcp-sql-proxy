package proxy

import (
	"context"
	"errors"
	"net/http"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

// instancesGetFn is a test injection point to simulate sqladmin.instances.get outcomes.
var instancesGetFn = func(ctx context.Context, httpClient *http.Client, project, instance string) error {
	svc, err := sqladmin.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}
	_, err = svc.Instances.Get(project, instance).Do()
	return err
}

// SetTestInstancesGet overrides the instances-get function for unit tests.
// It returns a restore function.
func SetTestInstancesGet(fn func(ctx context.Context, httpClient *http.Client, project, instance string) error) (restore func()) {
	prev := instancesGetFn
	if fn != nil {
		instancesGetFn = fn
	}
	return func() {
		instancesGetFn = prev
	}
}

// VerifyAccess performs a pre-flight permission check via SQL Admin instances.get.
// It fails fast with typed errors mapped by UserFacingError.
func VerifyAccess(ctx context.Context, client *http.Client, instance string) error {
	parts, err := ParseCloudSQLInstance(instance)
	if err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if client == nil {
		return &PermissionCheckUnavailableError{Cause: errors.New("missing http client")}
	}

	err = instancesGetFn(ctx, client, parts.Project, parts.Instance)
	if err == nil {
		return nil
	}

	var gerr *googleapi.Error
	if errors.As(err, &gerr) {
		switch gerr.Code {
		case http.StatusForbidden:
			return MissingCloudSQLOrIAPRolesError{MissingRoles: defaultLikelyMissingRoles()}
		case http.StatusNotFound:
			return InstanceNotFoundOrNotAccessibleError{Cause: err}
		default:
			// Any other SQL Admin error is treated as indeterminate for UX purposes.
			return &PermissionCheckUnavailableError{Cause: err}
		}
	}

	// Transport errors, timeouts, and unknown failures.
	return &PermissionCheckUnavailableError{Cause: err}
}
