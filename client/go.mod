module mandatory-exercise-2/client

go 1.17

replace mandatory-exercise-2/service => ../service

replace mandatory-exercise-2/utils => ../utils

require (
	google.golang.org/grpc v1.42.0
	mandatory-exercise-2/service v0.0.0-00010101000000-000000000000
	mandatory-exercise-2/utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/protobuf v1.5.0 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)
