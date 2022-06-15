package saga

import (
	"errors"
	"fmt"
	"github.com/uibricks/studio-engine/internal/pkg/logger"
	"github.com/uibricks/studio-engine/internal/pkg/saga/storage"
	"golang.org/x/net/context"
	"reflect"
	"runtime"
	"strings"
	"time"
)

const LogPrefix = "saga_"

// Saga presents current execute transaction.
// A Saga constituted by small sub-transactions.
type Saga struct {
	context  context.Context
	ctxValue reflect.Value
	sec      *ExecutionCoordinator
}

func (s *Saga) Storage() storage.Storage {
	return s.sec.storage
}

func (s *Saga) startSaga() {
	log := &Log{
		Type: SagaStart,
		Time: time.Now(),
	}
	s.Storage().AppendLog(log)
}

// ExecSub executes a sub-transaction for given subTxID(which define in SEC initialize) and arguments.
// it returns current Saga.
func (s *Saga) ExecSub(subTxID string, args ...interface{}) (*Saga, []reflect.Value) {
	subTxDef := s.sec.MustFindSubTxDef(subTxID)
	log := &Log{
		Type:       ActionStart,
		SubTxID:    subTxID,
		Time:       time.Now(),
		Params:     MarshalParam(s.sec, args),
		Compensate: subTxDef.compensate != reflect.ValueOf(nil),
	}
	s.Storage().AppendLog(log)

	params := make([]reflect.Value, 0, len(args)+1)
	//params = append(params, reflect.ValueOf(s.context))

	for _, arg := range args {
		params = append(params, reflect.ValueOf(arg))
	}

	//ctxMap[subTxID] = params[0]

	result := subTxDef.action.Call(params)
	if isReturnError(result) {
		s.Abort()
		return s, result
	}

	log = &Log{
		Type:    ActionEnd,
		SubTxID: subTxID,
		Time:    time.Now(),
	}
	// err = LogStorage().AppendLog(s.logID, log.mustMarshal())
	s.Storage().AppendLog(log)

	return s, result
}

// EndSaga finishes a Saga's execution.
func (s *Saga) EndSaga() {
	log := &Log{
		Type: SagaEnd,
		Time: time.Now(),
	}
	s.Storage().AppendLog(log)
	// sets s.data to empty, to prevent execution of previous saga's rollback
	s.Storage().Cleanup()
}

// Abort stop and compensate to rollback to start situation.
// This method will stop continue sub-transaction and do Compensate for executed sub-transaction.
// SubTx will call this method internal.
func (s *Saga) Abort() {
	//logs, err := s.Storage().Lookup(s.logID)
	//if err != nil {
	//	panic("Abort Panic")
	//}
	alog := &Log{
		Type: SagaAbort,
		Time: time.Now(),
	}
	// err = LogStorage().AppendLog(s.logID, alog.mustMarshal())
	s.Storage().AppendLog(alog)

	logs := s.Storage().Lookup()

	for i := len(logs) - 1; i >= 0; i-- {
		logData := logs[i]
		log := logData.(*Log) //mustUnmarshalLog(logData)
		if log.Type == ActionStart && log.Compensate {
			if err := s.compensate(*log); err != nil {
				panic("Compensate Failure..")
			}
		}
	}
}

func (s *Saga) compensate(tlog Log) error {
	clog := &Log{
		Type:    CompensateStart,
		SubTxID: tlog.SubTxID,
		Time:    time.Now(),
	}
	// err := LogStorage().AppendLog(s.logID, clog.mustMarshal())
	s.Storage().AppendLog(clog)

	//typ := s.sec.MustFindParamName(reflect.ValueOf(s.context).Type())

	args := UnmarshalParam(s.sec, tlog.Params)
	// fmt.Println(args)

	params := make([]reflect.Value, 0, len(args)+1)
	// params = append(params, reflect.ValueOf(s.context))
	params = append(params, args...)

	subDef := s.sec.MustFindSubTxDef(tlog.SubTxID)

	compFuncName := runtime.FuncForPC(subDef.compensate.Pointer()).Name()
	compFuncName = strings.Split(compFuncName, ".")[len(strings.Split(compFuncName, "."))-1]
	compFuncName = strings.Split(compFuncName, "-")[0]

	params[0] = s.ctxValue

	// ToDo - lock??
	result := subDef.compensate.Call(params)

	clog = &Log{
		Type:    CompensateEnd,
		SubTxID: tlog.SubTxID,
		Time:    time.Now(),
	}
	// err = LogStorage().AppendLog(s.logID, clog.mustMarshal())
	s.Storage().AppendLog(clog)

	if isReturnError(result) {
		logger.WithContext(s.context).Errorf("Failed to execute compensatory function : %s, error : %v", compFuncName, result[0])
		return errors.New(fmt.Sprintf("Failed to execute compensatory function : %s, error : %v", compFuncName, result[0]))
	}

	return nil
}

func isReturnError(result []reflect.Value) bool {
	if len(result) > 0 && !result[0].IsNil() {
		return true
	}
	return false
}
