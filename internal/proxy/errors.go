package proxy

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/auth"
)

// MissingCloudSQLOrIAPRolesError indicates startup-time permission gaps.
// Missing roles are provided as likely/target roles for DevOps to grant.
type MissingCloudSQLOrIAPRolesError struct {
	MissingRoles []string
}

func (e MissingCloudSQLOrIAPRolesError) Error() string {
	return "missing Cloud SQL / IAP permissions"
}

type PermissionCheckUnavailableError struct {
	Cause error
}

func (e PermissionCheckUnavailableError) Error() string {
	return "permission check unavailable"
}

type InstanceNotFoundOrNotAccessibleError struct {
	Cause error
}

func (e InstanceNotFoundOrNotAccessibleError) Error() string {
	return "instance not found or not accessible"
}

func defaultLikelyMissingRoles() []string {
	return []string{
		"roles/cloudsql.client",
		"roles/iap.tunnelResourceAccessor",
	}
}

// UserFacingError converts known proxy errors into actionable, human-readable messages.
// It intentionally avoids leaking low-level details such as stack traces or secret material.
func UserFacingError(err error) string {
	if err == nil {
		return ""
	}

	if errors.Is(err, auth.ErrMissingPrincipalEmail) {
		return "Unable to resolve the authenticated user email required for connection instructions. Re-run authentication to obtain the needed identity scopes."
	}

	var denied MissingCloudSQLOrIAPRolesError
	if errors.As(err, &denied) {
		missing := denied.MissingRoles
		if len(missing) == 0 {
			missing = defaultLikelyMissingRoles()
		}

		roles := strings.Join(missing, " and ")
		return fmt.Sprintf(
			"Missing required permissions to connect to Cloud SQL via IAP. Please contact DevOps to grant %s.",
			roles,
		)
	}

	var unavailable PermissionCheckUnavailableError
	if errors.As(err, &unavailable) {
		return "Permission check unavailable due to network/API uncertainty. Please verify connectivity and retry."
	}

	var notFound InstanceNotFoundOrNotAccessibleError
	if errors.As(err, &notFound) {
		return "Instance not found or not accessible. Verify the --instance value format project:region:instance and your access."
	}

	return fmt.Sprintf("Startup failed: %v", err)
}
