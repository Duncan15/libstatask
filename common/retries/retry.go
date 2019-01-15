package retries

////RetryError customize error type for retry mechanism
//type RetryError struct {
//	desc string
//	cause error
//}
//
//func (err *RetryError)Error() string {
//	return fmt.Sprintf("%s, cause:%v", err.desc, err.cause)
//}
////NewRetryError create new retry error
//func NewRetryError(desc string, err error) *RetryError {
//	return &RetryError{
//		desc:desc,
//		cause:err,
//	}
//}
//Retry basic retry method
func Retry(count int, op func() error) error {
	var err error
	for i := 0; i < count; i++ {
		if err = op(); err == nil {
			return nil
		}
	}
	return err
}
