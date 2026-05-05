package proxy

// DialerPlan describes which connector defaults must be enforced.
type DialerPlan struct {
	UseIAMAuthN  bool
	UsePrivateIP bool
}

// DefaultDialerPlan enforces IAM auth and defaults to public IP unless specified.
func DefaultDialerPlan(usePrivateIP bool) DialerPlan {
	return DialerPlan{
		UseIAMAuthN:  true,
		UsePrivateIP: usePrivateIP,
	}
}
