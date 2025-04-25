package entity

type SendopostRunNotificator interface {
	NotifyListeners(msg string)
	AddListener(listener Listener)
	RemoveListener(listener Listener)
	StopNotificate()
}
