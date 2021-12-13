build:
		cd PVS_API
		mkdir	-p functions
		GOBIN=${PWD}/functions go install ./functions/...