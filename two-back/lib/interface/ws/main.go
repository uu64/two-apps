package ws

import (
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

// Disconnect disconnects the connection to the specified user
func Disconnect(svc *apigatewaymanagementapi.ApiGatewayManagementApi, endpoint string, connectionID string) {
	svc.Endpoint = endpoint
	svc.DeleteConnection(&apigatewaymanagementapi.DeleteConnectionInput{
		ConnectionId: &connectionID,
	})
}

// Send send a message to the specified users
func Send(svc *apigatewaymanagementapi.ApiGatewayManagementApi, endpoint string, connectionIDs []string, message string) {
	svc.Endpoint = endpoint

	for _, id := range connectionIDs {
		svc.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: &id,
			Data:         []byte(message),
		})
	}
}
