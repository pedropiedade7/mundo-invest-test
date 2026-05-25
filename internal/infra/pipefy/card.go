package pipefy

import "strconv"

const (
	DefaultPipeID = "123456"

	CreateCardMutation = `mutation CreateClientCard($pipeId: ID!, $clienteNome: String!, $clienteEmail: String!, $tipoSolicitacao: String!, $valorPatrimonio: String!) {
  createCard(input: {
    pipe_id: $pipeId
    fields_attributes: [
      { field_id: "cliente_nome", field_value: $clienteNome }
      { field_id: "cliente_email", field_value: $clienteEmail }
      { field_id: "tipo_solicitacao", field_value: $tipoSolicitacao }
      { field_id: "valor_patrimonio", field_value: $valorPatrimonio }
    ]
  }) {
    card {
      id
    }
  }
}`
)

type CreateCardInput struct {
	ClientName     string
	ClientEmail    string
	RequestType    string
	PortfolioValue int64
}

func (c *Client) BuildCreateCardRequest(input CreateCardInput) GraphQLRequest {
	return GraphQLRequest{
		Query: CreateCardMutation,
		Variables: map[string]any{
			"pipeId":          c.pipeID,
			"clienteNome":     input.ClientName,
			"clienteEmail":    input.ClientEmail,
			"tipoSolicitacao": input.RequestType,
			"valorPatrimonio": formatMoneyField(input.PortfolioValue),
		},
	}
}

func formatMoneyField(value int64) string {
	return strconv.FormatInt(value, 10)
}
