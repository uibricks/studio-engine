package utils

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"github.com/uibricks/studio-engine/internal/pkg/constants"
	"github.com/uibricks/studio-engine/internal/pkg/logger"
	"github.com/uibricks/studio-engine/internal/pkg/request"
	"github.com/uibricks/studio-engine/internal/pkg/types"
	"google.golang.org/grpc/metadata"
	"reflect"
	"time"
	"unsafe"
)

type contextData struct {
	Values      map[interface{}]interface{}
	HasCancel   bool
	HasDeadline bool
	Deadline    time.Time
}

type SerializeOpts struct {
	RetainCancel    bool
	RetainDeadline  bool
	IgnoreFunctions bool
}

var (
	cancelCtxType reflect.Type
	timeCtxType   reflect.Type
)

func init() {
	cancelCtx, c := context.WithCancel(context.Background())
	c()
	cancelCtxType = reflect.ValueOf(cancelCtx).Elem().Type()
	timeCtx, c := context.WithDeadline(context.Background(), time.Time{})
	c()
	timeCtxType = reflect.ValueOf(timeCtx).Elem().Type()
}

// PrepareSharableContext - This will read the context data to share between the services in case
// if there is no way to share the context directly.
// Receiving service will need to recontruct this byte to a context to use
func PrepareSharableContext(ctx *context.Context, opts ...SerializeOpts) ([]byte, error) {
	buf := new(bytes.Buffer)
	e := gob.NewEncoder(buf)

	s := contextData{
		HasCancel: false,
		Deadline:  time.Time{},
	}

	// *** enforcing timeout in context if not present***
	serialized := buildMap(*ctx, s)
	if !serialized.HasCancel {
		ctxTemp, _ := context.WithTimeout(*ctx, constants.Saga_Timeout_In_Minutes*time.Minute)
		ctx = &ctxTemp
		serialized = buildMap(*ctx, s)
	}

	if len(opts) > 0 {

		if !opts[0].RetainCancel {
			serialized.HasCancel = false
		}
		if !opts[0].RetainDeadline {
			serialized.HasDeadline = false
		}
	}

	err := e.Encode(serialized)
	return buf.Bytes(), err
}

func buildMap(ctx context.Context, s contextData) contextData {

	rs := reflect.ValueOf(ctx).Elem()
	if rs.Type() == reflect.ValueOf(context.Background()).Elem().Type() {
		return s
	}

	rf := rs.FieldByName("key")
	if !rf.IsValid() {
		if rs.Type() == cancelCtxType {
			s.HasCancel = true
		}
		if rs.Type() == timeCtxType {
			deadline := rs.FieldByName("deadline")
			deadline = reflect.NewAt(deadline.Type(), unsafe.Pointer(deadline.UnsafeAddr())).Elem()
			deadlineTime := deadline.Convert(reflect.TypeOf(time.Time{})).Interface().(time.Time)
			if s.HasDeadline && deadlineTime.Before(s.Deadline) {
				s.Deadline = deadlineTime
			} else {
				s.HasDeadline = true
				s.Deadline = deadlineTime
			}
		}
	}

	parent := rs.FieldByName("Context")
	if parent.IsValid() && !parent.IsNil() {
		return buildMap(parent.Interface().(context.Context), s)
	}
	return s
}

func ReconstructContext(sharedContext types.SharedContext) context.Context {

	if sharedContext.ContextData == nil {
		return context.Background()
	}

	var ctxData []byte
	b, _ := json.Marshal(sharedContext.ContextData)
	json.Unmarshal(b, &ctxData)
	ctx, _, err := deserializeCtx(ctxData)

	if err != nil {
		logger.Sugar.Errorf("Error white reconstructing the context - %v", err)
		return context.Background()
	}

	// Partial reconstruction of context from project service
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs(request.RequestIDKey, sharedContext.RequestID))

	return ctx

}

// DeserializeCtx inflates the byte-array output of SerializeCtx into a context and optional CancelFunc
// The options specified during serialization dictate whether CancelFunc is non-nil
func deserializeCtx(ser []byte) (context.Context, context.CancelFunc, error) {
	dec := gob.NewDecoder(bytes.NewReader(ser))
	data := contextData{}
	err := dec.Decode(&data)
	if err != nil {
		return context.Background(), func() {}, err
	}

	// make a new base context
	ctx := context.Background()

	// get back the values
	for key, val := range data.Values {
		ctx = context.WithValue(ctx, key, val)
	}

	// get back the cancel
	var c context.CancelFunc
	if data.HasCancel {
		ctx, c = context.WithCancel(ctx)
	}

	// get back the deadline
	if data.HasDeadline {
		ctx, c = context.WithDeadline(ctx, data.Deadline)
	}

	return ctx, c, nil
}
