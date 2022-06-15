package rabbitmq

import "fmt"

type ConnectionError struct {
	err error
}

func (c *ConnectionError) Error() string {
	return fmt.Sprintf("Failed to connect. Error - %v", c.err)
}

type OpenChannelError struct {
	err error
}

func (c *OpenChannelError) Error() string {
	return fmt.Sprintf("Failed to open a channel. Error - %v", c.err)
}

type DeclareQueueError struct {
	name string
	err  error
}

func (d *DeclareQueueError) Error() string {
	return fmt.Sprintf("Failed to declare queue(%s). Error - %v", d.name, d.err)
}

type ConsumeError struct {
	name string
	err  error
}

func (c *ConsumeError) Error() string {
	return fmt.Sprintf("Failed to register a consumer(%s). Error - %v", c.name, c.err)
}
