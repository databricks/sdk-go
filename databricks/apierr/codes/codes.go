// Package codes defines error codes for API errors and their retry semantics.
package codes

// Code is a numeric code for an error.
type Code uint32

const (
	// Unknown indicates an error that cannot be classified.
	//
	// This code might be used for malformed error responses or error responses
	// using an error code that cannot be mapped to a code in this package.
	Unknown Code = iota

	// OK indicates the operation completed successfully.
	OK

	// Canceled indicates the operation was canceled (typically by the caller).
	Canceled

	// InvalidArgument indicates the client specified an invalid argument.
	//
	// By contrast with FailedPrecondition, InvalidArgument indicates arguments
	// that are problematic regardless of the state of the system. For example,
	// a malformed request parameter.
	InvalidArgument

	// DeadlineExceeded means the operation expired before completion.
	//
	// For operations that modify the state of the system, this error may be
	// returned even if the operation has completed successfully. For
	// example, a successful response from a server could have been delayed
	// long enough for the deadline to expire.
	DeadlineExceeded

	// NotFound means a requested entity (e.g. a resource or a file) was
	// not found.
	NotFound

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	AlreadyExists

	// PermissionDenied indicates the caller does not have permission to
	// execute the specified operation.
	//
	// This is different from an error returned when the user has exhausted
	// some resource (e.g. too many requests) which is a ResourceExhausted
	// error.
	PermissionDenied

	// ResourceExhausted indicates some resource has been exhausted, perhaps
	// a per-user quota, or perhaps the entire file system is out of space.
	ResourceExhausted

	// FailedPrecondition indicates the operation was rejected because the
	// system is not in a state required for the operation's execution.
	// For example, directory to be deleted may be non-empty, an rmdir
	// operation is applied to a non-directory, etc.
	FailedPrecondition

	// Aborted indicates the operation was aborted, typically due to a
	// concurrency issue like sequencer check failures, transaction aborts,
	// etc.
	Aborted

	// OutOfRange means the operation was attempted past the valid range.
	// E.g., seeking or reading past end of file.
	//
	// Unlike InvalidArgument, this error indicates a problem that may
	// be fixed if the system state changes. For example, a 32-bit file
	// system will generate InvalidArgument if asked to read at an
	// offset that is not in the range [0,2^32-1], but it will generate
	// OutOfRange if asked to read from an offset past the current
	// file size.
	//
	// There is a fair bit of overlap between FailedPrecondition and
	// OutOfRange. We recommend using OutOfRange (the more specific
	// error) when it applies so that callers who are iterating through
	// a space can easily look for an OutOfRange error to detect when
	// they are done.
	OutOfRange

	// Unimplemented indicates the operation is not implemented or not
	// supported/enabled in this service.
	Unimplemented

	// Internal indicates an internal error. This means some invariants
	// expected by the underlying system have been broken. If you see
	// this error, something is very broken.
	Internal

	// Unavailable indicates the service is currently unavailable.
	//
	// This is most likely a transient condition and may be corrected by
	// retrying. Though, this might not always be safe to retry if the
	// operation is non-idempotent.
	//
	// The Databricks SDK will generally automatically retry the request
	// with a backoff when encountering this error.
	Unavailable

	// DataLoss indicates unrecoverable data loss or corruption.
	DataLoss

	// Unauthenticated indicates the request does not have valid
	// authentication credentials for the operation.
	Unauthenticated

	nCodes // dummy code to count the number of codes in testing
)

func (c Code) String() string {
	switch c {
	case OK:
		return "OK"
	case Canceled:
		return "CANCELLED"
	case InvalidArgument:
		return "INVALID_ARGUMENT"
	case DeadlineExceeded:
		return "DEADLINE_EXCEEDED"
	case NotFound:
		return "NOT_FOUND"
	case AlreadyExists:
		return "ALREADY_EXISTS"
	case PermissionDenied:
		return "PERMISSION_DENIED"
	case ResourceExhausted:
		return "RESOURCE_EXHAUSTED"
	case FailedPrecondition:
		return "FAILED_PRECONDITION"
	case Aborted:
		return "ABORTED"
	case OutOfRange:
		return "OUT_OF_RANGE"
	case Unimplemented:
		return "UNIMPLEMENTED"
	case Internal:
		return "INTERNAL"
	case Unavailable:
		return "UNAVAILABLE"
	case DataLoss:
		return "DATA_LOSS"
	case Unauthenticated:
		return "UNAUTHENTICATED"
	default:
		return "UNKNOWN"
	}
}

// FromString converts a string representation of an error code to its
// corresponding Code value. If the string does not match any known code,
// Unknown is returned.
func FromString(s string) Code {
	switch s {
	case "OK":
		return OK
	case "CANCELLED":
		return Canceled
	case "INVALID_ARGUMENT":
		return InvalidArgument
	case "DEADLINE_EXCEEDED":
		return DeadlineExceeded
	case "NOT_FOUND":
		return NotFound
	case "ALREADY_EXISTS":
		return AlreadyExists
	case "PERMISSION_DENIED":
		return PermissionDenied
	case "RESOURCE_EXHAUSTED":
		return ResourceExhausted
	case "FAILED_PRECONDITION":
		return FailedPrecondition
	case "ABORTED":
		return Aborted
	case "OUT_OF_RANGE":
		return OutOfRange
	case "UNIMPLEMENTED":
		return Unimplemented
	case "INTERNAL":
		return Internal
	case "UNAVAILABLE":
		return Unavailable
	case "DATA_LOSS":
		return DataLoss
	case "UNAUTHENTICATED":
		return Unauthenticated
	default:
		return Unknown
	}
}
