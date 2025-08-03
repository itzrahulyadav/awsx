package util

import _ "github.com/aws/aws-sdk-go-v2/aws"

// ContainsProtocol checks if the protocol matches the target or is "all" (-1)
func ContainsProtocol(protocol *string, target string) bool {
	return protocol != nil && (*protocol == target || *protocol == "-1")
}