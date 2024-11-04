package interface_smtp

type Ismtp interface {
	SendNotificationWithEmailOtp(otp int, recieverEmail string, recieverName string) error
	SendRestPasswordEmailOtp(otp int, recieverEmail string) error
}
