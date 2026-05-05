package proxy

// DialerPlan describes which connector defaults must be enforced.
type DialerPlan struct {
	UseIAMAuthN  bool
	UsePrivateIP bool
}

// DefaultDialerPlan enforces the security baseline for this feature:
// IAM database authentication and private IP connectivity.
func DefaultDialerPlan() DialerPlan {
	return DialerPlan{
		UseIAMAuthN:  true,
		UsePrivateIP: true,
	}
}
