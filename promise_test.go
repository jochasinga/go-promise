package promise

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreatePromise(t *testing.T) {
	p := New(func() (interface{}, error) {
		<-time.After(100)
		return "hello", nil
	})
	assert.NotNil(t, p)
	assert.IsType(t, p, (*Promise)(nil))
	assert.Implements(t, (*Thenable)(nil), p)
}

func TestPromiseResolved(t *testing.T) {
	p := New(func() (interface{}, error) {
		return "hello", nil
	})
	assert.IsType(t, p, (*Promise)(nil))
	p.Then(func(result interface{}) {
		assert.Equal(t, result, "hello")
	})
}

func TestPromiseRejected(t *testing.T) {
	p := New(func() (interface{}, error) {
		return nil, errors.New("Too bad")
	})
	p.Then(func(result interface{}) {
		assert.Nil(t, result)
	}, func(err error) {
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "Too bad")
	})
}

func TestPromiseResolveWithDelay(t *testing.T) {
	a := time.Now()
	_ = New(func() (interface{}, error) {
		<-time.After(100)
		return "hello", nil
	}).Then(func(result interface{}) {
		assert.WithinDuration(t, a, time.Now(), 150*time.Millisecond)
	})
}

func TestPromiseIsConcurrent(t *testing.T) {
	assert := assert.New(t)
	var res interface{}
	var _err error
	_ = New(func() (interface{}, error) {
		<-time.After(100)
		return struct{}{}, nil
	}).Then(func(result interface{}) {
		res = result
	}, func(err error) {
		_err = err
	})

	a := 1 + 2
	b := "hello " + "world"
	assert.Equal(a, 3)
	assert.Equal(b, "hello world")

	// the main routine does not wait
	assert.Nil(res)
	assert.Nil(_err)
}

func TestCreateNewPromiseFromChannel(t *testing.T) {
	assert := assert.New(t)
	resultC := make(chan interface{})
	p := From(resultC)
	assert.IsType(p, (*Promise)(nil))
	assert.NotNil(p)
	assert.NotNil(p.resolved)
	assert.Nil(p.rejected)
}

func TestCreateNewPromiseFromChannels(t *testing.T) {
	assert := assert.New(t)
	result := make(chan interface{})
	err := make(chan error)
	p := From(result, err)
	p.Then(func(result interface{}) {
		assert.Equal(result, "hello")
	})
	assert.NotNil(p)
	assert.IsType(p, (*Promise)(nil))
	assert.NotNil(p.resolved)
	assert.NotNil(p.rejected)
}

func TestConvertPromiseToChannels(t *testing.T) {
	p1 := New(func() (interface{}, error) {
		return "hello", nil
	})
	result, err := p1.To()
	select {
	case r := <-result:
		assert.Equal(t, r, "hello")
	case <-err:
		assert.Fail(t, "Error should not be emitted")
	}

	p2 := New(func() (interface{}, error) {
		return nil, errors.New("Too bad")
	})
	result, err = p2.To()
	select {
	case <-result:
		assert.Fail(t, "Result should not be emitted")
	case e := <-err:
		assert.Equal(t, e.Error(), "Too bad")
	}
}
