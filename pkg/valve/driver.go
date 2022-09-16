package valve

import (
	"log"
	"sync"
	"time"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

type Driver interface {
	SwitchOn() error
	SwitchOff() error
	Status() ValveStatus
}

func New(pin1, pin2 string) Driver {
	r := raspi.NewAdaptor()

	r1 := gpio.NewRelayDriver(r, pin1)
	r2 := gpio.NewRelayDriver(r, pin2)

	return &driver{r1: r1, r2: r2}
}

type driver struct {
	r1 *gpio.RelayDriver
	r2 *gpio.RelayDriver

	m      sync.Mutex
	status ValveStatus
}

type ValveStatus bool

const (
	ValveOpen  ValveStatus = true
	ValveClose ValveStatus = false
)

func (d *driver) SwitchOn() error {
	log.Println("Switching on...")
	return d.switchRelay(d.r1, ValveOpen)
}

func (d *driver) SwitchOff() error {
	log.Println("Switching off...")
	return d.switchRelay(d.r2, ValveClose)
}

func (d *driver) Status() ValveStatus {
	return d.status
}

func (d *driver) switchRelay(r *gpio.RelayDriver, next ValveStatus) error {
	d.m.Lock()
	defer d.m.Unlock()

	if err := d.neutral(); err != nil {
		return err
	}

	if err := r.On(); err != nil {
		return err
	}

	time.Sleep(300 * time.Millisecond)

	if err := r.Off(); err != nil {
		return err
	}
	d.status = next
	return nil
}

func (d *driver) neutral() error {
	if err := d.r1.Off(); err != nil {
		return err
	}

	if err := d.r2.Off(); err != nil {
		return err
	}

	return nil
}
