package errorx

var (
	// CommonErrors is a namespace for general purpose errors designed for universal use.
	// These errors should typically be used in opaque manner, implying no handing in user code.
	// When handling is required, it is best to use custom error types with both standard and custom traits.
	CommonErrors = NewNamespace("common")

	IllegalArgument      = CommonErrors.NewType("illegal_argument")
	IllegalState         = CommonErrors.NewType("illegal_state")
	IllegalFormat        = CommonErrors.NewType("illegal_format")
	InitializationFailed = CommonErrors.NewType("initialization_failed")
	DataUnavailable      = CommonErrors.NewType("data_unavailable")
	UnsupportedOperation = CommonErrors.NewType("unsupported_operation")
	RejectedOperation    = CommonErrors.NewType("rejected_operation")
	Interrupted          = CommonErrors.NewType("interrupted")
	AssertionFailed      = CommonErrors.NewType("assertion_failed")
	InternalError        = CommonErrors.NewType("internal_error")
	ExternalError        = CommonErrors.NewType("external_error")
	ConcurrentUpdate     = CommonErrors.NewType("concurrent_update")
	TimeoutElapsed       = CommonErrors.NewType("timeout", Timeout())
	NotImplemented       = UnsupportedOperation.NewSubtype("not_implemented")
	UnsupportedVersion   = UnsupportedOperation.NewSubtype("version")
)
