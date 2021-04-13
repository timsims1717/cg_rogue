package run

var CurrentRun encounterRun

type encounterRun struct {
	Level int
}

func StartRun() {
	CurrentRun = encounterRun{}
}
