package pipefy

const UpdateCardFieldsMutation = `mutation UpdateClientCard($cardId: ID!, $status: [String!], $prioridade: [String!]) {
  status: updateCardField(input: {
    card_id: $cardId
    field_id: "status_cliente"
    new_value: $status
  }) {
    success
    card {
      id
    }
  }
  prioridade: updateCardField(input: {
    card_id: $cardId
    field_id: "prioridade_cliente"
    new_value: $prioridade
  }) {
    success
    card {
      id
    }
  }
}`

func (c *Client) BuildUpdateCardFieldsRequest(cardID, status, prioridade string) GraphQLRequest {
	return GraphQLRequest{
		Query: UpdateCardFieldsMutation,
		Variables: map[string]any{
			"cardId":     cardID,
			"status":     []string{status},
			"prioridade": []string{prioridade},
		},
	}
}
