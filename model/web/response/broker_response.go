package response

import "backend/model/entity"

type GetBrokerResponse struct {
	Message string           `json:"message"`
	Data    []BrokerResponse `json:"data"`
}

type BrokerResponse struct {
	Code string `json:"broker_code"`
	Name string `json:"name"`
}

func ToBrokerResponse(broker *entity.Broker) BrokerResponse {
	return BrokerResponse{
		Code: broker.Broker_Code,
		Name: broker.Broker_Code + " - " + broker.Broker_Name,
	}
}

func ToBrokerResponses(listBroker []entity.Broker) []BrokerResponse {
	brokerResponses := make([]BrokerResponse, 0, len(listBroker))
	for _, broker := range listBroker {
		brokerResponses = append(brokerResponses, ToBrokerResponse(&broker))
	}
	return brokerResponses
}
