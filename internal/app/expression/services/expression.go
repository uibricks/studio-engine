package expression

import (
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	"github.com/uibricks/studio-engine/internal/app/expression/utils"
	expressionpb "github.com/uibricks/studio-engine/internal/pkg/proto/expression"
	"github.com/uibricks/studio-engine/internal/pkg/rabbitmq"
)

type Server struct {
	Rabbit     *rabbitmq.Rabbit
	ReplyQueue amqp.Queue
	Channel    *rabbitmq.Channel
}

func ProvideExpressionServer(rabbit *rabbitmq.Rabbit, replyQ amqp.Queue, ch *rabbitmq.Channel) *Server {
	return &Server{
		Rabbit:     rabbit,
		ReplyQueue: replyQ,
		Channel:    ch,
	}
}

// EvalExpression is a helper function for conversion of expressions to their evaluated-data-form
func (s *Server) EvalExpression(_ context.Context, expressionReq *expressionpb.EvalExpressionRequest) (*expressionpb.EvalExpressionResponse, error) {
	expressionMetaData := map[string]utils.ExpressionMD{}
	b, _ := json.Marshal(expressionReq.GetExpressions())
	if err := json.Unmarshal(b, &expressionMetaData); err != nil {
		return nil, err
	}

	expressions := make([]utils.Expression, 0)
	b, _ = json.Marshal(expressionReq.GetExpressionMenu())
	if err := json.Unmarshal(b, &expressions); err != nil {
		return nil, err
	}

	//data := `[{"id":"1","name":"Mira Lioma","address":{"street":"4000 Edison Ave","city":"Sacramento","state":"CA","zip":"95821"},"classes":[{"number":"09","teacher":"Grazyna","students":[{"id":"1","name":"Rahul","marks":[{"sub":"math","mark":77}]},{"id":"2","name":"John","marks":[{"sub":"math","mark":66},{"sub":"science","mark":71}]},{"id":"3","name":"Anwesha","marks":[{"sub":"math","mark":88},{"sub":"science","mark":91}]}]},{"number":"10","teacher":"Kristin Baker","students":[{"id":"1","name":"Benjamin","marks":[{"sub":"math","mark":55},{"sub":"science","mark":61}]},{"id":"2","name":"Oliver","marks":[{"sub":"math","mark":44},{"sub":"science","mark":51}]},{"id":"3","name":"William","marks":[{"sub":"math","mark":33},{"sub":"science","mark":41}]}]}]},{"id":"2","name":"Arcade Middle School","address":{"street":"3500 Edison Ave","city":"Sacramento","state":"CA","zip":"95821"},"classes":[{"number":"07","teacher":"Julio Alvarez","students":[{"id":"1","name":"Rohit","marks":[{"sub":"math","mark":81},{"sub":"science","mark":77}]},{"id":"2","name":"Peggy","marks":[{"sub":"math","mark":71},{"sub":"science","mark":66}]},{"id":"3","name":"Sara","marks":[{"sub":"math","mark":61},{"sub":"science","mark":55}]}]},{"number":"08","teacher":"Matthew","students":[{"id":"1","name":"Levi","marks":[{"sub":"math","mark":51},{"sub":"science","mark":44}]},{"id":"2","name":"Mateo","marks":[{"sub":"math","mark":41},{"sub":"science","mark":33}]},{"id":"3","name":"David","marks":[{"sub":"math","mark":91},{"sub":"science","mark":88}]}]}]}]`

	res, err := utils.GetExpressionVal(expressionMetaData, expressions, expressionReq.GetData())
	//fmt.Println(utils.FlattenJson(expressionReq.GetData(), []string{}))
	//res, err := utils.GetExpressionVal(expressionMetaData, expressions, data, keys)
	b, _ = json.Marshal(res)

	return &expressionpb.EvalExpressionResponse{ExpressionResult: string(b)}, err
}
